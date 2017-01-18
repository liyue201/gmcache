package broker

import (
	"github.com/codinl/go-logger"
	"github.com/liyue201/gmcache/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"errors"
	"github.com/liyue201/gmcache/broker/rpc"
)

var ServerInternalError  = errors.New("Server internal error")

type IRpcServer interface {
	Run() error
	Stop() error
}

type RpcServer struct {
	addr     string
	listener net.Listener
	s        *grpc.Server
}

func NewRpcServer(addr string) IRpcServer {
	s := grpc.NewServer()
	rs := &RpcServer{
		addr:    addr,
		s:       s,
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
		rpc.CloseClientConn()
		this.listener = nil
	}
	return err
}

func (this *RpcServer) Set(ctx context.Context, arg *proto.SetOptArg) (*proto.SetOptRet, error) {
	logger.Debugf("RpcServer::Set(): in = %#v", arg)

	ret := &proto.SetOptRet{Code: proto.RCODE_SUCCESS}

	conn, err := rpc.GetClientConn()
	if err != nil {
		logger.Errorf("RpcServer Set: NewRpcServiceClient error")
		ret.Code = proto.RCODE_FAILURE
		return ret, ServerInternalError
	}
	client := proto.NewRpcServiceClient(conn)
	ret, err = client.Set(rpc.NewRpcContext(arg.Key), arg)
	if err != nil {
		logger.Errorf("RpcServer Set:", err)
	}
	return ret, err
}

func (this *RpcServer) Get(ctx context.Context, arg *proto.GetOptArg) (*proto.GetOptRet, error) {
	logger.Debugf("RpcServer::Get(): in = %#v", arg)

	ret := &proto.GetOptRet{Code: proto.RCODE_SUCCESS}

	conn, err := rpc.GetClientConn()
	if err != nil {
		logger.Errorf("RpcServer Get: NewRpcServiceClient error")
		ret.Code = proto.RCODE_FAILURE
		return ret, ServerInternalError
	}
	client := proto.NewRpcServiceClient(conn)
	ret, err = client.Get(rpc.NewRpcContext(arg.Key), arg)
	if err != nil {
		logger.Errorf("RpcServer Set:", err)
	}

	return ret, err
}

func (this *RpcServer) Delete(ctx context.Context, arg *proto.DelOptArg) (*proto.DelOptRet, error) {
	logger.Debugf("RpcServer::Delete(): arg = %#v", arg)

	ret := &proto.DelOptRet{Code: proto.RCODE_SUCCESS}

	conn, err := rpc.GetClientConn()
	if err != nil {
		logger.Errorf("RpcServer Delete: NewRpcServiceClient error")
		ret.Code = proto.RCODE_FAILURE
		return ret, ServerInternalError
	}
	client := proto.NewRpcServiceClient(conn)

	ret, err = client.Delete(rpc.NewRpcContext(arg.Key), arg)
	if err != nil {
		logger.Errorf("RpcServer Delete:", err)
	}
	return ret, nil
}
