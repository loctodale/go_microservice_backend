package account

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/model"
	"go_microservice_backend_api/internal/service_user/service"
	"go_microservice_backend_api/pkg/response"
)

var LoginController = new(cLoginController)

type cLoginController struct{}

// Verify OTP Login By User
// @Summary      Verify OTP Login By User
// @Description  Verify OTP Login By User
// @Tags         account management
// @Accept       json
// @Produce      json
// @Param        payload body model.VerifyInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /user/verify_account [post]
func (c *cLoginController) VerifyOTP(ctx *gin.Context) {
	var params model.VerifyInput
	if err := ctx.ShouldBind(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid)
	}
	result, err := service.UserLogin().VerifyOTP(ctx, &params)
	if err != nil {
		global.Logger.Error("Error verifying OTP", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrInvalidOTP)
		return
	}
	response.SuccessResponse(ctx, response.CodeSuccess, result)
}

// User Registration documentation
// @Summary      Register user
// @Description  When user registered send OTP to email
// @Tags         account management
// @Accept       json
// @Produce      json
// @Param        payload body model.RegisterInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /user/register [post]
func (c *cLoginController) Register(ctx *gin.Context) {
	var params model.RegisterInput
	if err := ctx.ShouldBind(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid)
		return
	}
	codeStatus, err := service.UserLogin().Register(ctx, &params)
	if err != nil {
		global.Logger.Error("Error registering user OTP", zap.Error(err))
		response.ErrorResponse(ctx, codeStatus)
		return
	}
	response.SuccessResponse(ctx, codeStatus, nil)
}

// Update User Password Register documentation
// @Summary      Update password register user
// @Description  After verify OTP update password
// @Tags         account management
// @Accept       json
// @Produce      json
// @Param        payload body model.UpdatePasswordRegister true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /user/update_pass_register  [post]
func (c *cLoginController) UpdatePasswordRegister(ctx *gin.Context) {
	var params model.UpdatePasswordRegister
	if err := ctx.ShouldBind(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid)
		return
	}
	codeStatus, err := service.UserLogin().UpdatePasswordRegister(ctx, params.UserToken, params.UserPassword)
	if err != nil {
		global.Logger.Error("Error update password register", zap.Error(err))
		response.ErrorResponse(ctx, codeStatus)
		return
	}
	response.SuccessResponse(ctx, codeStatus, nil)
}

// User Login
// @Summary      User Login
// @Description  User Login
// @Tags         account management
// @Accept       json
// @Produce      json
// @Param        payload body model.LoginInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /user/login  [post]
func (c *cLoginController) Login(ctx *gin.Context) {
	var params model.LoginInput
	if err := ctx.ShouldBind(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid)
	}
	codeResult, dataResult, err := service.UserLogin().Login(ctx, params)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken)
	}
	response.SuccessResponse(ctx, codeResult, dataResult)
}
