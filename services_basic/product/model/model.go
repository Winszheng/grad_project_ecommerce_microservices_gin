package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        int32          `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"` // 软删除功能
	IsDeleted bool           `gorm:"column:is_deleted" json:"-"`
}

type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type Item struct {
	Model
	Name string `json:"name" gorm:"type:varchar(25);not null"`

	ParentItemID int32 `json:"parent"`
	ParentItem   *Item `json:"-"`

	Sub []*Item `gorm:"foreignKey:ParentItemID;references:ID" json:"sub_item"`

	Level int32 `gorm:"type:int;not null;default:1" json:"level"`
}

// Brand 品牌
type Brand struct {
	Model
	Name    string `json:"name" gorm:"type:varchar(30);not null;unique"`
	LogoUrl string `json:"logo" gorm:"type:varchar(200);default:''"`
}

type ItemBrand struct {
	Model

	ItemID int32 `gorm:"index:idx_item_brand,unique"`
	Item   Item

	BrandID int32 `gorm:"index:idx_item_brand,unique"`
	Brand   Brand
}

// Product 定义商品表结构
type Product struct {
	Model

	ItemID  int32 `gorm:"type:int32;not null"`
	Item    Item
	BrandID int32 `gorm:"type:int32;not null"`
	Brand   Brand

	IsOnSale   bool `gorm:"default:false"`
	IsShipFree bool `gorm:"default:false"` // 是否包邮

	IsNew bool `gorm:"default:false;not null"`

	Name string `gorm:"not null"`

	ArticleNum string `gorm:"type:varchar(30);not null"`

	// 销量
	SoldNum int32 `gorm:"default:0"`

	FavoriteNum int32 `gorm:"default:0"`
	// 平时价
	NormalPrice float32 `gorm:"not null"`

	ProPrice float32 `gorm:"not null"`
	// 商品简介
	Brief string `gorm:"type:varchar(80);not null"`

	Images     GormList `gorm:"type:varchar(1000)"`
	DescImages GormList `gorm:"type:varchar(1000)"` // 商品描述的图片
	FrontImage string   `gorm:"type:varchar(200)"`
}
