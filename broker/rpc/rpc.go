package rpc

import (
	"github.com/codinl/go-logger"
	"google.golang.org/grpc"
	"sync"
)


var lock sync.Mutex
var clientConn *grpc.ClientConn

var LOCAL_TEST = true
var LOCAL_RPC_ADDR = "127.0.0.1:55555"

func dial() (*grpc.ClientConn, error) {
	var c *grpc.ClientConn
	var err error

	if LOCAL_TEST {
		c, err = grpc.Dial(LOCAL_RPC_ADDR, grpc.WithInsecure())
	}else {
		//etcd := viper.GetString("iot_rpc.etcd")
		//registryDir := viper.GetString("iot_rpc.registry_dir")
		//serviceName := viper.GetString("iot_rpc.service_name")

		//logger.Info("etcd =", etcd)
		//logger.Info("registryDir =", registryDir)
		//logger.Info("serviceName =", serviceName)

		//r := NewResolver(registryDir, serviceName)
		//b := grpc.RoundRobin(r)
		//c, err = grpc.Dial(etcd, grpc.WithInsecure(), grpc.WithBalancer(b))
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
