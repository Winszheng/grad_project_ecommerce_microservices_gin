package global

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/config"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/proto"
	ut "github.com/go-playground/universal-translator"
)

var (
	ServiceConfig config.ServiceConfig

	Trans ut.Translator // 全局翻译器

	UserSrvClient proto.UserClient
)
