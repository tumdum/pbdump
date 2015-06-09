package pbdump

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"testing"
)

func TestReadVarintOfOneByte(t *testing.T) {
	for i := byte(0); i != 128; i++ {
		b := bytes.NewBuffer([]byte{i})
		v, err := readVarint(b)
		if err != nil {
			t.Fatalf("Failed to decode '%d': %s", i, err)
		}
		if v != uint64(i) {
			t.Fatalf("Expected '%d' got '%d'", i, v)
		}
	}
}

func TestFoo(t *testing.T) {
	var msg EmptyMessage
	proto.Marshal(&msg)
}
