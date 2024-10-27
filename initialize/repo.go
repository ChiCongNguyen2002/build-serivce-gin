package initialize

import (
	"build-service-gin/common/mongodb"
	postgres "build-service-gin/common/postgresql"
	"build-service-gin/repositories/mongotx"
	"build-service-gin/repositories/user_transaction_history"
	"build-service-gin/repositories/user_transaction_history_postgresql"
)

var (
	repositories *Repositories
)

type Repositories struct {
	IUserTransactionHistoryRepo         user_transaction_history.IUserTransactionHistoryRepo
	IUserTransactionHistoryPostgresRepo user_transaction_history_postgresql.IUserTransactionHistoryPostgresSQLRepo
	IMongoTxRepository                  mongotx.IMongoTxRepository
}

func NewRepositories(dbStorage *mongodb.DatabaseStorage, postgres *postgres.DatabasePostgresql) *Repositories {
	repositories = &Repositories{
		IUserTransactionHistoryRepo:         user_transaction_history.NewRepoUserTransactionHistory(dbStorage),
		IUserTransactionHistoryPostgresRepo: user_transaction_history_postgresql.NewRepoUserTransactionHistoryPostgresql(postgres),
		IMongoTxRepository:                  mongotx.IMongoTxRepository(dbStorage),
	}
	return repositories
}
