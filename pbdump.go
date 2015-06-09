package pbdump

import (
	"encoding/binary"
	"fmt"
	"io"
)

type decoder func(io.ByteReader) (fmt.Stringer, error)

var (
	decoders = []decoder{decodeVarint}
)

type key struct {
	Tag, Type int
}

type StringerString string

func (s StringerString) String() string {
	return string(s)
}

type StringerVarint uint64

func (s StringerVarint) String() string {
	return fmt.Sprint(uint64(s))
}

func Dump(r io.ByteReader) (map[int]fmt.Stringer, error) {
	m := make(map[int]fmt.Stringer)
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
		m[k.Tag] = v
	}
	return m, nil
}

func decodeVarint(r io.ByteReader) (fmt.Stringer, error) {
	if n, err := binary.ReadUvarint(r); err != nil {
		return nil, err
	} else {
		return StringerVarint(n), nil
	}
}

func readKey(r io.ByteReader) (key, error) {
	n, err := binary.ReadUvarint(r)
	if err != nil {
		return key{}, err
	}
	return key{int(n >> 3), int(n & 0x7)}, nil
}
