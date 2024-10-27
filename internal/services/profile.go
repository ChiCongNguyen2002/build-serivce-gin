package services

import (
	"build-service-gin/common/logger"
	"build-service-gin/config"
	modelsServ "build-service-gin/internal/domains"
	"build-service-gin/pkg/helpers/adapters"
	"build-service-gin/pkg/helpers/resp"
	"build-service-gin/repositories/mongotx"
	"build-service-gin/repositories/user_transaction_history"
	"build-service-gin/repositories/user_transaction_history_postgresql"
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileService struct {
	conf                  *config.SystemConfig
	mongoRepo             mongotx.IMongoTxRepository
	profileRepo           user_transaction_history.IUserTransactionHistoryRepo
	profileRepoPostgresql user_transaction_history_postgresql.IUserTransactionHistoryPostgresSQLRepo
	//redisClient           *redis.Client
}

type IProfileService interface {
	GetUserHistoryByProfile(ctx context.Context, req modelsServ.GetUserTransactionHistoryReq) ([]modelsServ.UserTransactionHistory, int64, *resp.CustomError)
	CreateUserTransactionHistory(ctx context.Context, order *modelsServ.UserTransactionHistory) (*modelsServ.UserTransactionHistory, *resp.CustomError)
	UpdateUserTransactionHistoryByProfile(ctx context.Context, order *modelsServ.UserTransactionHistory, profileID string) (*modelsServ.UserTransactionHistory, *resp.CustomError)
	DeleteUserTransactionHistoryByProfile(ctx context.Context, profileID string) *resp.CustomError
	CreateOrderTransactionPoint(ctx context.Context, order *modelsServ.OrderSuccessEvent) *resp.CustomError
	CompleteTransactionPoint(ctx context.Context, order *modelsServ.EarnPointOrderEvent) *resp.CustomError

	GetUserHistoryByProfilePostgresql(ctx context.Context, req modelsServ.GetUserTransactionHistoryReq) ([]modelsServ.UserTransactionHistory, int64, *resp.CustomError)
	CreateUserTransactionHistoryPostgresql(ctx context.Context, order *modelsServ.UserTransactionHistory) (*modelsServ.UserTransactionHistory, *resp.CustomError)
}

func NewProfileService(conf *config.SystemConfig, mongoRepo mongotx.IMongoTxRepository, profileRepo user_transaction_history.IUserTransactionHistoryRepo, profileRepoPostgresql user_transaction_history_postgresql.IUserTransactionHistoryPostgresSQLRepo) IProfileService {
	return &ProfileService{
		conf:                  conf,
		mongoRepo:             mongoRepo,
		profileRepo:           profileRepo,
		profileRepoPostgresql: profileRepoPostgresql,
		//redisClient:           redisClient,
	}
}

func (p *ProfileService) GetUserHistoryByProfile(ctx context.Context, req modelsServ.GetUserTransactionHistoryReq) ([]modelsServ.UserTransactionHistory, int64, *resp.CustomError) {
	var profileID string
	var txTypes []string
	var recentMonth time.Time

	if req.ProfileID != "" {
		profileID = req.ProfileID
	}

	if req.RecentMonth > 0 {
		now := time.Now()
		recentMonth = now.AddDate(0, -req.RecentMonth, 0)
	}

	//upper case txType and status
	req.TxType = strings.ToUpper(req.TxType)
	req.Status = strings.ToUpper(req.Status)

	//get info user history
	userHistoryTxs, totalTxs, err := p.profileRepo.GetUserTransactionHistoryByProfile(ctx, profileID, txTypes, recentMonth, req.Offset, req.Limit, req.Status)
	if err != nil {
		return nil, 0, &resp.CustomError{ErrorCode: resp.ErrNotFound, Description: err.Error()}
	}

	pointTxsServ := adapters.AdapterProfile{}.ConvRepo2DomainServArrayUserTransactionHistoryTx(userHistoryTxs)
	return pointTxsServ, totalTxs, nil
}

func (p *ProfileService) CreateUserTransactionHistory(ctx context.Context, order *modelsServ.UserTransactionHistory) (*modelsServ.UserTransactionHistory, *resp.CustomError) {
	if order == nil {
		return nil, &resp.CustomError{ErrorCode: resp.ErrNotFound, Description: "order cannot be nil"}
	}
	pointTxsServ := adapters.AdapterProfile{}.ConvDomainToRepo(order)
	userHistory, err := p.profileRepo.CreateUserTransactionHistory(ctx, &pointTxsServ)
	if err != nil {
		return nil, &resp.CustomError{ErrorCode: resp.ErrNotFound, Description: err.Error()}
	}
	pointTxsServToDomain := adapters.AdapterProfile{}.ConvRepoToDomain(userHistory)
	return pointTxsServToDomain, nil
}

func (p *ProfileService) UpdateUserTransactionHistoryByProfile(ctx context.Context, order *modelsServ.UserTransactionHistory, profileID string) (*modelsServ.UserTransactionHistory, *resp.CustomError) {
	pointTxsServ := adapters.AdapterProfile{}.ConvDomainToRepo(order)
	userHistory, err := p.profileRepo.UpdateUserTransactionHistoryByProfile(ctx, &pointTxsServ, profileID)
	if err != nil {
		return order, nil
	}
	pointTxsServToDomain := adapters.AdapterProfile{}.ConvRepoToDomain(userHistory)
	return pointTxsServToDomain, nil
}

func (p *ProfileService) DeleteUserTransactionHistoryByProfile(ctx context.Context, profileID string) *resp.CustomError {
	err := p.profileRepo.DeleteUserTransactionHistoryByProfile(ctx, profileID)
	if err != nil {
		return nil
	}
	return nil
}

func (p *ProfileService) CreateOrderTransactionPoint(ctx context.Context, order *modelsServ.OrderSuccessEvent) *resp.CustomError {
	log := logger.GetLogger().AddTraceInfoContextRequest(ctx)
	log.Info().Interface("order", order).Msg("CreateOrderTransactionPoint - Start")

	err := p.mongoRepo.ExecTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		// Convert order to repo models
		newOrder := &modelsServ.OrderSuccessEvent{}
		newOrder.BuildCreateOrderTransaction(order)
		orderRepoModel := adapters.AdapterProfile{}.ConvertOrderCreateDomainToRepo(newOrder)

		if err := p.profileRepo.UpsertCreateOrderTransaction(sessionCtx, orderRepoModel); err != nil {
			log.Error().Err(err).Msg("CreateOrderTransactionPoint - Failed to upsert order")
			return nil, err
		}

		log.Info().Msg("CreateOrderTransactionPoint - Order upsert successfully")
		return nil, nil
	})

	if err != nil {
		return &resp.CustomError{ErrorCode: resp.ErrSystem, Description: err.Error()}
	}

	return nil
}

func (p *ProfileService) CompleteTransactionPoint(ctx context.Context, order *modelsServ.EarnPointOrderEvent) *resp.CustomError {
	log := logger.GetLogger().AddTraceInfoContextRequest(ctx)
	log.Info().Interface("order", order).Msg("CompleteOrderEarnPoint - Start")

	err := p.mongoRepo.ExecTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		orderRepoModel := adapters.AdapterProfile{}.ConvertCompleteOrderDomainToRepo(order)

		if err := p.profileRepo.UpsertCompleteOrderTransaction(sessionCtx, orderRepoModel); err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		log.Error().Err(err).Msg("CompleteOrderEarnPoint - Transaction failed")
		return &resp.CustomError{ErrorCode: resp.ErrSystem, Description: err.Error()}
	}

	log.Info().Msg("CompleteOrderEarnPoint - Order upsert successfully")
	return nil
}

func (p *ProfileService) GetUserHistoryByProfilePostgresql(ctx context.Context, req modelsServ.GetUserTransactionHistoryReq) ([]modelsServ.UserTransactionHistory, int64, *resp.CustomError) {
	var profileID string
	var txTypes []string
	var recentMonth time.Time

	// Validate and set profile ID from the request
	if req.ProfileID != "" {
		profileID = req.ProfileID
	}

	// Cache key for user history
	//cacheKey := "user_history:" + profileID

	// Step 1: Try to get user history from cache
	//var userHistory []modelsServ.UserTransactionHistory
	//if err := p.redisClient.GetDataCache(ctx, cacheKey, &userHistory); err == nil {
	//	return userHistory, int64(len(userHistory)), nil // Return cached history if found
	//}

	// Step 2: Determine the recent month filter if applicable
	if req.RecentMonth > 0 {
		now := time.Now()
		recentMonth = now.AddDate(0, -req.RecentMonth, 0)
	}

	// Convert txType and status to uppercase
	req.TxType = strings.ToUpper(req.TxType)
	req.Status = strings.ToUpper(req.Status)

	// Step 3: Query user transaction history from the PostgreSQL database
	userHistoryTxs, totalTxs, err := p.profileRepoPostgresql.GetUserTransactionHistoryByProfile(ctx, profileID, txTypes, recentMonth, req.Offset, req.Limit, req.Status)
	if err != nil {
		return nil, 0, &resp.CustomError{ErrorCode: resp.ErrNotFound, Description: err.Error()}
	}

	// Step 4: Convert repository model to service model
	pointTxsServ := adapters.AdapterProfile{}.ConvRepo2DomainServArrayUserTransactionHistoryTxPostgresql(userHistoryTxs)

	// Step 5: Cache the retrieved user history
	//if err := p.redisClient.SetDataCache(ctx, cacheKey, pointTxsServ, 0); err != nil {
	//	logger.GetLogger().Err(err).Msg("failed to set cache")
	//}

	return pointTxsServ, totalTxs, nil
}

func (p *ProfileService) CreateUserTransactionHistoryPostgresql(ctx context.Context, order *modelsServ.UserTransactionHistory) (*modelsServ.UserTransactionHistory, *resp.CustomError) {
	if order == nil {
		return nil, &resp.CustomError{ErrorCode: resp.ErrNotFound, Description: "order cannot be nil"}
	}
	pointTxsServ := adapters.AdapterProfile{}.ConvDomainToRepoPostgresql(order)
	userHistory, err := p.profileRepoPostgresql.CreateUserTransactionHistory(ctx, &pointTxsServ)
	if err != nil {
		return nil, &resp.CustomError{ErrorCode: resp.ErrNotFound, Description: err.Error()}
	}
	pointTxsServToDomain := adapters.AdapterProfile{}.ConvRepoToDomainPostgresql(userHistory)
	return pointTxsServToDomain, nil
}
