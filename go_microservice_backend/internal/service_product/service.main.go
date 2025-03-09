package service_product

import (
	"go_microservice_backend_api/internal/service_product/grpc"
	"go_microservice_backend_api/internal/service_product/private_config"
	"go_microservice_backend_api/internal/service_product/service"
	"log"
)

func ProductServiceMain() {
	config := private_config.NewProductConfig()
	config.InitProductSql()
	config.InitKafkaProducer()
	service.InitProductServiceInterface()
	//r := router.NewInitRouter().InitRouterProduct()
	//return r
	err := grpc.ProductGrpcServer()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
