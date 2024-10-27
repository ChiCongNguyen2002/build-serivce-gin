package middlewares

import (
	"bufio"
	"build-service-gin/common/logger"
	"build-service-gin/common/utils"
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/sjson"
)

// AddExtraDataForRequestContext middleware to add extra data to the request context
func AddExtraDataForRequestContext(c *gin.Context) {
	reqID := c.Request.Header.Get("X-Request-ID")
	if reqID == "" {
		reqID = generateNewRequestID()
	}

	// Set request_id to request
	c.Request.Header.Set("X-Request-ID", reqID)

	// Set request_id to response
	c.Writer.Header().Set("X-Request-ID", reqID)

	// Set trace_info to context
	traceInfo := utils.TraceInfo{RequestID: reqID}
	ctx := c.Request.Context()
	ctxTraceInfo := context.WithValue(ctx, utils.KeyTraceInfo, traceInfo)
	c.Request = c.Request.WithContext(ctxTraceInfo)

	c.Next()
}

func generateNewRequestID() string {
	return uuid.New().String()
}

func Logging(c *gin.Context) {
	start := time.Now()
	req := c.Request

	// Request
	reqBody := []byte{}
	if req.Body != nil { // Read
		reqBody, _ = io.ReadAll(req.Body)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset

	if !json.Valid(reqBody) && len(reqBody) > 256 {
		reqBody = append(reqBody[:256], []byte("***")...)
	}

	// Response
	resBody := new(bytes.Buffer)
	writer := &bodyDumpResponseWriter{Writer: resBody, ResponseWriter: c.Writer}
	c.Writer = writer

	c.Next() // Proceed to the next handler

	// After handler processing
	res := c.Writer
	latency := time.Since(start)
	latencyInMs := float64(latency.Nanoseconds()) / 1000000.0
	statusCode := res.Status()
	method := req.Method
	path := req.URL.Path
	if path == "" {
		path = "/"
	}
	if strings.Contains(path, "/health") {
		return
	}

	requestBody := string(reqBody)
	responseBody := resBody.String()

	ctx := c.Request.Context()

	if newReqBody := ctx.Value(utils.KeyRequestBody); newReqBody != nil {
		if str, err := utils.AnyToString(newReqBody); err == nil {
			requestBody = str
		}
	}

	if newResBody := ctx.Value(utils.KeyResponseBody); newResBody != nil {
		if str, err := utils.AnyToString(newResBody); err == nil {
			if v, err := sjson.Set(responseBody, "data", str); err == nil {
				responseBody = v
			}
		}
	}

	log := logger.GetLogger().AddTraceInfoContextRequest(req.Context())

	var newLog logger.Logger
	newLog = *log

	var eventLog *logger.Event
	if statusCode >= 500 {
		eventLog = newLog.Error()
	} else {
		eventLog = newLog.Info()
	}

	eventLog.Str("method", method).
		Str("path", path).
		Str("ip", c.ClientIP()). // Use ClientIP() for Gin
		Str("user_agent", req.UserAgent()).
		Str("request_id", req.Header.Get("X-Request-ID")).
		Int("statusCode", statusCode).
		Float64("latency", latencyInMs).
		Interface("params", c.Request.URL.Query()).
		Str("request_body", requestBody).
		Str("response_body", responseBody).Msg("request income")
}

type bodyDumpResponseWriter struct {
	io.Writer
	gin.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	if n, err := w.Writer.Write(b); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
