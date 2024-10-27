package models

type OrderRequest struct {
	OrderNumber string  `json:"orderNumber"`
	CreateTime  int64   `json:"createTime"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	VGAUserID   string  `json:"vgaUserId"`
	SourceType  string  `json:"sourceType"`
}
