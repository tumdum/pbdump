package pbdump

import (
	"fmt"
	"io"
)

func Dump(r io.Reader) (map[int]fmt.Stringer, error) {
	return nil, nil
}

func readVarint(r io.Reader) (uint64, error) {
	b := []byte{0}
	n, err := r.Read(b)
	if err != nil {
		return 0, err
	} else if n != 1 {
		return 0, fmt.Errorf("Expected to read 1 byte, read '%d' bytes", n)
	}
	return uint64(b[0]), nil
}
