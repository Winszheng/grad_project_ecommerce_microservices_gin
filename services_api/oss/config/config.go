package config

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name       string       `mapstructure:"name" json:"name"`
	Host       string       `mapstructure:"host" json:"host"`
	Tags       []string     `mapstructure:"tags" json:"tags"`
	Port       int          `mapstructure:"port" json:"port"`
	JWTInfo    JWTConfig    `mapstructure:"jwt" json:"jwt"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`
	OssInfo    OssConfig    `mapstructure:"oss_service" json:"oss"`
}

type OssConfig struct {
	AccessID     string `mapstructure:"access_id" json:"key"`
	AccessSecret string `mapstructure:"access_secret" json:"secrect"`
	Host         string `mapstructure:"host" json:"host"`
	UploadDir    string `mapstructure:"upload_dir" json:"upload_dir"`
}
