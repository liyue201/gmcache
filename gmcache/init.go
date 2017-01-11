package gmcache

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/liyue201/gmcache/utils"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	RpcPort int //grpc listening port

	logDir   string
	logFile  string
	logLevel int
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
		RpcPort:  viper.GetInt("server.rpc_port"),
		logDir:   utils.AbsPath(viper.GetString("log.dir")),
		logFile:  viper.GetString("log.file"),
		logLevel: viper.GetInt("log.level"),
	}
	return nil
}

func initLog() error {
	if !utils.PathExist(AppConfig.logDir) {
		os.MkdirAll(AppConfig.logDir, os.ModePerm)
	}

	err := logger.Init(AppConfig.logDir, AppConfig.logFile, AppConfig.logLevel)
	if err != nil {
		fmt.Println("logger init error err=", err)
		return err
	}

	logger.SetConsole(true)

	fmt.Println("logger init success")
	return nil
}
