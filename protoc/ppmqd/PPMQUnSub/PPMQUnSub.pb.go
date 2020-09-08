// Code generated by protoc-gen-go. DO NOT EDIT.
// source: PPMQUnSub.proto

package PPMQUnSub

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reqest
type Req struct {
	Topics               []string `protobuf:"bytes,1,rep,name=topics,proto3" json:"topics"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Req) Reset()         { *m = Req{} }
func (m *Req) String() string { return proto.CompactTextString(m) }
func (*Req) ProtoMessage()    {}
func (*Req) Descriptor() ([]byte, []int) {
	return fileDescriptor_4e799cda449b5eac, []int{0}
}

func (m *Req) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Req.Unmarshal(m, b)
}
func (m *Req) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Req.Marshal(b, m, deterministic)
}
func (m *Req) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Req.Merge(m, src)
}
func (m *Req) XXX_Size() int {
	return xxx_messageInfo_Req.Size(m)
}
func (m *Req) XXX_DiscardUnknown() {
	xxx_messageInfo_Req.DiscardUnknown(m)
}

var xxx_messageInfo_Req proto.InternalMessageInfo

func (m *Req) GetTopics() []string {
	if m != nil {
		return m.Topics
	}
	return nil
}

// Response
type Resp struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Resp) Reset()         { *m = Resp{} }
func (m *Resp) String() string { return proto.CompactTextString(m) }
func (*Resp) ProtoMessage()    {}
func (*Resp) Descriptor() ([]byte, []int) {
	return fileDescriptor_4e799cda449b5eac, []int{1}
}

func (m *Resp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Resp.Unmarshal(m, b)
}
func (m *Resp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Resp.Marshal(b, m, deterministic)
}
func (m *Resp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Resp.Merge(m, src)
}
func (m *Resp) XXX_Size() int {
	return xxx_messageInfo_Resp.Size(m)
}
func (m *Resp) XXX_DiscardUnknown() {
	xxx_messageInfo_Resp.DiscardUnknown(m)
}

var xxx_messageInfo_Resp proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Req)(nil), "PPMQUnSub.Req")
	proto.RegisterType((*Resp)(nil), "PPMQUnSub.Resp")
}

func init() { proto.RegisterFile("PPMQUnSub.proto", fileDescriptor_4e799cda449b5eac) }

var fileDescriptor_4e799cda449b5eac = []byte{
	// 118 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x0f, 0x08, 0xf0, 0x0d,
	0x0c, 0xcd, 0x0b, 0x2e, 0x4d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x84, 0x0b, 0x28,
	0xc9, 0x72, 0x31, 0x07, 0xa5, 0x16, 0x0a, 0x89, 0x71, 0xb1, 0x95, 0xe4, 0x17, 0x64, 0x26, 0x17,
	0x4b, 0x30, 0x2a, 0x30, 0x6b, 0x70, 0x06, 0x41, 0x79, 0x4a, 0x6c, 0x5c, 0x2c, 0x41, 0xa9, 0xc5,
	0x05, 0x46, 0xc6, 0x5c, 0xac, 0x01, 0x01, 0x41, 0x05, 0xc9, 0x42, 0x5a, 0x5c, 0xac, 0x01, 0xce,
	0x89, 0x39, 0x39, 0x42, 0x7c, 0x7a, 0x08, 0x53, 0x83, 0x52, 0x0b, 0xa5, 0xf8, 0x51, 0xf8, 0xc5,
	0x05, 0x4a, 0x0c, 0x49, 0x6c, 0x60, 0xdb, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x88, 0x0a,
	0x99, 0xbc, 0x80, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PPRpcClient is the client API for PPRpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PPRpcClient interface {
	PCall(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Resp, error)
}

type pPRpcClient struct {
	cc *grpc.ClientConn
}

func NewPPRpcClient(cc *grpc.ClientConn) PPRpcClient {
	return &pPRpcClient{cc}
}

func (c *pPRpcClient) PCall(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Resp, error) {
	out := new(Resp)
	err := c.cc.Invoke(ctx, "/PPMQUnSub.PPRpc/PCall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PPRpcServer is the server API for PPRpc service.
type PPRpcServer interface {
	PCall(context.Context, *Req) (*Resp, error)
}

func RegisterPPRpcServer(s *grpc.Server, srv PPRpcServer) {
	s.RegisterService(&_PPRpc_serviceDesc, srv)
}

func _PPRpc_PCall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PPRpcServer).PCall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PPMQUnSub.PPRpc/PCall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PPRpcServer).PCall(ctx, req.(*Req))
	}
	return interceptor(ctx, in, info, handler)
}

var _PPRpc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "PPMQUnSub.PPRpc",
	HandlerType: (*PPRpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PCall",
			Handler:    _PPRpc_PCall_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "PPMQUnSub.proto",
}