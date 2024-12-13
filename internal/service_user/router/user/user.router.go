package user

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/internal/service_user/controller/account"
)

type UserRouter struct{}

func (ur *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	//userController, _ := wire.InitUserRouterHandle()
	userRouterPublic := Router.Group("/users")
	{
		userRouterPublic.POST("/register", account.LoginController.Register)
		userRouterPublic.POST("/verify_account", account.LoginController.VerifyOTP)
		userRouterPublic.POST("/otp")
		userRouterPublic.POST("/login", account.LoginController.Login)
		userRouterPublic.POST("/update_pass_register", account.LoginController.UpdatePasswordRegister)
	}

	userRouterPrivate := Router.Group("/user")
	//userRouterPrivate.Use(limit)
	//userRouterPrivate.Use(Authen)
	//userRouterPrivate.Use(Permission)
	{
		userRouterPrivate.GET("/get_info/:id")
	}
}
