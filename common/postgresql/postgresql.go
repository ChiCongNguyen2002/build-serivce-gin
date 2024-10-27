package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

type DatabasePostgresql struct {
	db *gorm.DB
}

var (
	dbStorage *gorm.DB
)

func ConnectPostgresql(ctx context.Context, cfg *PostgresqlConfig) (*DatabasePostgresql, error) {
	if dbStorage != nil {
		return &DatabasePostgresql{db: dbStorage}, nil
	}

	if strings.HasPrefix(cfg.Host, "https://") {
		cfg.Host = cfg.Host[8:]
	}

	if strings.HasPrefix(cfg.Host, "http://") {
		cfg.Host = cfg.Host[7:]
	}

	if cfg.Host[len(cfg.Host)-1:] == "/" {
		cfg.Host = cfg.Host[:len(cfg.Host)-1]
	}

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(func() logger.LogLevel {
			switch cfg.LogLevel {
			case int(logger.Silent):
				return logger.Silent
			case int(logger.Error):
				return logger.Error
			case int(logger.Warn):
				return logger.Warn
			case int(logger.Info):
				return logger.Info
			default:
				return logger.Silent
			}
		}()),
	})

	if err != nil {
		return nil, err
	}

	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// SetConnMaxIdleTime sets the maximum amount of time a connection may be idle.
	sqlDB.SetConnMaxIdleTime(2 * time.Hour)

	dbStorage = db

	return &DatabasePostgresql{db: dbStorage}, nil
}

func (d *DatabasePostgresql) GetDB() *gorm.DB {
	return dbStorage
}
