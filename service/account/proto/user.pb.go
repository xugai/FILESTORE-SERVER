// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.4
// source: user.proto

package proto

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

type ReqSignup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *ReqSignup) Reset() {
	*x = ReqSignup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqSignup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqSignup) ProtoMessage() {}

func (x *ReqSignup) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqSignup.ProtoReflect.Descriptor instead.
func (*ReqSignup) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{0}
}

func (x *ReqSignup) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *ReqSignup) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type RespSignup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *RespSignup) Reset() {
	*x = RespSignup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RespSignup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RespSignup) ProtoMessage() {}

func (x *RespSignup) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RespSignup.ProtoReflect.Descriptor instead.
func (*RespSignup) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{1}
}

func (x *RespSignup) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *RespSignup) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ReqSignin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *ReqSignin) Reset() {
	*x = ReqSignin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqSignin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqSignin) ProtoMessage() {}

func (x *ReqSignin) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqSignin.ProtoReflect.Descriptor instead.
func (*ReqSignin) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{2}
}

func (x *ReqSignin) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *ReqSignin) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type RespSignin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Token   string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *RespSignin) Reset() {
	*x = RespSignin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RespSignin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RespSignin) ProtoMessage() {}

func (x *RespSignin) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RespSignin.ProtoReflect.Descriptor instead.
func (*RespSignin) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{3}
}

func (x *RespSignin) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *RespSignin) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *RespSignin) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ReqUserInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *ReqUserInfo) Reset() {
	*x = ReqUserInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqUserInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqUserInfo) ProtoMessage() {}

func (x *ReqUserInfo) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqUserInfo.ProtoReflect.Descriptor instead.
func (*ReqUserInfo) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{4}
}

func (x *ReqUserInfo) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

type RespUserInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code         int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message      string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Username     string `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
	Email        string `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	Phone        string `protobuf:"bytes,5,opt,name=phone,proto3" json:"phone,omitempty"`
	Signup       string `protobuf:"bytes,6,opt,name=signup,proto3" json:"signup,omitempty"`
	LastActiveAt string `protobuf:"bytes,7,opt,name=lastActiveAt,proto3" json:"lastActiveAt,omitempty"`
	Status       int32  `protobuf:"varint,8,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *RespUserInfo) Reset() {
	*x = RespUserInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RespUserInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RespUserInfo) ProtoMessage() {}

func (x *RespUserInfo) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RespUserInfo.ProtoReflect.Descriptor instead.
func (*RespUserInfo) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{5}
}

func (x *RespUserInfo) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *RespUserInfo) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *RespUserInfo) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *RespUserInfo) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *RespUserInfo) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *RespUserInfo) GetSignup() string {
	if x != nil {
		return x.Signup
	}
	return ""
}

func (x *RespUserInfo) GetLastActiveAt() string {
	if x != nil {
		return x.LastActiveAt
	}
	return ""
}

func (x *RespUserInfo) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

var File_user_proto protoreflect.FileDescriptor

var file_user_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x43, 0x0a, 0x09,
	0x52, 0x65, 0x71, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x22, 0x3a, 0x0a, 0x0a, 0x52, 0x65, 0x73, 0x70, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x12,
	0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x43, 0x0a,
	0x09, 0x52, 0x65, 0x71, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x22, 0x50, 0x0a, 0x0a, 0x52, 0x65, 0x73, 0x70, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x29, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x22,
	0xd8, 0x01, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x70, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x69, 0x67, 0x6e, 0x75, 0x70,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x12, 0x22,
	0x0a, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x41, 0x74, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65,
	0x41, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0x82, 0x01, 0x0a, 0x0b, 0x55,
	0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x23, 0x0a, 0x06, 0x53, 0x69,
	0x67, 0x6e, 0x75, 0x70, 0x12, 0x0a, 0x2e, 0x52, 0x65, 0x71, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70,
	0x1a, 0x0b, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x22, 0x00, 0x12,
	0x23, 0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x12, 0x0a, 0x2e, 0x52, 0x65, 0x71, 0x53,
	0x69, 0x67, 0x6e, 0x69, 0x6e, 0x1a, 0x0b, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x53, 0x69, 0x67, 0x6e,
	0x69, 0x6e, 0x22, 0x00, 0x12, 0x29, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x0c, 0x2e, 0x52, 0x65, 0x71, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x0d,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x00, 0x42,
	0x07, 0x5a, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_proto_rawDescOnce sync.Once
	file_user_proto_rawDescData = file_user_proto_rawDesc
)

func file_user_proto_rawDescGZIP() []byte {
	file_user_proto_rawDescOnce.Do(func() {
		file_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_proto_rawDescData)
	})
	return file_user_proto_rawDescData
}

var file_user_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_user_proto_goTypes = []interface{}{
	(*ReqSignup)(nil),    // 0: ReqSignup
	(*RespSignup)(nil),   // 1: RespSignup
	(*ReqSignin)(nil),    // 2: ReqSignin
	(*RespSignin)(nil),   // 3: RespSignin
	(*ReqUserInfo)(nil),  // 4: ReqUserInfo
	(*RespUserInfo)(nil), // 5: RespUserInfo
}
var file_user_proto_depIdxs = []int32{
	0, // 0: UserService.Signup:input_type -> ReqSignup
	2, // 1: UserService.Signin:input_type -> ReqSignin
	4, // 2: UserService.UserInfo:input_type -> ReqUserInfo
	1, // 3: UserService.Signup:output_type -> RespSignup
	3, // 4: UserService.Signin:output_type -> RespSignin
	5, // 5: UserService.UserInfo:output_type -> RespUserInfo
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_user_proto_init() }
func file_user_proto_init() {
	if File_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqSignup); i {
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
		file_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RespSignup); i {
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
		file_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqSignin); i {
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
		file_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RespSignin); i {
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
		file_user_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqUserInfo); i {
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
		file_user_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RespUserInfo); i {
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
			RawDescriptor: file_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_proto_goTypes,
		DependencyIndexes: file_user_proto_depIdxs,
		MessageInfos:      file_user_proto_msgTypes,
	}.Build()
	File_user_proto = out.File
	file_user_proto_rawDesc = nil
	file_user_proto_goTypes = nil
	file_user_proto_depIdxs = nil
}
