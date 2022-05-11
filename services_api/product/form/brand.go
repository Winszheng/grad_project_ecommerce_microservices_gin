package form

type BrandForm struct {
	Name string `form:"name" json:"name" binding:"required,min=2,max=20"`
	Logo string `form:"logo" json:"logo" binding:"url"`
}

type ItemBrandForm struct {
	ItemId  int `form:"item_id" json:"item_id" binding:"required"`
	BrandId int `form:"brand_id" json:"brand_id" binding:"required"`
}
