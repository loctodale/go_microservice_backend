package service_user

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/internal/service_user/private_config"
	"go_microservice_backend_api/internal/service_user/router"
	"go_microservice_backend_api/internal/service_user/service"
)

func ServiceUserMain() *gin.Engine {
	config := private_config.NewUserConfig()
	config.InitKafkaProducer()
	config.InitUserMysql()
	service.InitUserServiceInterface()
	r := router.NewInitRouterUser().InitRouterUser()
	return r
}
