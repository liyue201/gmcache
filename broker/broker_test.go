package broker

import (

	"github.com/liyue201/gmcache/broker/rpc"
	"github.com/liyue201/gmcache/proto"
	"testing"
	"github.com/liyue201/gmcache/broker/config"
	"log"
)

func initTest()  {
	if err := config.InitConfig(`.`); err != nil {
		log.Print("Init config:", err)
		return
	}

	if err := InitLog(); err != nil {
		return
	}
}

//go test -v -run="TestSet"
func TestSet(t *testing.T)  {
	initTest()

	conn, err := rpc.GetClientConn()
	if err != nil {
		t.Error(err)
		return
	}
	c := proto.NewRpcServiceClient(conn)
	if c == nil {
		t.Error("NewRpcServiceClient error")
	}

	SetArg := proto.SetOptArg{
		Key: "test",
		Val: []byte("bbbbbbb"),
		Ttl: 1000,
	}

	ret, err := c.Set(rpc.NewRpcContext(SetArg.Key), &SetArg)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("ret=%#v", ret)
}

//go test -v -run="TestGet"
func TestGet(t *testing.T)  {
	initTest()

	conn, err := rpc.GetClientConn()
	if err != nil {
		t.Error(err)
		return
	}
	c := proto.NewRpcServiceClient(conn)
	if c == nil {
		t.Error("NewRpcServiceClient error")
	}

	arg := proto.GetOptArg{
		Key: "test",
	}

	ret, err := c.Get(rpc.NewRpcContext(arg.Key), &arg)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("ret=%#v", ret)
}

//go test -v -run="TestDelete"
func TestDelete(t *testing.T)  {
	initTest()

	conn, err := rpc.GetClientConn()
	if err != nil {
		t.Error(err)
		return
	}
	c := proto.NewRpcServiceClient(conn)
	if c == nil {
		t.Error("NewRpcServiceClient error")
	}

	delArg := proto.DelOptArg{
		Key: "test",
	}

	ret, err := c.Delete(rpc.NewRpcContext(delArg.Key), &delArg)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("del ret=%#v", ret)
}