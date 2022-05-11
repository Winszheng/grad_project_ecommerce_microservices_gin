package global

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig() {
	viper.AutomaticEnv()
	isDev := viper.GetBool("E_DEV")
	var configFilename string
	if isDev {
		configFilename = "user/config-dev.yaml"
		zap.S().Info("以开发环境启动")
	} else {
		configFilename = "user/config-pro.yaml"
		zap.S().Info("以生产环境启动")
	}

	v := viper.New()

	v.SetConfigFile(configFilename)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	err := v.Unmarshal(&ServiceConfig)

	if err != nil {
		panic(err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Info("检测到配置信息被修改")
		v.ReadInConfig() // 重新读配置
		v.Unmarshal(ServiceConfig)
		zap.S().Infof("配置信息修改后: %v", ServiceConfig)
	})

}
