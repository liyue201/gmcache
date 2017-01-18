package broker

import (
	"github.com/liyue201/gmcache/broker/config"
	"github.com/liyue201/gmcache/utils"
	"log"
	"github.com/codinl/go-logger"
	"github.com/judwhite/go-svc/svc"
	"syscall"
	"fmt"
)

func Main() {
	app := &broker{}
	if err := svc.Run(app, syscall.SIGINT, syscall.SIGTERM); err != nil {
		logger.Println(err)
	}
}

type broker struct {
	utils.WaitGroupWrapper
	rpcServer  IRpcServer
}

func (this *broker) Init(env svc.Environment) error {
	AppDir := utils.GetAppDir()

	if err := config.InitConfig(AppDir); err != nil {
		log.Print("Init config:", err)
		return err
	}

	if err := InitLog(); err != nil {
		return err
	}
	addr := fmt.Sprintf("0.0.0.0:%d", config.AppConfig.Service.RpcPort)
	this.rpcServer = NewRpcServer(addr)

	return nil
}

func (this *broker) Start() error {
	logger.Println("broker start")

	this.Wrap(func() {
		this.rpcServer.Run()
	})
	return nil
}

func (this *broker) Stop() error {
	this.rpcServer.Stop()

	this.Wait()

	logger.Println("broker stopped")
	return nil
}

