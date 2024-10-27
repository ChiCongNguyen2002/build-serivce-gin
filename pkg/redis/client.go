package redis

import (
	"build-service-gin/common/logger"
	"build-service-gin/config"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type Client struct {
	client     redis.UniversalClient
	expDefault time.Duration
}

var redisClient *redis.Client

func ConnectRedis(ctx context.Context, cfg config.RedisConfig) (*Client, error) {
	log := logger.GetLogger()
	if redisClient != nil {
		_, err := redisClient.Ping(ctx).Result()
		if err == nil {
			return &Client{
				client:     redisClient,
				expDefault: time.Second * 60,
			}, nil
		}
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       0,
		Password: cfg.Password,
		Username: cfg.User,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Err(err).Msg("ping redis failed")
		return nil, err
	}

	log.Info().Msg("connect redis successfully")
	return &Client{
		client:     redisClient,
		expDefault: time.Second * 60,
	}, nil
}

func (c *Client) GetInstance() redis.UniversalClient {
	return c.client
}

func (c *Client) GetDataCache(ctx context.Context, key string, rs interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	} else {
		err = json.Unmarshal(data, &rs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) SetDataCache(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if exp == 0 {
		exp = c.expDefault
	}
	return c.client.Set(ctx, key, data, exp).Err()
}

func (c *Client) IncrementDataCache(ctx context.Context, key string) error {
	return c.client.Incr(ctx, key).Err()
}

func (c *Client) DecrementDataCache(ctx context.Context, key string) error {
	value, err := c.client.Get(ctx, key).Int64()
	if err == nil && value >= 0 {
		return c.client.Incr(ctx, key).Err()
	}

	return nil
}

func (c *Client) RemoteDataCache(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Client) AcquireLock(ctx context.Context, lockKey string, lockTimeout time.Duration) (bool, error) {
	isSet, err := c.client.SetNX(ctx, lockKey, 1, lockTimeout).Result()
	if err != nil {
		return false, err
	}

	return isSet, err
}

func (c *Client) ReleaseLock(ctx context.Context, lockKey string) error {
	_, err := c.client.Del(ctx, lockKey).Result()
	if err != nil {
		return err
	}

	return nil
}
