package global

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/order/config"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/order/proto"
	"gorm.io/gorm"
)

var (
	DB            *gorm.DB
	ServiceConfig config.ServiceConfig
	StockClient   proto.StockClient
	ProductClient proto.ProductClient
)
