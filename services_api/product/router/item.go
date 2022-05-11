package router

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/api/item"
	"github.com/gin-gonic/gin"
)

func InitItemRouter(Router *gin.RouterGroup) {
	CategoryRouter := Router.Group("item")
	{
		CategoryRouter.GET("", item.List)
		CategoryRouter.DELETE("/:id", item.Delete)
		CategoryRouter.GET("/:id", item.GetItem)
		CategoryRouter.POST("", item.Create)
		CategoryRouter.PUT("/:id", item.Update)
	}
}
