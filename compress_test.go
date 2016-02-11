package mt

import (
	_ "bytes"
	_ "encoding/binary"
	"fmt"
	"testing"
)

var _ = fmt.Sprintf(" ")

/*
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
	if cblock.prev_time != t1 {
		t.Error("cblock.prev_time != t1")
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
		var t2 = Time(60 - t1)
		cblock := NewCompressedBlock()
		cblock.StartTime = 1

		cblock.compressTime(t1)
		fmt.Println("delta: ", cblock.prev_delta)
		cblock.compressTime(t2)

		if cblock.bitNum != MAX_BIT {
			t.Error("2 cblock.bitNum!=MAX_BIT", cblock.bitNum)
		}

		if cblock.data[9] != 151 { // 1001 0111
			t.Error("2  cblock.data[9] != 151", cblock.data[9])
		}
	}

	{ //D is between [-255, 255]
		var t1 = Time(1)
		var t2 = Time(254 - t1)
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
		//1100 0111  1110 0000
		//|----199-|   224
		if cblock.data[9] != 199 {
			t.Error("cblock.data[9]!=199", cblock.data[9])
		}

		if cblock.data[10] != 224 {
			t.Error("cblock.data[9]!=224", cblock.data[10])
		}
	}

	{ //D is between  [-2047, 2048]
		var t1 = Time(1)
		var t2 = Time(2045)
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

		if checkBit(cblock.data[9], MAX_BIT-3) {
			t.Error("2 checkBit(cblock.data[0],MAX_BIT-3)", cblock.data[9])
		}
		//1110 0011  1111 1000
		//|----231-|   252
		if cblock.data[9] != 227 {
			t.Error("cblock.data[9]!=239", cblock.data[9])
		}

		if cblock.data[10] != 254 {
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

		if cblock.byteNum != 8+5 {
			t.Error("2 cblock.byteNum!=8+4 ", cblock.byteNum, 8+5)
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
		//1111 0101  0101 1101 1011 1101 1011
		//|----245-|   252
		if cblock.data[9] != 245 {
			t.Error("cblock.data[9]!=245", cblock.data[9])
		}

		if cblock.data[10] != 93 {
			t.Error("cblock.data[10]!=93", cblock.data[10])
		}

		if cblock.data[11] != 219 {
			t.Error("cblock.data[11]!=219", cblock.data[11])
		}
	}
}

func TestCompressTimeRead(t *testing.T) {
	var t1 = Time(2)
	{
		cblock := NewCompressedBlock()
		cblock.StartTime = 1

		cblock.compressTime(t1)
		cblock.bitNum = 0
		cblock.byteNum = 0

		tm := cblock.readTime(0)
		if tm != t1 {
			t.Error("tm!=t1", tm, t1)
		}
	}

	{ // D=0
		var t2 = Time(t1)
		cblock := NewCompressedBlock()
		cblock.StartTime = 1

		cblock.compressTime(t1)
		cblock.compressTime(t2)
		cblock.bitNum = 0
		cblock.byteNum = 0

		tm := cblock.readTime(0)
		if tm != t1 {
			t.Error("tm!=t1", tm, t1)
		}

		tm2 := cblock.readTime(tm)

		if tm2 != t2 {
			t.Error("tm2!=t2", tm2, t2)
		}
	}

	{ // D= [-63, 64],
		var t2 = Time(64 - t1)
		cblock := NewCompressedBlock()
		cblock.StartTime = 1

		cblock.compressTime(t1)
		cblock.compressTime(t2)
		cblock.bitNum = 0
		cblock.byteNum = 0

		tm := cblock.readTime(0)
		if tm != t1 {
			t.Error("tm!=t1", tm, t1)
		}

		tm = cblock.readTime(tm)
		if tm != t2 {
			t.Error("tm!=t2", tm, t2)
		}
	}

	{ // D= [-255, 256],
		var t2 = Time(t1 + 100)
		cblock := NewCompressedBlock()
		cblock.StartTime = 1

		cblock.compressTime(t1)
		cblock.compressTime(t2)
		cblock.bitNum = 0
		cblock.byteNum = 0

		tm := cblock.readTime(0)
		if tm != t1 {
			t.Error("tm!=t1", tm, t1)
		}

		tm = cblock.readTime(tm)
		if tm != t2 {
			t.Error("tm!=t2", tm, t2)
		}
	}

	{ // D= [-2047, 2048]
		var t2 = Time(t1 + 997)
		cblock := NewCompressedBlock()
		cblock.StartTime = 1

		cblock.compressTime(t1)
		cblock.compressTime(t2)
		cblock.bitNum = 0
		cblock.byteNum = 0

		tm := cblock.readTime(0)
		if tm != t1 {
			t.Error("tm!=t1", tm, t1)
		}

		tm = cblock.readTime(tm)
		if tm != t2 {
			t.Error("tm!=t2", tm, t2)
		}
	}

	{ // D= > 2048
		var t2 = Time(t1 + 2048)
		cblock := NewCompressedBlock()
		cblock.StartTime = 1

		cblock.compressTime(t1)
		cblock.compressTime(t2)
		cblock.bitNum = 0
		cblock.byteNum = 0

		tm := cblock.readTime(0)
		if tm != t1 {
			t.Error("tm!=t1", tm, t1)
		}

		tm = cblock.readTime(tm)
		if tm != t2 {
			t.Error("tm!=t2", tm, t2)
		}
	}

}
*/
func TestCompressTimeManyAppends(t *testing.T) {
	cblock := NewCompressedBlock()
	cblock.StartTime = 1

	deltaI := Time(1)
	times := []Time{}

	for i := Time(1); i < 10000; i += deltaI {
		fmt.Println("i:", i)
		//		fmt.Println(cblock.data[0:150])
		times = append(times, i)
		cblock.compressTime(i)

		//		old_byte := cblock.byteNum
		//		old_bit := cblock.bitNum
		//		readed_time := Time(0)
		//		for j, v := range times {
		//			cblock.bitNum = 0
		//			cblock.byteNum = 0
		//			readed_time = cblock.readTime(readed_time)
		//			if readed_time != v {
		//				t.Error("readed_time!=v ", readed_time, v, " j=", j, times)
		//				return
		//			}

		//		}
		//		cblock.bitNum = old_bit
		//		cblock.byteNum = old_byte
		//deltaI += 25
	}
	fmt.Println("count: ", len(times), times)
	cblock.bitNum = 0
	cblock.byteNum = 0
	readed_time := Time(0)
	for i, v := range times {

		readed_time = cblock.readTime(readed_time)
		if readed_time != v {
			t.Error("readed_time!=v ", readed_time, v, " i=", i)
		}
	}
}
