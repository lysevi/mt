package mt

import (
	"bytes"
	"encoding/binary"
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

func TestCompressTimePrevTimes(t *testing.T) {
	var t1 = Time(63)
	var t2 = Time(64)
	cblock := NewCompressedBlock()
	cblock.StartTime = 1

	cblock.compressTime(t1)
	if cblock.prev_time != t1 || cblock.prev_delta != (t1-cblock.StartTime) {
		t.Error("cblock.prev_time != t1 || cblock.prev_delta != (t1-cblock.StartTime)")
	}

	cblock.compressTime(t2)
	if cblock.prev_time != t2 {
		t.Error("cblock.prev_time != t2")
	}
}

func TestCompressTimeAddFirst(t *testing.T) {
	var t1 = Time(63)
	cblock := NewCompressedBlock()
	cblock.StartTime = 1

	cblock.compressTime(t1)

	b := cblock.data[0:len(cblock.data)]
	buf := bytes.NewBuffer(b)
	var readed_t Time
	binary.Read(buf, binary.LittleEndian, &readed_t)

	if readed_t+cblock.StartTime != t1 {
		t.Error("readed_t != t1", readed_t, t1)
	}

	if cblock.byteNum != 9 {
		t.Error("cblock.byteNum!=9")
	}
}

func TestCompressTimeAddSecond(t *testing.T) {
	{ //D is zero
		var t1 = Time(1)
		var t2 = Time(1)
		cblock := NewCompressedBlock()
		cblock.StartTime = 1

		cblock.compressTime(t1)
		cblock.compressTime(t2)

		if cblock.bitNum != 6 {
			t.Error("1 cblock.bitNum!=6", cblock.bitNum)
		}
		if checkBit(cblock.data[0], MAX_BIT) {
			t.Error("1 checkBit(cblock.data[0],MAX_BIT)")
		}
	}

	{ //D is between [-63, 64]
		var t1 = Time(1)
		var t2 = Time(10)
		cblock := NewCompressedBlock()
		cblock.StartTime = 1

		cblock.compressTime(t1)
		cblock.compressTime(t2)

		if cblock.bitNum != 7 {
			t.Error("2 cblock.bitNum!=0", cblock.bitNum)
		}

		if cblock.data[9] != 164 { // 10 1001 00
			t.Error("2 checkBit(cblock.data[0],164)", cblock.data[9])
		}
	}

	{ //D is between [-255, 255]
		var t1 = Time(1)
		var t2 = Time(230)
		cblock := NewCompressedBlock()
		cblock.StartTime = 1

		cblock.compressTime(t1)
		cblock.compressTime(t2)

		if cblock.byteNum != 8+2 {
			t.Error("2 cblock.byteNum!=8+2 ", cblock.byteNum, 8+2)
		}

		if !checkBit(cblock.data[9], MAX_BIT) {
			t.Error("2 checkBit(cblock.data[0],MAX_BIT)", cblock.data[9])
		}

		if !checkBit(cblock.data[9], MAX_BIT-1) {
			t.Error("2 checkBit(cblock.data[0],MAX_BIT-1)", cblock.data[9])
		}

		if checkBit(cblock.data[9], MAX_BIT-2) {
			t.Error("2 checkBit(cblock.data[0],MAX_BIT-2)", cblock.data[9])
		}
		//110 11100  1000
		//|----212-|   128
		if cblock.data[9] != 212 {
			t.Error("cblock.data[9]!=212", cblock.data[9])
		}

		if cblock.data[10] != 128 {
			t.Error("cblock.data[9]!=128", cblock.data[10])
		}
	}

	{ //D is between  [-2047, 2048]
		var t1 = Time(1)
		var t2 = Time(2045)
		cblock := NewCompressedBlock()
		cblock.StartTime = 1

		cblock.compressTime(t1)
		cblock.compressTime(t2)

		if cblock.byteNum != 8+3 {
			t.Error("2 cblock.byteNum!=8+3 ", cblock.byteNum, 8+3)
		}

		if !checkBit(cblock.data[9], MAX_BIT) {
			t.Error("2 checkBit(cblock.data[0],MAX_BIT)", cblock.data[9])
		}

		if !checkBit(cblock.data[9], MAX_BIT-1) {
			t.Error("2 checkBit(cblock.data[0],MAX_BIT-1)", cblock.data[9])
		}

		if !checkBit(cblock.data[9], MAX_BIT-2) {
			t.Error("2 checkBit(cblock.data[0],MAX_BIT-2)", cblock.data[9])
		}

		if checkBit(cblock.data[9], MAX_BIT-3) {
			t.Error("2 checkBit(cblock.data[0],MAX_BIT-3)", cblock.data[9])
		}
		//1110 0011  1111 1000
		//|----231-|   252
		if cblock.data[9] != 227 {
			t.Error("cblock.data[9]!=239", cblock.data[9])
		}

		if cblock.data[10] != 248 {
			t.Error("cblock.data[9]!=248", cblock.data[10])
		}
	}

	{ //D is > 2048
		var t1 = Time(1)
		var t2 = Time(900011)
		cblock := NewCompressedBlock()
		cblock.StartTime = 1

		cblock.compressTime(t1)
		cblock.compressTime(t2)

		if cblock.byteNum != 8+4 {
			t.Error("2 cblock.byteNum!=8+4 ", cblock.byteNum, 8+4)
		}

		if !checkBit(cblock.data[9], MAX_BIT) {
			t.Error("2 checkBit(cblock.data[0],MAX_BIT)", cblock.data[9])
		}

		if !checkBit(cblock.data[9], MAX_BIT-1) {
			t.Error("2 checkBit(cblock.data[0],MAX_BIT-1)", cblock.data[9])
		}

		if !checkBit(cblock.data[9], MAX_BIT-2) {
			t.Error("2 checkBit(cblock.data[0],MAX_BIT-2)", cblock.data[9])
		}

		if !checkBit(cblock.data[9], MAX_BIT-3) {
			t.Error("2 checkBit(cblock.data[0],MAX_BIT-3)", cblock.data[9])
		}
		//1111 0101  0111 0111 1011 0000 0000
		//|----245-|   252
		if cblock.data[9] != 245 {
			t.Error("cblock.data[9]!=245", cblock.data[9])
		}

		if cblock.data[10] != 119 {
			t.Error("cblock.data[9]!=119", cblock.data[10])
		}

		if cblock.data[11] != 176 {
			t.Error("cblock.data[9]!=248", cblock.data[10])
		}
	}
}

func TestCompressTimeRead(t *testing.T) {
	var t1 = Time(256)
	{
		cblock := NewCompressedBlock()
		cblock.StartTime = 1

		cblock.compressTime(t1)
		cblock.bitNum = 0
		cblock.byteNum = 0

		tm := cblock.readTime()
		if tm != t1 {
			t.Error("tm!=t1", tm, t1)
		}
	}

	{
		var t2 = Time(t1)
		cblock := NewCompressedBlock()
		cblock.StartTime = 1

		cblock.compressTime(t1)
		cblock.compressTime(t2)
		cblock.bitNum = 0
		cblock.byteNum = 0

		tm := cblock.readTime()
		if tm != t1 {
			t.Error("tm!=t1", tm, t1)
		}

		tm = cblock.readTime()
		if tm != t2 {
			t.Error("tm!=t2", tm, t2)
		}
	}

}
