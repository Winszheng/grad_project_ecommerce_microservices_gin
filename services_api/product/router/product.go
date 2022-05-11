package router

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/api/product"
	middlewares "github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/middleware"
	"github.com/gin-gonic/gin"
)

func InitProductRouter(router *gin.RouterGroup) {
	productRouter := router.Group("product")

	{
		productRouter.GET("", product.FilterList)
		productRouter.POST("", middlewares.JWTAuth, middlewares.IsAdminAuth, product.CreateProduct)

		productRouter.GET("/:id", product.GetProduct)
		productRouter.DELETE("/:id", product.DeleteProduct)
		productRouter.GET("/:id/stock", product.GetStock)

		productRouter.PUT("/:id", middlewares.JWTAuth, middlewares.IsAdminAuth, product.Update)
		productRouter.PATCH("/:id", middlewares.JWTAuth, middlewares.IsAdminAuth, product.UpdateStatus)
	}
}
