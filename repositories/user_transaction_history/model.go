package user_transaction_history

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserTransactionHistory struct {
	TransactionID        string     `bson:"transaction_id"`
	TransactionType      string     `bson:"transaction_type"`
	ProfileID            string     `bson:"profile_id"`
	Status               string     `bson:"status"`
	PointAmount          int64      `bson:"point_amount"`
	PointType            int64      `bson:"point_type"`
	TotalAmount          float64    `bson:"total_amount"`
	Currency             string     `bson:"currency"`
	PaymentTransactionID string     `bson:"payment_transaction_id"`
	Source               string     `bson:"source"`
	SourceTime           *time.Time `bson:"source_time"`
	SourceType           string     `bson:"source_type"`
	CreatedAt            *time.Time `bson:"created_at"`
	UpdatedAt            *time.Time `bson:"updated_at"`
}

func (r UserTransactionHistory) CollectionName() string {
	return "user_transaction_history"
}

func (r UserTransactionHistory) IndexModels() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: FUserTransactionHistoryTransactionID, Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
}
