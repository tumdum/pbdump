package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tumdum/pbdump"
)

func main() {
	all, err := ioutil.ReadAll(os.Stdin)
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "Failed to read stdin: '%v'", err)
		os.Exit(1)
	}
	b, err := ToBytes(string(all))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode stdin: '%v'", err)
		os.Exit(1)
	}
	msg, err := pbdump.Dump(bytes.NewBuffer(b))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode message: '%v'", err)
		os.Exit(1)
	}
	fmt.Println(msg.String())
}

func ToBytes(s string) ([]byte, error) {
	trimmed := strings.Trim(s, "\t\n\r ")
	withoutSpace := strings.Replace(trimmed, " ", "", -1)
	if b, err := hex.DecodeString(withoutSpace); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}
