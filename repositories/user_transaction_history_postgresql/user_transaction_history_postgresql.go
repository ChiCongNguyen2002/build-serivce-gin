package user_transaction_history_postgresql

import (
	postgres "build-service-gin/common/postgresql"
	"context"
	"fmt"
	"time"
)

type UserTransactionHistoryPostgresSQLRepo struct {
	*postgres.Repository[UserTransactionHistory]
}

type IUserTransactionHistoryPostgresSQLRepo interface {
	GetUserTransactionHistoryByProfile(ctx context.Context, profileID string, txTypes []string, recentMonth time.Time, skip, limit int64, status string) ([]*UserTransactionHistory, int64, error)
	CreateUserTransactionHistory(ctx context.Context, data *UserTransactionHistory) (*UserTransactionHistory, error)
	UpdateUserTransactionHistoryByProfile(ctx context.Context, data *UserTransactionHistory, profileID string) (*UserTransactionHistory, error)
	DeleteUserTransactionHistoryByProfile(ctx context.Context, profileID string) error
	//UpsertCreateOrderTransaction(ctx context.Context, data *UserTransactionHistory) error
	//UpsertCompleteOrderTransaction(ctx context.Context, data *UserTransactionHistory) error
}

func NewRepoUserTransactionHistoryPostgresql(dbStorage *postgres.DatabasePostgresql) IUserTransactionHistoryPostgresSQLRepo {
	return &UserTransactionHistoryPostgresSQLRepo{
		Repository: postgres.NewRepository[UserTransactionHistory](dbStorage),
	}
}

func (r *UserTransactionHistoryPostgresSQLRepo) GetUserTransactionHistoryByProfile(ctx context.Context, profileID string, txTypes []string, recentMonth time.Time, skip, limit int64, status string) ([]*UserTransactionHistory, int64, error) {
	var rs []*UserTransactionHistory
	var total int64

	query := r.GetDB().Model(&UserTransactionHistory{}).Where("profile_id = ?", profileID)

	if len(txTypes) > 0 {
		query = query.Where("transaction_type IN ?", txTypes)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if !recentMonth.IsZero() {
		query = query.Where("created_at >= ?", recentMonth)
	}

	err := query.Offset(int(skip)).Limit(int(limit)).Order("created_at desc").Find(&rs).Error
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return rs, total, nil
}

func (r *UserTransactionHistoryPostgresSQLRepo) CreateUserTransactionHistory(ctx context.Context, data *UserTransactionHistory) (*UserTransactionHistory, error) {
	t := time.Now()
	data.CreatedAt = &t
	data.UpdatedAt = &t

	if err := r.GetDB().Create(data).Error; err != nil {
		fmt.Printf("Failed to insert data: %+v\nError: %v\n", data, err)
		return nil, err
	}
	fmt.Printf("Insert successful: %+v\n", data)
	return data, nil
}

func (r *UserTransactionHistoryPostgresSQLRepo) UpdateUserTransactionHistoryByProfile(ctx context.Context, data *UserTransactionHistory, profileID string) (*UserTransactionHistory, error) {
	update := map[string]interface{}{
		"source":                 data.Source,
		"source_time":            data.SourceTime,
		"profile_id":             profileID,
		"source_type":            data.SourceType,
		"payment_transaction_id": data.PaymentTransactionID,
		"transaction_id":         data.TransactionID,
		"transaction_type":       data.TransactionType,
		"point_amount":           data.PointAmount,
		"total_amount":           data.TotalAmount,
		"point_type":             data.PointType,
		"currency":               data.Currency,
		"status":                 data.Status,
		"updated_at":             time.Now(),
	}

	if err := r.GetDB().Model(&UserTransactionHistory{}).Where("profile_id = ?", profileID).Updates(update).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *UserTransactionHistoryPostgresSQLRepo) DeleteUserTransactionHistoryByProfile(ctx context.Context, profileID string) error {
	if err := r.GetDB().Where("profile_id = ?", profileID).Delete(&UserTransactionHistory{}).Error; err != nil {
		return err
	}
	return nil
}

//func (r *UserTransactionHistoryPostgresSQLRepo) UpsertCreateOrderTransaction(ctx context.Context, data *UserTransactionHistory) error {
//	return r.GetDB().Clauses(gorm.OnConflict{
//		Columns:   []string{"transaction_id"},
//		DoUpdates: gorm.AssignmentColumns([]string{"source", "source_time", "source_type", "payment_transaction_id", "updated_at"}),
//	}).Create(data).Error
//}
//
//func (r *UserTransactionHistoryPostgresSQLRepo) UpsertCompleteOrderTransaction(ctx context.Context, data *UserTransactionHistory) error {
//	return r.GetDB().Clauses(gorm.OnConflict{
//		Columns:   []string{"transaction_id"}, // column to check for conflict
//		DoUpdates: gorm.AssignmentColumns([]string{"point_amount", "transaction_type", "updated_at", "point_type", "status"}),
//	}).Create(data).Error
//}
