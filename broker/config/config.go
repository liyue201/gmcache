package config

import (
	"github.com/liyue201/gmcache/utils"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	//log
	Service   ServiceCfg
	Discovery DiscoveryCfg
	Log       LogCfg
}

type ServiceCfg struct {
	HttpPort int
	RpcPort  int
}

type DiscoveryCfg struct {
	Etcd        string
	RegistryDir string
	ServiceName string
}

type LogCfg struct {
	Dir   string
	File  string
	Level int
}

//Default Configure
var AppConfig = &Config{
	Service: ServiceCfg{
		HttpPort: 8001,
		RpcPort: 8002,
	},
	Discovery: DiscoveryCfg{
		Etcd:        "127.0.0.1:4001",
		RegistryDir: "/dev",
		ServiceName: "gmcache",
	},
	Log: LogCfg{
		Dir:   "./log",
		File:  "broker.log",
		Level: 0,
	},
}

func InitConfig(path string) error {
	log.Println("InitConfig: path=", path)

	viper.AddConfigPath(path)
	//viper.AddConfigPath(path)
	viper.SetConfigFile("broker.conf")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	SetString := func(in *string, key string) {
		if viper.IsSet(key) {
			*in = viper.GetString(key)
		}
	}
	SetInt := func(in *int, key string) {
		if viper.IsSet(key) {
			*in = viper.GetInt(key)
		}
	}
	SetInt(&AppConfig.Service.HttpPort, "service.http_port")
	SetInt(&AppConfig.Service.RpcPort, "service.rpc_port")

	SetString(&AppConfig.Discovery.Etcd, "discovery.etcd")
	SetString(&AppConfig.Discovery.RegistryDir, "discovery.registry_dir")
	SetString(&AppConfig.Discovery.ServiceName, "discovery.service_name")

	SetString(&AppConfig.Log.Dir, "log.dir")
	AppConfig.Log.Dir = utils.AbsPath(AppConfig.Log.Dir)
	SetString(&AppConfig.Log.File, "log.file")
	SetInt(&AppConfig.Log.Level, "log.level")

	return nil
}

func CheckConfig() error {

	return nil
}
