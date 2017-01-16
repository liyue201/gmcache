package config

import (
	"errors"
	"github.com/liyue201/gmcache/utils"
	"github.com/spf13/viper"
)

type Config struct {
	//grpc server
	RpcPort int //listening port

	//storage
	MemoryLimit   int // in M
	BucketNum     int
	CleanInterval int //in ms

	//log
	LogDir   string
	LogFile  string
	LogLevel int
}

//Default Configure
var AppConfig = &Config{
	RpcPort:       55555,
	MemoryLimit:   1024,
	BucketNum:     10,
	CleanInterval: 20,

	LogDir:   "./log",
	LogFile:  "gmcache.log",
	LogLevel: 0,
}

func InitConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigFile("gmcache.conf")
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
	SetInt(&AppConfig.MemoryLimit, "storage.memory_limit")
	SetInt(&AppConfig.BucketNum, "storage.bucket_num")
	SetInt(&AppConfig.CleanInterval, "storage.clean_interval")

	SetString(&AppConfig.LogDir, "log.dir")
	AppConfig.LogDir = utils.AbsPath(AppConfig.LogDir)
	SetString(&AppConfig.LogFile, "log.file")
	SetInt(&AppConfig.LogLevel, "log.level")

	return  checkConfig()
}

func checkConfig() error {
	if AppConfig.BucketNum <= 0 {
		return errors.New("Bucket number must bee greater than 0")
	}
	return nil
}
