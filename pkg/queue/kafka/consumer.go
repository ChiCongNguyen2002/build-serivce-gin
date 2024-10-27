package kafka

import (
	"build-service-gin/common/logger"
	"build-service-gin/common/utils"
	"build-service-gin/config"
	"build-service-gin/pkg/queue"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ConsumerInterface interface {
	OnEvent(handler queue.OnEventHandler)
	Start(ctx context.Context) error
	Shutdown(ctx context.Context)
}

type Consumer struct {
	cs      *kafka.Consumer
	topics  []string
	handler queue.OnEventHandler

	started bool
	mu      sync.RWMutex
}

func NewConsumer(cfg config.KafkaConfig, topics []string) *Consumer {
	log := logger.GetLogger()
	cs, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.BootstrapServers,
		"group.id":          cfg.GroupID,
		"auto.offset.reset": cfg.AutoOffsetReset,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("init kafka consumer failed")
	}

	log.Info().Msgf("init kafka consumer success : TOPIC = %v", topics)
	return &Consumer{
		cs:     cs,
		topics: topics,
	}
}

func (s *Consumer) OnEvent(handler queue.OnEventHandler) {
	s.mu.Lock()
	if handler != nil {
		s.handler = handler
	}
	s.mu.Unlock()
}

func (s *Consumer) Start(ctx context.Context) error {
	log := logger.GetLogger().AddTraceInfoContextRequest(ctx)
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return queue.ErrAlreadyStarted
	}
	if s.handler == nil {
		s.mu.Unlock()
		return queue.ErrNilEventHandler
	}
	s.started = true
	s.mu.Unlock()

	if err := s.cs.SubscribeTopics(s.topics, nil); err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	for {
		msg, err := s.cs.ReadMessage(10 * time.Second)
		if err != nil {
			var kerr kafka.Error
			if errors.As(err, &kerr); kerr.Code() == kafka.ErrTimedOut {
				continue
			}
			log.Warn().Err(err).Msg("kafka read message failed")
			continue
		}

		traceInfoExisted := false
		newCtx := context.Background()
		for _, h := range msg.Headers {
			if h.Key == utils.KeyTraceInfo {
				traceInfo := utils.TraceInfo{}
				if err = json.Unmarshal(h.Value, &traceInfo); err != nil {
					break
				}
				newCtx = context.WithValue(newCtx, utils.KeyTraceInfo, traceInfo)
				traceInfoExisted = true
			}
		}

		if !traceInfoExisted {
			newCtx = context.WithValue(newCtx, utils.KeyTraceInfo, utils.GetRequestIdByContext(newCtx))
		}

		if err = s.handler(newCtx, msg.Key, msg.Value); err != nil {
			log.Err(err).Msg("kafka handler failed")
			continue
		}

		if _, err = s.cs.CommitMessage(msg); err != nil {
			log.Err(err).Msg("kafka commit failed")
		}
	}
}

func (s *Consumer) Shutdown() {
}
