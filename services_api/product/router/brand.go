package router

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/api/brand"
	ib "github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/api/item_brand"
	"github.com/gin-gonic/gin"
)

func InitBrandRouter(Router *gin.RouterGroup) {
	BrandRouter := Router.Group("brand")
	{
		BrandRouter.GET("", brand.List)
		BrandRouter.GET("/:id", brand.Get)
		BrandRouter.DELETE("/:id", brand.Delete)
		BrandRouter.POST("", brand.Create)
		BrandRouter.PUT("/:id", brand.Update)
	}

	CategoryBrandRouter := Router.Group("item_brand")
	{
		CategoryBrandRouter.GET("", ib.List)
		CategoryBrandRouter.POST("", ib.Create)
		CategoryBrandRouter.DELETE("/:id", ib.Delete)
		CategoryBrandRouter.PUT("/:id", ib.Update)
		CategoryBrandRouter.GET("/:id", ib.GetBrandByItem)
	}
}
