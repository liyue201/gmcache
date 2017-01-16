package broker

import (

	"github.com/liyue201/gmcache/broker/rpc"
	"github.com/liyue201/gmcache/proto"
	"golang.org/x/net/context"
	"testing"
)


//go test -v -run="TestSet"
func TestSet(t *testing.T)  {

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

	ret, err := c.Set(context.Background(), &SetArg)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("ret=%#v", ret)
}

//go test -v -run="TestGet"
func TestGet(t *testing.T)  {

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

	ret, err := c.Get(context.Background(), &arg)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("ret=%#v", ret)
}


//go test -v -run="TestDelete"
func TestDelete(t *testing.T)  {

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

	ret, err := c.Delete(context.Background(), &delArg)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("del ret=%#v", ret)
}