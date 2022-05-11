package main

import (
	"fmt"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/router"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	global.InitLogger()
	global.InitConfig()
	if err := global.InitTrans("zh"); err != nil {
		panic(err)
	}
	global.InitGrpcClient()
	dev := viper.GetBool("E_DEV")
	if !dev {
		port, err := global.GetFreePort()
		if err == nil {
			global.ServiceConfig.Port = port
		}
	}

	Router := router.InitRouter()

	zap.S().Debugf("启动服务器[api-order] 端口: %d\n", global.ServiceConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServiceConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}
