// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.4
// source: inspire/proto/inspire.proto

package inspire_proto

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

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid    string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Prompt  string `protobuf:"bytes,2,opt,name=prompt,proto3" json:"prompt,omitempty"`
	Context string `protobuf:"bytes,3,opt,name=context,proto3" json:"context,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inspire_proto_inspire_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_inspire_proto_inspire_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_inspire_proto_inspire_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *Request) GetPrompt() string {
	if x != nil {
		return x.Prompt
	}
	return ""
}

func (x *Request) GetContext() string {
	if x != nil {
		return x.Context
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid    string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Prompt  string `protobuf:"bytes,2,opt,name=prompt,proto3" json:"prompt,omitempty"`
	Context string `protobuf:"bytes,3,opt,name=context,proto3" json:"context,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inspire_proto_inspire_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_inspire_proto_inspire_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_inspire_proto_inspire_proto_rawDescGZIP(), []int{1}
}

func (x *Response) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *Response) GetPrompt() string {
	if x != nil {
		return x.Prompt
	}
	return ""
}

func (x *Response) GetContext() string {
	if x != nil {
		return x.Context
	}
	return ""
}

var File_inspire_proto_inspire_proto protoreflect.FileDescriptor

var file_inspire_proto_inspire_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x69, 0x6e, 0x73, 0x70, 0x69, 0x72, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x69, 0x6e, 0x73, 0x70, 0x69, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x69,
	0x6e, 0x73, 0x70, 0x69, 0x72, 0x65, 0x22, 0x4f, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x22, 0x50, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x72, 0x6f, 0x6d, 0x70,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x32, 0x3c, 0x0a, 0x0a, 0x41, 0x64, 0x64,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x49, 0x6e, 0x73, 0x70, 0x69,
	0x72, 0x65, 0x12, 0x10, 0x2e, 0x69, 0x6e, 0x73, 0x70, 0x69, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x69, 0x6e, 0x73, 0x70, 0x69, 0x72, 0x65, 0x2e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x6f, 0x6a, 0x69, 0x6e, 0x2d, 0x6b, 0x72, 0x2f, 0x66,
	0x69, 0x62, 0x65, 0x72, 0x2d, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x69, 0x6e, 0x73, 0x70, 0x69, 0x72,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x69, 0x6e, 0x73, 0x70, 0x69, 0x72, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_inspire_proto_inspire_proto_rawDescOnce sync.Once
	file_inspire_proto_inspire_proto_rawDescData = file_inspire_proto_inspire_proto_rawDesc
)

func file_inspire_proto_inspire_proto_rawDescGZIP() []byte {
	file_inspire_proto_inspire_proto_rawDescOnce.Do(func() {
		file_inspire_proto_inspire_proto_rawDescData = protoimpl.X.CompressGZIP(file_inspire_proto_inspire_proto_rawDescData)
	})
	return file_inspire_proto_inspire_proto_rawDescData
}

var file_inspire_proto_inspire_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_inspire_proto_inspire_proto_goTypes = []interface{}{
	(*Request)(nil),  // 0: inspire.Request
	(*Response)(nil), // 1: inspire.Response
}
var file_inspire_proto_inspire_proto_depIdxs = []int32{
	0, // 0: inspire.AddService.Inspire:input_type -> inspire.Request
	1, // 1: inspire.AddService.Inspire:output_type -> inspire.Response
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_inspire_proto_inspire_proto_init() }
func file_inspire_proto_inspire_proto_init() {
	if File_inspire_proto_inspire_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_inspire_proto_inspire_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_inspire_proto_inspire_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
			RawDescriptor: file_inspire_proto_inspire_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_inspire_proto_inspire_proto_goTypes,
		DependencyIndexes: file_inspire_proto_inspire_proto_depIdxs,
		MessageInfos:      file_inspire_proto_inspire_proto_msgTypes,
	}.Build()
	File_inspire_proto_inspire_proto = out.File
	file_inspire_proto_inspire_proto_rawDesc = nil
	file_inspire_proto_inspire_proto_goTypes = nil
	file_inspire_proto_inspire_proto_depIdxs = nil
}
