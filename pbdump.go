package pbdump

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

type decoder func(io.ByteReader) (fmt.Stringer, error)

var (
	decoders []decoder
)

func init() {
	decoders = []decoder{decodeVarint, decodeDouble, decodeLengthDelimited}
}

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
		tmp[i] = "'" + v.String() + "'"
	}
	return "{" + strings.Join(tmp, ";") + "}"
}

type StringerMessage map[int]StringerRepeated

func (s StringerMessage) String() string {
	buf := ""
	for k, v := range s {
		buf += fmt.Sprint(k) + " -> " + fmt.Sprint(v) + ", "
	}
	return buf
}

type StringerDouble float64

func (s StringerDouble) String() string {
	return fmt.Sprint(float64(s))
}

func Dump(r io.ByteReader) (StringerMessage, error) {
	if m, err := decodeMessage(r); err != nil {
		return nil, err
	} else {
		return m.(StringerMessage), nil
	}
}

func decodeMessage(r io.ByteReader) (fmt.Stringer, error) {
	m := make(StringerMessage)
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

func decodeLengthDelimited(r io.ByteReader) (fmt.Stringer, error) {
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
			return StringerString(string(b)), nil
		} else {
			return msg, nil
		}
	}
}

func decodeDouble(r io.ByteReader) (fmt.Stringer, error) {
	fullReader := ByteReaderReader{r}
	var d float64
	if err := binary.Read(&fullReader, binary.LittleEndian, &d); err != nil {
		return nil, err
	} else {
		return StringerDouble(d), nil
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
