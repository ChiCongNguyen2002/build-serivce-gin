package adapters

import (
	modelsHandler "build-service-gin/api/http/models"
	model2 "build-service-gin/api/msgbroker/models"
	modelsServ "build-service-gin/internal/domains"
	modelsRepo "build-service-gin/repositories/user_transaction_history"
	modelRepoPostgres "build-service-gin/repositories/user_transaction_history_postgresql"
)

type AdapterProfile struct{}

func (a AdapterProfile) ConvReq2ServUserTransactionHistoryTx(d modelsHandler.GetUserTransactionHistoryReq) (data *modelsServ.GetUserTransactionHistoryReq) {
	data = &modelsServ.GetUserTransactionHistoryReq{
		TxType:      d.TxType,
		RecentMonth: d.RecentMonth,
		Offset:      d.Offset,
		Limit:       d.Limit,
		Status:      d.Status,
		ProfileID:   d.ProfileID,
	}
	return data
}

func (a AdapterProfile) ConvModelToDomainUserTransactionHistoryTx(d modelsHandler.UserTransactionHistory) (data *modelsServ.UserTransactionHistory) {
	data = &modelsServ.UserTransactionHistory{
		TransactionID:        d.TransactionID,
		TransactionType:      d.TransactionType,
		ProfileID:            d.ProfileID,
		Status:               d.Status,
		PointAmount:          d.PointAmount,
		PointType:            d.PointType,
		TotalAmount:          d.TotalAmount,
		Currency:             d.Currency,
		PaymentTransactionID: d.PaymentTransactionID,
		Source:               d.Source,
		SourceTime:           d.SourceTime,
		SourceType:           d.SourceType,
		CreatedAt:            d.CreatedAt,
		UpdatedAt:            d.UpdatedAt,
	}
	return data
}

func (a AdapterProfile) ConvRepo2DomainServArrayUserTransactionHistoryTx(listDataDomain []*modelsRepo.UserTransactionHistory) (data []modelsServ.UserTransactionHistory) {
	for _, item := range listDataDomain {
		converted := a.ConvRepo2ServPoint(item)
		data = append(data, *converted)
	}
	return data
}

func (a AdapterProfile) ConvDomainToRepo(d *modelsServ.UserTransactionHistory) (data modelsRepo.UserTransactionHistory) {
	data = modelsRepo.UserTransactionHistory{
		TransactionID:        d.TransactionID,
		TransactionType:      d.TransactionType,
		ProfileID:            d.ProfileID,
		Status:               d.Status,
		PointAmount:          d.PointAmount,
		PointType:            d.PointType,
		TotalAmount:          d.TotalAmount,
		Currency:             d.Currency,
		PaymentTransactionID: d.PaymentTransactionID,
		Source:               d.Source,
		SourceTime:           d.SourceTime,
		SourceType:           d.SourceType,
		CreatedAt:            d.CreatedAt,
		UpdatedAt:            d.UpdatedAt,
	}
	return data
}

func (a AdapterProfile) ConvRepoToDomain(d *modelsRepo.UserTransactionHistory) (data *modelsServ.UserTransactionHistory) {
	data = &modelsServ.UserTransactionHistory{
		TransactionID:        d.TransactionID,
		TransactionType:      d.TransactionType,
		ProfileID:            d.ProfileID,
		Status:               d.Status,
		PointAmount:          d.PointAmount,
		PointType:            d.PointType,
		TotalAmount:          d.TotalAmount,
		Currency:             d.Currency,
		PaymentTransactionID: d.PaymentTransactionID,
		Source:               d.Source,
		SourceTime:           d.SourceTime,
		SourceType:           d.SourceType,
		CreatedAt:            d.CreatedAt,
		UpdatedAt:            d.UpdatedAt,
	}
	return data
}

func (a AdapterProfile) ConvRepo2ServPoint(d *modelsRepo.UserTransactionHistory) (data *modelsServ.UserTransactionHistory) {
	data = &modelsServ.UserTransactionHistory{
		TransactionID:        d.TransactionID,
		TransactionType:      d.TransactionType,
		ProfileID:            d.ProfileID,
		Status:               d.Status,
		PointAmount:          d.PointAmount,
		PointType:            d.PointType,
		TotalAmount:          d.TotalAmount,
		Currency:             d.Currency,
		PaymentTransactionID: d.PaymentTransactionID,
		Source:               d.Source,
		SourceTime:           d.SourceTime,
		SourceType:           d.SourceType,
		CreatedAt:            d.CreatedAt,
		UpdatedAt:            d.UpdatedAt,
	}
	return data
}

func (a AdapterProfile) ConvertOrderCreateDomainToRepo(data *modelsServ.OrderSuccessEvent) *modelsRepo.UserTransactionHistory {
	return &modelsRepo.UserTransactionHistory{
		ProfileID:            data.ProfileID,
		TotalAmount:          data.TotalAmount,
		Status:               data.Status,
		Currency:             data.Currency,
		Source:               data.Source,
		SourceTime:           data.SourceTime,
		SourceType:           data.SourceType,
		PaymentTransactionID: data.PaymentTransactionID,
	}
}

func (a AdapterProfile) ConvertOrderUpdateDomainToRepo(data *modelsServ.OrderSuccessEvent) *modelsRepo.UserTransactionHistory {
	return &modelsRepo.UserTransactionHistory{
		Source:               data.Source,
		SourceTime:           data.SourceTime,
		SourceType:           data.SourceType,
		PaymentTransactionID: data.PaymentTransactionID,
	}
}

func (a AdapterProfile) ConvertCompleteOrderDomainToRepo(data *modelsServ.EarnPointOrderEvent) *modelsRepo.UserTransactionHistory {
	return &modelsRepo.UserTransactionHistory{
		TransactionID:   data.TransactionID,
		TransactionType: data.TransactionType,
		Status:          data.Status,
		PointAmount:     data.PointAmount,
		PointType:       data.PointType,
		TotalAmount:     data.TotalAmount,
		Currency:        data.Currency,
		ProfileID:       data.ProfileID,
	}
}

func (a AdapterProfile) ConvertEventCreateOrderToDomain(event *model2.OrderSuccessEvent) *modelsServ.OrderSuccessEvent {
	return &modelsServ.OrderSuccessEvent{
		ProfileID:            event.ProfileID,
		TotalAmount:          event.TotalAmount,
		Currency:             event.Currency,
		Status:               event.Status,
		Region:               event.Region,
		Source:               event.Source,
		SourceTime:           event.SourceTime,
		SourceType:           event.SourceType,
		PaymentTransactionID: event.PaymentTransactionID,
	}
}

func (a AdapterProfile) ConvertEventEarnPointToDomain(event *model2.EarnPointOrderEvent) *modelsServ.EarnPointOrderEvent {
	return &modelsServ.EarnPointOrderEvent{
		TransactionID:   event.TransactionID,
		ReferenceCode:   event.ReferenceCode,
		Region:          event.Region,
		TransactionType: event.TransactionType,
		PointType:       event.PointType,
		PointAmount:     event.PointAmount,
		TotalAmount:     event.TotalAmount,
		Currency:        event.Currency,
		ProfileID:       event.ProfileID,
	}
}

func (a AdapterProfile) ConvDomainToRepoPostgresql(d *modelsServ.UserTransactionHistory) (data modelRepoPostgres.UserTransactionHistory) {
	data = modelRepoPostgres.UserTransactionHistory{
		TransactionID:        d.TransactionID,
		TransactionType:      d.TransactionType,
		ProfileID:            d.ProfileID,
		Status:               d.Status,
		PointAmount:          d.PointAmount,
		PointType:            d.PointType,
		TotalAmount:          d.TotalAmount,
		Currency:             d.Currency,
		PaymentTransactionID: d.PaymentTransactionID,
		Source:               d.Source,
		SourceTime:           d.SourceTime,
		SourceType:           d.SourceType,
		CreatedAt:            d.CreatedAt,
		UpdatedAt:            d.UpdatedAt,
	}
	return data
}

func (a AdapterProfile) ConvRepoToDomainPostgresql(d *modelRepoPostgres.UserTransactionHistory) (data *modelsServ.UserTransactionHistory) {
	data = &modelsServ.UserTransactionHistory{
		TransactionID:        d.TransactionID,
		TransactionType:      d.TransactionType,
		ProfileID:            d.ProfileID,
		Status:               d.Status,
		PointAmount:          d.PointAmount,
		PointType:            d.PointType,
		TotalAmount:          d.TotalAmount,
		Currency:             d.Currency,
		PaymentTransactionID: d.PaymentTransactionID,
		Source:               d.Source,
		SourceTime:           d.SourceTime,
		SourceType:           d.SourceType,
		CreatedAt:            d.CreatedAt,
		UpdatedAt:            d.UpdatedAt,
	}
	return data
}

func (a AdapterProfile) ConvRepo2DomainServArrayUserTransactionHistoryTxPostgresql(listDataDomain []*modelRepoPostgres.UserTransactionHistory) (data []modelsServ.UserTransactionHistory) {
	for _, item := range listDataDomain {
		converted := a.ConvRepo2ServPointPostgresql(item)
		data = append(data, *converted)
	}
	return data
}

func (a AdapterProfile) ConvRepo2ServPointPostgresql(d *modelRepoPostgres.UserTransactionHistory) (data *modelsServ.UserTransactionHistory) {
	data = &modelsServ.UserTransactionHistory{
		TransactionID:        d.TransactionID,
		TransactionType:      d.TransactionType,
		ProfileID:            d.ProfileID,
		Status:               d.Status,
		PointAmount:          d.PointAmount,
		PointType:            d.PointType,
		TotalAmount:          d.TotalAmount,
		Currency:             d.Currency,
		PaymentTransactionID: d.PaymentTransactionID,
		Source:               d.Source,
		SourceTime:           d.SourceTime,
		SourceType:           d.SourceType,
		CreatedAt:            d.CreatedAt,
		UpdatedAt:            d.UpdatedAt,
	}
	return data
}
