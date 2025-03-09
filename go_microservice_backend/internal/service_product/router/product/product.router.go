package product

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/internal/service_product/controller"
)

type rProductRouter struct{}

func NewProductRouter() *rProductRouter {
	return &rProductRouter{}
}

func (r *rProductRouter) InitProductRouter(router *gin.RouterGroup) {
	productPrivateRouter := router.Group("/product/private")
	{
		productPrivateRouter.POST("/", controller.NewProductController().CreateNewProduct)
		productPrivateRouter.POST("/sku/", controller.NewProductController().CreateNewSKU)
	}
}
