package router

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/api"
	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha)
	}

}
