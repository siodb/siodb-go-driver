// Code generated by protoc-gen-go. DO NOT EDIT.
// source: CommonTypes.proto

package siodb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

//* Single message from server or IO manager.
type StatusMessage struct {
	//* Message status code.
	StatusCode int32 `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	//* Message text.
	Text                 string   `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusMessage) Reset()         { *m = StatusMessage{} }
func (m *StatusMessage) String() string { return proto.CompactTextString(m) }
func (*StatusMessage) ProtoMessage()    {}
func (*StatusMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_c2115d584b0d8b29, []int{0}
}

func (m *StatusMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusMessage.Unmarshal(m, b)
}
func (m *StatusMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusMessage.Marshal(b, m, deterministic)
}
func (m *StatusMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusMessage.Merge(m, src)
}
func (m *StatusMessage) XXX_Size() int {
	return xxx_messageInfo_StatusMessage.Size(m)
}
func (m *StatusMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusMessage.DiscardUnknown(m)
}

var xxx_messageInfo_StatusMessage proto.InternalMessageInfo

func (m *StatusMessage) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *StatusMessage) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

//* Structured column data type attribute description
type AttributeDescription struct {
	//* Column name.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	//* Column type.
	Type ColumnDataType `protobuf:"varint,2,opt,name=type,proto3,enum=siodb.ColumnDataType" json:"type,omitempty"`
	//* Column can have null values.
	IsNull bool `protobuf:"varint,3,opt,name=is_null,json=isNull,proto3" json:"is_null,omitempty"`
	//* Attributes of a structured data type
	Attribute            []*AttributeDescription `protobuf:"bytes,4,rep,name=attribute,proto3" json:"attribute,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *AttributeDescription) Reset()         { *m = AttributeDescription{} }
func (m *AttributeDescription) String() string { return proto.CompactTextString(m) }
func (*AttributeDescription) ProtoMessage()    {}
func (*AttributeDescription) Descriptor() ([]byte, []int) {
	return fileDescriptor_c2115d584b0d8b29, []int{1}
}

func (m *AttributeDescription) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AttributeDescription.Unmarshal(m, b)
}
func (m *AttributeDescription) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AttributeDescription.Marshal(b, m, deterministic)
}
func (m *AttributeDescription) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AttributeDescription.Merge(m, src)
}
func (m *AttributeDescription) XXX_Size() int {
	return xxx_messageInfo_AttributeDescription.Size(m)
}
func (m *AttributeDescription) XXX_DiscardUnknown() {
	xxx_messageInfo_AttributeDescription.DiscardUnknown(m)
}

var xxx_messageInfo_AttributeDescription proto.InternalMessageInfo

func (m *AttributeDescription) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AttributeDescription) GetType() ColumnDataType {
	if m != nil {
		return m.Type
	}
	return ColumnDataType_COLUMN_DATA_TYPE_BOOL
}

func (m *AttributeDescription) GetIsNull() bool {
	if m != nil {
		return m.IsNull
	}
	return false
}

func (m *AttributeDescription) GetAttribute() []*AttributeDescription {
	if m != nil {
		return m.Attribute
	}
	return nil
}

//* Describes column returned by server
type ColumnDescription struct {
	//* Column name.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	//* Column type.
	Type ColumnDataType `protobuf:"varint,2,opt,name=type,proto3,enum=siodb.ColumnDataType" json:"type,omitempty"`
	//* Column can have null values.
	IsNull bool `protobuf:"varint,3,opt,name=is_null,json=isNull,proto3" json:"is_null,omitempty"`
	//* Attributes of a structured data type
	Attribute            []*AttributeDescription `protobuf:"bytes,4,rep,name=attribute,proto3" json:"attribute,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *ColumnDescription) Reset()         { *m = ColumnDescription{} }
func (m *ColumnDescription) String() string { return proto.CompactTextString(m) }
func (*ColumnDescription) ProtoMessage()    {}
func (*ColumnDescription) Descriptor() ([]byte, []int) {
	return fileDescriptor_c2115d584b0d8b29, []int{2}
}

func (m *ColumnDescription) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ColumnDescription.Unmarshal(m, b)
}
func (m *ColumnDescription) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ColumnDescription.Marshal(b, m, deterministic)
}
func (m *ColumnDescription) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ColumnDescription.Merge(m, src)
}
func (m *ColumnDescription) XXX_Size() int {
	return xxx_messageInfo_ColumnDescription.Size(m)
}
func (m *ColumnDescription) XXX_DiscardUnknown() {
	xxx_messageInfo_ColumnDescription.DiscardUnknown(m)
}

var xxx_messageInfo_ColumnDescription proto.InternalMessageInfo

func (m *ColumnDescription) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ColumnDescription) GetType() ColumnDataType {
	if m != nil {
		return m.Type
	}
	return ColumnDataType_COLUMN_DATA_TYPE_BOOL
}

func (m *ColumnDescription) GetIsNull() bool {
	if m != nil {
		return m.IsNull
	}
	return false
}

func (m *ColumnDescription) GetAttribute() []*AttributeDescription {
	if m != nil {
		return m.Attribute
	}
	return nil
}

func init() {
	proto.RegisterType((*StatusMessage)(nil), "siodb.StatusMessage")
	proto.RegisterType((*AttributeDescription)(nil), "siodb.AttributeDescription")
	proto.RegisterType((*ColumnDescription)(nil), "siodb.ColumnDescription")
}

func init() { proto.RegisterFile("CommonTypes.proto", fileDescriptor_c2115d584b0d8b29) }

var fileDescriptor_c2115d584b0d8b29 = []byte{
	// 250 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x91, 0x31, 0x4b, 0xc4, 0x40,
	0x10, 0x85, 0xd9, 0xbb, 0xdc, 0x69, 0xf6, 0x50, 0x70, 0x39, 0x31, 0x68, 0x61, 0x48, 0x15, 0x9b,
	0x14, 0x67, 0x65, 0xa9, 0x49, 0x61, 0xa3, 0xc5, 0x6a, 0x7f, 0x6c, 0x2e, 0x83, 0x2c, 0x24, 0x3b,
	0x21, 0x33, 0x0b, 0xde, 0x1f, 0x12, 0xfc, 0x97, 0x92, 0x8d, 0x22, 0x82, 0x3f, 0xe0, 0xba, 0xe1,
	0xbd, 0xc7, 0x7b, 0x1f, 0xbb, 0xf2, 0xac, 0xc4, 0xae, 0x43, 0xf7, 0xba, 0xef, 0x81, 0x8a, 0x7e,
	0x40, 0x46, 0xb5, 0x20, 0x8b, 0x4d, 0x7d, 0xb9, 0x2e, 0xb1, 0xf5, 0x9d, 0xab, 0x0c, 0x9b, 0xd1,
	0x9d, 0xcc, 0xac, 0x92, 0x27, 0x2f, 0x6c, 0xd8, 0xd3, 0x13, 0x10, 0x99, 0x37, 0x50, 0xd7, 0x72,
	0x45, 0x41, 0xd8, 0xee, 0xb0, 0x81, 0x44, 0xa4, 0x22, 0x5f, 0x68, 0x39, 0x49, 0x25, 0x36, 0xa0,
	0x94, 0x8c, 0x18, 0xde, 0x39, 0x99, 0xa5, 0x22, 0x8f, 0x75, 0xb8, 0xb3, 0x4f, 0x21, 0xd7, 0xf7,
	0xcc, 0x83, 0xad, 0x3d, 0x43, 0x05, 0xb4, 0x1b, 0x6c, 0xcf, 0x16, 0xdd, 0x18, 0x76, 0xa6, 0x9b,
	0x6a, 0x62, 0x1d, 0x6e, 0x75, 0x23, 0x23, 0xde, 0xf7, 0x10, 0x0a, 0x4e, 0x37, 0xe7, 0x45, 0xc0,
	0x2b, 0xfe, 0xd2, 0xe9, 0x10, 0x51, 0x17, 0xf2, 0xc8, 0xd2, 0xd6, 0xf9, 0xb6, 0x4d, 0xe6, 0xa9,
	0xc8, 0x8f, 0xf5, 0xd2, 0xd2, 0xb3, 0x6f, 0x5b, 0x75, 0x27, 0x63, 0xf3, 0xb3, 0x97, 0x44, 0xe9,
	0x3c, 0x5f, 0x6d, 0xae, 0xbe, 0x8b, 0xfe, 0xe3, 0xd0, 0xbf, 0xe9, 0xec, 0x43, 0x8c, 0x8f, 0x14,
	0xc6, 0x0e, 0x1a, 0xf4, 0x61, 0xf6, 0x28, 0xea, 0x65, 0xf8, 0xa5, 0xdb, 0xaf, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xfb, 0xdf, 0xe2, 0xeb, 0xd7, 0x01, 0x00, 0x00,
}
