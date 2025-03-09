package service

import (
	"go_microservice_backend_api/internal/service_shop/database"
	"go_microservice_backend_api/internal/service_shop/local"
	"go_microservice_backend_api/internal/service_shop/service/implement"
)

func InitShopServiceInterface() {
	userQueries := database.New(local.ShopDb)
	InitShopRegisterService(implement.NewShopRegister(userQueries))
}
