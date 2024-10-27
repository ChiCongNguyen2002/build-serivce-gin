package domains

import "time"

type UserTransactionHistory struct {
	TransactionID        string     `json:"transactionID"`
	TransactionType      string     `json:"transactionType"`
	ProfileID            string     `json:"profileID"`
	Status               string     `json:"status"`
	PointAmount          int64      `json:"pointAmount"`
	PointType            int64      `json:"pointType"`
	TotalAmount          float64    `json:"totalAmount"`
	Currency             string     `json:"currency"`
	PaymentTransactionID string     `json:"paymentTransactionID"`
	Source               string     `json:"source"`
	SourceTime           *time.Time `json:"sourceTime"`
	SourceType           string     `json:"sourceType"`
	CreatedAt            *time.Time `json:"createdAt"`
	UpdatedAt            *time.Time `json:"updatedAt"`
}

type GetUserTransactionHistoryReq struct {
	Offset      int64  `json:"offset"`
	Limit       int64  `json:"limit"`
	ProfileID   string `json:"profileID"`
	TxType      string `json:"txType"`
	Status      string `json:"status"`
	RecentMonth int    `json:"recentMonth"`
}

type QueueMessage struct {
	RawData   string `json:"rawData"`
	Signature string `json:"signature"`
}

type PublishMessage struct {
	EventType string      `json:"eventType"`
	Data      interface{} `json:"data"`
}

type Order struct {
	OrderNumber string  `json:"orderNumber"`
	CreateTime  int64   `json:"createTime"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	VGAUserID   string  `json:"vgaUserId"`
	SourceCode  string  `json:"sourceCode"`
	SourceType  string  `json:"sourceType"`
}

type OrderPoint struct {
	OrderNumber string  `json:"orderNumber"`
	CreateTime  int64   `json:"createTime"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	VGAUserID   string  `json:"vgaUserId"`
	Region      string  `json:"region"`
}

type OrderMessage struct {
	SourceType string `json:"sourceType"`
	RawData    string `json:"value"`
	Signature  string `json:"sign"`
}

type OrderSuccessEvent struct {
	ProfileID            string     `json:"profile_id"`
	TotalAmount          float64    `json:"total_amount"`
	Currency             string     `json:"currency"`
	Status               string     `json:"status"`
	Region               string     `json:"region"`
	Source               string     `json:"source"`
	SourceTime           *time.Time `json:"source_time"`
	SourceType           string     `json:"source_type"`
	PaymentTransactionID string     `json:"payment_transaction_id"`
}

func (r *OrderSuccessEvent) BuildCreateOrderTransaction(req *OrderSuccessEvent) {
	r.TotalAmount = req.TotalAmount
	r.Currency = req.Currency
	r.Status = req.Status
	r.Source = req.Source
	r.Region = req.Region
	r.SourceType = req.SourceType
	r.SourceTime = req.SourceTime
	r.PaymentTransactionID = req.PaymentTransactionID
}

type EarnPointOrderEvent struct {
	TransactionID   string  `json:"transactionId"`
	ReferenceCode   string  `json:"referenceCode"`
	TransactionType string  `json:"transactionType"`
	Status          string  `json:"status"`
	Region          string  `json:"region"`
	PointAmount     int64   `json:"pointAmount"`
	PointType       int64   `json:"pointType"`
	TotalAmount     float64 `json:"totalAmount"`
	Currency        string  `json:"currency"`
	ProfileID       string  `json:"profileID"`
}
