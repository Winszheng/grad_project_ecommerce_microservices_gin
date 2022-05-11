package initialize

import (
	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/oss/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	viper.AutomaticEnv()
	isDev := viper.GetBool("E_DEV")
	var configFilename string
	if isDev {
		configFilename = "oss/config-dev.yaml"
	} else {
		configFilename = "oss/config-pro.yaml"
	}

	zap.S().Info("configFilename:", configFilename)

	v := viper.New()

	v.SetConfigFile(configFilename)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	err := v.Unmarshal(&global.ServerConfig)

	if err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息: %+v", global.ServerConfig)

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Info("检测到配置信息被修改")
		v.ReadInConfig()
		v.Unmarshal(&global.ServerConfig)
		zap.S().Infof("配置信息修改后: %v", global.ServerConfig)
	})

}
