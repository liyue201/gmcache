package gmcache

import (
	"fmt"
	"github.com/codinl/go-logger"
	"github.com/judwhite/go-svc/svc"
	"github.com/liyue201/gmcache/gmcache/config"
	"github.com/liyue201/gmcache/gmcache/registry"
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
	rpcServer  IRpcServer
	storage    *StorageManager
	etcdClient *registry.EtcdReigistryClient
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
	this.storage = NewStorageManager(config.AppConfig.Storage.BucketNum, int64(config.AppConfig.Storage.MemoryLimit)*1024,
		time.Duration(time.Duration(config.AppConfig.Storage.CleanInterval)*time.Second))

	addr := fmt.Sprintf("0.0.0.0:%d", config.AppConfig.RpcPort)
	this.rpcServer = NewRpcServer(addr, this.storage)

	var err error
	this.etcdClient, err = registry.NewClient(config.AppConfig.Reg.ETCD,
		config.AppConfig.Reg.RegistryDir,
		config.AppConfig.Reg.ServiceName,
		config.AppConfig.Reg.NodeName,
		config.AppConfig.Reg.NodeAddr,
		time.Duration(time.Duration(config.AppConfig.Reg.TTL)*time.Second),
	)

	return err
}

func (this *gmcache) Start() error {
	logger.Println("gmcache start")

	this.Wrap(func() {
		this.rpcServer.Run()
	})
	this.Wrap(func() {
		this.storage.Run()
	})
	this.Wrap(func() {
		this.etcdClient.Register()
	})
	return nil
}

func (this *gmcache) Stop() error {
	this.rpcServer.Stop()
	this.storage.Stop()
	this.etcdClient.Unregister()

	this.Wait()

	logger.Println("gmcache stopped")
	return nil
}
