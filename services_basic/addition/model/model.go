package model

type Feedback struct {
	Model

	UserID int32  `gorm:"type:int;index"`
	Type   int32  `gorm:"type:int comment '留言类型: 1(留言),2(投诉),3(询问),4(售后),5(求购)'"`
	Title  string `gorm:"type:varchar(100)"`

	Msg  string
	File string `gorm:"type:varchar(200)"`
}

type Address struct {
	Model

	UserID   int32  `gorm:"type:int;index"`
	Province string `gorm:"type:varchar(10)"`
	City     string `gorm:"type:varchar(10)"`
	District string `gorm:"type:varchar(20)"`
	Address  string `gorm:"type:varchar(100)"`
	Name     string `gorm:"type:varchar(20)"`
	Mobile   string `gorm:"type:varchar(11)"`
}
