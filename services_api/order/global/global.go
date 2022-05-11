package global

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/config"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/proto"
	ut "github.com/go-playground/universal-translator"
)

var (
	ServiceConfig config.ServiceConfig
	Trans         ut.Translator

	OrderClient   proto.OrderClient
	StockClient   proto.StockClient
	ProductClient proto.ProductClient
)
