package global

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/product/config"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/product/proto"
	"gorm.io/gorm"
)

var (
	DB            *gorm.DB
	ServiceConfig config.ServiceConfig
	StockClient   proto.StockClient
)
