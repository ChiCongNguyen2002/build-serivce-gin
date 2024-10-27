package resp

import "fmt"

type Resp struct {
	ErrorCode   int64       `json:"errorCode"`
	Description string      `json:"description"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data"`
	Paging      *Paging     `json:"paging,omitempty"`
}

type CustomError struct {
	ErrorCode   int64  `json:"errorCode"`
	Description string `json:"description"`
}

func (c *CustomError) Error() string {
	return fmt.Sprintf("%v - %v", c.ErrorCode, c.Description)
}

func BuildErrorResp(errCode int64, description, lang string) Resp {
	err := GetMappingError(errCode, lang)
	if description == "" {
		description = err.Message
	}
	return Resp{ErrorCode: err.ErrorCode, Message: err.Message, Description: description}
}

func BuildSuccessResp(lang string, data interface{}) Resp {
	err := GetMappingError(StatusOK, lang)
	return Resp{Message: err.Message, Data: data, ErrorCode: StatusOK}
}
