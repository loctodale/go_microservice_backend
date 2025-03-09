package router

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/service_user/router/user"
)

//
//type RouterGroup struct {
//	User user.UserRouterGroup
//}
//
//var RouterGroupApp = new(RouterGroup)

type InitRouter struct{}

func NewInitRouterUser() *InitRouter {
	return &InitRouter{}
}

func (ir *InitRouter) InitRouterUser() *gin.Engine {
	var r *gin.Engine
	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}
	MainGroup := r.Group("/api")
	router := user.UserRouterGroup{}
	userRouter := router.UserRouter
	{
		userRouter.InitUserRouter(MainGroup)
	}
	productRouter := router.ProductRouter
	{
		productRouter.InitProductRouter(MainGroup)
	}
	return r
}
