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
	} else if m.Message() == nil {
		t.Fatalf("Expected to decode message, got '%v'", m)
	} else if v, ok := m.Message().Get(1); !ok {
		t.Fatalf("Missing required field '1': '%v'", m)
	} else if !HasVarints(v, 42) {
		t.Fatalf("Incorrect value for field, expected '%v', got '%v'", 42, v)
	} else if v[0].Start() != 1 {
		t.Fatalf("Incorrect start offset, expected 1, got '%v'", v[0].Start())
	}
}

func TestMessageWithRepeatedInt(t *testing.T) {
	msg := MessageWithRepeatedInt{Ids: []int32{1, 2, 333, 456789}}
	buf := MustMarshal(&msg)
	expectedOffsets := []uint64{1, 3, 5, 8}
	expectedLengths := []uint64{1, 1, 2, 3}
	if m, err := Dump(buf); err != nil {
		t.Fatalf("Failed to dump: '%v'", err)
	} else if v, ok := m.Message().Get(1); !ok {
		t.Fatalf("Missing filed repeated field '1': '%v'", m)
	} else if !HasVarints(v, 1, 2, 333, 456789) {
		t.Fatalf("Missing values for tag '1': '%v'", v)
	} else {
		HasCorrectStarts(t, v, expectedOffsets)
		HasCorrectLength(t, v, expectedLengths)
	}
}

func TestMessageWithString(t *testing.T) {
	name := "name"
	msg := MessageWithString{Name: &name}
	buf := MustMarshal(&msg)
	if m, err := Dump(buf); err != nil {
		t.Fatalf("Failed to dump: '%v'", err)
	} else if v, ok := m.Message().Get(1); !ok {
		t.Fatalf("Missing required field '1': '%v'", m)
	} else if !HasStrings(v, name) {
		t.Fatalf("Incorrect value, expected '%s', got '%s'", name, v)
	} else if v[0].Start() != 1 {
		t.Fatalf("Expected offset 1, got %d", v[0].Start())
	} else if v[0].Length() != uint64(len(name)+1) {
		t.Fatalf("Incorrect length")
	}
}

func TestMessageWithEmbeddedRepeatedMessageWithString(t *testing.T) {
	name1 := "name1"
	msg1 := MessageWithString{Name: &name1}
	name2 := "name2"
	msg2 := MessageWithString{Name: &name2}
	msg := MessageWithEmbeddedRepeatedMessageWithString{
		Messages: []*MessageWithString{
			&msg1, &msg2,
		},
	}
	buf := MustMarshal(&msg)
	m, err := Dump(buf)
	if err != nil {
		t.Fatalf("Failed to dump: '%v'", err)
	}
	v, _ := m.Message().Get(1)
	if len(v) != 2 {
		t.Fatalf("Expected to have to repeated messages, got: '%d' (%v)", len(v), v)
	}
	if m0 := v[0].Message(); m0 == nil {
		t.Fatalf("Expected message, found '%#v'", v[0])
	} else if !HasStrings(mustGet(m0, 1), name1) {
		t.Fatal("First message expected to have string '%v', got '%v'", name1, m0)
	} else if v[0].Start() != 1 {
		t.Fatalf("First message expected to start at 1, got %d", v[0].Start())
	} else if mustGet(m0, 1)[0].Start() != 3 {
		t.Fatalf("String from first message expected to start at 3, got %d", mustGet(m0, 1)[0].Start())
	}
	if m1 := v[1].Message(); m1 == nil {
		t.Fatalf("Expected message, found '%#v'", v[1])
	} else if !HasStrings(mustGet(m1, 1), name2) {
		t.Fatalf("Second message expected to have string '%v', got '%v'", name2, m1)
	} else if v[1].Start() != 10 {
		t.Fatalf("Second message expected to start at 10, got %d", v[1].Start())
	} else if mustGet(m1, 1)[0].Start() != 12 {
		t.Fatalf("String from second message expected to start at 12, got %d", mustGet(m1, 1)[0].Start())
	}
}

func TestMessageWithDouble(t *testing.T) {
	d := float64(3.14159)
	msg := MessageWithDouble{D: &d}
	b, err := proto.Marshal(&msg)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(b)
	if m, err := Dump(buf); err != nil {
		t.Fatalf("Failed to dump: '%v'", err)
	} else if v, ok := m.Message().Get(1); !ok {
		t.Fatalf("Missing required field '1': '%v'", m)
	} else if !HasDouble(v, d) {
		t.Fatalf("Incorrect value, expected '%v', got '%v'", d, v)
	} else if v[0].Start() != 1 {
		t.Fatalf("Expected that double starts at 1, got %d", v[0].Start())
	} else if v[0].Length() != 8 {
		t.Fatalf("Incorrect length")
	}
}

func HasDouble(actual []Value, expected ...float64) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i, _ := range actual {
		if d := actual[i].Double(); d == nil {
			return false
		} else if *d != expected[i] {
			return false
		}
	}
	return true
}

func HasVarints(actual []Value, expected ...uint64) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i, _ := range actual {
		if vp := actual[i].Varint(); vp == nil {
			return false
		} else if *vp != expected[i] {
			return false
		}
	}
	return true
}

func HasStrings(actual []Value, expected ...string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i, _ := range actual {
		if s := actual[i].String(); s == nil {
			return false
		} else if *s != expected[i] {
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

func mustGet(m *Message, id int) []Value {
	if v, ok := m.Get(id); ok {
		return v
	} else {
		panic("Failed to get field from message")
	}
}

func HasCorrectStarts(t *testing.T, v []Value, expected []uint64) {
	for i, start := range expected {
		if v[i].Start() != start {
			t.Fatalf("Field %d: expected start %d, got %d", i, start, v[i].Start())
		}
	}
}

func HasCorrectLength(t *testing.T, v []Value, expected []uint64) {
	for i, length := range expected {
		if v[i].Length() != length {
			t.Fatalf("Field %d: expected length %d, got %d", i, length, v[i].Length())
		}
	}
}
