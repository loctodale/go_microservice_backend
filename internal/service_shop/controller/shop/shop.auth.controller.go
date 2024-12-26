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
// @Param        payload body model.RegisterInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /shops/auth/public/register [post]
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

// Shop register
// @Summary      Shop verify OTP
// @Description  Shop verify OTP
// @Tags         shop account management
// @Accept       json
// @Produce      json
// @Param        payload body model.VerifyInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /shops/auth/public/verifyOTP [post]
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

// Shop register
// @Summary      Shop Change Password After Register
// @Description  Shop Change Password After Register
// @Tags         shop account management
// @Accept       json
// @Produce      json
// @Param        payload body model.ShopChangePasswordRegisterInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /shops/auth/private/change_password
func (c *cShopAuthController) ChangePasswordRegister(ctx *gin.Context) {
	username := ctx.GetHeader("X-Consumer-Username")
	var params model.ShopChangePasswordRegisterInput
	if err := ctx.ShouldBind(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid)
	}
	result, err := service.ShopRegisterService().ChangePasswordRegister(ctx, username, params.Password)
	if err != nil {
		global.Logger.Error(err.Error(), zap.Error(err))
	}
	response.SuccessResponse(ctx, 1, result)
	return
}

// Shop Login
// @Summary      Shop Login
// @Description  Shop Login
// @Tags         shop account management
// @Accept       json
// @Produce      json
// @Param        payload body model.ShopLoginInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /shops/auth/public/login
func (c *cShopAuthController) Login(ctx *gin.Context) {
	var params model.ShopLoginInput
	if err := ctx.ShouldBind(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid)
	}
	result, err := service.ShopRegisterService().LoginShop(ctx, params)
	if err != nil {
		global.Logger.Error(err.Error(), zap.Error(err))
		response.ErrorResponse(ctx, 400)
	} else {
		response.SuccessResponse(ctx, 1, result)
	}
	return
}
