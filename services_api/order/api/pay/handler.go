package pay

import (
	"context"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/proto"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"net/http"
)

func Notify(ctx *gin.Context) {
	zap.S().Info("调用[pav.Notify]")

	client, err := alipay.New(global.ServiceConfig.AliPayInfo.AppID, global.ServiceConfig.AliPayInfo.AppPrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = client.LoadAliPayPublicKey(global.ServiceConfig.AliPayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝的公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	noti, err := client.GetTradeNotification(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	_, err = global.OrderClient.UpdateOrderStatus(context.Background(), &proto.OrderStatusRequest{
		OrderNo: noti.OutTradeNo,
		Status:  string(noti.TradeStatus),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	ctx.String(http.StatusOK, "success")
}
