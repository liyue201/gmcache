package gmcache

import (
	"github.com/liyue201/gmcache/utils"
	"syscall"
	"github.com/judwhite/go-svc/svc"
	"log"
)

func Main() {

	app := &gmcache{}
	if err := svc.Run(app, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}

type gmcache struct {
	utils.WaitGroupWrapper
}

func (p *gmcache) Init(env svc.Environment) error {
	AppDir := utils.GetAppDir()
	err := initConfig(AppDir)
	return  err
}

func (p *gmcache) Start() error {
	log.Println("gmcache start")

	return nil
}

func (p *gmcache) Stop() error {
	log.Println("gmcache stopped")
	return nil
}
