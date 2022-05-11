package form

type ProductInfoForm struct {
	Name        string  `form:"name" json:"name" binding:"required,min=2,max=100"`
	NormalPrice float32 `form:"normal_price" json:"normal_price" binding:"required,min=0"`
	Stock       int32   `form:"stocks" json:"stocks"`
	ArticleNum  string  `form:"article_num" json:"article_num"`                      // 货号
	ProPrice    float32 `form:"pro_price" json:"pro_price" binding:"required,min=0"` // 销售价
	IsShipFree  *bool   `form:"is_ship_free" json:"is_ship_free" binding:"required"`

	ItemId     int32    `form:"item_id" json:"item_id"`
	Brand      int32    `form:"brand_id" json:"brand_id"`
	Images     []string `form:"images" json:"images"`
	DescImages []string `form:"desc_images" json:"desc_images"`

	FrontImage string `form:"front_image" json:"front_image"`

	Brief string `form:"brief" json:"brief" binding:"required,min=2"`
}
