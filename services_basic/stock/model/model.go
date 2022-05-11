package model

type Stock struct {
	//Model
	ProductID int32 `gorm:"index" json:"product_id"`
	Num       int32 `json:"num"`
}

type Order struct {
	No     string    `gorm:"type:varchar(200);index:idx_order_no,unique;"` // 订单号，不然没必要专门弄个No出来
	Status int32     `gorm:"type:varchar(200)"`                            // 1 表示已扣减 2. 表示已归还
	Detail StockList `gorm:"type:varchar(200)"`
}
