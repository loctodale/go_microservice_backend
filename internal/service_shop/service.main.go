package service_shop

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/internal/service_shop/private_config"
	"go_microservice_backend_api/internal/service_shop/router"
	"go_microservice_backend_api/internal/service_shop/service"
)

func ServiceShopMain() *gin.Engine {
	config := private_config.NewShopConfig()
	config.InitShopKafkaProducer()
	config.InitShopMysql()
	service.InitShopServiceInterface()
	r := router.NewInitRouterShop().InitRouterShop()

	return r
}
