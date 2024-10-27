package models

import "time"

type GetUserTransactionHistoryReq struct {
	ProfileID   string `form:"profileID" query:"profileID"`
	Offset      int64  `form:"offset" query:"offset"`
	Limit       int64  `form:"limit" query:"limit"`
	TxType      string `form:"txType" query:"txType"`
	Status      string `form:"status" query:"status"`
	RecentMonth int    `form:"recentMonth" query:"recentMonth"`
}

type GetUserTransactionHistoryByProfileReq struct {
	ProfileID string `form:"profileID" query:"profileID"`
}

type UserTransactionHistory struct {
	TransactionID        string     `json:"transactionID" form:"transactionID"`
	TransactionType      string     `json:"transactionType" form:"transactionType"`
	ProfileID            string     `json:"profileID" form:"profileID"`
	Status               string     `json:"status" form:"status"`
	PointAmount          int64      `json:"pointAmount" form:"pointAmount"`
	PointType            int64      `json:"pointType" form:"pointType"`
	TotalAmount          float64    `json:"totalAmount" form:"totalAmount"`
	Currency             string     `json:"currency" form:"currency"`
	PaymentTransactionID string     `json:"paymentTransactionID" form:"paymentTransactionID"`
	Source               string     `json:"source" form:"source"`
	SourceTime           *time.Time `json:"sourceTime" form:"sourceTime"`
	SourceType           string     `json:"sourceType" form:"sourceType"`
	CreatedAt            *time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt            *time.Time `json:"updatedAt" form:"updatedAt"`
}
