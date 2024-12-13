package shop

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/internal/service_shop/controller/shop"
)

type AuthRouter struct{}

func (ar *AuthRouter) InitShopAuth(Router *gin.RouterGroup) {
	shopPublicRouter := Router.Group("/shops")
	{
		shopPublicRouter.POST("/register", shop.NewAuthController.Register)
	}
}
