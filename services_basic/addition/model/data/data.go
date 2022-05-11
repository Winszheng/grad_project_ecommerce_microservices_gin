package data

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/product/model"
	"gorm.io/gorm"
	"math/rand"
)

var (
	DB *gorm.DB
)

func CreatePhoneBrand() {
	brand := model.Brand{
		Model: model.Model{},
	}
	brand.Name = "Apple"
	DB.Create(&brand)

	brand.ID = int32(rand.Intn(100))
	brand.Name = "华为"
	DB.Create(&brand)

	brand.ID = int32(rand.Intn(100))
	brand.Name = "小米"
	DB.Create(&brand)

	brand.ID = int32(rand.Intn(100))
	brand.Name = "荣耀"
	DB.Create(&brand)

	brand.ID = int32(rand.Intn(100))
	brand.Name = "三星"
	DB.Create(&brand)

	brand.ID = 2
	brand.Name = "林氏果业"
	DB.Create(&brand)
}

func InitItem() {
	//DB.AutoMigrate(&model.Item{})
	//item1 := model.Item{
	//	//Model:    model.Model{},
	//	Name: "美妆个护/宠物",
	//	//ParentID: 0,
	//	Parent: nil,
	//	//Sub:   nil,
	//	Level: 1,
	//	//IsTab: false,
	//}
	//DB.Save(&item1)

}
