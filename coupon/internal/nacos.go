package internal

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

type nacosConfig struct {
	IpAddr              string `mapstructure:"IpAddr"`
	Port                uint64 `mapstructure:"Port"`
	TimeoutMs           uint64 `mapstructure:"TimeOutMs"`
	NamespaceId         string `mapstructure:"NamespaceId"`
	DataId              string `mapstructure:"DataId"`
	CacheDir            string `mapstructure:"CacheDir"`
	NotLoadCacheAtStart bool   `mapstructure:"NotLoadCacheAtStart"`
	LogDir              string `mapstructure:"LogDir"`
	LogLevel            string `mapstructure:"LogLevel"`
}

var AppConf AppConfig
var nacosConf nacosConfig
var fileName = "./config_center/dev-config.yaml"

func init() {
	InitNacos()
	InitDB()
	InitRedis()
}

func InitTest() {

}

func InitNacos() {
	v := viper.New()
	v.SetConfigFile(fileName)
	v.ReadInConfig()
	v.Unmarshal(&nacosConf)
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: nacosConf.IpAddr,
			Port:   nacosConf.Port,
		},
	}

	clientConfig := constant.ClientConfig{
		TimeoutMs:           nacosConf.TimeoutMs,
		NamespaceId:         nacosConf.NamespaceId,
		CacheDir:            nacosConf.CacheDir,
		NotLoadCacheAtStart: nacosConf.NotLoadCacheAtStart,
		LogDir:              nacosConf.LogDir,
		LogLevel:            nacosConf.LogLevel,
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: nacosConf.DataId,
		Group:  "coupon-dev",
	})
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(content), &AppConf)
}
