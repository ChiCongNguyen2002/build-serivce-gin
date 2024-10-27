package models

import (
	"time"
)

type OrderSuccessEvent struct {
	ProfileID            string     `json:"profile_id"`
	TotalAmount          float64    `json:"total_amount"`
	Currency             string     `json:"currency"`
	Status               string     `json:"status"`
	Region               string     `json:"region"`
	ReferenceCode        string     `json:"reference_code"`
	Source               string     `json:"source"`
	SourceTime           *time.Time `json:"source_time"`
	SourceType           string     `json:"source_type"`
	PaymentTransactionID string     `json:"payment_transaction_id"`
}

type EarnPointOrderEvent struct {
	TransactionID   string  `json:"transactionID"`
	TransactionType string  `json:"transactionType"`
	PointType       int64   `json:"pointType"`
	ReferenceCode   string  `json:"referenceCode"`
	Region          string  `json:"region"`
	PointAmount     int64   `json:"pointAmount"`
	TotalAmount     float64 `json:"totalAmount"`
	Currency        string  `json:"currency"`
	ProfileID       string  `json:"profileId"`
	Status          string  `json:"status"`
}

type CallbackMessage struct {
	EventType string      `json:"eventType"`
	EventData interface{} `json:"data"`
}

type OrderEventData struct {
	OrderNumber   string `json:"orderNumber"`
	ReferenceCode string `json:"referCode"`
	CreateTime    int64  `json:"createTime"`
	Signature     string `json:"signature"`
	ProfileID     string `json:"profileID"`
	Amount        int64  `json:"amount"`
	Region        string `json:"region"`
	Currency      string `json:"currency"`
	Source        string `json:"source"`
	SourceType    string `json:"sourceType"`
	EventType     string `json:"eventType"`
	RawData       string `json:"rawData"`
	Status        string `json:"status"`
}

type RawData struct {
	PaymentTransactionID string `json:"paymentTransID"`
}
