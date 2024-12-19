package shop

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/internal/service_shop/controller/shop"
)

type AuthRouter struct{}

func (ar *AuthRouter) InitShopAuth(Router *gin.RouterGroup) {
	shopPublicRouter := Router.Group("/auth/public")
	{
		shopPublicRouter.POST("/register", shop.NewAuthController.Register)
		shopPublicRouter.POST("/verifyOTP", shop.NewAuthController.VerifyOTP)
	}

	shopPrivateRouter := Router.Group("/auth/private")
	{
		shopPrivateRouter.POST("/change_password", shop.NewAuthController.ChangePasswordRegister)
	}
}
