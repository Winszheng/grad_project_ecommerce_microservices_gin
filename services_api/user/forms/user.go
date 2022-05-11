package forms

type PasswordLoginForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`

	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`

	CaptchaAnswer string `form:"captcha_answer" json:"captcha_answer" binding:"required,min=5,max=5"`
	CaptchaId     string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

type RegisterUserForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required"`
	Nickname string `form:"nickname" json:"nickname" binding:"required"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
}

type UpdateUserForm struct {
	Name     string `form:"name" json:"name" binding:"required,min=1,max=10"`
	Gender   string `form:"gender" json:"gender" binding:"required,oneof=female male"`
	Birthday string `form:"birthday" json:"birthday" binding:"required,datetime=2006-01-02"`
}
