package router

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/addition/api/address"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/addition/api/message"
	middlewares "github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/addition/middleware"
	"github.com/gin-gonic/gin"
)

func InitAdditionRouter(Router *gin.RouterGroup) {
	AddressRouter := Router.Group("address")
	{
		AddressRouter.GET("", middlewares.JWTAuth, address.List)
		AddressRouter.DELETE("/:id", middlewares.JWTAuth, address.Delete)
		AddressRouter.POST("", middlewares.JWTAuth, address.New)
		AddressRouter.PUT("/:id", middlewares.JWTAuth, address.Update)
	}
	MessageRouter := Router.Group("message").Use(middlewares.JWTAuth)
	{
		MessageRouter.GET("", message.List)
		MessageRouter.POST("", message.New)
	}
}
