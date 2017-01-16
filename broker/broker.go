package broker

import (
	//"github.com/apsdehal/go-logger"
	"github.com/liyue201/gmcache/broker/config"
	//"github.com/liyue201/gmcache/broker/rpc"
	//"github.com/liyue201/gmcache/proto"
	"github.com/liyue201/gmcache/utils"
	//"golang.org/x/net/context"
	"log"
	//"time"
)

func Main() {
	AppDir := utils.GetAppDir()

	if err := config.InitConfig(AppDir); err != nil {
		log.Print("Init config:", err)
		return
	}

	if err := InitLog(); err != nil {
		return
	}
}
