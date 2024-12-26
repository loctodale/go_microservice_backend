package service

import (
	"context"
	"go_microservice_backend_api/internal/model"
)

type IProductService interface {
	AddNewProduct(ctx context.Context, in model.CreateProductInput, shopId string) (statusCode int, err error)
	AddNewSKUProduct(ctx context.Context, in model.CreateProductSKUInput) (statusCode int, err error)
}

var (
	localProductService IProductService
)

func ProductService() IProductService {
	if localProductService == nil {
		panic("Product Service Must Be Initialize First")
	}
	return localProductService
}

func InitProductService(i IProductService) {
	localProductService = i
}
