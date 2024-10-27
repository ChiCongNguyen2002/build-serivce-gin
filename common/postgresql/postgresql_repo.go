package postgresql

import (
	"build-service-gin/common/logger"
	"gorm.io/gorm"
)

type ModelInterface interface {
	TableName() string
}

type Repository[T ModelInterface] struct {
	*gorm.DB
	tableName string
	db        *gorm.DB
}

func NewRepository[T ModelInterface](dbStorage *DatabasePostgresql) *Repository[T] {
	log := logger.GetLogger()

	var t T
	err := dbStorage.db.Table(t.TableName()).AutoMigrate(t)
	if err != nil {
		log.Fatal().Err(err).Msgf("create table_name=%s failed", t.TableName())
	}

	return &Repository[T]{
		DB:        dbStorage.db.Table(t.TableName()),
		tableName: t.TableName(),
		db:        dbStorage.db,
	}
}

func (r *Repository[T]) WithTableName(tx *gorm.DB) *gorm.DB {
	return tx.Table(r.tableName)
}

func (r *Repository[T]) GetDB() *gorm.DB {
	return r.db
}
