package router

import (
	middlewares "github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRouter() *gin.Engine {
	zap.S().Infof("进入router模块")

	Router := gin.Default()
	Router.Use(middlewares.Cors)

	RouterGroup := Router.Group("/p/v1")
	InitProductRouter(RouterGroup)
	InitItemRouter(RouterGroup)
	InitBrandRouter(RouterGroup)
	return Router
}
