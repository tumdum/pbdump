package pbdump

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestFoo(t *testing.T) {
	var msg EmptyMessage
	b, _ := proto.Marshal(&msg)
	buf := bytes.NewBuffer(b)
	m, _ := Dump(buf)
	fmt.Println(m)
}

func TestMessageWithInt(t *testing.T) {
	msg := MessageWithInt{Id: proto.Int32(42)}
	b, err := proto.Marshal(&msg)
	if err != nil {
		t.Fail()
	}
	buf := bytes.NewBuffer(b)
	m, err := Dump(buf)
	if err != nil {
		t.Fatalf("Failed to dump: '%v'", err)
	}
	v, ok := m[1]
	if !ok {
		t.Fatalf("Missing required field '1': '%v'", m)
	}
	if n, ok := v.(StringerVarint); !ok {
		t.Fatalf("Incorrect type under key: '%v'", v)
	} else if uint64(n) != 42 {
		t.Fatalf("Incorrect value under key: '%v'", n)
	}
}
