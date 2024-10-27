package redis

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Producer struct {
	redisClient redis.UniversalClient
	topic       string

	pr *redis.PubSub
}

func NewProducer(redisClient redis.UniversalClient, topic string) *Producer {
	return &Producer{
		redisClient: redisClient,
		topic:       topic,
	}
}

func (s *Producer) Publish(ctx context.Context, key, value interface{}) error {
	valueData, _ := json.Marshal(value)
	return s.redisClient.Publish(ctx, s.topic, valueData).Err()
}
