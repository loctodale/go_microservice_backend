package controller

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/internal/model"
	"go_microservice_backend_api/internal/service_product/service"
	"go_microservice_backend_api/pkg/response"
)

type cCategoryController struct {
}

func NewCategoryConntroller() *cCategoryController {
	return &cCategoryController{}
}

// Shop Create Category
// @Summary      Shop Operation Create Category
// @Description  Shop Operation Create Category
// @Tags         shop op
// @Accept       json
// @Produce      json
// @Param        payload body model.CreateCategoryInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /shop-op/category/private
func (c *cCategoryController) CreateNewCategory(ctx *gin.Context) {
	var params model.CreateCategoryInput

	if err := ctx.ShouldBind(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid)
		return
	}

	result, err := service.CategoryService().AddNewCategory(ctx, params)
	if err != nil {
		response.ErrorResponse(ctx, result)
		return
	} else {
		response.SuccessResponse(ctx, 1, result)
	}
}
