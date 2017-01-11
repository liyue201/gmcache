package gmcache

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/judwhite/go-svc/svc"
	"github.com/liyue201/gmcache/utils"
	"syscall"
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
}

func (this *gmcache) Init(env svc.Environment) error {
	AppDir := utils.GetAppDir()

	if err := initConfig(AppDir); err != nil {
		return err
	}
	if err := initLog(); err != nil {
		return err
	}

	addr := fmt.Sprintf("0.0.0.0:%d", AppConfig.RpcPort)
	//log.Println("rpc addr:", addr)
	this.rpcServer = NewRpcServer(addr)
	return nil
}

func (this *gmcache) Start() error {
	logger.Println("gmcache start")

	this.Wrap(func() {
		this.rpcServer.Run()
	})
	return nil
}

func (this *gmcache) Stop() error {
	err := this.rpcServer.Stop()
	this.Wait()

	logger.Println("gmcache stopped")
	return err
}
