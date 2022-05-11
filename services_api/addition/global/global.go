package global

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/addition/config"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/addition/proto"
	ut "github.com/go-playground/universal-translator"
)

var (
	ServiceConfig config.ServiceConfig

	Trans ut.Translator

	Client proto.AddressClient
)
