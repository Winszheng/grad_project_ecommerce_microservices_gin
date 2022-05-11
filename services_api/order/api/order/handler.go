package order

import (
	"context"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/api"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/form"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/model"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/proto"
	"github.com/gin-gonic/gin"
	alipay "github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func GetList(ctx *gin.Context) {

	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")

	request := proto.OrderFilterRequest{}

	model := claims.(*model.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserID = int32(userId.(uint))
	}

	pages := ctx.DefaultQuery("pageNum", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Page = int32(pagesInt)

	perNums := ctx.DefaultQuery("pageSize", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNum = int32(perNumsInt)

	request.Page = int32(pagesInt)
	request.PagePerNum = int32(perNumsInt)

	rsp, err := global.OrderClient.GetOrderList(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("获取订单列表失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}
	orderList := make([]interface{}, 0)

	for _, item := range rsp.Data {
		tmpMap := map[string]interface{}{}

		tmpMap["id"] = item.Id
		tmpMap["status"] = item.Status
		tmpMap["user"] = item.UserID
		tmpMap["post"] = item.Comment
		tmpMap["total"] = item.Amount
		tmpMap["address"] = item.Address
		tmpMap["name"] = item.Name
		tmpMap["mobile"] = item.Mobile
		tmpMap["order_sn"] = item.No
		tmpMap["id"] = item.Id
		tmpMap["add_time"] = item.CreateTime

		orderList = append(orderList, tmpMap)
	}
	reMap["data"] = orderList
	ctx.JSON(http.StatusOK, reMap)
}

func Create(ctx *gin.Context) {
	orderForm := form.CreateOrderForm{}
	if err := ctx.ShouldBindJSON(&orderForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderClient.CreateOrder(context.WithValue(context.Background(), "ginContext", ctx), &proto.OrderRequest{
		UserID:  int32(userId.(uint)),
		Name:    orderForm.Name,
		Mobile:  orderForm.Mobile,
		Address: orderForm.Address,
		Comment: orderForm.Post,
	})
	if err != nil {
		zap.S().Errorw("新建订单失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	config := global.ServiceConfig
	client, err := alipay.New(global.ServiceConfig.AliPayInfo.AppID, global.ServiceConfig.AliPayInfo.AppPrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = client.LoadAliPayPublicKey(config.AliPayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝的公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = config.AliPayInfo.NotifyURL
	p.ReturnURL = config.AliPayInfo.ReturnURL
	p.Subject = "商城订单-" + rsp.No
	p.OutTradeNo = rsp.No
	p.TotalAmount = strconv.FormatFloat(float64(rsp.Amount), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)

	if err != nil {
		zap.S().Errorw("生成支付url失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	_, err = global.OrderClient.UpdateOrderStatus(context.Background(), &proto.OrderStatusRequest{
		OrderNo: rsp.No,
		Status:  "TRADE_SUCCESS",
	})
	if err != nil {
		zap.S().Info(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         rsp.Id,
		"alipay_url": url.String(),
	})
}

// Delete
func Delete(ctx *gin.Context) {
	id := ctx.Param("id") // 订单的id
	orderID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url参数格式出错",
		})
		return
	}

	rsp, err := global.OrderClient.GetOrderDetail(context.Background(), &proto.OrderRequest{
		OrderID: int32(orderID),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	global.OrderClient.UpdateOrderStatus(context.Background(), &proto.OrderStatusRequest{
		OrderNo: rsp.OrderBasicInfo.No,
		Status:  "TRADE_CLOSED",
	})
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "成功关闭订单",
	})
}

// GetOrder
func GetOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, _ := ctx.Get("userId")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	request := proto.OrderRequest{
		OrderID: int32(i),
	}
	claims, _ := ctx.Get("claims")
	model := claims.(*model.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserID = int32(userId.(uint))
	}

	rsp, err := global.OrderClient.GetOrderDetail(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("获取订单详情失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	reMap := gin.H{}
	temp := rsp.OrderBasicInfo
	reMap["id"] = temp.Id
	reMap["status"] = temp.Status
	reMap["user"] = temp.UserID

	reMap["post"] = temp.Comment
	reMap["total"] = temp.Amount
	reMap["address"] = temp.Address
	reMap["name"] = temp.Name
	reMap["mobile"] = temp.Mobile
	reMap["order_sn"] = temp.No

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.ProductList {
		tmpMap := gin.H{
			"id":    item.GoodsId,
			"name":  item.GoodsName,
			"image": item.GoodsImage,
			"price": item.GoodsPrice,
			"nums":  item.Nums,
		}

		goodsList = append(goodsList, tmpMap)
	}
	reMap["goods"] = goodsList

	ctx.JSON(http.StatusOK, reMap)
}
