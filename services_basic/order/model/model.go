package model

import "time"

type Cart struct {
	Model
	UserID    int32 `gorm:"index"`
	ProductID int32 `gorm:"index"`
	Num       int32
	Checked   bool
}

type Order struct {
	Model

	UserID int32  `gorm:"index"`
	No     string `gorm:"varchar(100);index"`

	Status  string     `grom:"type:varchar(20)"`
	TradeNo string     `gorm:"varchar(100)"`
	Amount  float32    // 订单金额
	PayTime *time.Time `gorm:"type:datetime"`

	Address string `grom:"type:varchar(100)"`
	Name    string `gorm:"type:varchar(20)"`
	Mobile  string `gorm:"type:varchar(11)"`  // 收货人电话号码
	Comment string `gorm:"type:varchar(100)"` // 订单备注
}

type OrderItem struct {
	Model

	OrderID int32 `gorm:"index"`
	UserID  int32 `gorm:"index"`

	ProductName string `gorm:"varchar(100)"`
	Image       string `gorm:"varchar(200)"` // url
	Price       float32
	Num         int32
}
