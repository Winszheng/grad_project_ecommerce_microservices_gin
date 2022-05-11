package global

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/oss/config"
	ut "github.com/go-playground/universal-translator"
)

var (
	Trans ut.Translator

	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	//NacosConfig *config.NacosConfig = &config.NacosConfig{}

)
