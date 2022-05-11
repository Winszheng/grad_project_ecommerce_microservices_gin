package router

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/api/cart"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/middleware"
	"github.com/gin-gonic/gin"
)

func InitCartRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("cart").Use(middleware.JWTAuth)
	{
		GoodsRouter.GET("", cart.GetList)
		GoodsRouter.DELETE("/:id", cart.Delete)
		GoodsRouter.POST("", cart.Create)
		GoodsRouter.PATCH("/:id", cart.Update)
	}
}
