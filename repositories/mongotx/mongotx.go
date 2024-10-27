package mongotx

import (
	"build-service-gin/common/mongodb"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTxRepository struct {
	db *mongodb.DatabaseStorage
}

type IMongoTxRepository interface {
	ExecTransaction(ctx context.Context, callback func(sessCtx mongo.SessionContext) (interface{}, error)) error
}

func NewMongoTxRepository(dbStorage *mongodb.DatabaseStorage) IMongoTxRepository {
	return &MongoTxRepository{
		db: dbStorage,
	}
}

func (r *MongoTxRepository) ExecTransaction(ctx context.Context, callback func(sessCtx mongo.SessionContext) (interface{}, error)) error {
	return r.db.ExecTransaction(ctx, callback)
}
