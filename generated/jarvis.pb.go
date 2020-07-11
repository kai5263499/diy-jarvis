// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.7.1
// source: jarvis.proto

package generated

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Type int32

const (
	Type_RegisterAudioSourceRequestType  Type = 0
	Type_OutputRequestType               Type = 1
	Type_TextRequestType                 Type = 2
	Type_RegisterAudioSourceResponseType Type = 3
	Type_OutputResponseType              Type = 4
	Type_TextResponseType                Type = 5
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "RegisterAudioSourceRequestType",
		1: "OutputRequestType",
		2: "TextRequestType",
		3: "RegisterAudioSourceResponseType",
		4: "OutputResponseType",
		5: "TextResponseType",
	}
	Type_value = map[string]int32{
		"RegisterAudioSourceRequestType":  0,
		"OutputRequestType":               1,
		"TextRequestType":                 2,
		"RegisterAudioSourceResponseType": 3,
		"OutputResponseType":              4,
		"TextResponseType":                5,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_jarvis_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_jarvis_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_jarvis_proto_rawDescGZIP(), []int{0}
}

type Code int32

const (
	Code_ERROR    Code = 0
	Code_ACCEPTED Code = 1
)

// Enum value maps for Code.
var (
	Code_name = map[int32]string{
		0: "ERROR",
		1: "ACCEPTED",
	}
	Code_value = map[string]int32{
		"ERROR":    0,
		"ACCEPTED": 1,
	}
)

func (x Code) Enum() *Code {
	p := new(Code)
	*p = x
	return p
}

func (x Code) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Code) Descriptor() protoreflect.EnumDescriptor {
	return file_jarvis_proto_enumTypes[1].Descriptor()
}

func (Code) Type() protoreflect.EnumType {
	return &file_jarvis_proto_enumTypes[1]
}

func (x Code) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Code.Descriptor instead.
func (Code) EnumDescriptor() ([]byte, []int) {
	return file_jarvis_proto_rawDescGZIP(), []int{1}
}

type Base struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      Type   `protobuf:"varint,1,opt,name=Type,proto3,enum=Type" json:"Type,omitempty"`
	Code      Code   `protobuf:"varint,2,opt,name=Code,proto3,enum=Code" json:"Code,omitempty"`
	Id        string `protobuf:"bytes,3,opt,name=Id,proto3" json:"Id,omitempty"`
	SourceId  string `protobuf:"bytes,4,opt,name=SourceId,proto3" json:"SourceId,omitempty"`
	Timestamp uint64 `protobuf:"varint,5,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	Text      string `protobuf:"bytes,6,opt,name=Text,proto3" json:"Text,omitempty"`
	SinkId    string `protobuf:"bytes,7,opt,name=SinkId,proto3" json:"SinkId,omitempty"`
}

func (x *Base) Reset() {
	*x = Base{}
	if protoimpl.UnsafeEnabled {
		mi := &file_jarvis_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Base) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Base) ProtoMessage() {}

func (x *Base) ProtoReflect() protoreflect.Message {
	mi := &file_jarvis_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Base.ProtoReflect.Descriptor instead.
func (*Base) Descriptor() ([]byte, []int) {
	return file_jarvis_proto_rawDescGZIP(), []int{0}
}

func (x *Base) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_RegisterAudioSourceRequestType
}

func (x *Base) GetCode() Code {
	if x != nil {
		return x.Code
	}
	return Code_ERROR
}

func (x *Base) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Base) GetSourceId() string {
	if x != nil {
		return x.SourceId
	}
	return ""
}

func (x *Base) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Base) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *Base) GetSinkId() string {
	if x != nil {
		return x.SinkId
	}
	return ""
}

var File_jarvis_proto protoreflect.FileDescriptor

var file_jarvis_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6a, 0x61, 0x72, 0x76, 0x69, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb2,
	0x01, 0x0a, 0x04, 0x42, 0x61, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x05, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x19, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x05, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x64, 0x12, 0x1a, 0x0a,
	0x08, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x78, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x65, 0x78, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x53,
	0x69, 0x6e, 0x6b, 0x49, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x69, 0x6e,
	0x6b, 0x49, 0x64, 0x2a, 0xa9, 0x01, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x22, 0x0a, 0x1e,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41, 0x75, 0x64, 0x69, 0x6f, 0x53, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x10, 0x00,
	0x12, 0x15, 0x0a, 0x11, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x54, 0x65, 0x78, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x10, 0x02, 0x12, 0x23, 0x0a, 0x1f,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41, 0x75, 0x64, 0x69, 0x6f, 0x53, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x10,
	0x03, 0x12, 0x16, 0x0a, 0x12, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x10, 0x04, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x65, 0x78,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x10, 0x05, 0x2a,
	0x1f, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52,
	0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x41, 0x43, 0x43, 0x45, 0x50, 0x54, 0x45, 0x44, 0x10, 0x01,
	0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x3b, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_jarvis_proto_rawDescOnce sync.Once
	file_jarvis_proto_rawDescData = file_jarvis_proto_rawDesc
)

func file_jarvis_proto_rawDescGZIP() []byte {
	file_jarvis_proto_rawDescOnce.Do(func() {
		file_jarvis_proto_rawDescData = protoimpl.X.CompressGZIP(file_jarvis_proto_rawDescData)
	})
	return file_jarvis_proto_rawDescData
}

var file_jarvis_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_jarvis_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_jarvis_proto_goTypes = []interface{}{
	(Type)(0),    // 0: Type
	(Code)(0),    // 1: Code
	(*Base)(nil), // 2: Base
}
var file_jarvis_proto_depIdxs = []int32{
	0, // 0: Base.Type:type_name -> Type
	1, // 1: Base.Code:type_name -> Code
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_jarvis_proto_init() }
func file_jarvis_proto_init() {
	if File_jarvis_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_jarvis_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Base); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_jarvis_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_jarvis_proto_goTypes,
		DependencyIndexes: file_jarvis_proto_depIdxs,
		EnumInfos:         file_jarvis_proto_enumTypes,
		MessageInfos:      file_jarvis_proto_msgTypes,
	}.Build()
	File_jarvis_proto = out.File
	file_jarvis_proto_rawDesc = nil
	file_jarvis_proto_goTypes = nil
	file_jarvis_proto_depIdxs = nil
}
