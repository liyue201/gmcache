package broker

import (
	"flag"
	"fmt"
	"github.com/codinl/go-logger"
	"github.com/judwhite/go-svc/svc"
	"github.com/liyue201/gmcache/broker/config"
	"github.com/liyue201/gmcache/utils"
	"github.com/liyue201/martini"
	"log"
	"path/filepath"
	"syscall"
)

var configPath *string = flag.String("c", "", "Use -c <config file path>")

func Main() {
	flag.Parse()
	app := &broker{}
	if err := svc.Run(app, syscall.SIGINT, syscall.SIGTERM); err != nil {
		logger.Println(err)
	}
}

type broker struct {
	utils.WaitGroupWrapper
	rpcServer  IRpcServer
	httpServer *martini.ClassicMartini
}

func (this *broker) Init(env svc.Environment) error {

	defaultConfigPath := utils.GetAppDir() + string(filepath.Separator) + "broker.conf"
	if *configPath == "" {
		*configPath = defaultConfigPath
	}
	if err := config.InitConfig(*configPath); err != nil {
		log.Print("Init config:", err)
		return err
	}

	if err := InitLog(); err != nil {
		return err
	}
	rpcAddr := fmt.Sprintf("0.0.0.0:%d", config.AppConfig.Service.RpcPort)
	this.rpcServer = NewRpcServer(rpcAddr)

	httpServer := martini.Classic()
	InitRouter(httpServer)
	this.httpServer = httpServer

	return nil
}

func (this *broker) Start() error {
	logger.Println("broker start")

	this.Wrap(func() {
		this.rpcServer.Run()
	})

	this.Wrap(func() {
		httpAddr := fmt.Sprintf("0.0.0.0:%d", config.AppConfig.Service.HttpPort)
		this.httpServer.RunOnAddr(httpAddr)
	})

	return nil
}

func (this *broker) Stop() error {
	this.rpcServer.Stop()
	this.httpServer.Stop()
	this.Wait()

	logger.Println("broker stopped")
	return nil
}
