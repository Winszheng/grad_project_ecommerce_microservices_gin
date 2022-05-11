package cart

import (
	"context"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/api"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/form"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/proto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func GetList(ctx *gin.Context) {
	zap.S().Info("调用[order.GetList]")
	userID, _ := ctx.Get("userId")

	id := userID.(uint)
	rsp, err := global.OrderClient.GetCartList(context.Background(), &proto.UserInfo{Id: int32(id)})
	if err != nil {
		zap.S().Info("[GetList]查询购物车失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	data := make([]interface{}, 0)
	for _, cart := range rsp.Data {
		temp := map[string]interface{}{}
		temp["id"] = cart.CartID
		temp["product_id"] = cart.ProductID
		temp["product_name"] = cart.ProductName
		temp["product_price"] = cart.Price
		temp["nums"] = cart.Num
		temp["checked"] = cart.Checked

		temp["image"] = cart.Image

		data = append(data, temp)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  data,
	})
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url参数格式出错",
		})
		return
	}
	userId, _ := ctx.Get("userId")
	_, err = global.OrderClient.DeleteCart(context.Background(), &proto.CartRequest{
		UserID:    int32(userId.(uint)),
		ProductID: int32(productID),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "成功删除购物车条目",
	})
}

func Create(ctx *gin.Context) {
	cart := form.CreateCartFrom{}
	if err := ctx.ShouldBindJSON(&cart); err != nil {
		api.HandleValidatorError(ctx, err)
	}

	zap.S().Info("before 1")
	if _, err := global.ProductClient.GetProductDetailByID(context.Background(), &proto.ProductID{Id: cart.ProductID}); err != nil {
		zap.S().Info("创建购物车条目时，发现商品不存在")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	zap.S().Info("before 2")
	stock, err := global.StockClient.GetStock(context.Background(), &proto.ProductInfo{ProductID: cart.ProductID})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	if stock.Num < cart.Num {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "库存不足",
		})
		return
	}

	zap.S().Info("before 3")
	userID, _ := ctx.Get("userId")
	rsp, err := global.OrderClient.CreateCart(context.Background(), &proto.CartRequest{

		UserID:    int32(userID.(uint)),
		ProductID: cart.ProductID,
		Num:       cart.Num,
		Checked:   false,
	})
	zap.S().Info("after 2")
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	zap.S().Info("before success")
	ctx.JSON(http.StatusOK, gin.H{"id": rsp.CartID})
}

func Update(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	itemForm := form.UpdateCartForm{}

	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")
	request := proto.CartRequest{
		UserID:    int32(userId.(uint)),
		ProductID: int32(i),
		Num:       itemForm.Num,
		Checked:   false,
	}

	if itemForm.Checked != nil {
		request.Checked = *itemForm.Checked
	}

	_, err = global.OrderClient.UpdateCart(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("更新购物车记录失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}
