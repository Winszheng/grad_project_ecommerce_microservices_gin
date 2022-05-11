package item

import (
	"context"
	"encoding/json"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/api"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/form"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/proto"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func List(ctx *gin.Context) {

	zap.S().Info("调用[item.List]")
	r, err := global.ProductClient.GetAllItemList(context.Background(), &empty.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	data := make([]interface{}, 0)
	err = json.Unmarshal([]byte(r.JsonData), &data)
	if err != nil {
		zap.S().Errorw("[List] 查询 【分类列表】失败： ", err.Error())
	}

	ctx.JSON(http.StatusOK, data)
}

func GetItem(ctx *gin.Context) {
	zap.S().Info("调用[product.GetItem]")
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	reMap := make(map[string]interface{})
	subCategorys := make([]interface{}, 0)
	if r, err := global.ProductClient.GetSubItem(context.Background(), &proto.ItemListRequest{
		Id: int32(i),
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	} else {

		for _, value := range r.SubItem {
			subCategorys = append(subCategorys, map[string]interface{}{
				"id":              value.Id,
				"name":            value.Name,
				"level":           value.Level,
				"parent_category": value.ParentItem,
				"is_tab":          value.IsTab,
			})
		}
		reMap["id"] = r.Info.Id
		reMap["name"] = r.Info.Name
		reMap["level"] = r.Info.Level
		reMap["parent_category"] = r.Info.ParentItem
		reMap["is_tab"] = r.Info.IsTab
		reMap["sub_item"] = subCategorys

		ctx.JSON(http.StatusOK, reMap)
	}
	return
}

func Create(ctx *gin.Context) {
	categoryForm := form.ItemForm{}
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.ProductClient.CreateItem(context.Background(), &proto.ItemInfoRequest{
		Name: categoryForm.Name,
		//IsTab:      *categoryForm.IsTab,
		Level:      categoryForm.Level,
		ParentItem: categoryForm.ParentCategory,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id
	request["name"] = rsp.Name
	request["parent"] = rsp.ParentItem
	request["level"] = rsp.Level
	request["is_tab"] = rsp.IsTab

	ctx.JSON(http.StatusOK, request)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	//1. 先查询出该分类写的所有子分类
	//2. 将所有的分类全部逻辑删除
	//3. 将该分类下的所有的商品逻辑删除
	_, err = global.ProductClient.DeleteItem(context.Background(), &proto.DeleteItemRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func Update(ctx *gin.Context) {
	categoryForm := form.UpdateCategoryForm{}
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	request := &proto.ItemInfoRequest{
		Id:   int32(i),
		Name: categoryForm.Name,
	}
	//if categoryForm.IsTab != nil {
	//	request.IsTab = *categoryForm.IsTab
	//}
	_, err = global.ProductClient.UpdateItem(context.Background(), request)
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
