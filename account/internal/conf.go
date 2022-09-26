package internal

type AppConfig struct {
	RedisConfig      RedisConfig      `mapstructure:"redis" json:"redis"`
	ConsulConfig     ConsulConfig     `mapstructure:"consul" json:"consul"`
	ProductSrvConfig AccountSrvConfig `mapstructure:"account_srv" json:"account_srv"`
	ProductWebConfig AccountWebConfig `mapstructure:"account_web" json:"account_web"`
	DBConfig         DBConfig         `mapstructure:"db" json:"db"`
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
