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
	} else if v, ok := m.attributes[1]; !ok {
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
	} else if v, ok := m.attributes[1]; !ok {
		t.Fatalf("Missing filed repeated field '1': '%v'", m)
	} else if !HasVarints(v, 1, 2, 333, 456789) {
		t.Fatalf("Missing values for tag '1': '%v'", v)
	}
}

func TestMessageWithString(t *testing.T) {
	name := "name"
	msg := MessageWithString{Name: &name}
	buf := MustMarshal(&msg)
	if m, err := Dump(buf); err != nil {
		t.Fatalf("Failed to dump: '%v'", err)
	} else if v, ok := m.attributes[1]; !ok {
		t.Fatalf("Missing required field '1': '%v'", m)
	} else if !HasStrings(v, name) {
		t.Fatalf("Incorrect value, expected '%s', got '%s'", name, v)
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
	v, ok := m.attributes[1]
	if !ok {
		t.Fatalf("Missing filed repeated field '1': '%v'", m)
	} else if len(v) != 2 {
		t.Fatalf("Expected to have to repeated messages, got: '%d' (%v)", len(v), v)
	}
	if v0, ok := v[0].(StringerMessage); !ok {
		t.Fatalf("Expected message, found '%#v'", v[0])
	} else if !HasStrings(v0.attributes[1], name1) {
		t.Fatal("First message expected to have string '%v', got '%v'", name1, v0)
	}
	if v1, ok := v[1].(StringerMessage); !ok {
		t.Fatalf("Expected message, found '%#v'", v[1])
	} else if !HasStrings(v1.attributes[1], name2) {
		t.Fatalf("Second message expected to have string '%v', got '%v'", name2, v1)
	}
}

func TestMessageWithDouble(t *testing.T) {
	d := float64(3.14159)
	msg := MessageWithDouble{D: &d}
	buf := MustMarshal(&msg)
	if m, err := Dump(buf); err != nil {
		t.Fatalf("Failed to dump: '%v'", err)
	} else if v, ok := m.attributes[1]; !ok {
		t.Fatalf("Missing required field '1': '%v'", m)
	} else if !HasDouble(v, d) {
		t.Fatalf("Incorrect value, expected '%v', got '%v'", d, v)
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

func HasStrings(actual StringerRepeated, expected ...string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i, _ := range actual {
		if v, ok := actual[i].(StringerString); !ok {
			return false
		} else if string(v) != expected[i] {
			return false
		}
	}
	return true
}

func HasDouble(actual StringerRepeated, expected ...float64) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i, _ := range actual {
		if v, ok := actual[i].(StringerDouble); !ok {
			return false
		} else if float64(v) != expected[i] {
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
