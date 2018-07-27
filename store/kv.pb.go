// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kv.proto

/*
Package store is a generated protocol buffer package.

It is generated from these files:
	kv.proto

It has these top-level messages:
	ProtoKVItem
*/
package store

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ProtoKVItem struct {
	Key              []byte `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Value            []byte `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *ProtoKVItem) Reset()                    { *m = ProtoKVItem{} }
func (m *ProtoKVItem) String() string            { return proto.CompactTextString(m) }
func (*ProtoKVItem) ProtoMessage()               {}
func (*ProtoKVItem) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ProtoKVItem) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *ProtoKVItem) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*ProtoKVItem)(nil), "store.ProtoKVItem")
}

func init() { proto.RegisterFile("kv.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 84 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xc8, 0x2e, 0xd3, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0x2e, 0xc9, 0x2f, 0x4a, 0x55, 0x32, 0xe5, 0xe2, 0x0e,
	0x00, 0xf1, 0xbd, 0xc3, 0x3c, 0x4b, 0x52, 0x73, 0x85, 0x04, 0xb8, 0x98, 0xb3, 0x53, 0x2b, 0x25,
	0x18, 0x15, 0x18, 0x35, 0x78, 0x82, 0x40, 0x4c, 0x21, 0x11, 0x2e, 0xd6, 0xb2, 0xc4, 0x9c, 0xd2,
	0x54, 0x09, 0x26, 0xb0, 0x18, 0x84, 0x03, 0x08, 0x00, 0x00, 0xff, 0xff, 0xb6, 0xf7, 0xfa, 0x56,
	0x48, 0x00, 0x00, 0x00,
}