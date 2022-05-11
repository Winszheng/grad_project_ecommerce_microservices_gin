package brand

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

func List(ctx *gin.Context) {
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("pnum", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := global.ProductClient.BrandList(context.Background(), &proto.BrandFilterRequest{
		Pages:       int32(pnInt),
		PagePerNums: int32(pSizeInt),
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	reMap := make(map[string]interface{})
	reMap["total"] = rsp.Total
	//for _, value := range rsp.Data[pnInt : pnInt*pSizeInt+pSizeInt] {
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["name"] = value.Name
		reMap["logo"] = value.Logo

		result = append(result, reMap)
	}

	reMap["data"] = result

	ctx.JSON(http.StatusOK, reMap)
}

func Create(ctx *gin.Context) {
	brandForm := form.BrandForm{}
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.ProductClient.CreateBrand(context.Background(), &proto.BrandRequest{
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id
	request["name"] = rsp.Name
	request["logo"] = rsp.Logo

	ctx.JSON(http.StatusOK, request)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.ProductClient.DeleteBrand(context.Background(), &proto.BrandRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func Update(ctx *gin.Context) {
	zap.S().Info("调用[brand.Update]")
	brandForm := form.BrandForm{}
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.ProductClient.UpdateBrand(context.Background(), &proto.BrandRequest{
		Id:   int32(i),
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "品牌信息修改完成"})
}

func Get(ctx *gin.Context) {
	idx := ctx.Param("id")
	id, err := strconv.ParseInt(idx, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound) // ctx报错
	}

	result, err := global.ProductClient.GetBrand(context.Background(), &proto.BrandRequest{Id: int32(id)})

	rsp := map[string]interface{}{
		"name": result.Name,
		"logo": result.Logo,
	}

	ctx.JSON(http.StatusOK, rsp)
}
