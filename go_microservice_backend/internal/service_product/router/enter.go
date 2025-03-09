package router

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/service_product/router/category"
	"go_microservice_backend_api/internal/service_product/router/product"
)

type InitRouter struct {
}

func NewInitRouter() *InitRouter {
	return &InitRouter{}
}

func (s *InitRouter) InitRouterProduct() *gin.Engine {
	var r *gin.Engine
	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	MainGroup := r.Group("/api/shop-op")
	productRouter := product.NewProductRouter().InitProductRouter
	{
		productRouter(MainGroup)
	}
	categoryRouter := category.NewCategoryRouter().InitCategoryRouter
	{
		categoryRouter(MainGroup)
	}
	return r
}
