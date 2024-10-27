package adapters

import (
	"build-service-gin/api/msgbroker/models"
	"time"
)

type AdapterOrderPoint struct{}

func (a AdapterOrderPoint) ConvertOrderEventDataToUserHistory(earnData *models.OrderEventData, rawData models.RawData) models.OrderSuccessEvent {
	sourceTime := time.UnixMilli(earnData.CreateTime)
	return models.OrderSuccessEvent{
		ProfileID:            earnData.ProfileID,
		TotalAmount:          float64(earnData.Amount),
		Currency:             earnData.Currency,
		Status:               earnData.Status,
		Source:               earnData.Source,
		SourceTime:           &sourceTime,
		SourceType:           earnData.SourceType,
		Region:               earnData.Region,
		ReferenceCode:        earnData.OrderNumber,
		PaymentTransactionID: rawData.PaymentTransactionID,
	}
}

func (a AdapterOrderPoint) ConvertEarnEventDataToOEarnPointSuccessEvent(earnData *models.EarnPointOrderEvent) models.EarnPointOrderEvent {
	return models.EarnPointOrderEvent{
		TransactionID:   earnData.TransactionID,
		ReferenceCode:   earnData.ReferenceCode,
		TransactionType: earnData.TransactionType,
		PointAmount:     earnData.PointAmount,
		PointType:       earnData.PointType,
		TotalAmount:     earnData.TotalAmount,
		Currency:        earnData.Currency,
		Region:          earnData.Region,
		Status:          earnData.Status,
		ProfileID:       earnData.ProfileID,
	}
}
