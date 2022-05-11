package form

type CreateCartFrom struct {
	ProductID int32 `json:"product_id" bind:"required"`
	Num       int32 `json:"num" bind:"required"`
}

type UpdateCartForm struct {
	Num     int32 `json:"nums" form:"nums" bind:"required"`
	Checked *bool `json:"checked" form:"checked"`
}
