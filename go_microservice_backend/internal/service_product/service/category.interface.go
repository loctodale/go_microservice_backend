package service

import (
	"context"
	"go_microservice_backend_api/internal/model"
)

type ICategoryService interface {
	AddNewCategory(ctx context.Context, in model.CreateCategoryInput) (statusCode int, err error)
}

var localCategoryService ICategoryService

func CategoryService() ICategoryService {
	if localCategoryService == nil {
		panic("Category Service Must Be Initialize First")
	}
	return localCategoryService
}

func InitCategoryService(i ICategoryService) {
	localCategoryService = i
}
