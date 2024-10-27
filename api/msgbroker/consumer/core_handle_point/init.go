package core_handle_point

//
//import (
//	"build-service-gin/api/msgbroker/models"
//	"build-service-gin/common/logger"
//	"build-service-gin/pkg/helpers/adapters"
//	queuekafka "build-service-gin/pkg/queue/kafka"
//	"context"
//	"encoding/json"
//)
//
//type ConsumerEarnPoint struct {
//	cs               *queuekafka.Consumer
//	corePointHandler *CorePointHandler
//}
//
//func NewConsumerEarnPoint(
//	cs *queuekafka.Consumer,
//	corePointHandler *CorePointHandler,
//) *ConsumerEarnPoint {
//	return &ConsumerEarnPoint{
//		cs:               cs,
//		corePointHandler: corePointHandler,
//	}
//}
//
//func (s *ConsumerEarnPoint) Start(ctx context.Context) error {
//	s.cs.OnEvent(func(ctx context.Context, key, value []byte) error {
//		log := logger.GetLogger().AddTraceInfoContextRequest(ctx)
//		var ev models.CallbackMessage
//		if err := json.Unmarshal(value, &ev); err != nil {
//			log.Err(err).Str("key", string(key)).Str("value", string(value)).Msg("decode failed")
//			return nil
//		}
//		log.Info().Any("core handler point success callback received", ev).Msg("event data")
//		eventDataByte, err := json.Marshal(ev.EventData)
//		if err != nil {
//			log.Err(err).Str("key", string(key)).Str("value", string(value)).Msg("event data marshaling failed")
//			return err
//		}
//
//		var earnData models.EarnPointOrderEvent
//		if err := json.Unmarshal(eventDataByte, &earnData); err != nil {
//			log.Err(err).Str("key", string(key)).Str("value", string(value)).Msg("event data unmarshaling failed")
//			return err
//		}
//
//		data := adapters.AdapterOrderPoint{}.ConvertEarnEventDataToOEarnPointSuccessEvent(earnData)
//		// Process the earn point event
//		err = s.corePointHandler.CorePointHandle(ctx, data)
//		if err != nil {
//			log.Err(err).Msg("handling earn point event failed")
//			return err
//		}
//
//		return nil
//	})
//
//	if err := s.cs.Start(ctx); err != nil {
//		return err
//	}
//
//	return nil
//}
