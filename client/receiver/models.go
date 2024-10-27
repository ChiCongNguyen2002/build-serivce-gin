package receiver

type OrderResp struct {
	Status    string `json:"status"`
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
}
