package private_config

import (
	"github.com/segmentio/kafka-go"
	"go_microservice_backend_api/internal/service_sendmail/local"
	"time"
)

type Config struct{}

func NewSendMailConfig() *Config {
	return &Config{}
}

func (c *Config) InitKafkaReader() {
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"kafka0:9092"},
		GroupID:        "group-verify-otp",
		Topic:          "otp-auth-topic",
		CommitInterval: time.Second,
	})
	local.KafkaSendMailConsumer = kafkaReader
}
