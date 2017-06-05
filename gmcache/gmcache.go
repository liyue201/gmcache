package gmcache

import (
	"flag"
	"fmt"
	"github.com/codinl/go-logger"
	etcd "github.com/coreos/etcd/client"
	"github.com/judwhite/go-svc/svc"
	"github.com/liyue201/gmcache/gmcache/config"
	"github.com/liyue201/gmcache/utils"
	registry "github.com/liyue201/grpc-lb/registry/etcd"
	"log"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var configPath *string = flag.String("c", "", "Use -c <config file path>")

func Main() {
	app := &gmcache{}
	if err := svc.Run(app, syscall.SIGINT, syscall.SIGTERM); err != nil {
		logger.Println(err)
	}
}

type gmcache struct {
	utils.WaitGroupWrapper
	rpcServer    IRpcServer
	storage      *StorageManager
	etcdRegistry *registry.EtcdReigistry
}

func (this *gmcache) Init(env svc.Environment) error {
	defaultConfigPath := utils.GetAppDir() + string(filepath.Separator) + "gmcache.conf"
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
	this.storage = NewStorageManager(config.AppConfig.Storage.BucketNum, int64(config.AppConfig.Storage.MemoryLimit)*1024*1024,
		time.Duration(time.Duration(config.AppConfig.Storage.CleanInterval)*time.Second))

	addr := fmt.Sprintf("0.0.0.0:%d", config.AppConfig.RpcPort)
	this.rpcServer = NewRpcServer(addr, this.storage)

	etcdConfig := etcd.Config{
		Endpoints: strings.Split(config.AppConfig.Reg.ETCD, ","),
	}

	var err error
	this.etcdRegistry, err = registry.NewRegistry(
		registry.Option{
			EtcdConfig:  etcdConfig,
			RegistryDir: config.AppConfig.Reg.RegistryDir,
			ServiceName: config.AppConfig.Reg.ServiceName,
			NodeName:    config.AppConfig.Reg.NodeName,
			NodeAddr:    config.AppConfig.Reg.NodeAddr,
			Ttl:         time.Duration(time.Duration(config.AppConfig.Reg.TTL) * time.Second),
		})

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
		this.etcdRegistry.Register()
	})
	return nil
}

func (this *gmcache) Stop() error {
	this.rpcServer.Stop()
	this.storage.Stop()
	this.etcdRegistry.Deregister()

	this.Wait()

	logger.Println("gmcache stopped")
	return nil
}
