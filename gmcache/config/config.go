package config

import (
	"errors"
	"github.com/liyue201/gmcache/utils"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	//grpc server
	RpcPort int //listening port
	Storage StorageCfg
	Log     LogCfg
	Reg     RegistryCfg
}

type StorageCfg struct {
	MemoryLimit   int // in M
	BucketNum     int
	CleanInterval int //in ms
}

type LogCfg struct {
	Dir   string
	File  string
	Level int
}

type RegistryCfg struct {
	ETCD        string
	RegistryDir string
	ServiceName string
	NodeName    string
	NodeAddr    string
	TTL         int
}

//Default Configure
var AppConfig = &Config{
	RpcPort: 55555,
	Storage: StorageCfg{
		MemoryLimit:   1024,
		BucketNum:     10,
		CleanInterval: 20,
	},
	Log: LogCfg{
		Dir:   "./log",
		File:  "gmcache.log",
		Level: 0,  //0 debug, 1 info, 2 warn, 3 error
	},
	Reg: RegistryCfg{
		ETCD:        "127.0.0.1:4001",
		RegistryDir: "/dev",
		ServiceName: "gmcache",
		NodeName:    "node-01",
		NodeAddr:    "127.0.0.1:55555",
		TTL:         5,
	},
}

func InitConfig(path string) error {
	log.Println("InitConfig:", path)

	viper.SetConfigFile(path)
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
	SetInt(&AppConfig.RpcPort, "server.rpc_port")
	SetInt(&AppConfig.Storage.MemoryLimit, "storage.memory_limit")
	SetInt(&AppConfig.Storage.BucketNum, "storage.bucket_num")
	SetInt(&AppConfig.Storage.CleanInterval, "storage.clean_interval")

	SetString(&AppConfig.Log.Dir, "log.dir")
	AppConfig.Log.Dir = utils.AbsPath(AppConfig.Log.Dir)
	SetString(&AppConfig.Log.File, "log.file")
	SetInt(&AppConfig.Log.Level, "log.level")

	SetString(&AppConfig.Reg.ETCD, "registry.etcd")
	SetString(&AppConfig.Reg.RegistryDir, "registry.registry_dir")
	SetString(&AppConfig.Reg.ServiceName, "registry.service_name")
	SetString(&AppConfig.Reg.NodeName, "registry.node_name")
	SetString(&AppConfig.Reg.NodeAddr, "registry.node_addr")
	SetInt(&AppConfig.Reg.TTL, "registry.ttl")

	return checkConfig()
}

func checkConfig() error {
	if AppConfig.Storage.BucketNum <= 0 {
		return errors.New("Bucket number must bee greater than 0")
	}
	//todo : check other items
	return nil
}
