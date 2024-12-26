package service_product

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/internal/service_product/private_config"
	"go_microservice_backend_api/internal/service_product/router"
	"go_microservice_backend_api/internal/service_product/service"
)

func ProductServiceMain() *gin.Engine {
	config := private_config.NewProductConfig()
	config.InitProductSql()
	config.InitKafkaProducer()
	service.InitProductServiceInterface()
	r := router.NewInitRouter().InitRouterProduct()
	return r
}
