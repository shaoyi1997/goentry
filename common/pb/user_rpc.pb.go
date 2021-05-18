// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: user_rpc.proto

package pb

import (
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

type RpcRequest_Method int32

const (
	RpcRequest_Login    RpcRequest_Method = 1
	RpcRequest_Logout   RpcRequest_Method = 2
	RpcRequest_Update   RpcRequest_Method = 3
	RpcRequest_Upload   RpcRequest_Method = 4
	RpcRequest_Register RpcRequest_Method = 5
)

// Enum value maps for RpcRequest_Method.
var (
	RpcRequest_Method_name = map[int32]string{
		1: "Login",
		2: "Logout",
		3: "Update",
		4: "Upload",
		5: "Register",
	}
	RpcRequest_Method_value = map[string]int32{
		"Login":    1,
		"Logout":   2,
		"Update":   3,
		"Upload":   4,
		"Register": 5,
	}
)

func (x RpcRequest_Method) Enum() *RpcRequest_Method {
	p := new(RpcRequest_Method)
	*p = x
	return p
}

func (x RpcRequest_Method) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RpcRequest_Method) Descriptor() protoreflect.EnumDescriptor {
	return file_user_rpc_proto_enumTypes[0].Descriptor()
}

func (RpcRequest_Method) Type() protoreflect.EnumType {
	return &file_user_rpc_proto_enumTypes[0]
}

func (x RpcRequest_Method) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *RpcRequest_Method) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = RpcRequest_Method(num)
	return nil
}

// Deprecated: Use RpcRequest_Method.Descriptor instead.
func (RpcRequest_Method) EnumDescriptor() ([]byte, []int) {
	return file_user_rpc_proto_rawDescGZIP(), []int{0, 0}
}

type LoginRegisterResponse_ErrorCode int32

const (
	LoginRegisterResponse_InvalidUsername     LoginRegisterResponse_ErrorCode = 1
	LoginRegisterResponse_InvalidPassword     LoginRegisterResponse_ErrorCode = 2
	LoginRegisterResponse_MissingCredentials  LoginRegisterResponse_ErrorCode = 3
	LoginRegisterResponse_InternalServerError LoginRegisterResponse_ErrorCode = 4
)

// Enum value maps for LoginRegisterResponse_ErrorCode.
var (
	LoginRegisterResponse_ErrorCode_name = map[int32]string{
		1: "InvalidUsername",
		2: "InvalidPassword",
		3: "MissingCredentials",
		4: "InternalServerError",
	}
	LoginRegisterResponse_ErrorCode_value = map[string]int32{
		"InvalidUsername":     1,
		"InvalidPassword":     2,
		"MissingCredentials":  3,
		"InternalServerError": 4,
	}
)

func (x LoginRegisterResponse_ErrorCode) Enum() *LoginRegisterResponse_ErrorCode {
	p := new(LoginRegisterResponse_ErrorCode)
	*p = x
	return p
}

func (x LoginRegisterResponse_ErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LoginRegisterResponse_ErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_user_rpc_proto_enumTypes[1].Descriptor()
}

func (LoginRegisterResponse_ErrorCode) Type() protoreflect.EnumType {
	return &file_user_rpc_proto_enumTypes[1]
}

func (x LoginRegisterResponse_ErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *LoginRegisterResponse_ErrorCode) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = LoginRegisterResponse_ErrorCode(num)
	return nil
}

// Deprecated: Use LoginRegisterResponse_ErrorCode.Descriptor instead.
func (LoginRegisterResponse_ErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_user_rpc_proto_rawDescGZIP(), []int{2, 0}
}

type LogoutResponse_ErrorCode int32

const (
	LogoutResponse_MissingUsername     LogoutResponse_ErrorCode = 1
	LogoutResponse_InternalServerError LogoutResponse_ErrorCode = 2
)

// Enum value maps for LogoutResponse_ErrorCode.
var (
	LogoutResponse_ErrorCode_name = map[int32]string{
		1: "MissingUsername",
		2: "InternalServerError",
	}
	LogoutResponse_ErrorCode_value = map[string]int32{
		"MissingUsername":     1,
		"InternalServerError": 2,
	}
)

func (x LogoutResponse_ErrorCode) Enum() *LogoutResponse_ErrorCode {
	p := new(LogoutResponse_ErrorCode)
	*p = x
	return p
}

func (x LogoutResponse_ErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LogoutResponse_ErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_user_rpc_proto_enumTypes[2].Descriptor()
}

func (LogoutResponse_ErrorCode) Type() protoreflect.EnumType {
	return &file_user_rpc_proto_enumTypes[2]
}

func (x LogoutResponse_ErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *LogoutResponse_ErrorCode) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = LogoutResponse_ErrorCode(num)
	return nil
}

// Deprecated: Use LogoutResponse_ErrorCode.Descriptor instead.
func (LogoutResponse_ErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_user_rpc_proto_rawDescGZIP(), []int{4, 0}
}

type RpcRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Method *RpcRequest_Method `protobuf:"varint,1,req,name=method,enum=RpcRequest_Method" json:"method,omitempty"`
	Token  *string            `protobuf:"bytes,2,opt,name=token" json:"token,omitempty"`
}

func (x *RpcRequest) Reset() {
	*x = RpcRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RpcRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RpcRequest) ProtoMessage() {}

func (x *RpcRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RpcRequest.ProtoReflect.Descriptor instead.
func (*RpcRequest) Descriptor() ([]byte, []int) {
	return file_user_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *RpcRequest) GetMethod() RpcRequest_Method {
	if x != nil && x.Method != nil {
		return *x.Method
	}
	return RpcRequest_Login
}

func (x *RpcRequest) GetToken() string {
	if x != nil && x.Token != nil {
		return *x.Token
	}
	return ""
}

type LoginRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username *string `protobuf:"bytes,1,req,name=username" json:"username,omitempty"`
	Password *string `protobuf:"bytes,2,req,name=password" json:"password,omitempty"`
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_user_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *LoginRequest) GetUsername() string {
	if x != nil && x.Username != nil {
		return *x.Username
	}
	return ""
}

func (x *LoginRequest) GetPassword() string {
	if x != nil && x.Password != nil {
		return *x.Password
	}
	return ""
}

type LoginRegisterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User  *User                            `protobuf:"bytes,1,opt,name=user" json:"user,omitempty"`
	Token *string                          `protobuf:"bytes,2,opt,name=token" json:"token,omitempty"`
	Error *LoginRegisterResponse_ErrorCode `protobuf:"varint,3,opt,name=error,enum=LoginRegisterResponse_ErrorCode" json:"error,omitempty"`
}

func (x *LoginRegisterResponse) Reset() {
	*x = LoginRegisterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_rpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginRegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRegisterResponse) ProtoMessage() {}

func (x *LoginRegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_rpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRegisterResponse.ProtoReflect.Descriptor instead.
func (*LoginRegisterResponse) Descriptor() ([]byte, []int) {
	return file_user_rpc_proto_rawDescGZIP(), []int{2}
}

func (x *LoginRegisterResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *LoginRegisterResponse) GetToken() string {
	if x != nil && x.Token != nil {
		return *x.Token
	}
	return ""
}

func (x *LoginRegisterResponse) GetError() LoginRegisterResponse_ErrorCode {
	if x != nil && x.Error != nil {
		return *x.Error
	}
	return LoginRegisterResponse_InvalidUsername
}

type LogoutRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username *string `protobuf:"bytes,1,req,name=username" json:"username,omitempty"`
	Token    *string `protobuf:"bytes,2,req,name=token" json:"token,omitempty"`
}

func (x *LogoutRequest) Reset() {
	*x = LogoutRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_rpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogoutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogoutRequest) ProtoMessage() {}

func (x *LogoutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_rpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogoutRequest.ProtoReflect.Descriptor instead.
func (*LogoutRequest) Descriptor() ([]byte, []int) {
	return file_user_rpc_proto_rawDescGZIP(), []int{3}
}

func (x *LogoutRequest) GetUsername() string {
	if x != nil && x.Username != nil {
		return *x.Username
	}
	return ""
}

func (x *LogoutRequest) GetToken() string {
	if x != nil && x.Token != nil {
		return *x.Token
	}
	return ""
}

type LogoutResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success *bool                     `protobuf:"varint,1,req,name=Success" json:"Success,omitempty"`
	Error   *LogoutResponse_ErrorCode `protobuf:"varint,2,opt,name=error,enum=LogoutResponse_ErrorCode" json:"error,omitempty"`
}

func (x *LogoutResponse) Reset() {
	*x = LogoutResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_rpc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogoutResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogoutResponse) ProtoMessage() {}

func (x *LogoutResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_rpc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogoutResponse.ProtoReflect.Descriptor instead.
func (*LogoutResponse) Descriptor() ([]byte, []int) {
	return file_user_rpc_proto_rawDescGZIP(), []int{4}
}

func (x *LogoutResponse) GetSuccess() bool {
	if x != nil && x.Success != nil {
		return *x.Success
	}
	return false
}

func (x *LogoutResponse) GetError() LogoutResponse_ErrorCode {
	if x != nil && x.Error != nil {
		return *x.Error
	}
	return LogoutResponse_MissingUsername
}

type UpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   *uint64 `protobuf:"varint,1,req,name=user_id,json=userId" json:"user_id,omitempty"`
	Nickname *string `protobuf:"bytes,2,req,name=nickname" json:"nickname,omitempty"`
}

func (x *UpdateRequest) Reset() {
	*x = UpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_rpc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRequest) ProtoMessage() {}

func (x *UpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_rpc_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRequest.ProtoReflect.Descriptor instead.
func (*UpdateRequest) Descriptor() ([]byte, []int) {
	return file_user_rpc_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateRequest) GetUserId() uint64 {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return 0
}

func (x *UpdateRequest) GetNickname() string {
	if x != nil && x.Nickname != nil {
		return *x.Nickname
	}
	return ""
}

type RegisterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username *string `protobuf:"bytes,1,req,name=username" json:"username,omitempty"`
	Password *string `protobuf:"bytes,2,req,name=password" json:"password,omitempty"`
	Nickname *string `protobuf:"bytes,3,opt,name=nickname" json:"nickname,omitempty"`
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_rpc_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_rpc_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_user_rpc_proto_rawDescGZIP(), []int{6}
}

func (x *RegisterRequest) GetUsername() string {
	if x != nil && x.Username != nil {
		return *x.Username
	}
	return ""
}

func (x *RegisterRequest) GetPassword() string {
	if x != nil && x.Password != nil {
		return *x.Password
	}
	return ""
}

func (x *RegisterRequest) GetNickname() string {
	if x != nil && x.Nickname != nil {
		return *x.Nickname
	}
	return ""
}

var File_user_rpc_proto protoreflect.FileDescriptor

var file_user_rpc_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x95, 0x01, 0x0a,
	0x0a, 0x52, 0x70, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x06, 0x6d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x02, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x52, 0x70,
	0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52,
	0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x45, 0x0a,
	0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x09, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e,
	0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x4c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x10, 0x02, 0x12, 0x0a,
	0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06, 0x55, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x10, 0x04, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x10, 0x05, 0x22, 0x46, 0x0a, 0x0c, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x02,
	0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0xe8, 0x01, 0x0a,
	0x15, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65,
	0x72, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x36, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22,
	0x66, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x13, 0x0a, 0x0f,
	0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x10,
	0x01, 0x12, 0x13, 0x0a, 0x0f, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x50, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x10, 0x02, 0x12, 0x16, 0x0a, 0x12, 0x4d, 0x69, 0x73, 0x73, 0x69, 0x6e,
	0x67, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x10, 0x03, 0x12, 0x17,
	0x0a, 0x13, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x04, 0x22, 0x41, 0x0a, 0x0d, 0x4c, 0x6f, 0x67, 0x6f, 0x75,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20,
	0x02, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x96, 0x01, 0x0a, 0x0e, 0x4c,
	0x6f, 0x67, 0x6f, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x02, 0x28, 0x08, 0x52, 0x07,
	0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x2f, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x4c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64,
	0x65, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x39, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x13, 0x0a, 0x0f, 0x4d, 0x69, 0x73, 0x73, 0x69, 0x6e, 0x67,
	0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x49, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x10, 0x02, 0x22, 0x44, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x02, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a,
	0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x02, 0x28, 0x09, 0x52,
	0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x65, 0x0a, 0x0f, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x02, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65,
	0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x2e, 0x67, 0x61, 0x72, 0x65, 0x6e, 0x61, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x73, 0x68, 0x61, 0x6f, 0x79, 0x69, 0x2e, 0x68, 0x6f, 0x6e, 0x67, 0x2f, 0x67,
	0x6f, 0x2d, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x2d, 0x74, 0x61, 0x73, 0x6b, 0x2f, 0x70, 0x62,
}

var (
	file_user_rpc_proto_rawDescOnce sync.Once
	file_user_rpc_proto_rawDescData = file_user_rpc_proto_rawDesc
)

func file_user_rpc_proto_rawDescGZIP() []byte {
	file_user_rpc_proto_rawDescOnce.Do(func() {
		file_user_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_rpc_proto_rawDescData)
	})
	return file_user_rpc_proto_rawDescData
}

var file_user_rpc_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_user_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_user_rpc_proto_goTypes = []interface{}{
	(RpcRequest_Method)(0),               // 0: RpcRequest.Method
	(LoginRegisterResponse_ErrorCode)(0), // 1: LoginRegisterResponse.ErrorCode
	(LogoutResponse_ErrorCode)(0),        // 2: LogoutResponse.ErrorCode
	(*RpcRequest)(nil),                   // 3: RpcRequest
	(*LoginRequest)(nil),                 // 4: LoginRequest
	(*LoginRegisterResponse)(nil),        // 5: LoginRegisterResponse
	(*LogoutRequest)(nil),                // 6: LogoutRequest
	(*LogoutResponse)(nil),               // 7: LogoutResponse
	(*UpdateRequest)(nil),                // 8: UpdateRequest
	(*RegisterRequest)(nil),              // 9: RegisterRequest
	(*User)(nil),                         // 10: User
}
var file_user_rpc_proto_depIdxs = []int32{
	0,  // 0: RpcRequest.method:type_name -> RpcRequest.Method
	10, // 1: LoginRegisterResponse.user:type_name -> User
	1,  // 2: LoginRegisterResponse.error:type_name -> LoginRegisterResponse.ErrorCode
	2,  // 3: LogoutResponse.error:type_name -> LogoutResponse.ErrorCode
	4,  // [4:4] is the sub-list for method output_type
	4,  // [4:4] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_user_rpc_proto_init() }
func file_user_rpc_proto_init() {
	if File_user_rpc_proto != nil {
		return
	}
	file_user_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_user_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RpcRequest); i {
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
		file_user_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginRequest); i {
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
		file_user_rpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginRegisterResponse); i {
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
		file_user_rpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogoutRequest); i {
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
		file_user_rpc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogoutResponse); i {
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
		file_user_rpc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateRequest); i {
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
		file_user_rpc_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterRequest); i {
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
			RawDescriptor: file_user_rpc_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_user_rpc_proto_goTypes,
		DependencyIndexes: file_user_rpc_proto_depIdxs,
		EnumInfos:         file_user_rpc_proto_enumTypes,
		MessageInfos:      file_user_rpc_proto_msgTypes,
	}.Build()
	File_user_rpc_proto = out.File
	file_user_rpc_proto_rawDesc = nil
	file_user_rpc_proto_goTypes = nil
	file_user_rpc_proto_depIdxs = nil
}
