package gmcache

import (
	"github.com/apsdehal/go-logger"
	"github.com/liyue201/gmcache/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"time"
)

type IRpcServer interface {
	Run() error
	Stop() error
}

type RpcServer struct {
	addr     string
	listener net.Listener
	s        *grpc.Server
	storage  IStorage
}

func NewRpcServer(addr string, storage IStorage) IRpcServer {
	s := grpc.NewServer()
	rs := &RpcServer{
		addr:    addr,
		s:       s,
		storage: storage,
	}
	return rs
}

func (this *RpcServer) Run() error {
	logger.Println("RpcServer is listening on:", this.addr)

	listener, err := net.Listen("tcp", this.addr)
	if err != nil {
		logger.Println("RpcServer liten error:", err)
		return err
	}

	this.listener = listener
	proto.RegisterRpcServiceServer(this.s, this)
	return this.s.Serve(listener)
}

func (this *RpcServer) Stop() (err error) {
	if this.listener != nil {
		err = this.listener.Close()
		this.listener = nil
	}
	return err
}


func (this *RpcServer) Set(ctx context.Context, arg *proto.SetOptArg) (*proto.SetOptRet, error) {
	logger.Debugf("RpcServer::Set(): in = %#v", arg)

	ret := &proto.SetOptRet{Code: proto.RCODE_SUCCESS}
	err := this.storage.Set(arg.Key, arg.Val, time.Duration(arg.Ttl*1000))
	if err != nil {
		ret.Code = proto.RCODE_FAILURE
	}
	return ret, nil
}

func (this *RpcServer) Get(ctx context.Context, arg *proto.GetOptArg) (*proto.GetOptRet, error) {
	logger.Debugf("RpcServer::Get(): in = %#v", arg)

	ret := &proto.GetOptRet{Code: proto.RCODE_SUCCESS}
	item, err := this.storage.Get(arg.Key)
	if err == nil {
		ret.Val = item.value
	} else {
		ret.Code = proto.RCODE_FAILURE
	}
	return ret, nil
}

func (this *RpcServer) Delete(ctx context.Context, arg *proto.DelOptArg) (*proto.DelOptRet, error) {
	logger.Debugf("RpcServer::Delete(): arg = %#v", arg)

	ret := &proto.DelOptRet{Code: proto.RCODE_SUCCESS}
	err := this.storage.Delete(arg.Key)
	if err != nil {
		ret.Code = proto.RCODE_FAILURE
	}
	return ret, nil
}
