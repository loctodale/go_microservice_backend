package shop

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/model"
	"go_microservice_backend_api/internal/service_shop/service"
	"go_microservice_backend_api/pkg/response"
)

type cShopAuthController struct {
}

var NewAuthController = new(cShopAuthController)

// Shop register
// @Summary      Shop register new account
// @Description  Shop register new account
// @Tags         shop account management
// @Accept       json
// @Produce      json
// @Param        payload body model.VerifyInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /user/verify_account [post]
func (c *cShopAuthController) Register(ctx *gin.Context) {
	var params model.RegisterInput
	if err := ctx.ShouldBind(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid)
	}

	result, err := service.ShopRegisterService().Register(ctx, params)
	if err != nil {
		global.Logger.Error(err.Error(), zap.Error(err))
		response.ErrorResponse(ctx, result)
		return
	}

	if result == 1 {
		response.SuccessResponse(ctx, response.CodeSuccess, nil)
		return
	} else {
		response.ErrorResponse(ctx, response.ErrCodeOTPNotExists)
		return
	}
}

func (c *cShopAuthController) VerifyOTP(ctx *gin.Context) {
	var params model.VerifyInput
	if err := ctx.ShouldBind(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid)
	}
	result, err := service.ShopRegisterService().VerifyOTP(ctx, params)
	if err != nil {
		global.Logger.Error(err.Error(), zap.Error(err))
		return
	}
	response.SuccessResponse(ctx, 1, result)
	return
}

func (c *cShopAuthController) ChangePasswordRegister(ctx *gin.Context) {
	result, err := service.ShopRegisterService().ChangePasswordRegister(ctx)
	if err != nil {
		global.Logger.Error(err.Error(), zap.Error(err))
	}
	response.SuccessResponse(ctx, 1, result)
	return
}
