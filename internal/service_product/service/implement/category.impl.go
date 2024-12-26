package implement

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/model"
	"go_microservice_backend_api/internal/service_product/database"
	"go_microservice_backend_api/pkg/response"
)

type sCategoryService struct {
	r *database.Queries
}

func NewCategoryService(r *database.Queries) *sCategoryService {
	return &sCategoryService{
		r: r,
	}
}

func (s *sCategoryService) AddNewCategory(ctx context.Context, in model.CreateCategoryInput) (statusCode int, err error) {
	result, err := s.r.AddNewCategory(ctx, database.AddNewCategoryParams{
		CategoryDescription: in.CategoryDescription,
		CategoryIcon:        in.CategoryIcon,
		CategoryName:        in.CategoryName,
		CategorySort:        int64(in.CategorySort),
		CategorySpuCount:    int64(in.CategorySPUCount),
		CategoryStatus:      int64(in.CategoryStatus),
		HasActiveChildren:   in.HasActiveChildren,
		ParentID:            int64(in.ParentId),
	})
	if err != nil {
		global.Logger.Error("Failed to add new category", zap.Error(err))
		return 0, fmt.Errorf("failed to add new category: %w", err)
	}
	lastestID, err := result.LastInsertId()
	if lastestID == 0 {
		return response.ErrCreateFailed, err
	}
	return 1, nil
}
