package initialize

import (
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/pkg/logger"
)

func InitLogger() {
	global.Logger = logger.NewLogger()

}
