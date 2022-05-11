package main

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/stock/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

func main() {
	dsn := "root:mysql2333@tcp(127.0.0.1:3306)/basic_stock?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // log leverl
			Colorful:      true,        // 禁用彩色打印
		})

	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名默认单数
		},
	})

	db.AutoMigrate(&model.Order{}, &model.Stock{})
}
