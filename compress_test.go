package main

import (
	_ "bytes"
	_ "encoding/binary"
	"fmt"
	"testing"
)

var _ = fmt.Sprintf(" ")

func TestCompressTimePanic(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("compressTime not panic")
		}
	}()

	cblock := NewCompressedBlock()
	cblock.StartTime = 1
	cblock.compressTime(0)
}
