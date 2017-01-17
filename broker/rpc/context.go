package rpc

import (
	"time"
	"golang.org/x/net/context"
)

//RpcContext is an implementation of context.Context
type RpcContext struct {
	key string
}

func NewRpcContext(key string) context.Context {
	return &RpcContext{key:key}
}

func (*RpcContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*RpcContext) Done() <-chan struct{} {
	return nil
}

func (*RpcContext) Err() error {
	return nil
}

func (this *RpcContext) Value(key interface{}) interface{} {
	return this.key
}
