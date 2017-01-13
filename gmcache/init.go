package gmcache

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/liyue201/gmcache/utils"
	"github.com/spf13/viper"
	"os"
	"errors"
)

var AppConfig *Config

type Config struct {
	//grpc server
	RpcPort int //listening port

	//storage
	BucketNum int

	//log
	logDir   string
	logFile  string
	logLevel int
}

func NewDefaultConfig() *Config {
	return &Config{
		RpcPort:   55555,
		BucketNum: 10,
		logDir:    "./log",
		logFile:   "gmcache.log",
		logLevel:  0,
	}
}

func InitConfig(path string) error {
	AppConfig = NewDefaultConfig()

	viper.AddConfigPath(path)
	viper.SetConfigFile("gmcache.conf")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	AppConfig = &Config{
		RpcPort:   viper.GetInt("server.rpc_port"),
		BucketNum: viper.GetInt("storage.bucket_num"),
		logDir:    utils.AbsPath(viper.GetString("log.dir")),
		logFile:   viper.GetString("log.file"),
		logLevel:  viper.GetInt("log.level"),
	}
	return nil
}

func CheckConfig() error {
	if AppConfig.BucketNum <= 0 {
		return errors.New("Bucket number must bee greater than 0")
	}
	return nil
}

func InitLog() error {
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
