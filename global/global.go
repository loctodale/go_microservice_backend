package global

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go_microservice_backend_api/pkg/logger"
	"go_microservice_backend_api/pkg/settings"
)

var (
	Config        settings.Config
	Logger        *logger.LoggerZap
	Rdb           *redis.Client
	Mdb           *sql.DB
	KafkaProducer *kafka.Writer
	KafkaConsumer *kafka.Reader
)
