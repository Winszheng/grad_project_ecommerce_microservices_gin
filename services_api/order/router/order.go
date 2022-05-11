package router

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/api/order"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/api/pay"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/middleware"
	"github.com/gin-gonic/gin"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	r := Router.Group("order").Use(middleware.JWTAuth)
	{
		r.GET("", order.GetList)
		r.POST("", order.Create)
		r.GET(":id", order.GetOrder)
		r.DELETE(":id", order.Delete)
	}

	p := Router.Group("pay")
	{
		p.POST("alipay/notify", pay.Notify)
	}
}
