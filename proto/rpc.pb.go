// Code generated by protoc-gen-go.
// source: rpc.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	rpc.proto

It has these top-level messages:
	SetOptArg
	SetOptRet
	GetOptArg
	GetOptRet
	DelOptArg
	DelOptRet
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

// return code
type RCODE int32

const (
	RCODE_SUCCESS RCODE = 0
	RCODE_FAILURE RCODE = 1
)

var RCODE_name = map[int32]string{
	0: "SUCCESS",
	1: "FAILURE",
}
var RCODE_value = map[string]int32{
	"SUCCESS": 0,
	"FAILURE": 1,
}

func (x RCODE) String() string {
	return proto1.EnumName(RCODE_name, int32(x))
}
func (RCODE) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type SetOptArg struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Val []byte `protobuf:"bytes,2,opt,name=val,proto3" json:"val,omitempty"`
	Ttl uint64 `protobuf:"varint,3,opt,name=ttl" json:"ttl,omitempty"`
}

func (m *SetOptArg) Reset()                    { *m = SetOptArg{} }
func (m *SetOptArg) String() string            { return proto1.CompactTextString(m) }
func (*SetOptArg) ProtoMessage()               {}
func (*SetOptArg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type SetOptRet struct {
	Code RCODE `protobuf:"varint,1,opt,name=code,enum=proto.RCODE" json:"code,omitempty"`
}

func (m *SetOptRet) Reset()                    { *m = SetOptRet{} }
func (m *SetOptRet) String() string            { return proto1.CompactTextString(m) }
func (*SetOptRet) ProtoMessage()               {}
func (*SetOptRet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type GetOptArg struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *GetOptArg) Reset()                    { *m = GetOptArg{} }
func (m *GetOptArg) String() string            { return proto1.CompactTextString(m) }
func (*GetOptArg) ProtoMessage()               {}
func (*GetOptArg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type GetOptRet struct {
	Code RCODE  `protobuf:"varint,1,opt,name=code,enum=proto.RCODE" json:"code,omitempty"`
	Val  []byte `protobuf:"bytes,2,opt,name=val,proto3" json:"val,omitempty"`
}

func (m *GetOptRet) Reset()                    { *m = GetOptRet{} }
func (m *GetOptRet) String() string            { return proto1.CompactTextString(m) }
func (*GetOptRet) ProtoMessage()               {}
func (*GetOptRet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type DelOptArg struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *DelOptArg) Reset()                    { *m = DelOptArg{} }
func (m *DelOptArg) String() string            { return proto1.CompactTextString(m) }
func (*DelOptArg) ProtoMessage()               {}
func (*DelOptArg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type DelOptRet struct {
	Code RCODE `protobuf:"varint,1,opt,name=code,enum=proto.RCODE" json:"code,omitempty"`
}

func (m *DelOptRet) Reset()                    { *m = DelOptRet{} }
func (m *DelOptRet) String() string            { return proto1.CompactTextString(m) }
func (*DelOptRet) ProtoMessage()               {}
func (*DelOptRet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func init() {
	proto1.RegisterType((*SetOptArg)(nil), "proto.SetOptArg")
	proto1.RegisterType((*SetOptRet)(nil), "proto.SetOptRet")
	proto1.RegisterType((*GetOptArg)(nil), "proto.GetOptArg")
	proto1.RegisterType((*GetOptRet)(nil), "proto.GetOptRet")
	proto1.RegisterType((*DelOptArg)(nil), "proto.DelOptArg")
	proto1.RegisterType((*DelOptRet)(nil), "proto.DelOptRet")
	proto1.RegisterEnum("proto.RCODE", RCODE_name, RCODE_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for RpcService service

type RpcServiceClient interface {
	Set(ctx context.Context, in *SetOptArg, opts ...grpc.CallOption) (*SetOptRet, error)
	Get(ctx context.Context, in *GetOptArg, opts ...grpc.CallOption) (*GetOptRet, error)
	Delete(ctx context.Context, in *DelOptArg, opts ...grpc.CallOption) (*DelOptRet, error)
}

type rpcServiceClient struct {
	cc *grpc.ClientConn
}

func NewRpcServiceClient(cc *grpc.ClientConn) RpcServiceClient {
	return &rpcServiceClient{cc}
}

func (c *rpcServiceClient) Set(ctx context.Context, in *SetOptArg, opts ...grpc.CallOption) (*SetOptRet, error) {
	out := new(SetOptRet)
	err := grpc.Invoke(ctx, "/proto.RpcService/Set", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rpcServiceClient) Get(ctx context.Context, in *GetOptArg, opts ...grpc.CallOption) (*GetOptRet, error) {
	out := new(GetOptRet)
	err := grpc.Invoke(ctx, "/proto.RpcService/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rpcServiceClient) Delete(ctx context.Context, in *DelOptArg, opts ...grpc.CallOption) (*DelOptRet, error) {
	out := new(DelOptRet)
	err := grpc.Invoke(ctx, "/proto.RpcService/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RpcService service

type RpcServiceServer interface {
	Set(context.Context, *SetOptArg) (*SetOptRet, error)
	Get(context.Context, *GetOptArg) (*GetOptRet, error)
	Delete(context.Context, *DelOptArg) (*DelOptRet, error)
}

func RegisterRpcServiceServer(s *grpc.Server, srv RpcServiceServer) {
	s.RegisterService(&_RpcService_serviceDesc, srv)
}

func _RpcService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetOptArg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.RpcService/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcServiceServer).Set(ctx, req.(*SetOptArg))
	}
	return interceptor(ctx, in, info, handler)
}

func _RpcService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOptArg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.RpcService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcServiceServer).Get(ctx, req.(*GetOptArg))
	}
	return interceptor(ctx, in, info, handler)
}

func _RpcService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelOptArg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.RpcService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcServiceServer).Delete(ctx, req.(*DelOptArg))
	}
	return interceptor(ctx, in, info, handler)
}

var _RpcService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.RpcService",
	HandlerType: (*RpcServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Set",
			Handler:    _RpcService_Set_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _RpcService_Get_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _RpcService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto1.RegisterFile("rpc.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 253 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x90, 0xd1, 0x4a, 0xc3, 0x40,
	0x10, 0x45, 0x5d, 0xd3, 0x56, 0x32, 0x16, 0x09, 0xfb, 0x14, 0x0a, 0x42, 0xcc, 0x53, 0x2a, 0xb6,
	0x0f, 0xf5, 0x03, 0x24, 0x24, 0x31, 0x08, 0x42, 0x61, 0x96, 0x7e, 0x80, 0xc6, 0x41, 0xc4, 0x40,
	0x96, 0x65, 0x28, 0xf8, 0x17, 0x7e, 0xb2, 0xcc, 0x6a, 0x62, 0x4b, 0x29, 0xf4, 0x29, 0x33, 0xf7,
	0x5e, 0xce, 0xcd, 0x0e, 0x84, 0xce, 0x36, 0x4b, 0xeb, 0x3a, 0xee, 0xf4, 0xd8, 0x7f, 0xd2, 0x1c,
	0x42, 0x43, 0xbc, 0xb6, 0x9c, 0xbb, 0x77, 0x1d, 0x41, 0xf0, 0x49, 0x5f, 0xb1, 0x4a, 0x54, 0x16,
	0xa2, 0x8c, 0xa2, 0x6c, 0x5f, 0xda, 0xf8, 0x3c, 0x51, 0xd9, 0x14, 0x65, 0x14, 0x85, 0xb9, 0x8d,
	0x83, 0x44, 0x65, 0x23, 0x94, 0x31, 0x5d, 0xf4, 0x08, 0x24, 0xd6, 0x09, 0x8c, 0x9a, 0xee, 0x8d,
	0x3c, 0xe3, 0x6a, 0x35, 0xfd, 0x2d, 0x5b, 0x62, 0xb1, 0x2e, 0x2b, 0xf4, 0x4e, 0x7a, 0x0d, 0x61,
	0x7d, 0xbc, 0x31, 0x7d, 0xe8, 0xed, 0x93, 0x68, 0x87, 0x3f, 0x28, 0xfc, 0x92, 0xda, 0xa3, 0xfc,
	0x45, 0x6f, 0x9f, 0xc4, 0xbf, 0xbd, 0x81, 0xb1, 0x5f, 0xf5, 0x25, 0x5c, 0x98, 0x4d, 0x51, 0x54,
	0xc6, 0x44, 0x67, 0xb2, 0x3c, 0xe6, 0x4f, 0xcf, 0x1b, 0xac, 0x22, 0xb5, 0xfa, 0x56, 0x00, 0x68,
	0x1b, 0x43, 0x6e, 0xfb, 0xd1, 0x90, 0x9e, 0x43, 0x60, 0x88, 0x75, 0xf4, 0x07, 0x1b, 0xae, 0x3b,
	0xdb, 0x57, 0xa4, 0x7e, 0x0e, 0x41, 0xbd, 0x13, 0xad, 0x0f, 0xa2, 0xff, 0x97, 0xb8, 0x83, 0x49,
	0x49, 0x2d, 0x31, 0x0d, 0xe9, 0xe1, 0x91, 0xb3, 0x7d, 0x05, 0x89, 0x5f, 0x27, 0x5e, 0xb8, 0xff,
	0x09, 0x00, 0x00, 0xff, 0xff, 0xbe, 0xdb, 0x6f, 0xff, 0xf0, 0x01, 0x00, 0x00,
}