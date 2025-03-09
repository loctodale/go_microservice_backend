package router

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/service_shop/router/shop"
)

type InitRouter struct {
}

func NewInitRouterShop() *InitRouter {
	return &InitRouter{}
}

func (ir *InitRouter) InitRouterShop() *gin.Engine {
	var r *gin.Engine
	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	MainGroup := r.Group("/api/shops")
	router := shop.ShopRouterGroup{}
	shopAuthRouter := router.AuthRouter
	{
		shopAuthRouter.InitShopAuth(MainGroup)
	}
	return r
}
