package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/internal/service"
	"go_microservice_backend_api/internal/vo"
	"go_microservice_backend_api/pkg/response"
)

type UserController struct {
	userService service.IUserService
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) Register(c *gin.Context) {
	var params vo.UserRegistrationRequest
	if err := c.ShouldBind(&params); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamInvalid)
	}
	fmt.Printf("Email params: %s", params.Email)
	result := uc.userService.Register(params.Email, params.Purpose)
	if result == 1 {
		response.SuccessResponse(c, response.CodeSuccess, result)

	}
	response.ErrorResponse(c, result)
}
