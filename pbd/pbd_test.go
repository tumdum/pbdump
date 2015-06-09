package main

import (
	"bytes"
	"testing"
)

func TestBytesConversion(t *testing.T) {
	actual, err := ToBytes("08 96 01")
	if err != nil {
		t.Fatalf("Error: '%v'", err)
	}
	expected := []byte{0x8, 0x96, 0x1}
	if bytes.Compare(actual, expected) != 0 {
		t.Fatalf("Expected '%#v' got '%#v'", expected, actual)
	}
}
