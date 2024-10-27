package order

import (
	"build-service-gin/api/msgbroker/models"
	"build-service-gin/common/logger"
	"build-service-gin/internal/services"
	"build-service-gin/pkg/helpers/adapters"
	"context"
	"errors"
)

type OrderHandler struct {
	profileService services.IProfileService
}

func NewOrderHandler(profileService services.IProfileService) *OrderHandler {
	return &OrderHandler{
		profileService: profileService,
	}
}

func (h *OrderHandler) OrderHandle(ctx context.Context, data models.OrderSuccessEvent) error {
	log := logger.GetLogger().AddTraceInfoContextRequest(ctx)
	dataServ := adapters.AdapterProfile{}.ConvertEventCreateOrderToDomain(&data)
	if err := h.profileService.CreateOrderTransactionPoint(ctx, dataServ); err != nil {
		log.Err(err).Msg("Failed to process order and earn points")
		return errors.New(err.Error())
	}
	log.Info().Msg("Order success event processed successfully")
	return nil
}
