package router

import (
	middlewares "github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/addition/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRouter() *gin.Engine {
	zap.S().Infof("进入router模块")

	Router := gin.Default()
	Router.Use(middlewares.Cors)

	RouterGroup := Router.Group("/up/v1")
	InitAdditionRouter(RouterGroup)
	return Router
}
