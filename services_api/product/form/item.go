package form

// 定义接收前端数据的表单

type ItemForm struct {
	Name           string `form:"name" json:"name" binding:"required,min=2,max=20"`
	ParentCategory int32  `form:"parent" json:"parent"`
	Level          int32  `form:"level" json:"level" binding:"required,oneof=1 2 3"`
}

type UpdateCategoryForm struct {
	Name string `form:"name" json:"name" binding:"required,min=3,max=20"`
}
