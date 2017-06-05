package client

import (
	"context"
	"github.com/liyue201/gmcache/proto"
	"google.golang.org/grpc"
)

type Client struct {
	addr string
	conn *grpc.ClientConn
}

func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

func (this *Client) Connect() error {
	c, err := grpc.Dial(this.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	this.conn = c
	return nil
}

func (this *Client) Disconnect() error {
	if this.conn != nil {
		err := this.conn.Close()
		this.conn = nil
		return err
	}
	return nil
}

func (this *Client) Set(key string, value []byte, ttl uint64) error {
	arg := &proto.SetOptArg{
		Key: key,
		Val: value,
		Ttl: ttl,
	}
	c := proto.NewRpcServiceClient(this.conn)
	_, err := c.Set(context.Background(), arg)
	return err
}

func (this *Client) Get(key string) (value []byte, err error) {
	arg := &proto.GetOptArg{
		Key: key,
	}
	c := proto.NewRpcServiceClient(this.conn)
	ret, err := c.Get(context.Background(), arg)
	return ret.Val, err
}

func (this *Client) Delete(key string) error {
	arg := &proto.DelOptArg{
		Key: key,
	}
	c := proto.NewRpcServiceClient(this.conn)
	_, err := c.Delete(context.Background(), arg)
	return err
}
