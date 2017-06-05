package broker

import (
	"fmt"
	"github.com/codinl/go-logger"
	"github.com/liyue201/gmcache/broker/config"
	"github.com/liyue201/gmcache/utils"
	"log"
	"os"
)

func InitLog() error {
	if !utils.PathExist(config.AppConfig.Log.Dir) {
		os.MkdirAll(config.AppConfig.Log.Dir, os.ModePerm)
	}

	err := logger.Init(config.AppConfig.Log.Dir, config.AppConfig.Log.File, config.AppConfig.Log.Level)
	if err != nil {
		fmt.Println("logger init error err=", err)
		return err
	}

	logger.SetConsole(true)

	log.Println("logger init success")
	return nil
}
