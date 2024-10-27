package kafka

import (
	"build-service-gin/common/logger"
	"build-service-gin/common/utils"
	"build-service-gin/config"
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ProducerInterface interface {
	Publish(ctx context.Context, key, value interface{}) error
}

type Producer struct {
	pr    *kafka.Producer
	topic string
}

func NewProducer(cfg config.KafkaConfig, topic string) *Producer {
	log := logger.GetLogger()
	pr, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  cfg.BootstrapServers,
		"enable.idempotence": true,
		"acks":               "all",
	})

	if err != nil {
		log.Fatal().Err(err).Msg("init kafka producer failed")
	}

	log.Info().Msgf("init kafka producer success : TOPIC = %v", topic)
	return &Producer{
		pr:    pr,
		topic: topic,
	}
}

func marshal(val interface{}) ([]byte, error) {
	switch v := val.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return json.Marshal(val)
	}
}

func (s *Producer) Publish(ctx context.Context, key, value interface{}) error {
	keyData, err := marshal(key)
	if err != nil {
		return err
	}
	valueData, err := marshal(value)
	if err != nil {
		return err
	}

	header, err := marshal(utils.GetRequestIdByContext(ctx))
	if err != nil {
		return err
	}

	msg := &kafka.Message{
		Key:   keyData,
		Value: valueData,
		Headers: []kafka.Header{
			{
				Key:   utils.KeyTraceInfo,
				Value: header,
			},
		},
		TopicPartition: kafka.TopicPartition{
			Topic:     &s.topic,
			Partition: kafka.PartitionAny,
		},
	}

	if err = s.pr.Produce(msg, nil); err != nil {
		return err
	}
	return nil
}

func (s *Producer) PublishBytes(key, value []byte) error {
	msg := &kafka.Message{
		Key:   key,
		Value: value,
		TopicPartition: kafka.TopicPartition{
			Topic:     &s.topic,
			Partition: kafka.PartitionAny,
		},
	}
	return s.pr.Produce(msg, nil)
}

func (s *Producer) PublishMessage(msg *kafka.Message) error {
	return s.pr.Produce(msg, nil)
}

func (s *Producer) PublishWithTopic(ctx context.Context, topic string, key, value interface{}) error {
	keyData, err := marshal(key)
	if err != nil {
		return err
	}
	valueData, err := marshal(value)
	if err != nil {
		return err
	}

	header, err := marshal(utils.GetRequestIdByContext(ctx))
	if err != nil {
		return err
	}

	msg := &kafka.Message{
		Key:   keyData,
		Value: valueData,
		Headers: []kafka.Header{
			{
				Key:   utils.KeyTraceInfo,
				Value: header,
			},
		},
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
	}
	if err = s.pr.Produce(msg, nil); err != nil {
		return err
	}
	return nil
}

func (s *Producer) GetTopicName() string {
	return s.topic
}
