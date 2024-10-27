package utils

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/segmentio/ksuid"
	"github.com/tidwall/gjson"
	"golang.org/x/text/unicode/norm"
)

const (
	KeyTraceInfo               = "trace_info"
	HeaderXMeProfile           = "X-Me-Profile"
	KeyRateLimit               = "key-limit"
	KeyRequestBody             = "request_body"
	KeyResponseBody            = "response_body"
	JwtSub                     = "sub"
	JwtExp                     = "exp"
	KeyEchoContextRequestBody  = "echo_context_request_body"
	KeyEchoContextResponseBody = "echo_context_response_body"
	KeyMongoMultiConnName      = "mongo_multi_conn_name"
	KeyRegion                  = "X-Client-Region"
	KeySignature               = "signature"
	KeyXTicketId               = "X-Ticket-Id"
)

const (
	TagNameEncrypt = "encrypt"
	TagValEncrypt  = "true"
)

const (
	IASTypeService  = "SERVICE"
	IASTypeClient   = "CLIENT"
	IASTypePublic   = "PUBLIC"
	IASTypeInternal = "INTERNAL"
)

const (
	IASTokenExpireTypeLimited   = "LIMITED"
	IASTokenExpireTypeUnlimited = "UNLIMITED"
)

const (
	VGREncryptKey = "VGR_ENCRYPT_KEY"
)

const (
	XIASCode     = "X-IAS-Code"
	XRequestData = "X-Request-Data"
)

type TraceInfo struct {
	RequestID string `json:"request_id"`
}

func GetRequestIdByContext(ctx context.Context) *TraceInfo {
	value := ctx.Value(KeyTraceInfo)
	traceInfo, ok := value.(TraceInfo)
	if !ok {
		return nil
	}
	return &traceInfo
}

func CombineErrors(errs []error) error {
	var errRes error
	for _, err := range errs {
		errRes = combine2Err(errRes, err)
	}
	return errRes
}

func combine2Err(errSrc, err error) error {
	if errSrc == nil {
		return err
	}

	if err == nil {
		return errSrc
	}

	return fmt.Errorf("%s | %s", errSrc.Error(), err.Error())
}

func AnyToString(value any) (string, error) {
	if value == nil {
		return "", nil
	}

	if str, ok := value.(string); ok {
		return str, nil
	}

	if str, ok := value.([]byte); ok {
		return string(str), nil
	}

	byteValue, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(byteValue), nil
}

func FindKeysWithRegex(data, regex string) []string {
	keyMap := make(map[string]bool)

	re := regexp.MustCompile(regex)
	matches := re.FindAllStringSubmatch(data, -1)

	for _, match := range matches {
		key := strings.TrimSpace(match[1])
		keyMap[key] = true
	}

	keys := make([]string, 0, len(keyMap))
	for key := range keyMap {
		keys = append(keys, key)
	}

	return keys
}

// APIPathMatch compares two paths, return true if they match
// Path can contain parameters, for example:
//
//	pattern: /api/admin/v1/users/:id/role
//	path1: /api/admin/v1/users/123123/role
//	path2: /api/admin/v1/users/abcxyz/12332342/role
//
// path1 matches pattern, path2 does not match pattern
func APIPathMatch(pattern, path string) bool {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i, part := range patternParts {
		if part != pathParts[i] && !strings.HasPrefix(part, ":") {
			return false
		}
	}

	return true
}

// ConcatenateStructValues concatenates all struct values into a single string
func ConcatenateStructValues(item interface{}) (string, error) {
	v := reflect.ValueOf(item)
	typeOfS := v.Type()

	mapFileName := make(map[string]string)
	var tags []string
	for i := 0; i < v.NumField(); i++ {
		tag := typeOfS.Field(i).Tag.Get("json")
		tags = append(tags, tag)
		mapFileName[tag] = typeOfS.Field(i).Name
	}

	// Sort tag alphabetically
	sort.Strings(tags)

	// Concatenate values based on the sorted field names
	var concatenatedValues []string
	for _, tag := range tags {
		fieldValue := v.FieldByName(mapFileName[tag]).Interface()
		if strValue, err := AnyToString(fieldValue); err != nil {
			return "", err
		} else {
			concatenatedValues = append(concatenatedValues, strValue)
		}
	}

	return strings.Join(concatenatedValues, ""), nil
}

//func BuildMessageFromJsonByte(jsonByte []byte, messageFields []string, sep string) (string, error) {
//	if len(messageFields) == 0 {
//		return "", fmt.Errorf("messageFields is empty")
//	}
//
//	var msgVals []string
//	for _, fieldName := range messageFields {
//		value := gjson.GetBytes(jsonByte, fieldName).Value()
//		str, err := AnyToString(value)
//		if err != nil {
//			return "", err
//		}
//
//		msgVals = append(msgVals, str)
//	}
//
//	return strings.Join(msgVals, sep), nil
//}

// IsMn reports whether the rune has the mark, nonspacing property.
func IsMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

// RemoveDiacritics removes diacritics from the input string.
func RemoveDiacritics(input string) string {
	// Normalize the input to NFD form.
	normInput := norm.NFD.String(input)

	// Create a rune slice to hold the result.
	var result []rune

	// Iterate over the normalized input.
	for _, r := range normInput {
		// Append the rune to the result if it is not a combining mark.
		if !IsMn(r) {
			result = append(result, r)
		}
	}

	// Convert the rune slice back to a string.
	return string(result)
}

func NewContextWithRequestId(ctx context.Context) (context.Context, TraceInfo) {
	traceInfo := GetRequestIdByContext(ctx)
	if traceInfo == nil {
		traceInfo = &TraceInfo{RequestID: RandString()}
	}
	return context.WithValue(ctx, KeyTraceInfo, *traceInfo), *traceInfo
}

func NewContextBackgroundFromCtx(ctx context.Context) context.Context {
	newCtx := context.Background()

	traceInfo := GetRequestIdByContext(ctx)
	if traceInfo != nil {
		newCtx = context.WithValue(newCtx, KeyTraceInfo, *traceInfo)
	}

	return newCtx
}

func RandString() string {
	id := ksuid.New()
	return id.String()
}

func BuildMessageFromJsonByte(jsonByte []byte, messageFields []string, sep string) (string, error) {
	if len(messageFields) == 0 {
		return "", fmt.Errorf("messageFields is empty")
	}

	var msgVals []string
	for _, fieldName := range messageFields {
		value := gjson.GetBytes(jsonByte, fieldName).Value()
		str, err := AnyToString(value)
		if err != nil {
			return "", err
		}

		msgVals = append(msgVals, str)
	}

	return strings.Join(msgVals, sep), nil
}

func Encrypt(plaintext, secretKeyHex string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	secretKey, err := hex.DecodeString(secretKeyHex)
	if err != nil {
		return "", err
	}

	iv := secretKey[:aes.BlockSize]

	bPlaintext := PKCS5Padding([]byte(plaintext), aes.BlockSize)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, bPlaintext)

	// encode base64 and return
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertextBase64, secretKeyHex string) (string, error) {
	if ciphertextBase64 == "" {
		return "", nil
	}

	secretKey, err := hex.DecodeString(secretKeyHex)
	if err != nil {
		return "", err
	}

	ciphertextByte, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	iv := secretKey[:aes.BlockSize]

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertextByte, ciphertextByte)

	return string(PKCS5UnPadding(ciphertextByte)), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
