package logger

import (
	"build-service-gin/common/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"io"
	"os"
	"reflect"
	"sync"
)

var (
	loggerInstance *Logger
	mu             sync.RWMutex
	keyEncrypt     *string
)

const (
	KeyServiceName = "service_name"
	KeyLogId       = "log_id"
	KeyFileError   = "file_error"
)

func InitLog(serviceName string) {
	mu.Lock()
	defer mu.Unlock()
	if loggerInstance != nil {
		return
	}

	if serviceName == "" {
		log.Fatal().Msg("service name is empty")
	}

	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	multiWriter := io.MultiWriter(file, os.Stdout)

	lg := log.Output(multiWriter).With().Str(KeyServiceName, serviceName).Logger()

	loggerInstance = &Logger{lg}
}

func SetKeyEncrypt(key string) {
	if keyEncrypt == nil {
		keyEncrypt = &key
	}
}

func GetLogger() *Logger {
	mu.RLock()
	defer mu.RUnlock()
	uid := uuid.NewString()
	lg := loggerInstance.logger.With().Str(KeyLogId, uid).Logger()
	return &Logger{lg}
}

func SetGinReqEncrLog(c *gin.Context, req interface{}) {
	if keyEncrypt == nil || *keyEncrypt == "" {
		return
	}

	ctx := c.Request.Context()
	if req != nil {
		if newReq, err := utils.StructEncryptTagInterface(req, *keyEncrypt, utils.TagNameEncrypt, utils.TagValEncrypt); err == nil {
			if str, err := utils.AnyToString(newReq); err == nil {
				ctx = context.WithValue(ctx, utils.KeyRequestBody, str)
				c.Request = c.Request.WithContext(ctx)
			}
		}
	}
}

func SetGinRespEncrLog(c *gin.Context, resp interface{}) {
	if keyEncrypt == nil || *keyEncrypt == "" {
		return
	}

	ctx := c.Request.Context()

	if resp == nil {
		return
	}

	v := reflect.ValueOf(resp)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		if data := v.FieldByName("Data"); data.IsValid() {
			if data.Kind() == reflect.Ptr {
				data = data.Elem()
			}

			if newRes, err := utils.InterfaceEncryptTagInterface(data.Interface(), *keyEncrypt, utils.TagNameEncrypt, utils.TagValEncrypt); err == nil {
				if str, err := utils.AnyToString(newRes); err == nil {
					ctx = context.WithValue(ctx, utils.KeyResponseBody, str)
					c.Request = c.Request.WithContext(ctx)
				}
			}
		}
	}
}

func Encrypt[T any](data T) (T, error) {
	if keyEncrypt == nil || *keyEncrypt == "" {
		return data, nil
	}

	switch v := interface{}(data).(type) {
	case string:
		res, err := utils.Encrypt(v, *keyEncrypt)
		if err != nil {
			return data, err
		}

		var result interface{} = res
		return result.(T), nil
	case *string:
		res, err := utils.Encrypt(*v, *keyEncrypt)
		if err != nil {
			return data, err
		}

		var result interface{} = &res
		return result.(T), nil
	}

	return utils.InterfaceEncryptTag(data, *keyEncrypt, utils.TagNameEncrypt, utils.TagValEncrypt)
}

func EncryptInterface(data interface{}) (interface{}, error) {
	if keyEncrypt == nil || *keyEncrypt == "" {
		return data, nil
	}

	switch v := data.(type) {
	case string:
		return utils.Encrypt(v, *keyEncrypt)
	case *string:
		return utils.Encrypt(*v, *keyEncrypt)
	}

	return utils.InterfaceEncryptTagInterface(data, *keyEncrypt, utils.TagNameEncrypt, utils.TagValEncrypt)
}
