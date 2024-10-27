package config

import (
	logger2 "build-service-gin/common/logger"
	"build-service-gin/common/mongodb"
	"build-service-gin/common/postgresql"
	"github.com/caarlos0/env/v7"
)

type SystemConfig struct {
	Env       string `env:"ENV,required,notEmpty"`
	HttpPort  uint64 `env:"HTTP_PORT,required,notEmpty"`
	ServiceID string `env:"SERVICE_ID,required,notEmpty"`

	KafkaConfig          KafkaConfig                 `envPrefix:"KAFKA_"`
	KafkaTopicConfig     KafkaTopicConfig            `envPrefix:"TOPICS_"`
	MongoDBConfig        mongodb.MongoDBConfig       `envPrefix:"MONGODB_"`
	InternalToken        string                      `env:"INTERNAL_TOKEN,required,notEmpty"`
	RewardIntegrationUrl string                      `env:"REWARD_INTEGRATION_URL,required,notEmpty"`
	PostgresConfig       postgresql.PostgresqlConfig `envPrefix:"POSTGRES_"`
	RedisConfig          RedisConfig                 `envPrefix:"REDIS_"`
}

var configSingletonObj *SystemConfig

func LoadConfig() (cf *SystemConfig, err error) {
	log := logger2.GetLogger()

	if configSingletonObj != nil {
		cf = configSingletonObj
		return
	}

	cf = &SystemConfig{}
	if err = env.Parse(cf); err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal config")
	}

	configSingletonObj = cf
	return
}

func GetInstance() *SystemConfig {
	return configSingletonObj
}

type RedisConfig struct {
	Addr     string `env:"ADDRESS,required,notEmpty"`
	Password string `env:"PASS,required,notEmpty"`
	User     string `env:"USER"`
}

type KafkaConfig struct {
	BootstrapServers string `env:"BOOTSTRAP_SERVERS"`
	GroupID          string `env:"GROUP_ID"`
	AutoOffsetReset  string `env:"AUTO_OFFSET_RESET"`
}

type KafkaTopicConfig struct {
	TopicsRewardsPoint                string `env:"REWARDS_POINT,required,notEmpty"`
	TopicsCoreTransactionPointSuccess string `env:"CORE_TRANSACTION_POINT_SUCCESS,required,notEmpty"`
}
