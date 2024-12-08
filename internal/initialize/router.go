package initialize

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/router"
)

func InitRouter() *gin.Engine {
	var r *gin.Engine

	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}
	//middlewares
	//r.Use() //logging
	//r.Use() //cross
	
	// router
	MainGroup := r.Group("/v1/2024")
	{
		MainGroup.GET("/checkStatus")
	}
	managerRouter := router.RouterGroupApp.Manager
	{
		managerRouter.InitUserRouter(MainGroup)
		managerRouter.InitAdminRouter(MainGroup)
	}
	userRouter := router.RouterGroupApp.User

	{
		userRouter.InitUserRouter(MainGroup)
		userRouter.InitProductRouter(MainGroup)

	}

	return r
}
