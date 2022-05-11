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

type ServiceConfig struct {
	Name                string               `mapstructure:"name"`
	Port                int                  `mapstructure:"port"`
	AdditionServiceInfo ProductServiceConfig `mapstructure:"addition_service"`
	JWTInfo             JWTConfig            `mapstructure:"jwt"`
	ConsulInfo          ConsulConfig         `mapstructure:"consul"`
}
