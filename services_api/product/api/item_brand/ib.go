package item_brand

import (
	"context"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/api"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/form"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/proto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func GetBrandByItem(ctx *gin.Context) {
	zap.S().Info("调用[ib.GetBrandByItem]")

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	rsp, err := global.ProductClient.GetBrandListByItem(context.Background(), &proto.ItemInfoRequest{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["name"] = value.Name
		reMap["logo"] = value.Logo

		result = append(result, reMap)
	}

	ctx.JSON(http.StatusOK, result)
}

func List(ctx *gin.Context) {
	rsp, err := global.ProductClient.ItemBrandList(context.Background(), &proto.ItemBrandFilterRequest{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	reMap := map[string]interface{}{
		"total": rsp.Total,
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		temp := make(map[string]interface{})
		temp["id"] = value.Id
		zap.S().Info("temp: id == ", value.Id)
		temp["category"] = map[string]interface{}{
			"id":   value.Category.Id,
			"name": value.Category.Name,
		}
		temp["brand"] = map[string]interface{}{
			"id":   value.Brand.Id,
			"name": value.Brand.Name,
			"logo": value.Brand.Logo,
		}

		result = append(result, temp)
	}

	reMap["data"] = result
	ctx.JSON(http.StatusOK, reMap)
}

func Create(ctx *gin.Context) {
	categoryBrandForm := form.ItemBrandForm{}
	if err := ctx.ShouldBindJSON(&categoryBrandForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.ProductClient.CreateItemBrand(context.Background(), &proto.ItemBrandRequest{
		ItemID:  int32(categoryBrandForm.ItemId),
		BrandID: int32(categoryBrandForm.BrandId),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	response := make(map[string]interface{})
	response["id"] = rsp.Id

	ctx.JSON(http.StatusOK, response)
}

func Update(ctx *gin.Context) {
	categoryBrandForm := form.ItemBrandForm{}
	if err := ctx.ShouldBindJSON(&categoryBrandForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.ProductClient.UpdateItemBrand(context.Background(), &proto.ItemBrandRequest{
		Id:      int32(i),
		ItemID:  int32(categoryBrandForm.ItemId),
		BrandID: int32(categoryBrandForm.BrandId),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.ProductClient.DeleteItemBrand(context.Background(), &proto.ItemBrandRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "")
}
