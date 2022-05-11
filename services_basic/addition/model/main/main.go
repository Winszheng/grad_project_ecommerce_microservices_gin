package main

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/product/model"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/product/model/data"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

func main() {
	dsn := "root:mysql2333@tcp(127.0.0.1:3306)/basic_product?charset=utf8mb4&parseTime=True&loc=Local"

	// 没事，不用在意，反正这个logger以后要改的
	newlogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // log leverl
			Colorful:      true,        // 禁用彩色打印
		})

	data.DB, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newlogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名默认单数
		},
	})

	data.DB.AutoMigrate(&model.Item{})
}
