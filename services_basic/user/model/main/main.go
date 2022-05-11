package main

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/user/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema" // 用来做命名策略
	"log"
	"os"
	"time"
)

func main() {
	dsn := "root:mysql2333@tcp(127.0.0.1:3306)/basic_user?charset=utf8mb4&parseTime=True&loc=Local"

	newlogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // log leverl
			Colorful:      true,        // 禁用彩色打印
		})
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newlogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名默认单数
		},
	})

	db.AutoMigrate(&model.User{})

}
