package global

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/addition/config"
	"gorm.io/gorm"
)

var (
	DB            *gorm.DB
	ServiceConfig config.ServiceConfig
)
