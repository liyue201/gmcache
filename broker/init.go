package broker

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/liyue201/gmcache/utils"
	"github.com/liyue201/gmcache/broker/config"
	"os"
)


func InitLog() error {
	if !utils.PathExist(config.AppConfig.LogDir) {
		os.MkdirAll(config.AppConfig.LogDir, os.ModePerm)
	}

	err := logger.Init(config.AppConfig.LogDir, config.AppConfig.LogFile, config.AppConfig.LogLevel)
	if err != nil {
		fmt.Println("logger init error err=", err)
		return err
	}

	logger.SetConsole(true)

	fmt.Println("logger init success")
	return nil
}
