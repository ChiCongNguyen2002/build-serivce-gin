package user_transaction_history

import (
	mongodb "build-service-gin/common/mongodb"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserTransactionHistoryRepo struct {
	*mongodb.Repository[UserTransactionHistory]
}

type IUserTransactionHistoryRepo interface {
	GetUserTransactionHistoryByProfile(ctx context.Context, profileID string, txTypes []string, recentMonth time.Time, skip, limit int64, status string) ([]*UserTransactionHistory, int64, error)
	CreateUserTransactionHistory(ctx context.Context, data *UserTransactionHistory) (*UserTransactionHistory, error)
	UpdateUserTransactionHistoryByProfile(ctx context.Context, data *UserTransactionHistory, profileID string) (*UserTransactionHistory, error)
	DeleteUserTransactionHistoryByProfile(ctx context.Context, profileID string) error
	UpsertCreateOrderTransaction(ctx context.Context, data *UserTransactionHistory) error
	UpsertCompleteOrderTransaction(ctx context.Context, data *UserTransactionHistory) error
}

func NewRepoUserTransactionHistory(dbStorage *mongodb.DatabaseStorage) IUserTransactionHistoryRepo {
	return &UserTransactionHistoryRepo{
		Repository: mongodb.NewRepository[UserTransactionHistory](dbStorage),
	}
}

func (r *UserTransactionHistoryRepo) R() *UserTransactionHistoryRepo {
	return &UserTransactionHistoryRepo{
		Repository: r.Repository.NewFilterPlayer(),
	}
}

func (r *UserTransactionHistoryRepo) byProfileID(profileID string) *UserTransactionHistoryRepo {
	filter := bson.M{
		FUserTransactionHistoryProfileID: profileID,
	}
	r.Append(filter)
	return r
}

func (r *UserTransactionHistoryRepo) byGTECreatedAt(date time.Time) *UserTransactionHistoryRepo {
	filter := bson.M{
		FUserTransactionHistoryCreatedAt: bson.M{
			"$gte": date,
		},
	}
	r.Append(filter)
	return r
}

func (r *UserTransactionHistoryRepo) byTxType(txTypes string) *UserTransactionHistoryRepo {
	filter := bson.M{
		FUserTransactionHistoryTransactionType: txTypes,
	}
	r.Append(filter)
	return r
}

func (r *UserTransactionHistoryRepo) byTxTypes(txTypes []string) *UserTransactionHistoryRepo {
	filter := bson.M{}
	if len(txTypes) > 0 {
		filter[FUserTransactionHistoryTransactionType] = bson.M{"$in": txTypes}
	}
	r.Append(filter)
	return r
}

func (r *UserTransactionHistoryRepo) byStatus(status string) *UserTransactionHistoryRepo {
	filter := bson.M{
		FUserTransactionHistoryStatus: status,
	}
	r.Append(filter)
	return r
}

func (r *UserTransactionHistoryRepo) CheckExistTxType(ctx context.Context, txType string) (*UserTransactionHistory, error) {
	rs, err := r.R().byTxType(txType).FindOneDoc(ctx)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}

func (r *UserTransactionHistoryRepo) sort(sort bson.M) *UserTransactionHistoryRepo {
	r.AppendSort(sort)
	return r
}

func (r *UserTransactionHistoryRepo) limit(limit int64) *UserTransactionHistoryRepo {
	r.SetLimit(limit)
	return r
}

func (r *UserTransactionHistoryRepo) skip(skip int64) *UserTransactionHistoryRepo {
	r.SetSkip(skip)
	return r
}

func (r *UserTransactionHistoryRepo) byTransactionId(transactionId string) *UserTransactionHistoryRepo {
	filter := bson.M{
		FUserTransactionHistoryTransactionID: transactionId,
	}
	r.Append(filter)
	return r
}

func (r *UserTransactionHistoryRepo) GetUserTransactionHistoryByProfile(ctx context.Context, profileID string, txTypes []string, recentMonth time.Time, skip, limit int64, status string) ([]*UserTransactionHistory, int64, error) {
	sort := bson.M{FUserTransactionHistoryCreatedAt: -1}

	queryBuilder := r.R().byProfileID(profileID)
	if len(txTypes) > 0 {
		queryBuilder = queryBuilder.byTxTypes(txTypes)
	}

	if status != "" {
		queryBuilder = queryBuilder.byStatus(status)
	}

	rs, err := queryBuilder.byGTECreatedAt(recentMonth).sort(sort).limit(limit).skip(skip).FindDocs(ctx)
	if err != nil {
		return nil, 0, err
	}

	total, err := queryBuilder.byGTECreatedAt(recentMonth).CountDocs(ctx)
	if err != nil {
		return nil, 0, err
	}

	return rs, total, nil
}

func (r *UserTransactionHistoryRepo) CreateUserTransactionHistory(ctx context.Context, data *UserTransactionHistory) (*UserTransactionHistory, error) {
	t := time.Now()
	data.CreatedAt = &t
	data.UpdatedAt = &t

	fmt.Printf("Inserting data: %+v\n", data)

	result, err := r.R().CreateOneDocument(ctx, data)
	if err != nil {
		fmt.Printf("Failed to insert data: %+v\nError: %v\n", data, err)
		return data, err
	}

	fmt.Printf("Insert successful: %+v\n", result)

	return data, nil
}

func (r *UserTransactionHistoryRepo) UpdateUserTransactionHistoryByProfile(ctx context.Context, data *UserTransactionHistory, profileID string) (*UserTransactionHistory, error) {
	update := bson.M{
		"$set": bson.M{
			FUserTransactionHistorySource:               data.Source,
			FUserTransactionHistorySourceTime:           data.SourceTime,
			FUserTransactionHistoryProfileID:            profileID,
			FUserTransactionHistorySourceType:           data.SourceType,
			FUserTransactionHistoryPaymentTransactionID: data.PaymentTransactionID,
			FUserTransactionHistoryTransactionID:        data.TransactionID,
			FUserTransactionHistoryTransactionType:      data.TransactionType,
			FUserTransactionHistoryPointAmount:          data.PointAmount,
			FUserTransactionHistoryTotalAmount:          data.TotalAmount,
			FUserTransactionHistoryPointType:            data.PointType,
			FUserTransactionHistoryCurrency:             data.Currency,
			FUserTransactionHistoryStatus:               data.Status,
			FUserTransactionHistoryUpdatedAt:            time.Now(),
		},
	}
	_, err := r.R().UpdateOneDoc(ctx, update)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *UserTransactionHistoryRepo) DeleteUserTransactionHistoryByProfile(ctx context.Context, profileID string) error {
	_, err := r.R().byProfileID(profileID).DeleteOneDoc(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserTransactionHistoryRepo) UpsertCreateOrderTransaction(ctx context.Context, data *UserTransactionHistory) error {
	now := time.Now()
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	updater := bson.M{
		"$set": bson.M{
			FUserTransactionHistorySource:               data.Source,
			FUserTransactionHistorySourceTime:           data.SourceTime,
			FUserTransactionHistorySourceType:           data.SourceType,
			FUserTransactionHistoryPaymentTransactionID: data.PaymentTransactionID,
			FUserTransactionHistoryUpdatedAt:            now,
		},
		"$setOnInsert": bson.M{
			FUserTransactionHistoryTransactionID:   data.TransactionID,
			FUserTransactionHistoryTransactionType: data.TransactionType,
			FUserTransactionHistoryPointAmount:     data.PointAmount,
			FUserTransactionHistoryTotalAmount:     data.TotalAmount,
			FUserTransactionHistoryPointType:       data.PointType,
			FUserTransactionHistoryCurrency:        data.Currency,
			FUserTransactionHistoryProfileID:       data.ProfileID,
			FUserTransactionHistoryCreatedAt:       now,
			FUserTransactionHistoryStatus:          data.Status,
		},
	}

	_, err := r.R().FindOneAndUpdateDoc(ctx, updater, opts)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserTransactionHistoryRepo) UpsertCompleteOrderTransaction(ctx context.Context, data *UserTransactionHistory) error {
	now := time.Now()
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	updater := bson.M{
		"$set": bson.M{
			FUserTransactionHistoryPointAmount:     data.PointAmount,
			FUserTransactionHistoryTransactionID:   data.TransactionID,
			FUserTransactionHistoryTransactionType: data.TransactionType,
			FUserTransactionHistoryUpdatedAt:       now,
			FUserTransactionHistoryPointType:       data.PointType,
			FUserTransactionHistoryStatus:          data.Status,
		},
		"$setOnInsert": bson.M{
			FUserTransactionHistorySource:               data.Source,
			FUserTransactionHistoryPaymentTransactionID: data.PaymentTransactionID,
			FUserTransactionHistorySourceTime:           data.SourceTime,
			FUserTransactionHistorySourceType:           data.SourceType,
			FUserTransactionHistoryTotalAmount:          data.TotalAmount,
			FUserTransactionHistoryCurrency:             data.Currency,
			FUserTransactionHistoryProfileID:            data.ProfileID,
			FUserTransactionHistoryCreatedAt:            now,
		},
	}

	_, err := r.R().FindOneAndUpdateDoc(ctx, updater, opts)
	if err != nil {
		return err
	}

	return nil
}
