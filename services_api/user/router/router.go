package router

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/api"
	middlewares "github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/middleware"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.RouterGroup) {
	r := router.Group("user")

	{
		r.GET("", middlewares.JWTAuth, middlewares.IsAdminAuth, api.GetUserList)
		r.POST("/pwd_login", api.PasswordLogin)
		r.POST("/pwd_register", api.RegisterUser)

		r.GET("detail", middlewares.JWTAuth, api.GetUserDetail)

		r.PATCH("/update", middlewares.JWTAuth, api.UpdateUser)
	}
}
