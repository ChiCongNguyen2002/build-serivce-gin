package client

import (
	"build-service-gin/common/logger"
	"build-service-gin/common/utils"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	neturl "net/url"
	"os"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Client struct {
	*resty.Client
}

type (
	ResponseMiddlewareFunc func(client *resty.Client, response *resty.Response) error
	RequestMiddlewareFunc  func(client *resty.Client, request *resty.Request) error

	TypeKeyRequestBodyEncrypt  string
	TypeKeyResponseBodyEncrypt string
	TypeKeyQueryParamsEncrypt  string
)

const (
	KeyRestRequestBodyEncrypt  = "rest_request_body_encrypt"
	KeyRestResponseBodyEncrypt = "rest_response_body_encrypt"
	KeyRestQueryParamsEncrypt  = "rest_query_params_encrypt"

	KeyOffLogRespBody = "off_log_resp_body"
)

func NewClient(baseURL string, timeout time.Duration, maxRetry int, waitTime time.Duration, fns ...interface{}) *Client {
	var (
		fnResp             []ResponseMiddlewareFunc
		fnReq              []RequestMiddlewareFunc
		reqBodyEncrypt     []TypeKeyRequestBodyEncrypt
		respBodyEncrypt    []TypeKeyResponseBodyEncrypt
		queryParamsEncrypt []TypeKeyQueryParamsEncrypt
	)

	for _, fn := range fns {
		switch fn.(type) {
		case []ResponseMiddlewareFunc:
			fnResp = append(fnResp, fn.([]ResponseMiddlewareFunc)...)
		case []RequestMiddlewareFunc:
			fnReq = append(fnReq, fn.([]RequestMiddlewareFunc)...)
		case []TypeKeyRequestBodyEncrypt:
			reqBodyEncrypt = append(reqBodyEncrypt, fn.([]TypeKeyRequestBodyEncrypt)...)
		case []TypeKeyResponseBodyEncrypt:
			respBodyEncrypt = append(respBodyEncrypt, fn.([]TypeKeyResponseBodyEncrypt)...)
		case []TypeKeyQueryParamsEncrypt:
			queryParamsEncrypt = append(queryParamsEncrypt, fn.([]TypeKeyQueryParamsEncrypt)...)
		}
	}

	// get key encrypt from env
	keyEncrypt := os.Getenv(utils.VGREncryptKey)

	client := resty.New()
	client.SetBaseURL(baseURL)
	client.SetTimeout(timeout)
	client.SetRetryCount(maxRetry)
	client.SetRetryWaitTime(waitTime)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.EnableTrace()
	client.SetDisableWarn(true)

	client.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		ctx := request.Context()
		ctxNew, traceInfo := utils.NewContextWithRequestId(ctx)
		request.SetContext(ctxNew)
		request.SetHeader("X-Request-ID", traceInfo.RequestID) // Changed to Gin's style

		// get log
		log := logger.GetLogger().AddTraceInfoContextRequest(request.Context())

		reqBody := request.Body

		reqBaseURL := client.BaseURL + request.URL
		reqQueryParams := copyQueryParams(request.QueryParam)

		if keyEncrypt != "" {
			// encrypt request body fields
			var reqBodyEnc []string
			for _, be := range reqBodyEncrypt {
				reqBodyEnc = append(reqBodyEnc, string(be))
			}

			if result, err := encryptBodyFields(reqBody, reqBodyEnc, keyEncrypt); err != nil {
				log.Error().Err(err).Msg("encrypt request body fields error")
			} else {
				reqBody = result
			}

			if value, ok := ctx.Value(KeyRestRequestBodyEncrypt).([]string); ok {
				if result, err := encryptBodyFields(reqBody, value, keyEncrypt); err != nil {
					log.Error().Err(err).Msg("encrypt request body fields error")
				} else {
					reqBody = result
				}
			}

			// -----------------------------------------------

			// encrypt request query params
			var queryParamsEnc []string
			for _, qe := range queryParamsEncrypt {
				queryParamsEnc = append(queryParamsEnc, string(qe))
			}

			if result, err := encryptQueryParams(reqQueryParams, queryParamsEnc, keyEncrypt); err != nil {
				log.Error().Err(err).Msg("encrypt request query params error")
			} else {
				reqQueryParams = result
			}

			if value, ok := ctx.Value(KeyRestQueryParamsEncrypt).([]string); ok {
				if result, err := encryptQueryParams(reqQueryParams, value, keyEncrypt); err != nil {
					log.Error().Err(err).Msg("encrypt request query params error")
				} else {
					reqQueryParams = result
				}
			}
		}

		// log trace
		log.Info().
			Str("url", reqBaseURL).
			Interface("query_params", reqQueryParams).
			Str("method", request.Method).
			Interface("request_body", reqBody).
			Msg("rest api request")

		for _, fn := range fnReq {
			_ = fn(client, request)
		}

		return nil
	})

	// add interceptor to log
	client.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		ctx := response.Request.Context()

		// get log
		log := logger.GetLogger().AddTraceInfoContextRequest(ctx)

		// get trace request
		trace := response.Request.TraceInfo()

		logTrace := log.Info()
		if response.StatusCode() != http.StatusOK {
			logTrace = log.Error()
		}

		reqBody := response.Request.Body

		var respBody interface{}

		valOffLogRespBody := ctx.Value(KeyOffLogRespBody)
		value, ok := valOffLogRespBody.(bool)
		if ok && value && len(response.Body()) > 0 {
			respBody = "******"
		} else {
			respBody = response.Body()
		}

		reqBaseURL, reqQueryParams, err := parseURL(response.Request.URL)
		if err != nil {
			log.Error().Err(err).Msg("parse request url error")
		}

		if keyEncrypt != "" {
			// encrypt request body fields
			var reqBodyEnc []string
			for _, be := range reqBodyEncrypt {
				reqBodyEnc = append(reqBodyEnc, string(be))
			}

			if result, err := encryptBodyFields(reqBody, reqBodyEnc, keyEncrypt); err != nil {
				log.Error().Err(err).Msg("encrypt request body fields error")
			} else {
				reqBody = result
			}

			if value, ok := ctx.Value(KeyRestRequestBodyEncrypt).([]string); ok {
				if result, err := encryptBodyFields(reqBody, value, keyEncrypt); err != nil {
					log.Error().Err(err).Msg("encrypt request body fields error")
				} else {
					reqBody = result
				}
			}

			// -----------------------------------------------

			// encrypt response body fields
			var respBodyEnc []string
			for _, be := range respBodyEncrypt {
				respBodyEnc = append(respBodyEnc, string(be))
			}

			if result, err := encryptBodyFields(respBody, respBodyEnc, keyEncrypt); err != nil {
				log.Error().Err(err).Msg("encrypt response body fields error")
			} else {
				respBody = result
			}

			if value, ok := ctx.Value(KeyRestResponseBodyEncrypt).([]string); ok {
				if result, err := encryptBodyFields(respBody, value, keyEncrypt); err != nil {
					log.Error().Err(err).Msg("encrypt response body fields error")
				} else {
					respBody = result
				}
			}

			// -----------------------------------------------

			// encrypt request query params
			var queryParamsEnc []string
			for _, qe := range queryParamsEncrypt {
				queryParamsEnc = append(queryParamsEnc, string(qe))
			}

			if result, err := encryptQueryParams(reqQueryParams, queryParamsEnc, keyEncrypt); err != nil {
				log.Error().Err(err).Msg("encrypt request query params error")
			} else {
				reqQueryParams = result
			}

			if value, ok := ctx.Value(KeyRestQueryParamsEncrypt).([]string); ok {
				if result, err := encryptQueryParams(reqQueryParams, value, keyEncrypt); err != nil {
					log.Error().Err(err).Msg("encrypt request query params error")
				} else {
					reqQueryParams = result
				}
			}
		}

		// log trace
		logTrace.
			Str("url", reqBaseURL).
			Interface("query_params", reqQueryParams).
			Str("method", response.Request.Method).
			Str("status", response.Status()).
			Int("status_code", response.StatusCode()).
			Str("latency", response.Time().String()).
			Str("protocol", response.Proto()).
			Time("received_at", response.ReceivedAt()).
			Int64("size", response.Size()).
			Str("conn_time", trace.ConnTime.String()).
			Str("tcp_conn_time", trace.TCPConnTime.String()).
			Str("server_time", trace.ServerTime.String()).
			Str("response_time", trace.ResponseTime.String()).
			Interface("request_body", reqBody).
			Interface("response_body", respBody).
			Msg("rest api response")

		for _, fn := range fnResp {
			_ = fn(client, response)
		}

		return nil
	})

	return &Client{
		client,
	}
}

func encryptBodyFields(body interface{}, bodyFields []string, secretKey string) (interface{}, error) {
	if len(bodyFields) == 0 {
		return body, nil
	}

	if secretKey == "" {
		return body, nil
	}

	bodyStr, err := utils.AnyToString(body)
	if err != nil {
		return "", err
	}

	if bodyStr == "" {
		return body, nil
	}

	mappingBodyFields := make(map[string]bool)
	for _, field := range bodyFields {
		if _, ok := mappingBodyFields[field]; ok {
			continue
		}

		mappingBodyFields[field] = true

		value := gjson.Get(bodyStr, field).Value()
		if value == nil {
			continue
		}

		encryptedValue, err := utils.Encrypt(fmt.Sprintf("%v", value), secretKey)
		if err != nil {
			return "", err
		}

		bodyStr, err = sjson.Set(bodyStr, field, encryptedValue)
		if err != nil {
			return "", err
		}
	}

	return bodyStr, nil
}

func parseURL(_url string) (string, neturl.Values, error) {
	if _url == "" {
		return "", nil, nil
	}

	url, err := neturl.Parse(_url)
	if err != nil {
		return "", nil, err
	}

	path := url.Scheme + "://" + url.Host + url.Path
	queryParams := url.Query()
	return path, queryParams, nil
}

func encryptQueryParams(queryParams neturl.Values, fields []string, secretKey string) (neturl.Values, error) {
	if len(queryParams) == 0 {
		return queryParams, nil
	}

	result := make(neturl.Values)

	mappingFields := make(map[string]bool)
	for _, field := range fields {
		mappingFields[field] = true
	}

	for key, value := range queryParams {
		if _, ok := mappingFields[key]; ok {
			encryptedValue, err := utils.Encrypt(value[0], secretKey)
			if err != nil {
				return nil, err
			}

			result[key] = []string{encryptedValue}
		} else {
			result[key] = value
		}
	}

	return result, nil
}

func copyQueryParams(queryParams neturl.Values) neturl.Values {
	result := make(neturl.Values)
	for key, value := range queryParams {
		result[key] = value
	}
	return result
}

// ----------------------------- Example -----------------------------

//package main
//
//import (
//"gitlab.gt.vng.vn/vng-rewards/share/vgr-be-common/client"
//logger "gitlab.gt.vng.vn/vng-rewards/share/vgr-be-common/log"
//"context"
//"time"
//)
//
//func main() {
//	ctx := context.Background()
//	logger.InitLog("test")
//
//	url := "https://api.openbrewerydb.org/v1/breweries"
//
//	var key client.TypeKeyEncrypt = "94e553b0fc7e934278ca996d84a15cf135fdfc1375aa188936624da1aeb57544"
//	queryParamsEnc := []client.TypeKeyQueryParamsEncrypt{"by_city"}
//	respBodyEnc := []client.TypeKeyResponseBodyEncrypt{"#.id", "#.name"}
//
//	c := client.NewClient(url, 30*time.Second, 0, 1*time.Second, key, queryParamsEnc, respBodyEnc)
//
//	queryParamsEncValueCtx := []string{"by_state"}
//	ctx = context.WithValue(ctx, client.KeyRestQueryParamsEncrypt, queryParamsEncValueCtx)
//
//	respBodyEncValueCtx := []string{"#.brewery_type", "#.address_1"}
//	ctx = context.WithValue(ctx, client.KeyRestResponseBodyEncrypt, respBodyEncValueCtx)
//
//	queryParams := map[string]string{
//		"by_city":  "san_diego",
//		"by_state": "california",
//		"by_type":  "micro",
//		"per_page": "3",
//	}
//
//	_, err := c.R().SetContext(ctx).SetQueryParams(queryParams).Get("")
//	if err != nil {
//		logger.GetLogger().Error().Err(err).Msg("get error")
//		return
//	}
//}
