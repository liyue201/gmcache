package config

import (
	"github.com/liyue201/gmcache/utils"
	"github.com/spf13/viper"
)


type Config struct {
	//grpc server
	//log
	LogDir   string
	LogFile  string
	LogLevel int
}

//Default Configure
var AppConfig = &Config{
	LogDir:   "./log",
	LogFile:  "broker.log",
	LogLevel: 0,
}

func InitConfig(path string) error {
	viper.AddConfigPath(path)
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

	SetString(&AppConfig.LogDir, "log.dir")
	AppConfig.LogDir = utils.AbsPath(AppConfig.LogDir)
	SetString(&AppConfig.LogFile, "log.file")
	SetInt(&AppConfig.LogLevel, "log.level")

	return nil
}

func CheckConfig() error {

	return nil
}