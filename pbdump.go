package pbdump

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
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

type StringerRepeated []fmt.Stringer

func (s StringerRepeated) String() string {
	tmp := make([]string, len(s))
	for i, v := range s {
		tmp[i] = v.String()
	}
	return "{" + strings.Join(tmp, ",") + "}"
}

func Dump(r io.ByteReader) (map[int]StringerRepeated, error) {
	m := make(map[int]StringerRepeated)
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
		if s, ok := m[k.Tag]; ok {
			m[k.Tag] = append(s, v)
		} else {
			m[k.Tag] = StringerRepeated{v}
		}
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
