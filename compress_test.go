package main

import (
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

func TestCompressTime_Delta_64(t *testing.T) {
	cblock := NewCompressedBlock()

	res := cblock.delta_64(1)
	if res != 257 {
		t.Error("res!= 257", res)
	}

	res = cblock.delta_64(64)
	if res != 320 {
		t.Error("res!= 320", res)
	}

	res = cblock.delta_64(63)
	if res != 319 {
		t.Error("res!= 319", res)
	}
}

func TestCompressTime_Delta_256(t *testing.T) {
	cblock := NewCompressedBlock()

	res := cblock.delta_256(256)
	if res != 3328 {
		t.Error("res!= 3328", res)
	}

	res = cblock.delta_256(255)
	if res != 3327 {
		t.Error("res!= 3327", res)
	}

	res = cblock.delta_256(65)
	if res != 3137 {
		t.Error("res!= 3137", res)
	}
}

func TestCompressTime_Delta_2048(t *testing.T) {
	cblock := NewCompressedBlock()

	res := cblock.delta_2048(2048)
	if res != 59392 {
		t.Error("res!= 59392", res)
	}

	res = cblock.delta_2048(257)
	if res != 57601 {
		t.Error("res!= 57601", res)
	}

	res = cblock.delta_2048(4095)
	if res != 61439 {
		t.Error("res!= 61439", res)
	}
}

func TestCompressTime_Delta_Big(t *testing.T) {
	cblock := NewCompressedBlock()

	res := cblock.delta_big(2049)
	if res != 64424511489 {
		t.Error("res!= 64424511489", res)
	}

	res = cblock.delta_big(65535)
	if res != 64424574975 {
		t.Error("res!= 64424574975", res)
	}

	res = cblock.delta_big(4095)
	if res != 64424513535 {
		t.Error("res!= 64424513535", res)
	}

	res = cblock.delta_big(4294967295)
	if res != 68719476735 {
		t.Error("res!= 68719476735", res)
	}
}
