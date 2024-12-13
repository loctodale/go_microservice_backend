package service

import (
	"context"
	"go_microservice_backend_api/internal/model"
)

type (
	IShopRegisterService interface {
		Register(ctx context.Context, in model.RegisterInput) (codeResult int, err error)
	}
)

var (
	localShopRegisterService IShopRegisterService
)

func ShopRegisterService() IShopRegisterService {
	if localShopRegisterService == nil {
		panic("Shop Register server is nil")
	}
	return localShopRegisterService
}

func InitShopRegisterService(i IShopRegisterService) {
	localShopRegisterService = i
}
