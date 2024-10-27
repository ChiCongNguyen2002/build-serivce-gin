package resp

import "fmt"

const (
	LangEN = "EN"
	LangVI = "VI"
)

const (
	StatusOK = 1
)

const (
	// common error
	ErrSystem = 1000 + iota
	ErrNotFound
	ErrDataInvalid
	ErrStoreDataFailed
	ErrENVInvalid
)

const (
	// Logic error
	ErrHandler1LogicA = 2000 + iota
	ErrHandleLogicParseData
	ErrProfileTerritoryInvalid
	ErrProfileGamesInvalid
	ErrHandleAmountInvalid
	ErrHandleCurrencyInvalid
	ErrHandleTier
	ErrHandleProfileInvalid
	ErrHandleBalanceNotFound
	ErrHandleTxTypeNotFound
	ErrHandleInvalidMonth
	ErrHandleOffsetInvalid
	ErrHandleLimitInvalid
	ErrHandleProfileIdNotFound
	ErrHandleTierClient
	ErrHandleOrderNotFound
)

const (
	//auth error
	ErrAuth = 4000 + iota
)

const (
	// Rate limit error
	ErrTooManyRequest = 5000 + iota
	ErrRateLimit
)

const (
	// call external services error
	ErrCallExternalService = 9000 + iota
)

// TODO get mapping error-message in file config json
var (
	mappingErrEN = map[int64]mappingError{
		StatusOK:           {StatusOK, "Success"},
		ErrSystem:          {ErrSystem, "System error"},
		ErrAuth:            {ErrAuth, "Authenticate fail"},
		ErrNotFound:        {ErrNotFound, "Data not found"},
		ErrDataInvalid:     {ErrDataInvalid, "Data invalid"},
		ErrStoreDataFailed: {ErrStoreDataFailed, "There was an error during the data saving process"},
		ErrENVInvalid:      {ErrENVInvalid, "ENV invalid"},

		ErrHandler1LogicA:       {ErrHandler1LogicA, "ErrHandler1LogicA"},
		ErrHandleLogicParseData: {ErrHandleLogicParseData, "parse data fail"},
		ErrCallExternalService:  {ErrCallExternalService, "There was an error during the call external services"},

		ErrTooManyRequest:          {ErrTooManyRequest, "Too many request"},
		ErrRateLimit:               {ErrRateLimit, "Rate limit"},
		ErrHandleAmountInvalid:     {ErrHandleAmountInvalid, "Amount invalid"},
		ErrHandleCurrencyInvalid:   {ErrHandleCurrencyInvalid, "Currency invalid"},
		ErrHandleTier:              {ErrHandleTier, "Error retrieve tier"},
		ErrHandleProfileInvalid:    {ErrHandleProfileInvalid, "Profile invalid"},
		ErrHandleBalanceNotFound:   {ErrHandleBalanceNotFound, "Balance not found"},
		ErrHandleTxTypeNotFound:    {ErrHandleTxTypeNotFound, "Tx type not found"},
		ErrHandleInvalidMonth:      {ErrHandleInvalidMonth, "Invalid month"},
		ErrHandleOffsetInvalid:     {ErrHandleOffsetInvalid, "Offset invalid"},
		ErrHandleLimitInvalid:      {ErrHandleLimitInvalid, "Limit invalid"},
		ErrHandleProfileIdNotFound: {ErrHandleProfileIdNotFound, "Profile Id not found"},
		ErrHandleTierClient:        {ErrHandleTierClient, "Error retrieve tier client"},
		ErrHandleOrderNotFound:     {ErrHandleOrderNotFound, "Order not found"},
		ErrProfileTerritoryInvalid: {ErrProfileTerritoryInvalid, "Profile territory invalid"},
		ErrProfileGamesInvalid:     {ErrProfileGamesInvalid, "Profile games invalid"},
	}
)

type mappingError struct {
	ErrorCode int64  `json:"errorCode"`
	Message   string `json:"message"`
}

func (c *mappingError) Error() string {
	return fmt.Sprintf("%v - %v", c.ErrorCode, c.Message)
}

func (c *mappingError) GetErrCode() int64 {
	return c.ErrorCode
}

func (c *mappingError) GetErrMsg() string {
	return c.Message
}

func GetMappingError(errCode int64, lang string) mappingError {
	switch lang {
	case LangEN:
		if err, ok := mappingErrEN[errCode]; ok {
			return err
		}
	default:
		if err, ok := mappingErrEN[errCode]; ok {
			return err
		}
	}
	return mappingErrEN[ErrSystem]
}
