// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/kv-server.proto

package kvstore_methods

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

type Value struct {
	Values               []byte   `protobuf:"bytes,1,opt,name=values,proto3" json:"values,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Value) Reset()         { *m = Value{} }
func (m *Value) String() string { return proto.CompactTextString(m) }
func (*Value) ProtoMessage()    {}
func (*Value) Descriptor() ([]byte, []int) {
	return fileDescriptor_1477b0e33c60977b, []int{0}
}

func (m *Value) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Value.Unmarshal(m, b)
}
func (m *Value) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Value.Marshal(b, m, deterministic)
}
func (m *Value) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Value.Merge(m, src)
}
func (m *Value) XXX_Size() int {
	return xxx_messageInfo_Value.Size(m)
}
func (m *Value) XXX_DiscardUnknown() {
	xxx_messageInfo_Value.DiscardUnknown(m)
}

var xxx_messageInfo_Value proto.InternalMessageInfo

func (m *Value) GetValues() []byte {
	if m != nil {
		return m.Values
	}
	return nil
}

type Values struct {
	Values               []*Value `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
	Cursor               int32    `protobuf:"varint,2,opt,name=Cursor,proto3" json:"Cursor,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Values) Reset()         { *m = Values{} }
func (m *Values) String() string { return proto.CompactTextString(m) }
func (*Values) ProtoMessage()    {}
func (*Values) Descriptor() ([]byte, []int) {
	return fileDescriptor_1477b0e33c60977b, []int{1}
}

func (m *Values) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Values.Unmarshal(m, b)
}
func (m *Values) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Values.Marshal(b, m, deterministic)
}
func (m *Values) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Values.Merge(m, src)
}
func (m *Values) XXX_Size() int {
	return xxx_messageInfo_Values.Size(m)
}
func (m *Values) XXX_DiscardUnknown() {
	xxx_messageInfo_Values.DiscardUnknown(m)
}

var xxx_messageInfo_Values proto.InternalMessageInfo

func (m *Values) GetValues() []*Value {
	if m != nil {
		return m.Values
	}
	return nil
}

func (m *Values) GetCursor() int32 {
	if m != nil {
		return m.Cursor
	}
	return 0
}

type Key struct {
	Key                  string   `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Key) Reset()         { *m = Key{} }
func (m *Key) String() string { return proto.CompactTextString(m) }
func (*Key) ProtoMessage()    {}
func (*Key) Descriptor() ([]byte, []int) {
	return fileDescriptor_1477b0e33c60977b, []int{2}
}

func (m *Key) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Key.Unmarshal(m, b)
}
func (m *Key) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Key.Marshal(b, m, deterministic)
}
func (m *Key) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Key.Merge(m, src)
}
func (m *Key) XXX_Size() int {
	return xxx_messageInfo_Key.Size(m)
}
func (m *Key) XXX_DiscardUnknown() {
	xxx_messageInfo_Key.DiscardUnknown(m)
}

var xxx_messageInfo_Key proto.InternalMessageInfo

func (m *Key) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type KVPair struct {
	Key                  *Key     `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                *Value   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KVPair) Reset()         { *m = KVPair{} }
func (m *KVPair) String() string { return proto.CompactTextString(m) }
func (*KVPair) ProtoMessage()    {}
func (*KVPair) Descriptor() ([]byte, []int) {
	return fileDescriptor_1477b0e33c60977b, []int{3}
}

func (m *KVPair) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KVPair.Unmarshal(m, b)
}
func (m *KVPair) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KVPair.Marshal(b, m, deterministic)
}
func (m *KVPair) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KVPair.Merge(m, src)
}
func (m *KVPair) XXX_Size() int {
	return xxx_messageInfo_KVPair.Size(m)
}
func (m *KVPair) XXX_DiscardUnknown() {
	xxx_messageInfo_KVPair.DiscardUnknown(m)
}

var xxx_messageInfo_KVPair proto.InternalMessageInfo

func (m *KVPair) GetKey() *Key {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *KVPair) GetValue() *Value {
	if m != nil {
		return m.Value
	}
	return nil
}

// whether the operation is ok
type OperationOK struct {
	Ok                   bool     `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OperationOK) Reset()         { *m = OperationOK{} }
func (m *OperationOK) String() string { return proto.CompactTextString(m) }
func (*OperationOK) ProtoMessage()    {}
func (*OperationOK) Descriptor() ([]byte, []int) {
	return fileDescriptor_1477b0e33c60977b, []int{4}
}

func (m *OperationOK) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OperationOK.Unmarshal(m, b)
}
func (m *OperationOK) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OperationOK.Marshal(b, m, deterministic)
}
func (m *OperationOK) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OperationOK.Merge(m, src)
}
func (m *OperationOK) XXX_Size() int {
	return xxx_messageInfo_OperationOK.Size(m)
}
func (m *OperationOK) XXX_DiscardUnknown() {
	xxx_messageInfo_OperationOK.DiscardUnknown(m)
}

var xxx_messageInfo_OperationOK proto.InternalMessageInfo

func (m *OperationOK) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

// the args of scan
type ScanArgs struct {
	Cursor               int32    `protobuf:"varint,1,opt,name=cursor,proto3" json:"cursor,omitempty"`
	UseKey               bool     `protobuf:"varint,4,opt,name=useKey,proto3" json:"useKey,omitempty"`
	Match                *Key     `protobuf:"bytes,2,opt,name=match,proto3" json:"match,omitempty"`
	Count                int32    `protobuf:"varint,3,opt,name=count,proto3" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ScanArgs) Reset()         { *m = ScanArgs{} }
func (m *ScanArgs) String() string { return proto.CompactTextString(m) }
func (*ScanArgs) ProtoMessage()    {}
func (*ScanArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_1477b0e33c60977b, []int{5}
}

func (m *ScanArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScanArgs.Unmarshal(m, b)
}
func (m *ScanArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScanArgs.Marshal(b, m, deterministic)
}
func (m *ScanArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScanArgs.Merge(m, src)
}
func (m *ScanArgs) XXX_Size() int {
	return xxx_messageInfo_ScanArgs.Size(m)
}
func (m *ScanArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_ScanArgs.DiscardUnknown(m)
}

var xxx_messageInfo_ScanArgs proto.InternalMessageInfo

func (m *ScanArgs) GetCursor() int32 {
	if m != nil {
		return m.Cursor
	}
	return 0
}

func (m *ScanArgs) GetUseKey() bool {
	if m != nil {
		return m.UseKey
	}
	return false
}

func (m *ScanArgs) GetMatch() *Key {
	if m != nil {
		return m.Match
	}
	return nil
}

func (m *ScanArgs) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

func init() {
	proto.RegisterType((*Value)(nil), "kvstore.methods.Value")
	proto.RegisterType((*Values)(nil), "kvstore.methods.Values")
	proto.RegisterType((*Key)(nil), "kvstore.methods.Key")
	proto.RegisterType((*KVPair)(nil), "kvstore.methods.KVPair")
	proto.RegisterType((*OperationOK)(nil), "kvstore.methods.OperationOK")
	proto.RegisterType((*ScanArgs)(nil), "kvstore.methods.ScanArgs")
}

func init() { proto.RegisterFile("proto/kv-server.proto", fileDescriptor_1477b0e33c60977b) }

var fileDescriptor_1477b0e33c60977b = []byte{
	// 349 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xc1, 0x4b, 0xeb, 0x40,
	0x10, 0xc6, 0x9b, 0xa4, 0x09, 0x7d, 0xd3, 0xc7, 0x7b, 0x32, 0xd4, 0xb6, 0x16, 0xc5, 0xb2, 0x07,
	0x29, 0xa2, 0x11, 0x2a, 0xde, 0x44, 0x10, 0x05, 0x0f, 0x39, 0xb4, 0xa4, 0xd0, 0xa3, 0x10, 0xe3,
	0x60, 0x4b, 0xda, 0x6e, 0xd9, 0xdd, 0x04, 0x0a, 0x5e, 0xfc, 0xcf, 0x25, 0x93, 0x14, 0x8a, 0x69,
	0xf1, 0x94, 0xfd, 0xb2, 0xdf, 0x37, 0xdf, 0x2f, 0x43, 0xe0, 0x78, 0xad, 0xa4, 0x91, 0x37, 0x49,
	0x76, 0xad, 0x49, 0x65, 0xa4, 0x7c, 0xd6, 0xf8, 0x3f, 0xc9, 0xb4, 0x91, 0x8a, 0xfc, 0x25, 0x99,
	0x99, 0x7c, 0xd7, 0xe2, 0x1c, 0xdc, 0x69, 0xb4, 0x48, 0x09, 0xdb, 0xe0, 0x65, 0xf9, 0x41, 0x77,
	0xad, 0xbe, 0x35, 0xf8, 0x1b, 0x96, 0x4a, 0x8c, 0xc1, 0x63, 0x83, 0x46, 0x7f, 0xc7, 0xe1, 0x0c,
	0x9a, 0xc3, 0xb6, 0xff, 0x63, 0x98, 0xcf, 0xc6, 0x6d, 0x32, 0x9f, 0xf8, 0x94, 0x2a, 0x2d, 0x55,
	0xd7, 0xee, 0x5b, 0x03, 0x37, 0x2c, 0x95, 0xe8, 0x80, 0x13, 0xd0, 0x06, 0x8f, 0xf8, 0xc1, 0x6d,
	0x7f, 0xc2, 0xfc, 0x28, 0x5e, 0xc1, 0x0b, 0xa6, 0xe3, 0x68, 0xae, 0xf0, 0x02, 0x9c, 0xa4, 0xbc,
	0x6b, 0x0e, 0x5b, 0x95, 0x9e, 0x80, 0x36, 0x61, 0x6e, 0xc0, 0x2b, 0x70, 0xb9, 0x8c, 0x1b, 0x0e,
	0x13, 0x15, 0x26, 0x71, 0x06, 0xcd, 0xd1, 0x9a, 0x54, 0x64, 0xe6, 0x72, 0x35, 0x0a, 0xf0, 0x1f,
	0xd8, 0x32, 0xe1, 0x8e, 0x46, 0x68, 0xcb, 0x44, 0x7c, 0x42, 0x63, 0x12, 0x47, 0xab, 0x47, 0xf5,
	0xc1, 0xec, 0x71, 0xc1, 0x6e, 0x15, 0xec, 0x85, 0xca, 0xdf, 0xa7, 0x9a, 0x72, 0xee, 0x3a, 0xe7,
	0x4a, 0x85, 0x97, 0xe0, 0x2e, 0x23, 0x13, 0xcf, 0x4a, 0x90, 0xfd, 0xc8, 0x85, 0x05, 0x5b, 0xe0,
	0xc6, 0x32, 0x5d, 0x99, 0xae, 0xc3, 0xa3, 0x0b, 0x31, 0xfc, 0xb2, 0x01, 0x82, 0xe9, 0x84, 0x54,
	0x36, 0x8f, 0x49, 0xe1, 0x1d, 0x38, 0x2f, 0x64, 0x70, 0xef, 0xa0, 0xde, 0x81, 0xef, 0x14, 0x35,
	0x7c, 0x00, 0x67, 0x42, 0x06, 0x3b, 0xd5, 0x18, 0x2f, 0xb6, 0x77, 0x5a, 0xb9, 0xd8, 0xd9, 0x08,
	0xe7, 0xbd, 0x67, 0x5a, 0x90, 0xa1, 0x03, 0xcd, 0xbf, 0xe5, 0xef, 0xa1, 0x9e, 0xef, 0x10, 0x4f,
	0x2a, 0xbe, 0xed, 0x6a, 0x7b, 0x9d, 0xfd, 0xf0, 0x5a, 0xd4, 0xde, 0x3c, 0xfe, 0x49, 0x6f, 0xbf,
	0x03, 0x00, 0x00, 0xff, 0xff, 0x44, 0xb1, 0xd7, 0x88, 0xbd, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// KVServicerClient is the client API for KVServicer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type KVServicerClient interface {
	Get(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Value, error)
	Set(ctx context.Context, in *KVPair, opts ...grpc.CallOption) (*OperationOK, error)
	Delete(ctx context.Context, in *Key, opts ...grpc.CallOption) (*OperationOK, error)
	Scan(ctx context.Context, in *ScanArgs, opts ...grpc.CallOption) (*Values, error)
}

type kVServicerClient struct {
	cc *grpc.ClientConn
}

func NewKVServicerClient(cc *grpc.ClientConn) KVServicerClient {
	return &kVServicerClient{cc}
}

func (c *kVServicerClient) Get(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Value, error) {
	out := new(Value)
	err := c.cc.Invoke(ctx, "/kvstore.methods.KVServicer/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVServicerClient) Set(ctx context.Context, in *KVPair, opts ...grpc.CallOption) (*OperationOK, error) {
	out := new(OperationOK)
	err := c.cc.Invoke(ctx, "/kvstore.methods.KVServicer/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVServicerClient) Delete(ctx context.Context, in *Key, opts ...grpc.CallOption) (*OperationOK, error) {
	out := new(OperationOK)
	err := c.cc.Invoke(ctx, "/kvstore.methods.KVServicer/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVServicerClient) Scan(ctx context.Context, in *ScanArgs, opts ...grpc.CallOption) (*Values, error) {
	out := new(Values)
	err := c.cc.Invoke(ctx, "/kvstore.methods.KVServicer/Scan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KVServicerServer is the server API for KVServicer service.
type KVServicerServer interface {
	Get(context.Context, *Key) (*Value, error)
	Set(context.Context, *KVPair) (*OperationOK, error)
	Delete(context.Context, *Key) (*OperationOK, error)
	Scan(context.Context, *ScanArgs) (*Values, error)
}

func RegisterKVServicerServer(s *grpc.Server, srv KVServicerServer) {
	s.RegisterService(&_KVServicer_serviceDesc, srv)
}

func _KVServicer_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVServicerServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvstore.methods.KVServicer/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServicerServer).Get(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVServicer_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KVPair)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVServicerServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvstore.methods.KVServicer/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServicerServer).Set(ctx, req.(*KVPair))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVServicer_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVServicerServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvstore.methods.KVServicer/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServicerServer).Delete(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVServicer_Scan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScanArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVServicerServer).Scan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvstore.methods.KVServicer/Scan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServicerServer).Scan(ctx, req.(*ScanArgs))
	}
	return interceptor(ctx, in, info, handler)
}

var _KVServicer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kvstore.methods.KVServicer",
	HandlerType: (*KVServicerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _KVServicer_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _KVServicer_Set_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _KVServicer_Delete_Handler,
		},
		{
			MethodName: "Scan",
			Handler:    _KVServicer_Scan_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/kv-server.proto",
}
