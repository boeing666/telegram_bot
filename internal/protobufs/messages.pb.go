// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.0
// 	protoc        v5.27.0--rc1
// source: messages.proto

package protobufs

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

type MessageID int32

const (
	MessageID_AddNewChannel MessageID = 0
	MessageID_MyChannels    MessageID = 1
	MessageID_AddNewKeyWord MessageID = 2
	MessageID_RemoveKeyWord MessageID = 3
	MessageID_NextChannels  MessageID = 4
	MessageID_PrevChannels  MessageID = 5
	MessageID_NextKeyWords  MessageID = 6
	MessageID_PrevKeyWords  MessageID = 7
	MessageID_ChannelInfo   MessageID = 8
	MessageID_Back          MessageID = 9
	MessageID_MainPage      MessageID = 10
	MessageID_RemoveChannel MessageID = 11
	MessageID_Spacer        MessageID = 666
)

// Enum value maps for MessageID.
var (
	MessageID_name = map[int32]string{
		0:   "AddNewChannel",
		1:   "MyChannels",
		2:   "AddNewKeyWord",
		3:   "RemoveKeyWord",
		4:   "NextChannels",
		5:   "PrevChannels",
		6:   "NextKeyWords",
		7:   "PrevKeyWords",
		8:   "ChannelInfo",
		9:   "Back",
		10:  "MainPage",
		11:  "RemoveChannel",
		666: "Spacer",
	}
	MessageID_value = map[string]int32{
		"AddNewChannel": 0,
		"MyChannels":    1,
		"AddNewKeyWord": 2,
		"RemoveKeyWord": 3,
		"NextChannels":  4,
		"PrevChannels":  5,
		"NextKeyWords":  6,
		"PrevKeyWords":  7,
		"ChannelInfo":   8,
		"Back":          9,
		"MainPage":      10,
		"RemoveChannel": 11,
		"Spacer":        666,
	}
)

func (x MessageID) Enum() *MessageID {
	p := new(MessageID)
	*p = x
	return p
}

func (x MessageID) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageID) Descriptor() protoreflect.EnumDescriptor {
	return file_messages_proto_enumTypes[0].Descriptor()
}

func (MessageID) Type() protoreflect.EnumType {
	return &file_messages_proto_enumTypes[0]
}

func (x MessageID) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageID.Descriptor instead.
func (MessageID) EnumDescriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{0}
}

type MessageHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time  uint64    `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	Msgid MessageID `protobuf:"varint,2,opt,name=msgid,proto3,enum=protobufs.MessageID" json:"msgid,omitempty"`
	Msg   []byte    `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *MessageHeader) Reset() {
	*x = MessageHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageHeader) ProtoMessage() {}

func (x *MessageHeader) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageHeader.ProtoReflect.Descriptor instead.
func (*MessageHeader) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{0}
}

func (x *MessageHeader) GetTime() uint64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *MessageHeader) GetMsgid() MessageID {
	if x != nil {
		return x.Msgid
	}
	return MessageID_AddNewChannel
}

func (x *MessageHeader) GetMsg() []byte {
	if x != nil {
		return x.Msg
	}
	return nil
}

type ButtonChanneInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId int64 `protobuf:"varint,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
}

func (x *ButtonChanneInfo) Reset() {
	*x = ButtonChanneInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ButtonChanneInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ButtonChanneInfo) ProtoMessage() {}

func (x *ButtonChanneInfo) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ButtonChanneInfo.ProtoReflect.Descriptor instead.
func (*ButtonChanneInfo) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{1}
}

func (x *ButtonChanneInfo) GetChannelId() int64 {
	if x != nil {
		return x.ChannelId
	}
	return 0
}

type ButtonMenuBack struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Newmenu MessageID `protobuf:"varint,1,opt,name=newmenu,proto3,enum=protobufs.MessageID" json:"newmenu,omitempty"`
	Msg     []byte    `protobuf:"bytes,3,opt,name=msg,proto3,oneof" json:"msg,omitempty"`
}

func (x *ButtonMenuBack) Reset() {
	*x = ButtonMenuBack{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ButtonMenuBack) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ButtonMenuBack) ProtoMessage() {}

func (x *ButtonMenuBack) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ButtonMenuBack.ProtoReflect.Descriptor instead.
func (*ButtonMenuBack) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{2}
}

func (x *ButtonMenuBack) GetNewmenu() MessageID {
	if x != nil {
		return x.Newmenu
	}
	return MessageID_AddNewChannel
}

func (x *ButtonMenuBack) GetMsg() []byte {
	if x != nil {
		return x.Msg
	}
	return nil
}

type ButtonRemoveKeyWord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	KeywordId int64 `protobuf:"varint,1,opt,name=keyword_id,json=keywordId,proto3" json:"keyword_id,omitempty"`
	GroupId   int64 `protobuf:"varint,2,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
}

func (x *ButtonRemoveKeyWord) Reset() {
	*x = ButtonRemoveKeyWord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ButtonRemoveKeyWord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ButtonRemoveKeyWord) ProtoMessage() {}

func (x *ButtonRemoveKeyWord) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ButtonRemoveKeyWord.ProtoReflect.Descriptor instead.
func (*ButtonRemoveKeyWord) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{3}
}

func (x *ButtonRemoveKeyWord) GetKeywordId() int64 {
	if x != nil {
		return x.KeywordId
	}
	return 0
}

func (x *ButtonRemoveKeyWord) GetGroupId() int64 {
	if x != nil {
		return x.GroupId
	}
	return 0
}

var File_messages_proto protoreflect.FileDescriptor

var file_messages_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x09, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x22, 0x61, 0x0a, 0x0d, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65,
	0x12, 0x2a, 0x0a, 0x05, 0x6d, 0x73, 0x67, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x49, 0x44, 0x52, 0x05, 0x6d, 0x73, 0x67, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03,
	0x6d, 0x73, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x31,
	0x0a, 0x10, 0x62, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49,
	0x64, 0x22, 0x5f, 0x0a, 0x0e, 0x62, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x4d, 0x65, 0x6e, 0x75, 0x42,
	0x61, 0x63, 0x6b, 0x12, 0x2e, 0x0a, 0x07, 0x6e, 0x65, 0x77, 0x6d, 0x65, 0x6e, 0x75, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x44, 0x52, 0x07, 0x6e, 0x65, 0x77, 0x6d,
	0x65, 0x6e, 0x75, 0x12, 0x15, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x48, 0x00, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x88, 0x01, 0x01, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x6d,
	0x73, 0x67, 0x22, 0x4f, 0x0a, 0x13, 0x62, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x52, 0x65, 0x6d, 0x6f,
	0x76, 0x65, 0x4b, 0x65, 0x79, 0x57, 0x6f, 0x72, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x6b, 0x65, 0x79,
	0x77, 0x6f, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x6b,
	0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x49, 0x64, 0x2a, 0xe5, 0x01, 0x0a, 0x09, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49,
	0x44, 0x12, 0x11, 0x0a, 0x0d, 0x41, 0x64, 0x64, 0x4e, 0x65, 0x77, 0x43, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x4d, 0x79, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x73, 0x10, 0x01, 0x12, 0x11, 0x0a, 0x0d, 0x41, 0x64, 0x64, 0x4e, 0x65, 0x77, 0x4b, 0x65,
	0x79, 0x57, 0x6f, 0x72, 0x64, 0x10, 0x02, 0x12, 0x11, 0x0a, 0x0d, 0x52, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x4b, 0x65, 0x79, 0x57, 0x6f, 0x72, 0x64, 0x10, 0x03, 0x12, 0x10, 0x0a, 0x0c, 0x4e, 0x65,
	0x78, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x10, 0x04, 0x12, 0x10, 0x0a, 0x0c,
	0x50, 0x72, 0x65, 0x76, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x10, 0x05, 0x12, 0x10,
	0x0a, 0x0c, 0x4e, 0x65, 0x78, 0x74, 0x4b, 0x65, 0x79, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x10, 0x06,
	0x12, 0x10, 0x0a, 0x0c, 0x50, 0x72, 0x65, 0x76, 0x4b, 0x65, 0x79, 0x57, 0x6f, 0x72, 0x64, 0x73,
	0x10, 0x07, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x6e, 0x66,
	0x6f, 0x10, 0x08, 0x12, 0x08, 0x0a, 0x04, 0x42, 0x61, 0x63, 0x6b, 0x10, 0x09, 0x12, 0x0c, 0x0a,
	0x08, 0x4d, 0x61, 0x69, 0x6e, 0x50, 0x61, 0x67, 0x65, 0x10, 0x0a, 0x12, 0x11, 0x0a, 0x0d, 0x52,
	0x65, 0x6d, 0x6f, 0x76, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x10, 0x0b, 0x12, 0x0b,
	0x0a, 0x06, 0x53, 0x70, 0x61, 0x63, 0x65, 0x72, 0x10, 0x9a, 0x05, 0x42, 0x0d, 0x5a, 0x0b, 0x2e,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_messages_proto_rawDescOnce sync.Once
	file_messages_proto_rawDescData = file_messages_proto_rawDesc
)

func file_messages_proto_rawDescGZIP() []byte {
	file_messages_proto_rawDescOnce.Do(func() {
		file_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_messages_proto_rawDescData)
	})
	return file_messages_proto_rawDescData
}

var file_messages_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_messages_proto_goTypes = []interface{}{
	(MessageID)(0),              // 0: protobufs.MessageID
	(*MessageHeader)(nil),       // 1: protobufs.MessageHeader
	(*ButtonChanneInfo)(nil),    // 2: protobufs.buttonChanneInfo
	(*ButtonMenuBack)(nil),      // 3: protobufs.buttonMenuBack
	(*ButtonRemoveKeyWord)(nil), // 4: protobufs.buttonRemoveKeyWord
}
var file_messages_proto_depIdxs = []int32{
	0, // 0: protobufs.MessageHeader.msgid:type_name -> protobufs.MessageID
	0, // 1: protobufs.buttonMenuBack.newmenu:type_name -> protobufs.MessageID
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_messages_proto_init() }
func file_messages_proto_init() {
	if File_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageHeader); i {
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
		file_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ButtonChanneInfo); i {
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
		file_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ButtonMenuBack); i {
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
		file_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ButtonRemoveKeyWord); i {
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
	file_messages_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_messages_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_messages_proto_goTypes,
		DependencyIndexes: file_messages_proto_depIdxs,
		EnumInfos:         file_messages_proto_enumTypes,
		MessageInfos:      file_messages_proto_msgTypes,
	}.Build()
	File_messages_proto = out.File
	file_messages_proto_rawDesc = nil
	file_messages_proto_goTypes = nil
	file_messages_proto_depIdxs = nil
}
