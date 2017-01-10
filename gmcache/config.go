package gmcache

import (
	"github.com/spf13/viper"
)

type Config struct {
	RpcPort int  //grpc listening port
}

var AppConfig *Config

func initConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigFile("gmcache.conf")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	AppConfig = &Config{
		RpcPort: viper.GetInt("server.rpc_port"),
	}
	return  nil
}
