// Code generated by protoc-gen-go.
// source: example.proto
// DO NOT EDIT!

/*
Package pbdump is a generated protocol buffer package.

It is generated from these files:
	example.proto

It has these top-level messages:
	EmptyMessage
	MessageWithInt
	MessageWithRepeatedInt
	MessageWithString
	MessageWithEmbeddedRepeatedMessageWithString
	MessageWithDouble
*/
package pbdump

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type EmptyMessage struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *EmptyMessage) Reset()         { *m = EmptyMessage{} }
func (m *EmptyMessage) String() string { return proto.CompactTextString(m) }
func (*EmptyMessage) ProtoMessage()    {}

type MessageWithInt struct {
	Id               *int32 `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *MessageWithInt) Reset()         { *m = MessageWithInt{} }
func (m *MessageWithInt) String() string { return proto.CompactTextString(m) }
func (*MessageWithInt) ProtoMessage()    {}

func (m *MessageWithInt) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

type MessageWithRepeatedInt struct {
	Ids              []int32 `protobuf:"varint,1,rep,name=ids" json:"ids,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MessageWithRepeatedInt) Reset()         { *m = MessageWithRepeatedInt{} }
func (m *MessageWithRepeatedInt) String() string { return proto.CompactTextString(m) }
func (*MessageWithRepeatedInt) ProtoMessage()    {}

func (m *MessageWithRepeatedInt) GetIds() []int32 {
	if m != nil {
		return m.Ids
	}
	return nil
}

type MessageWithString struct {
	Name             *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MessageWithString) Reset()         { *m = MessageWithString{} }
func (m *MessageWithString) String() string { return proto.CompactTextString(m) }
func (*MessageWithString) ProtoMessage()    {}

func (m *MessageWithString) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

type MessageWithEmbeddedRepeatedMessageWithString struct {
	Messages         []*MessageWithString `protobuf:"bytes,1,rep,name=messages" json:"messages,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *MessageWithEmbeddedRepeatedMessageWithString) Reset() {
	*m = MessageWithEmbeddedRepeatedMessageWithString{}
}
func (m *MessageWithEmbeddedRepeatedMessageWithString) String() string {
	return proto.CompactTextString(m)
}
func (*MessageWithEmbeddedRepeatedMessageWithString) ProtoMessage() {}

func (m *MessageWithEmbeddedRepeatedMessageWithString) GetMessages() []*MessageWithString {
	if m != nil {
		return m.Messages
	}
	return nil
}

type MessageWithDouble struct {
	D                *float64 `protobuf:"fixed64,1,req,name=d" json:"d,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *MessageWithDouble) Reset()         { *m = MessageWithDouble{} }
func (m *MessageWithDouble) String() string { return proto.CompactTextString(m) }
func (*MessageWithDouble) ProtoMessage()    {}

func (m *MessageWithDouble) GetD() float64 {
	if m != nil && m.D != nil {
		return *m.D
	}
	return 0
}

func init() {
}
