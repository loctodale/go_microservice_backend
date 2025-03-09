package service_shop

import (
	"go_microservice_backend_api/internal/service_shop/grpc"
	"go_microservice_backend_api/internal/service_shop/private_config"
	"go_microservice_backend_api/internal/service_shop/service"
	"log"
)

func ServiceShopMain() {
	config := private_config.NewShopConfig()
	config.InitShopKafkaProducer()
	config.InitShopMysql()
	service.InitShopServiceInterface()
	err := grpc.ShopGrpcServer()
	if err != nil {
		log.Fatalf(err.Error())
	}
	//r := router.NewInitRouterShop().InitRouterShop()
	//go func() {
	//	err := grpc.ShopGrpcServer()
	//	if err != nil {
	//		log.Fatalf(err.Error())
	//	}
	//}()
	return
}
