package gmcache

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"github.com/liyue201/gmcache/proto"
	"log"
)

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

func (rs *RpcServer) Run()  error {
	log.Println("RpcServer is listening on:", rs.addr)

	listener, err := net.Listen("tcp", rs.addr)
	if err != nil {
		log.Println("RpcServer liten error:", err)
		return err
	}

	rs.listener = listener
	return rs.s.Serve(listener)
}

func (rs *RpcServer) Stop() (err error) {
	if rs.listener != nil {
		err = rs.listener.Close()
		rs.listener = nil
	}
	return  err
}

func (rs *RpcServer) Set(ctx context.Context, in *proto.SetOptArg, opts ...grpc.CallOption) (*proto.SetOptRet, error) {
	ret := &proto.SetOptRet{}
	return  ret, nil
}

func (rs *RpcServer) Get(ctx context.Context, in *proto.GetOptArg, opts ...grpc.CallOption) (*proto.GetOptRet, error) {
	ret := &proto.GetOptRet{}
	return  ret, nil
}

func (rs *RpcServer) Delete(ctx context.Context, in *proto.DelOptArg, opts ...grpc.CallOption) (*proto.DelOptRet, error) {
	ret := &proto.DelOptRet{}
	return  ret, nil
}