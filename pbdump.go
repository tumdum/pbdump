package pbdump

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	// "strings"
)

type decoder func(io.ByteReader) (Value, error)

var (
	decoders []decoder
)

func init() {
	decoders = []decoder{decodeVarint, decodeDouble, decodeLengthDelimited}
}

type Message map[int][]Value

type Value interface {
	// Start returns offset in bytes from beggining of message to start to this
	// value (tag preceding value is not part of the value hence is counted
	// towards this offset.
	Start() uint64
	Payload() []byte
	String() *string
	Varint() *uint64
	Double() *float64
	Message() *Message
}

type variant struct {
	start   uint64
	str     *string
	varint  *uint64
	double  *float64
	message *Message
	payload []byte
}

func (v variant) Start() uint64 {
	return v.start
}

func (v variant) String() *string {
	return v.str
}

func (v variant) Varint() *uint64 {
	return v.varint
}

func (v variant) Double() *float64 {
	return v.double
}

func (v variant) Message() *Message {
	return v.message
}

func (v variant) Payload() []byte {
	return v.payload
}

type key struct {
	Tag, Type int
}

func Dump(r io.ByteReader) (Value, error) {
	if val, err := decodeMessage(r); err != nil {
		return nil, err
	} else {
		return val, nil
	}
}

func decodeMessage(r io.ByteReader) (Value, error) {
	tmp := make(Message)
	m := variant{message: &tmp}
	for {
		k, err := readKey(r)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		v, err := decoders[k.Type](r)
		if err != nil {
			return nil, err
		}
		if s, ok := (*m.Message())[k.Tag]; ok {
			(*m.Message())[k.Tag] = append(s, v)
		} else {
			(*m.Message())[k.Tag] = []Value{v}
		}
	}
	return m, nil
}

func decodeVarint(r io.ByteReader) (Value, error) {
	if n, err := binary.ReadUvarint(r); err != nil {
		return nil, err
	} else {
		return variant{varint: &n}, nil
	}
}

type ByteReaderReader struct {
	r io.ByteReader
}

func (r *ByteReaderReader) Read(b []byte) (int, error) {
	for i := 0; i < len(b); i++ {
		c, err := r.r.ReadByte()
		if err != nil {
			return i, err
		}
		b[i] = c
	}
	return len(b), nil
}

func decodeLengthDelimited(r io.ByteReader) (Value, error) {
	if l, err := binary.ReadUvarint(r); err != nil {
		return nil, err
	} else {
		b := make([]byte, l)
		fullReader := ByteReaderReader{r}
		if n, err := fullReader.Read(b); err != nil {
			return nil, err
		} else if uint64(n) != l {
			return nil, fmt.Errorf("Too little data")
		}
		buf := bytes.NewBuffer(b)
		if msg, err := decodeMessage(buf); err != nil {
			s := string(b)
			return variant{str: &s}, nil
		} else {
			return msg, nil
		}
	}
}

func decodeDouble(r io.ByteReader) (Value, error) {
	fullReader := ByteReaderReader{r}
	var d float64
	if err := binary.Read(&fullReader, binary.LittleEndian, &d); err != nil {
		return nil, err
	} else {
		return variant{double: &d}, nil
	}
}

func readKey(r io.ByteReader) (key, error) {
	n, err := binary.ReadUvarint(r)
	if err != nil {
		return key{}, err
	}
	k := key{int(n >> 3), int(n & 0x7)}
	if !isSupportedWireType(k.Type) {
		return key{}, fmt.Errorf("Unsupported wire type '%d'", k.Type)
	}
	return k, nil
}

func isSupportedWireType(t int) bool {
	return (t >= 0) && (t < len(decoders)) && (decoders[t] != nil)
}
