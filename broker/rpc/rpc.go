package rpc

import (
	"github.com/codinl/go-logger"
	"google.golang.org/grpc"
	"sync"
	"github.com/liyue201/gmcache/broker/config"
)


var lock sync.Mutex
var clientConn *grpc.ClientConn

var LOCAL_TEST = false
var LOCAL_RPC_ADDR = "127.0.0.1:55555"

func dial() (*grpc.ClientConn, error) {
	var c *grpc.ClientConn
	var err error

	if LOCAL_TEST {
		c, err = grpc.Dial(LOCAL_RPC_ADDR, grpc.WithInsecure())
	}else {
		logger.Info("etcd =", config.AppConfig.Discovery.Etcd)
		logger.Info("registryDir =", config.AppConfig.Discovery.RegistryDir)
		logger.Info("serviceName =", config.AppConfig.Discovery.ServiceName)

		r := NewResolver(config.AppConfig.Discovery.RegistryDir, config.AppConfig.Discovery.ServiceName)
		b := NewKetamaBalancer(r)
		c, err = grpc.Dial(config.AppConfig.Discovery.Etcd, grpc.WithInsecure(), grpc.WithBalancer(b))
	}

	if err != nil {
		logger.Errorf("grpc dial: %s", err.Error())
	}

	return c, err
}

func GetClientConn() (c *grpc.ClientConn, err error) {
	lock.Lock()
	defer lock.Unlock()

	if clientConn == nil {
		clientConn, err = dial()
	}
	c = clientConn
	return c, err
}

func CloseClientConn() {
	lock.Lock()
	defer lock.Unlock()

	if clientConn != nil {
		clientConn.Close()
		clientConn = nil
	}
}
