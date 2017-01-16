package config

import (
	"github.com/liyue201/gmcache/utils"
	"github.com/spf13/viper"
)

type Config struct {
	//log
	Log LogCfg
}

type LogCfg struct {
	Dir   string
	File  string
	Level int
}

//Default Configure
var AppConfig = &Config{
	Log: LogCfg{
		Dir:   "./log",
		File:  "broker.log",
		Level: 0,
	},
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

	SetString(&AppConfig.Log.Dir, "log.dir")
	AppConfig.Log.Dir = utils.AbsPath(AppConfig.Log.Dir)
	SetString(&AppConfig.Log.File, "log.file")
	SetInt(&AppConfig.Log.Level, "log.level")

	return nil
}

func CheckConfig() error {

	return nil
}