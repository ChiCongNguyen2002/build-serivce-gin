package services

import (
	"build-service-gin/client/receiver"
	"build-service-gin/common/logger"
	"build-service-gin/common/utils"
	"build-service-gin/config"
	modelsServ "build-service-gin/internal/domains"
	"build-service-gin/pkg/helpers/resp"
	"context"
	"encoding/json"
	"time"
)

type PointService struct {
	conf           *config.SystemConfig
	receiverClient receiver.IReceiverClient
}

type IPointService interface {
	CreatePointTransaction(ctx context.Context, order *modelsServ.Order) *resp.CustomError
}

func NewPointService(
	conf *config.SystemConfig,
	receiverClient receiver.IReceiverClient,
) IPointService {
	return &PointService{
		conf:           conf,
		receiverClient: receiverClient,
	}
}

func (s *PointService) CreatePointTransaction(ctx context.Context, order *modelsServ.Order) *resp.CustomError {
	log := logger.GetLogger().AddTraceInfoContextRequest(ctx)
	log.Info().Interface("transaction", order).Msg("CreateAdminPointTransaction")

	region := ctx.Value(utils.KeyRegion).(string)

	orderPoint := &modelsServ.OrderPoint{
		OrderNumber: utils.GetIdGenerate().GetIDStringV2(),
		CreateTime:  time.Now().UnixMilli(),
		Amount:      order.Amount,
		Currency:    order.Currency,
		VGAUserID:   order.VGAUserID,
		Region:      region,
	}

	orderPointJSON, err := json.Marshal(orderPoint)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal OrderPointMessage")
		return &resp.CustomError{ErrorCode: resp.ErrSystem, Description: "Failed to process order data"}
	}

	//signature, err := s.chain.GetContractLPPoint().GenSig(string(orderPointJSON))
	//if err != nil {
	//	log.Error().Err(err).Msg("Failed to generate signature")
	//	return &resp.CustomError{ErrorCode: resp.ErrSystem, Description: "Failed to sign order data"}
	//}

	signedOrder := modelsServ.OrderMessage{
		SourceType: order.SourceType,
		RawData:    string(orderPointJSON),
		//Signature:  signature,
	}

	_, err = s.receiverClient.PostOrder(ctx, signedOrder)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send order request")
		return &resp.CustomError{ErrorCode: resp.ErrSystem, Description: "error create order point"}
	}
	log.Info().Msg("Order created successfully")
	return nil

}
