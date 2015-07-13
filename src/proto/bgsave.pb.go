// Code generated by protoc-gen-go.
// source: bgsave.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	bgsave.proto

It has these top-level messages:
	BgSave
*/
package proto

import proto1 "github.com/golang/protobuf/proto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal

type BgSave struct {
}

func (m *BgSave) Reset()         { *m = BgSave{} }
func (m *BgSave) String() string { return proto1.CompactTextString(m) }
func (*BgSave) ProtoMessage()    {}

type BgSave_Key struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *BgSave_Key) Reset()         { *m = BgSave_Key{} }
func (m *BgSave_Key) String() string { return proto1.CompactTextString(m) }
func (*BgSave_Key) ProtoMessage()    {}

type BgSave_Keys struct {
	Names []string `protobuf:"bytes,1,rep,name=names" json:"names,omitempty"`
}

func (m *BgSave_Keys) Reset()         { *m = BgSave_Keys{} }
func (m *BgSave_Keys) String() string { return proto1.CompactTextString(m) }
func (*BgSave_Keys) ProtoMessage()    {}

type BgSave_NullResult struct {
}

func (m *BgSave_NullResult) Reset()         { *m = BgSave_NullResult{} }
func (m *BgSave_NullResult) String() string { return proto1.CompactTextString(m) }
func (*BgSave_NullResult) ProtoMessage()    {}

// Client API for BgSaveService service

type BgSaveServiceClient interface {
	MarkDirty(ctx context.Context, in *BgSave_Key, opts ...grpc.CallOption) (*BgSave_NullResult, error)
	MarkDirties(ctx context.Context, in *BgSave_Keys, opts ...grpc.CallOption) (*BgSave_NullResult, error)
}

type bgSaveServiceClient struct {
	cc *grpc.ClientConn
}

func NewBgSaveServiceClient(cc *grpc.ClientConn) BgSaveServiceClient {
	return &bgSaveServiceClient{cc}
}

func (c *bgSaveServiceClient) MarkDirty(ctx context.Context, in *BgSave_Key, opts ...grpc.CallOption) (*BgSave_NullResult, error) {
	out := new(BgSave_NullResult)
	err := grpc.Invoke(ctx, "/proto.BgSaveService/MarkDirty", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bgSaveServiceClient) MarkDirties(ctx context.Context, in *BgSave_Keys, opts ...grpc.CallOption) (*BgSave_NullResult, error) {
	out := new(BgSave_NullResult)
	err := grpc.Invoke(ctx, "/proto.BgSaveService/MarkDirties", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for BgSaveService service

type BgSaveServiceServer interface {
	MarkDirty(context.Context, *BgSave_Key) (*BgSave_NullResult, error)
	MarkDirties(context.Context, *BgSave_Keys) (*BgSave_NullResult, error)
}

func RegisterBgSaveServiceServer(s *grpc.Server, srv BgSaveServiceServer) {
	s.RegisterService(&_BgSaveService_serviceDesc, srv)
}

func _BgSaveService_MarkDirty_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(BgSave_Key)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(BgSaveServiceServer).MarkDirty(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _BgSaveService_MarkDirties_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(BgSave_Keys)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(BgSaveServiceServer).MarkDirties(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _BgSaveService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.BgSaveService",
	HandlerType: (*BgSaveServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MarkDirty",
			Handler:    _BgSaveService_MarkDirty_Handler,
		},
		{
			MethodName: "MarkDirties",
			Handler:    _BgSaveService_MarkDirties_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}
