package category

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/internal/service_product/controller"
)

type rCategoryRouter struct{}

func NewCategoryRouter() *rCategoryRouter {
	return &rCategoryRouter{}
}

func (r *rCategoryRouter) InitCategoryRouter(router *gin.RouterGroup) {
	categoryPrivateRouter := router.Group("/category/private")
	{
		categoryPrivateRouter.POST("/", controller.NewCategoryConntroller().CreateNewCategory)
	}
}
