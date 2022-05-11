package config

type ProductServiceConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"secret"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type OrderService struct {
	Name string `mapstructure:"name"`
}

type AliPayConfig struct {
	AppID         string `mapstructure:"app_id"`            // 应用id
	AppPrivateKey string `mapstructure:"app_private_key"`   // 应用私钥
	AliPublicKey  string `mapstructure:"alipay_public_key"` // 支付宝公钥
	NotifyURL     string `mapstructture:"notify_url"`       // 回调url
	ReturnURL     string `mapstructture:"return_url"`       // 支付成功后重定向到该地址
}

type ServiceConfig struct {
	Name           string               `mapstructure:"name"`
	Port           int                  `mapstructure:"port"`
	JWTInfo        JWTConfig            `mapstructure:"jwt"`
	ConsulInfo     ConsulConfig         `mapstructure:"consul"`
	OrderSrvConfig ProductServiceConfig `mapstructure:"order_srv"`

	ProductInfo ProductServiceConfig `mapstructure:"product_srv" json:"product_srv"`
	StockInfo   ProductServiceConfig `mapstructure:"stock_srv" json:"stock_srv"`
	AliPayInfo  AliPayConfig         `mapstructure:"alipay" json:"alipay"`
}
