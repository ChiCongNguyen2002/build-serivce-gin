package redis

import (
	"build-service-gin/pkg/queue"
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
)

type Consumer struct {
	redisClient redis.UniversalClient
	topic       string
	handler     queue.OnEventHandler

	cs *redis.PubSub
	mu sync.Mutex
}

func NewConsumer(redisClient redis.UniversalClient, topic string) *Consumer {
	return &Consumer{redisClient: redisClient, topic: topic}
}

func (s *Consumer) OnEvent(handler queue.OnEventHandler) {
	s.mu.Lock()
	if handler != nil {
		s.handler = handler
	}
	s.mu.Unlock()
}

func (s *Consumer) Start(ctx context.Context) error {
	s.mu.Lock()
	s.cs = s.redisClient.Subscribe(ctx, s.topic)
	s.mu.Unlock()

	for msg := range s.cs.Channel() {
		_ = s.handler(ctx, nil, []byte(msg.Payload))
	}

	return nil
}

func (s *Consumer) Shutdown(ctx context.Context) {
	// TODO implement me
	panic("implement me")
}
