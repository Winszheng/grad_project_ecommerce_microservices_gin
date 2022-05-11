package global

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/config"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/proto"
	ut "github.com/go-playground/universal-translator"
)

var (
	ServiceConfig config.ServiceConfig

	Trans ut.Translator

	ProductClient proto.ProductClient
)
