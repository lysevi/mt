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

func TestCompressIncBytePanic(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("incBit not panic")
		}
	}()

	cblock := NewCompressedBlock()
	cblock.byteNum = MAX_BLOCK_SIZE
	cblock.incByte()
}

func TestCompress2String(t *testing.T) {
	cblock := NewCompressedBlock()
	cblock.data[0] = 1
	str := cblock.String()
	if len(str) == 0 {
		t.Error("empty ", str)
	}
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

func TestCompressTime_Write_Delta_64(t *testing.T) {
	cblock := NewCompressedBlock()
	cblock.write_64(257)
	//	fmt.Println(cblock.String())
	if cblock.bitNum != 6 || cblock.byteNum != 1 {
		t.Error("cblock.bitNum != 1 || cblock.byteNum != 1", cblock.bitNum, cblock.byteNum)
	}
	//1 0000 0001
	if cblock.data[0] != 128 || cblock.data[1] != 128 {
		t.Error("cblock.data[0] != 1 || cblock.data[2] != 1: ", cblock.data[0], cblock.data[1])
	}
	//1 0100 0000
	cblock.write_64(320)
	if cblock.data[1] != 208 || cblock.data[2] != 0 {
		t.Error("cblock.data[1] != 208 || cblock.data[2] != 1: ", cblock.data[1], cblock.data[2])
	}
	//1 0011 1111
	cblock.write_64(319)
	if cblock.data[2] != 39 || cblock.data[3] != 224 {
		t.Error("cblock.data[2] != 39 || cblock.data[3] != 224: ", cblock.data[2], cblock.data[3])
	}
}

func TestCompressTime_Write_Delta_256(t *testing.T) {
	cblock := NewCompressedBlock()
	cblock.write_256(3328) //1101 0000 0000
	if cblock.data[0] != 208 || cblock.data[1] != 0 {
		t.Error("first error: ", cblock.data[0], cblock.data[1])
	}

	cblock.write_256(3327) //1100 11111111
	if cblock.data[1] != 12 || cblock.data[2] != 255 {
		t.Error("second error: ", cblock.data[1], cblock.data[2])
	}

	cblock.write_256(3137)
	if cblock.data[3] != 196 || cblock.data[4] != 16 {
		t.Error("second error: ", cblock.data[3], cblock.data[4])
	}
}

func TestCompressTime_Write_Delta_2048(t *testing.T) {
	cblock := NewCompressedBlock()
	cblock.write_2048(59392) //1110100000000000
	if cblock.data[0] != 232 || cblock.data[1] != 0 {
		t.Error("first error: ", cblock.data[0], cblock.data[1])
	}
	cblock.write_2048(57601) //1110000100000001
	if cblock.data[2] != 225 || cblock.data[3] != 1 {
		t.Error("second error: ", cblock.data[2], cblock.data[3])
	}

	cblock.write_2048(61439) //1110111111111111
	if cblock.data[4] != 239 || cblock.data[5] != 255 {
		t.Error("second error: ", cblock.data[4], cblock.data[5])
	}
}

func TestCompressTime_Write_Delta_big(t *testing.T) {
	cblock := NewCompressedBlock()
	cblock.write_big(64424511489) //111100000000000000000000100000000001
	if cblock.data[0] != 240 || cblock.data[1] != 0 || cblock.data[2] != 0 || cblock.data[3] != 128 || cblock.data[4] != 16 {
		t.Error(cblock.String())
	}
	cblock.write_big(64424574975) //111100000000000000001111111111111111
	if cblock.data[4] != 31 || cblock.data[5] != 0 || cblock.data[6] != 0 || cblock.data[7] != 255 || cblock.data[8] != 255 {
		t.Error(cblock.String())
	}
	cblock.write_big(64424513535) //111100000000000000000000111111111111
	if cblock.data[9] != 240 || cblock.data[10] != 0 || cblock.data[11] != 0 || cblock.data[12] != 255 || cblock.data[13] != 240 {
		t.Error(cblock.String())
	}
	cblock.write_big(68719476735) //111111111111111111111111111111111111
	if cblock.data[14] != 255 || cblock.data[15] != 255 || cblock.data[16] != 255 || cblock.data[16] != 255 || cblock.data[17] != 255 {
		t.Error(cblock.String())
	}
}
