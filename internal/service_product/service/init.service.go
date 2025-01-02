package service

import (
	"go_microservice_backend_api/internal/service_product/database"
	"go_microservice_backend_api/internal/service_product/local"
	"go_microservice_backend_api/internal/service_product/service/implement"
)

func InitProductServiceInterface() {
	ProductDb := database.New(local.ProductDb)
	InitProductService(implement.NewProductService(ProductDb))
	InitCategoryService(implement.NewCategoryService(ProductDb))
}
