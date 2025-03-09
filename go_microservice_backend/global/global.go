package global

import (
	"github.com/redis/go-redis/v9"
	"go_microservice_backend_api/pkg/logger"
	"go_microservice_backend_api/pkg/settings"
)

var (
	Config settings.Config
	Logger *logger.LoggerZap
	Rdb    *redis.Client
)
