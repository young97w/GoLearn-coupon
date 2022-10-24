package internal

type AppConfig struct {
	RedisConfig      RedisConfig      `mapstructure:"redis" json:"redis"`
	ConsulConfig     ConsulConfig     `mapstructure:"consul" json:"consul"`
	AccountSrvConfig AccountSrvConfig `mapstructure:"coupon_srv" json:"coupon_srv"`
	AccountWebConfig AccountWebConfig `mapstructure:"coupon_web" json:"coupon_web"`
	DBConfig         DBConfig         `mapstructure:"db" json:"db"`
	JWTKey           JWTKey           `mapstructure:"jwt_key"json:"jwt_key"`
	Debug            bool             `mapstructure:"debug" json:"debug"`
}

type AccountSrvConfig struct {
	SrvName string   `mapstructure:"srvName" json:"srvName"`
	Host    string   `mapstructure:"host" json:"host"`
	Port    int      `mapstructure:"port" json:"port"`
	Tags    []string `mapstructure:"tags" json:"tags"`
}

type AccountWebConfig struct {
	SrvName string   `mapstructure:"srvName" json:"srvName"`
	Host    string   `mapstructure:"host" json:"host"`
	Port    int      `mapstructure:"port" json:"port"`
	Tags    []string `mapstructure:"tags" json:"tags"`
}

type JWTKey struct {
	SigningKey string `mastructure:"signing_key" json:"signing_key"`
}
