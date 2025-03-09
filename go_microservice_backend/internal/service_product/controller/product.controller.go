package controller

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/internal/model"
	"go_microservice_backend_api/internal/service_product/service"
	"go_microservice_backend_api/pkg/response"
)

type cProductController struct {
}

func NewProductController() *cProductController {
	return &cProductController{}
}

// Shop Create Product
// @Summary      Shop Operation Create Product
// @Description  Shop Operation Create Product
// @Tags         shop op
// @Accept       json
// @Produce      json
// @Param        payload body model.CreateProductInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /shop-op/product/private
func (c *cProductController) CreateNewProduct(ctx *gin.Context) {
	var params model.CreateProductInput
	shopId := ctx.GetHeader("X-Consumer-Custom-ID")
	if err := ctx.ShouldBind(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid)
	}
	result, err := service.ProductService().AddNewProduct(ctx, params, shopId)
	if err != nil {
		response.ErrorResponse(ctx, result)
		return
	}

	response.SuccessResponse(ctx, response.CodeSuccess, nil)
	return
}

func (c *cProductController) CreateNewSKU(ctx *gin.Context) {
	var params model.CreateProductSKUInput
	if err := ctx.ShouldBind(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid)
	}

	result, err := service.ProductService().AddNewSKUProduct(ctx, params)
	if err != nil {
		response.ErrorResponse(ctx, result)
		return
	}
	response.SuccessResponse(ctx, response.CodeSuccess, nil)
}
