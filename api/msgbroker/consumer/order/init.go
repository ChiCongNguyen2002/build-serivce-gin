package order

//import (
//	"build-service-gin/api/msgbroker/models"
//	"build-service-gin/common/logger"
//	"build-service-gin/pkg/helpers/adapters"
//	queue_kafka "build-service-gin/pkg/queue/kafka"
//	"context"
//	"encoding/json"
//)
//
//type ConsumerOrder struct {
//	cs           *queue_kafka.Consumer
//	orderHandler *OrderHandler
//}
//
//func NewConsumerOrder(
//	cs *queue_kafka.Consumer,
//	orderHandler *OrderHandler,
//) *ConsumerOrder {
//	return &ConsumerOrder{
//		cs:           cs,
//		orderHandler: orderHandler,
//	}
//}
//
//func (s *ConsumerOrder) Start(ctx context.Context) error {
//	s.cs.OnEvent(func(ctx context.Context, key, value []byte) error {
//		log := logger.GetLogger().AddTraceInfoContextRequest(ctx)
//		var ev models.CallbackMessage
//		if err := json.Unmarshal(value, &ev); err != nil {
//			log.Err(err).Str("key", string(key)).Str("value", string(value)).Msg("decode failed")
//			return nil
//		}
//		log.Info().Any("order success callback received", ev).Msg("event data")
//
//		eventDataByte, err := json.Marshal(ev.EventData)
//		if err != nil {
//			log.Err(err).Str("key", string(key)).Str("value", string(value)).Msg("event data marshaling failed")
//			return err
//		}
//
//		// Unmarshal JSON byte slice to EarnEventData struct
//		var earnData models.OrderEventData
//		if err := json.Unmarshal(eventDataByte, &earnData); err != nil {
//			log.Err(err).Str("key", string(key)).Str("value", string(value)).Msg("event data unmarshal failed")
//			return err
//		}
//
//		var rawData models.RawData
//		if err := json.Unmarshal([]byte(earnData.RawData), &rawData); err != nil {
//			log.Err(err).Str("key", string(key)).Str("value", string(value)).Msg("Failed to unmarshal raw data: %v")
//			return err
//		}
//
//		data := adapters.AdapterOrderPoint{}.ConvertOrderEventDataToUserHistory(earnData, rawData)
//		if err := json.Unmarshal(eventDataByte, &data); err != nil {
//			log.Err(err).Str("key", string(key)).Str("value", string(value)).Msg("event data unmarshal failed")
//			return err
//		}
//
//		err = s.orderHandler.OrderHandle(ctx, data)
//		if err != nil {
//			log.Err(err).Msg("handling order success event failed")
//			return err
//		}
//		return nil
//	})
//
//	if err := s.cs.Start(ctx); err != nil {
//		return err
//	}
//	return nil
//}
