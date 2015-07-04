package pbdump

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type decoder func(*countingByteReader) (Value, error)

var (
	decoders []decoder
)

func init() {
	decoders = []decoder{decodeVarint, decodeDouble, decodeLengthDelimited}
}

type Message struct {
	m map[int][]Value
}

func (m Message) Get(id int) ([]Value, bool) {
	v, ok := m.m[id]
	return v, ok
}

func (m *Message) add(id int, v ...Value) {
	old, ok := m.m[id]
	if ok {
		m.m[id] = append(old, v...)
	} else {
		m.m[id] = v
	}
}

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

type countingByteReader struct {
	r     io.ByteReader
	total uint64
}

func (r *countingByteReader) ReadByte() (byte, error) {
	b, err := r.r.ReadByte()
	if err == nil {
		r.total += 1
	}
	return b, err
}

func Dump(r io.ByteReader) (Value, error) {
	reader := countingByteReader{r, 0}
	if val, err := decodeMessage(&reader); err != nil {
		return nil, err
	} else {
		return val, nil
	}
}

func decodeMessage(r *countingByteReader) (Value, error) {
	tmp := Message{make(map[int][]Value)}
	m := variant{message: &tmp, start: r.total}
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
		if _, ok := m.Message().Get(k.Tag); ok {
			m.Message().add(k.Tag, v)
		} else {
			m.Message().add(k.Tag, v)
		}
	}
	return m, nil
}

func decodeVarint(r *countingByteReader) (Value, error) {
	start := r.total
	if n, err := binary.ReadUvarint(r); err != nil {
		return nil, err
	} else {
		return variant{varint: &n, start: start}, nil
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

func decodeLengthDelimited(r *countingByteReader) (Value, error) {
	start := r.total
	if l, err := binary.ReadUvarint(r); err != nil {
		return nil, err
	} else {
		sub := r.total
		b := make([]byte, l)
		fullReader := ByteReaderReader{r}
		if n, err := fullReader.Read(b); err != nil {
			return nil, err
		} else if uint64(n) != l {
			return nil, fmt.Errorf("Too little data")
		}
		buf := bytes.NewBuffer(b)
		reader := countingByteReader{buf, sub}
		if msg, err := decodeMessage(&reader); err != nil {
			s := string(b)
			return variant{str: &s, start: start}, nil
		} else {
			m := msg.(variant)
			m.start = start
			return m, nil
		}
	}
}

func decodeDouble(r *countingByteReader) (Value, error) {
	start := r.total
	fullReader := ByteReaderReader{r}
	var d float64
	if err := binary.Read(&fullReader, binary.LittleEndian, &d); err != nil {
		return nil, err
	} else {
		return variant{double: &d, start: start}, nil
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
