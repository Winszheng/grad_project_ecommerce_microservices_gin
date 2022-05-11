package form

type AddressForm struct {
	Province     string `form:"province" json:"province" binding:"required"`
	City         string `form:"city" json:"city" binding:"required"`
	District     string `form:"district" json:"district" binding:"required"`
	Address      string `form:"address" json:"address" binding:"required"`
	SignerName   string `form:"signer_name" json:"signer_name" binding:"required"`
	SignerMobile string `form:"signer_mobile" json:"signer_mobile" binding:"required"`
}

type Feedback struct {
	Type    int32  `form:"type" json:"type" binding:"required,oneof=1 2 3 4 5"`
	Subject string `form:"subject" json:"subject" binding:"required"`
	Msg     string `form:"message" json:"message" binding:"required"`
	File    string `form:"file" json:"file"`
}
