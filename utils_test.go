package main

import (
	_ "fmt"
	"testing"
)

func TestMeasInTimeInterval(t *testing.T) {
	if !inTimeInterval(0, 10, 1) {
		t.Error("!inTimeInterval(0,10,1)")
	}

	if inTimeInterval(0, 10, 11) {
		t.Error("inTimeInterval(0,10,11)")
	}

	if !inTimeInterval(10, 0, 1) {
		t.Error("inTimeInterval(10,0,1)")
	}
}

func TestMeasIdFltr(t *testing.T) {
	var ids []Id

	if !idFltr(ids, 1) {
		t.Error("!idFltr(ids,1)")
	}
	ids = append(ids, 1)

	if idFltr(ids, 2) {
		t.Error("idFltr(ids,2)")
	}

	if !idFltr(ids, 1) {
		t.Error("!idFltr(ids,1)")
	}
}

func TestMeasFlagFltr(t *testing.T) {
	if !flagFltr(0, 1) {
		t.Error("!flagFltr(0,1)")
	}

	if flagFltr(2, 1) {
		t.Error("flagFltr(2,1)")
	}

	if !flagFltr(0, 1) {
		t.Error("flagFltr(0,1)")
	}

	if !flagFltr(1, 1) {
		t.Error("flagFltr(1,1)")
	}
}

func TestBitOperations(t *testing.T) {
	var value byte = 0
	value = setBit(value, 0, 1)

	if value != 1 {
		t.Error("value!=1", value)
	}

	if !checkBit(value, 0) {
		t.Error("!checkBit(value,0)")
	}

	value = setBit(value, 0, 0)
	if value != 0 {
		t.Error("value!=0", value)
	}

	if checkBit(value, 0) || value != 0 {
		t.Error("checkBit(value,0) || value != 0 ")
	}

	value = 0
	for i := 1; i < 8; i += 2 {
		value = setBit(value, uint8(i), uint8(1))
		if !checkBit(value, uint8(i)) {
			t.Errorf("!checkBit(value,i)")
		}
	}
	if value != 170 {
		t.Error("value!=170", value)
	}

	value = 0
	for i := 0; i < 8; i++ {
		value = setBit(value, uint8(i), uint8(1))
		if !checkBit(value, uint8(i)) {
			t.Errorf("!checkBit(value,i)")
		}
	}

	if value != 255 {
		t.Error("value!=255", value)
	}

	value = setBit(255, 7, 0)
	value = setBit(value, 0, 0)

	if value != 126 {
		t.Error("value!=126", value)
	}

	value = 64
	if checkBit(value, 7) || !checkBit(value, 6) {
		t.Error("64 test error")
	}
}

func TestBitOperations16(t *testing.T) {
	value := uint16(0)
	for i := 0; i < 16; i++ {
		value = setBit16(value, uint8(i), uint8(1))
		if !checkBit16(value, uint8(i)) {
			t.Errorf("!checkBit(value,i)")
		}
	}
	if value != 65535 {
		t.Error("value!=65535", value)
	}
}

func TestBitOperations32(t *testing.T) {
	value := uint32(0)
	for i := 0; i < 32; i++ {
		value = setBit32(value, uint8(i), uint8(1))
		if !checkBit32(value, uint8(i)) {
			t.Errorf("!checkBit(value,i)")
		}
	}
	if value != 4294967295 {
		t.Error("value!=4294967295", value)
	}
}

func TestBitOperations64(t *testing.T) {
	value := uint64(0)
	for i := 0; i < 64; i++ {
		value = setBit64(value, uint8(i), uint8(1))
		if !checkBit64(value, uint8(i)) {
			t.Errorf("!checkBit(value,i)")
		}
	}
}
