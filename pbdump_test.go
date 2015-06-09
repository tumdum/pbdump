package pbdump

import (
	"bytes"
	"io"
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestMessageWithInt(t *testing.T) {
	msg := MessageWithInt{Id: proto.Int32(42)}
	buf := MustMarshal(&msg)
	if m, err := Dump(buf); err != nil {
		t.Fatalf("Failed to dump: '%v'", err)
	} else if v, ok := m[1]; !ok {
		t.Fatalf("Missing required field '1': '%v'", m)
	} else if !HasVarints(v, 42) {
		t.Fatalf("Incorrect value for field, expected '%v', got '%v'", 42, v)
	}
}

func TestMessageWithRepeatedInt(t *testing.T) {
	msg := MessageWithRepeatedInt{Ids: []int32{1, 2, 333, 456789}}
	buf := MustMarshal(&msg)
	if m, err := Dump(buf); err != nil {
		t.Fatalf("Failed to dump: '%v'", err)
	} else if v, ok := m[1]; !ok {
		t.Fatalf("Missing filed repeated field '1': '%v'", m)
	} else if !HasVarints(v, 1, 2, 333, 456789) {
		t.Fatalf("Missing values for tag '1': '%v'", v)
	}
}

func HasVarints(actual StringerRepeated, expected ...uint64) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i, _ := range actual {
		if v, ok := actual[i].(StringerVarint); !ok {
			return false
		} else if uint64(v) != expected[i] {
			return false
		}
	}
	return true
}

func MustMarshal(msg proto.Message) io.ByteReader {
	b, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(b)
}
