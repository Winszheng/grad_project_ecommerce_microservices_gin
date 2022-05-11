package main

import (
	"fmt"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/router"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/utils"
	mxvalidator "github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/validator"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
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
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServiceConfig.Port = port
		}
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", mxvalidator.ValidateMobile)
		v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "非法的手机号码", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	Router := router.InitRouters()

	zap.S().Infof("启动服务器 端口: %d\n", global.ServiceConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServiceConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}
