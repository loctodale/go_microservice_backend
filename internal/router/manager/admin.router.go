package manager

import "github.com/gin-gonic/gin"

type AdminRouter struct{}

func (ar *AdminRouter) InitAdminRouter(Router *gin.RouterGroup) {
	adminRouterPublic := Router.Group("/admin")
	{
		adminRouterPublic.POST("/login")
	}

	adminRouterPrivate := Router.Group("/admin/user")
	{
		adminRouterPrivate.POST("/active_user")
	}
	
}
