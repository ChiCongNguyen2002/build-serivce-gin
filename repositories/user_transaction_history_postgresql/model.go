package user_transaction_history_postgresql

import "time"

type UserTransactionHistory struct {
	TransactionID        string     `gorm:"column:transaction_id;primaryKey"`
	TransactionType      string     `gorm:"column:transaction_type"`
	ProfileID            string     `gorm:"column:profile_id"`
	Status               string     `gorm:"column:status"`
	PointAmount          int64      `gorm:"column:point_amount"`
	PointType            int64      `gorm:"column:point_type"`
	TotalAmount          float64    `gorm:"column:total_amount"`
	Currency             string     `gorm:"column:currency"`
	PaymentTransactionID string     `gorm:"column:payment_transaction_id"`
	Source               string     `gorm:"column:source"`
	SourceTime           *time.Time `gorm:"column:source_time"`
	SourceType           string     `gorm:"column:source_type"`
	CreatedAt            *time.Time `gorm:"column:created_at"`
	UpdatedAt            *time.Time `gorm:"column:updated_at"`
}

func (UserTransactionHistory) TableName() string {
	return "user_transaction_history"
}
