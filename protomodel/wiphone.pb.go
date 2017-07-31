// Code generated by protoc-gen-go. DO NOT EDIT.
// source: wiphone.proto

/*
Package protomodel is a generated protocol buffer package.

It is generated from these files:
	wiphone.proto

It has these top-level messages:
	Credentials
	ConsumptionResponse
	AnonymousConsumptionRequest
*/
package protomodel

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

type Credentials struct {
	Operator string `protobuf:"bytes,1,opt,name=Operator" json:"Operator,omitempty"`
	Username string `protobuf:"bytes,2,opt,name=Username" json:"Username,omitempty"`
	Password string `protobuf:"bytes,3,opt,name=Password" json:"Password,omitempty"`
}

func (m *Credentials) Reset()                    { *m = Credentials{} }
func (m *Credentials) String() string            { return proto.CompactTextString(m) }
func (*Credentials) ProtoMessage()               {}
func (*Credentials) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Credentials) GetOperator() string {
	if m != nil {
		return m.Operator
	}
	return ""
}

func (m *Credentials) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *Credentials) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type ConsumptionResponse struct {
	InternetTotal    int64  `protobuf:"varint,1,opt,name=InternetTotal" json:"InternetTotal,omitempty"`
	InternetConsumed int64  `protobuf:"varint,2,opt,name=InternetConsumed" json:"InternetConsumed,omitempty"`
	CallTotal        int32  `protobuf:"varint,3,opt,name=CallTotal" json:"CallTotal,omitempty"`
	CallConsumed     int32  `protobuf:"varint,4,opt,name=CallConsumed" json:"CallConsumed,omitempty"`
	PeriodStart      int32  `protobuf:"varint,5,opt,name=periodStart" json:"periodStart,omitempty"`
	PeriodEnd        int32  `protobuf:"varint,6,opt,name=periodEnd" json:"periodEnd,omitempty"`
	UpdatedAt        int32  `protobuf:"varint,7,opt,name=updatedAt" json:"updatedAt,omitempty"`
	PhoneNumber      string `protobuf:"bytes,8,opt,name=phoneNumber" json:"phoneNumber,omitempty"`
}

func (m *ConsumptionResponse) Reset()                    { *m = ConsumptionResponse{} }
func (m *ConsumptionResponse) String() string            { return proto.CompactTextString(m) }
func (*ConsumptionResponse) ProtoMessage()               {}
func (*ConsumptionResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ConsumptionResponse) GetInternetTotal() int64 {
	if m != nil {
		return m.InternetTotal
	}
	return 0
}

func (m *ConsumptionResponse) GetInternetConsumed() int64 {
	if m != nil {
		return m.InternetConsumed
	}
	return 0
}

func (m *ConsumptionResponse) GetCallTotal() int32 {
	if m != nil {
		return m.CallTotal
	}
	return 0
}

func (m *ConsumptionResponse) GetCallConsumed() int32 {
	if m != nil {
		return m.CallConsumed
	}
	return 0
}

func (m *ConsumptionResponse) GetPeriodStart() int32 {
	if m != nil {
		return m.PeriodStart
	}
	return 0
}

func (m *ConsumptionResponse) GetPeriodEnd() int32 {
	if m != nil {
		return m.PeriodEnd
	}
	return 0
}

func (m *ConsumptionResponse) GetUpdatedAt() int32 {
	if m != nil {
		return m.UpdatedAt
	}
	return 0
}

func (m *ConsumptionResponse) GetPhoneNumber() string {
	if m != nil {
		return m.PhoneNumber
	}
	return ""
}

type AnonymousConsumptionRequest struct {
	DeviceId    string       `protobuf:"bytes,1,opt,name=DeviceId" json:"DeviceId,omitempty"`
	Credentials *Credentials `protobuf:"bytes,2,opt,name=Credentials" json:"Credentials,omitempty"`
	PhoneNumber string       `protobuf:"bytes,3,opt,name=phoneNumber" json:"phoneNumber,omitempty"`
}

func (m *AnonymousConsumptionRequest) Reset()                    { *m = AnonymousConsumptionRequest{} }
func (m *AnonymousConsumptionRequest) String() string            { return proto.CompactTextString(m) }
func (*AnonymousConsumptionRequest) ProtoMessage()               {}
func (*AnonymousConsumptionRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *AnonymousConsumptionRequest) GetDeviceId() string {
	if m != nil {
		return m.DeviceId
	}
	return ""
}

func (m *AnonymousConsumptionRequest) GetCredentials() *Credentials {
	if m != nil {
		return m.Credentials
	}
	return nil
}

func (m *AnonymousConsumptionRequest) GetPhoneNumber() string {
	if m != nil {
		return m.PhoneNumber
	}
	return ""
}

func init() {
	proto.RegisterType((*Credentials)(nil), "protomodel.Credentials")
	proto.RegisterType((*ConsumptionResponse)(nil), "protomodel.ConsumptionResponse")
	proto.RegisterType((*AnonymousConsumptionRequest)(nil), "protomodel.AnonymousConsumptionRequest")
}

func init() { proto.RegisterFile("wiphone.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 317 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0x86, 0x69, 0xd7, 0xd6, 0x36, 0xb5, 0x20, 0xf1, 0x60, 0x50, 0x0f, 0x65, 0xf1, 0x20, 0x1e,
	0x7a, 0xd0, 0x93, 0xc7, 0x52, 0x3d, 0xf4, 0xa2, 0x12, 0xf5, 0x01, 0x52, 0x33, 0xe0, 0xc2, 0x6e,
	0x26, 0x26, 0xb3, 0x16, 0x5f, 0xc4, 0x77, 0xf0, 0x2d, 0x25, 0x09, 0xbb, 0xed, 0xda, 0x53, 0xf8,
	0xbf, 0x3f, 0xfc, 0x33, 0xcc, 0xcf, 0xa6, 0x9b, 0xc2, 0x7e, 0xa0, 0x81, 0xb9, 0x75, 0x48, 0xc8,
	0x59, 0x7c, 0x2a, 0xd4, 0x50, 0xe6, 0x8a, 0x4d, 0x96, 0x0e, 0x34, 0x18, 0x2a, 0x54, 0xe9, 0xf9,
	0x19, 0x1b, 0x3d, 0x59, 0x70, 0x8a, 0xd0, 0x89, 0xde, 0xac, 0x77, 0x35, 0x96, 0xad, 0x0e, 0xde,
	0x9b, 0x07, 0x67, 0x54, 0x05, 0xa2, 0x9f, 0xbc, 0x46, 0x07, 0xef, 0x59, 0x79, 0xbf, 0x41, 0xa7,
	0x45, 0x96, 0xbc, 0x46, 0xe7, 0xbf, 0x7d, 0x76, 0xb2, 0x44, 0xe3, 0xeb, 0xca, 0x52, 0x81, 0x46,
	0x82, 0xb7, 0x68, 0x3c, 0xf0, 0x4b, 0x36, 0x5d, 0x19, 0x02, 0x67, 0x80, 0x5e, 0x91, 0x54, 0x19,
	0x07, 0x66, 0xb2, 0x0b, 0xf9, 0x35, 0x3b, 0x6e, 0x40, 0x0a, 0x01, 0x1d, 0xa7, 0x67, 0x72, 0x8f,
	0xf3, 0x0b, 0x36, 0x5e, 0xaa, 0xb2, 0x4c, 0x69, 0x61, 0x8d, 0x81, 0xdc, 0x02, 0x9e, 0xb3, 0xa3,
	0x20, 0xda, 0x94, 0x83, 0xf8, 0xa1, 0xc3, 0xf8, 0x8c, 0x4d, 0x2c, 0xb8, 0x02, 0xf5, 0x0b, 0x29,
	0x47, 0x62, 0x10, 0xbf, 0xec, 0xa2, 0x30, 0x23, 0xc9, 0x07, 0xa3, 0xc5, 0x30, 0xcd, 0x68, 0x41,
	0x70, 0x6b, 0xab, 0x15, 0x81, 0x5e, 0x90, 0x38, 0x4c, 0x6e, 0x0b, 0x62, 0x7a, 0xe8, 0xe1, 0xb1,
	0xae, 0xd6, 0xe0, 0xc4, 0x28, 0x1e, 0x6a, 0x17, 0xe5, 0x3f, 0x3d, 0x76, 0xbe, 0x30, 0x68, 0xbe,
	0x2b, 0xac, 0x7d, 0xe7, 0x68, 0x9f, 0x35, 0x78, 0x0a, 0x77, 0xbe, 0x87, 0xaf, 0xe2, 0x1d, 0x56,
	0xba, 0xe9, 0xa7, 0xd1, 0xfc, 0xae, 0x53, 0x65, 0x3c, 0xd2, 0xe4, 0xe6, 0x74, 0xbe, 0x2d, 0x7b,
	0xbe, 0x63, 0xcb, 0x4e, 0xed, 0xff, 0x16, 0xcb, 0xf6, 0x16, 0x5b, 0x0f, 0x63, 0xcc, 0xed, 0x5f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x13, 0x3f, 0x20, 0x35, 0x4b, 0x02, 0x00, 0x00,
}
