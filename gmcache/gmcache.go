package gmcache

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/judwhite/go-svc/svc"
	"github.com/liyue201/gmcache/gmcache/config"
	"github.com/liyue201/gmcache/utils"
	"log"
	"syscall"
	"time"
)

func Main() {
	app := &gmcache{}
	if err := svc.Run(app, syscall.SIGINT, syscall.SIGTERM); err != nil {
		logger.Println(err)
	}
}

type gmcache struct {
	utils.WaitGroupWrapper
	rpcServer IRpcServer
	storage   *StorageManager
}

func (this *gmcache) Init(env svc.Environment) error {
	AppDir := utils.GetAppDir()

	if err := config.InitConfig(AppDir); err != nil {
		log.Print("Init config:", err)
		return err
	}

	if err := InitLog(); err != nil {
		return err
	}
	this.storage = NewStorageManager(config.AppConfig.BucketNum, int64(config.AppConfig.MemoryLimit)*1024,
		time.Duration(config.AppConfig.CleanInterval*1000))

	addr := fmt.Sprintf("0.0.0.0:%d", config.AppConfig.RpcPort)
	//log.Println("rpc addr:", addr)
	this.rpcServer = NewRpcServer(addr, this.storage)
	return nil
}

func (this *gmcache) Start() error {
	logger.Println("gmcache start")

	this.Wrap(func() {
		this.rpcServer.Run()
	})
	this.Wrap(func() {
		this.storage.Run()
	})
	return nil
}

func (this *gmcache) Stop() error {
	this.rpcServer.Stop()
	this.storage.Stop()

	this.Wait()

	logger.Println("gmcache stopped")
	return nil
}
