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
	cblock.writeTime(0)
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
	if cblock.data[14] != 255 || cblock.data[15] != 255 || cblock.data[16] != 255 || cblock.data[17] != 255 {
		t.Error(cblock.String())
	}
}

func TestCompressRead_64(t *testing.T) {
	cblock := NewCompressedBlock()
	res := cblock.delta_64(1)
	cblock.write_64(res)

	res = cblock.delta_64(64)
	cblock.write_64(res)

	res = cblock.delta_64(63)
	cblock.write_64(res)

	cblock.byteNum = 0
	cblock.bitNum = MAX_BIT

	if sr := cblock.readTime(0); sr != 1 {
		t.Error("sr:", sr, cblock.String())
	}

	if sr := cblock.readTime(0); sr != 64 {
		t.Error("sr:", sr)
	}

	if sr := cblock.readTime(0); sr != 63 {
		t.Error("sr:", sr)
	}
}

func TestCompressRead_256(t *testing.T) {
	cblock := NewCompressedBlock()
	res := cblock.delta_256(256)
	cblock.write_256(res) //110 1 0000 0000

	res = cblock.delta_256(255)
	cblock.write_256(res)

	res = cblock.delta_256(65)
	cblock.write_256(res)

	cblock.byteNum = 0
	cblock.bitNum = MAX_BIT

	if sr := cblock.readTime(0); sr != 256 {
		t.Error("sr:", sr, cblock.String())
	}

	if sr := cblock.readTime(0); sr != 255 {
		t.Error("sr:", sr)
	}

	if sr := cblock.readTime(0); sr != 65 {
		t.Error("sr:", sr)
	}
}

func TestCompressRead_2048(t *testing.T) {
	cblock := NewCompressedBlock()
	res := cblock.delta_2048(2048)
	cblock.write_2048(res) //1110 1000 0000 0000

	res = cblock.delta_2048(257)
	cblock.write_2048(res)

	res = cblock.delta_2048(4095)
	cblock.write_2048(res)

	cblock.byteNum = 0
	cblock.bitNum = MAX_BIT

	if sr := cblock.readTime(0); sr != 2048 {
		t.Error("sr:", sr, cblock.String())
	}

	if sr := cblock.readTime(0); sr != 257 {
		t.Error("sr:", sr)
	}

	if sr := cblock.readTime(0); sr != 4095 {
		t.Error("sr:", sr)
	}
}

func TestCompressRead_big(t *testing.T) {
	cblock := NewCompressedBlock()
	res := cblock.delta_big(2049)
	cblock.write_big(res) //111100000000000000000000100000000001

	res = cblock.delta_big(65535)
	cblock.write_big(res)

	res = cblock.delta_big(4095)
	cblock.write_big(res)

	cblock.byteNum = 0
	cblock.bitNum = MAX_BIT

	if sr := cblock.readTime(0); sr != 2049 {
		t.Error("sr:", sr, cblock.String())
	}

	if sr := cblock.readTime(0); sr != 65535 {
		t.Error("sr:", sr)
	}

	if sr := cblock.readTime(0); sr != 4095 {
		t.Error("sr:", sr)
	}
}

func TestCompressReadAll(t *testing.T) {
	cblock := NewCompressedBlock()
	iterations := 100
	for i := 0; i < iterations; i++ {
		cblock.write_64(cblock.delta_64(1))
		cblock.write_64(cblock.delta_64(64))
		cblock.write_64(cblock.delta_64(64))

		cblock.write_256(cblock.delta_256(256))
		cblock.write_256(cblock.delta_256(255))
		cblock.write_256(cblock.delta_256(65))

		cblock.write_2048(cblock.delta_2048(2048))
		cblock.write_2048(cblock.delta_2048(257))
		cblock.write_2048(cblock.delta_2048(4095))

		cblock.write_big(cblock.delta_big(2049))
		cblock.write_big(cblock.delta_big(65535))
		cblock.write_big(cblock.delta_big(4095))
	}

	cblock.byteNum = 0
	cblock.bitNum = MAX_BIT

	for i := 0; i < iterations; i++ {
		t_1 := cblock.readTime(0)
		t_64 := cblock.readTime(0)
		t_63 := cblock.readTime(0)
		if t_1 != 1 || t_64 != 64 || t_63 != 64 {
			t.Error("d64 read error i:", i, t_1, t_64, t_63)
			fmt.Print(cblock.String())
			return
		}

		t_256 := cblock.readTime(0)
		t_255 := cblock.readTime(0)
		t_65 := cblock.readTime(0)

		if t_256 != 256 || t_255 != 255 || t_65 != 65 {
			t.Error("d256 read error i:", i, t_256, t_255, t_65)
			fmt.Print(cblock.String())
			return
		}

		t_2048 := cblock.readTime(0)
		t_257 := cblock.readTime(0)
		t_4095 := cblock.readTime(0)
		if t_2048 != 2048 || t_257 != 257 || t_4095 != 4095 {
			t.Error("2048 error:", t_2048, t_257, t_4095, cblock.String())
		}

		t_2049 := cblock.readTime(0)
		t_65535 := cblock.readTime(0)
		t_4095 = cblock.readTime(0)
		if t_2049 != 2049 || t_65535 != 65535 || t_4095 != 4095 {
			t.Error("2048 error:", t_2048, t_257, t_4095, cblock.String())
		}
	}
}

func TestCompressTimeWrite(t *testing.T) {
	cblock := NewCompressedBlock()
	cblock.StartTime = 1
	iterations := 20
	tm := Time(1)
	times := []Time{}
	for i := 0; i < iterations; i++ {
		cblock.writeTime(tm)
		times = append(times, tm)
		tm *= 2
	}

	cblock.bitNum = MAX_BIT
	cblock.byteNum = 0

	readed_t := Time(cblock.StartTime)
	for _, tm = range times {
		readed_t = cblock.readTime(readed_t)
		if readed_t != tm {
			t.Error("read error:", readed_t, tm)
		}
	}
}

func TestCompressValue(t *testing.T) {
	{
		cblock := NewCompressedBlock()
		cblock.compressValue(10, 10, 0, 0)
		if cblock.data[0] != 0 || cblock.bitNum != MAX_BIT-1 {
			t.Error("value compress on equals: ", cblock.data[0], cblock.String())
		}
	}
	{
		cblock := NewCompressedBlock()
		zeros := cblock.leadingZeros(1080863910568919040) //0001....
		if zeros != 4 {
			t.Error("cblock.leadingZeros:", zeros)
		}

		zeros = cblock.tailngZeros(240) // 11110000
		if zeros != 4 {
			t.Error("cblock.tailngZeros:", zeros)
		}
	}
	{
		cblock := NewCompressedBlock()
		v1 := uint64(240) //1111 0000
		v2 := uint64(224) //1110 0000

		cblock.compressValue(v1, v2, 0, 0)
		if cblock.data[0] != 251 {
			t.Error("copmress error", cblock.data[0], cblock.String())
		}

		if cblock.data[1] != 18 {
			t.Error("copmress error", cblock.data[1], cblock.String())
		}

		if cblock.prevValue != v2 {
			t.Error("cblock.prevValue!=v2", cblock.prevValue)
		}
	}
	{
		cblock := NewCompressedBlock()
		v1 := uint64(240) //1111 0000
		v2 := uint64(224) //1110 0000

		cblock.compressValue(v1, v2, 59, 4)
		if cblock.data[0] != 160 {
			t.Error("copmress error", cblock.data[0], cblock.String())
		}
	}
	{
		cblock := NewCompressedBlock()
		v1 := uint64(240) //1111 0000
		v2 := uint64(224) //1110 0000
		cblock.writeValue(v1)
		if cblock.prevValue != v1 || cblock.data[0] != 0 {
			t.Error(cblock.prevValue, cblock.data[0])
		}
		cblock.writeValue(v2)

		if cblock.data[0] != 251 || cblock.data[1] != 18 {
			t.Error(cblock.data[0], cblock.data[1], cblock.String())
		}
	}

	{ // cur==prev
		cblock := NewCompressedBlock()
		v1 := uint64(240)
		v2 := uint64(240)
		cblock.writeValue(v1)
		cblock.writeValue(v2)
		cblock.bitNum = MAX_BIT
		cblock.byteNum = 0

		res := cblock.readValue(cblock.prevValue)
		if res != v2 {
			t.Error("res!=v2", res, v2)
		}
	}

	{ // cur!=prev
		cblock := NewCompressedBlock()
		v1 := uint64(240)
		v2 := uint64(96)
		v3 := uint64(176)
		cblock.writeValue(32)
		cblock.writeValue(v1)
		cblock.writeValue(v2)
		cblock.writeValue(v3)
		cblock.bitNum = MAX_BIT
		cblock.byteNum = 0

		res := cblock.readValue(cblock.startValue)
		if res != v1 {
			t.Error("res!=v2", res, v1)
		}

		res = cblock.readValue(res)
		if res != v2 {
			t.Error("res!=v2", res, v2)
		}

		res = cblock.readValue(res)
		if res != v3 {
			t.Error("res!=v3", res, v3)
		}
	}

	{ // tail/lead is equals
		cblock := NewCompressedBlock()
		v1 := uint64(3840) //111100000000
		v2 := uint64(3356) //111100010000

		cblock.writeValue(1)
		cblock.writeValue(v1)
		cblock.writeValue(v2)
		cblock.bitNum = MAX_BIT
		cblock.byteNum = 0

		res := cblock.readValue(cblock.startValue)
		if res != v1 {
			t.Error("res!=cblock.startValue", res, cblock.startValue)
		}

		res = cblock.readValue(res)
		if res != v2 {
			t.Error("res!=v2", res, v2)
		}
	}

	{ // tail/lead not equals
		cblock := NewCompressedBlock()
		v1 := uint64(3840) //111100000000
		v2 := uint64(3328) //110100000000

		cblock.writeValue(0)
		cblock.writeValue(v1)
		cblock.writeValue(v2)
		cblock.bitNum = MAX_BIT
		cblock.byteNum = 0

		res := cblock.readValue(cblock.startValue)
		if res != v1 {
			t.Error("res!=cblock.startValue", res, cblock.startValue)
		}

		res = cblock.readValue(res)
		if res != v2 {
			t.Error("res!=v2", res, v2)
		}
	}

	{
		cblock := NewCompressedBlock()
		values := []uint64{}
		delta := uint64(1)
		for i := uint64(0); i < 50; i++ {
			v := i * delta
			cblock.writeValue(v)
			values = append(values, v)
			delta *= 2
		}

		if cblock.startValue != 0 {
			t.Error("cblock.startValue != 0", cblock.startValue)
		}

		cblock.bitNum = MAX_BIT
		cblock.byteNum = 0

		readed_v := cblock.startValue
		for i, v := range values[1:] {
			readed_v = cblock.readValue(readed_v)
			if readed_v != v {
				t.Error("readed_v!=v", i, readed_v, v)
			}
		}
	}
}

func TestCompressFlag(t *testing.T) {
	{
		cblock := NewCompressedBlock()

		cblock.writeFlag(10)
		cblock.firstValue = false
		cblock.writeFlag(10)

		if cblock.data[0] != 0 || cblock.byteNum != 0 || cblock.bitNum != MAX_BIT-1 {
			t.Error(cblock.data[0], cblock.byteNum, cblock.bitNum)
		}

	}

	{
		cblock := NewCompressedBlock()

		cblock.writeFlag(10)
		cblock.firstValue = false
		cblock.writeFlag(11)

		if cblock.data[0] != 128 || cblock.byteNum != 8 {
			t.Error(cblock.data[0], cblock.byteNum)
		}

	}
	{
		cblock := NewCompressedBlock()

		flags := []Flag{}
		cblock.writeFlag(0)
		cblock.firstValue = false
		for i := Flag(1); i < 10; i++ {
			cblock.writeFlag(i)

			cblock.writeFlag(i)
			flags = append(flags, i)
		}

		cblock.bitNum = MAX_BIT
		cblock.byteNum = 0
		readed_flag := cblock.prevFlag

		for i, v := range flags {
			readed_flag = cblock.readFlag(readed_flag)
			if readed_flag != v {
				t.Error("readed_flag!=v", i, readed_flag, v, cblock.String())
			}

			readed_flag = cblock.readFlag(readed_flag)
			if readed_flag != v {
				t.Error("readed_flag!=v", i, readed_flag, v, cblock.String())
			}

		}
	}
}

func TestCompressMeas(t *testing.T) {
	{
		cblock := NewCompressedBlock()

		m1 := NewMeas(1, Time(1), int64(2), Flag(3))
		cblock.Add(m1)

		m2 := NewMeas(1, Time(2), int64(5), Flag(4))
		cblock.Add(m2)

		readed := cblock.ReadAll()

		if len(readed) != 2 {
			t.Error("len(readed)!=2", len(readed), readed)
		} else {
			if !measEqual(readed[0], m1) {
				t.Error("measEqual(readed[0], m1)", readed[0].String(), m1.String())
			}
			if !measEqual(readed[1], m2) {
				t.Error("measEqual(readed[1], m2)", readed[1].String(), m2.String())
			}
		}
	}
	cblock := NewCompressedBlock()
	iterations := 100
	for i := 0; i < iterations; i++ {
		m := NewMeas(1, Time(i), int64(i), Flag(i))
		cblock.Add(m)
	}

	meases := cblock.ReadAll()

	if len(meases) != iterations {
		t.Error("len(meases)!=iterations", len(meases), iterations)
	}
	for i, v := range meases {
		if v.Id != 1 || v.Flg != Flag(i) || v.Tstamp != Time(i) || v.Value != int64(i) {
			t.Error("meas read error: ", v.String())
		}
	}
}

func TestCompresseBlockAddOne(t *testing.T) {
	lc := NewCompressedBlock()

	m := NewMeas(1, 10, 10, 0)
	lc.Add(m)
	res := lc.ReadAll()
	if len(res) != 1 || !measEqual(res[0], m) {
		t.Error("add error: ", res)
	}
}

func TestCompresseBlockMeasAdd(t *testing.T) {
	lc := NewCompressedBlock()
	checkCompressWriterAdd(t, lc)
}

func TestCompresseBlockMeasAddRange(t *testing.T) {
	lc := NewCompressedBlock()
	checkWriterAddRange(t, lc)
}

func TestCompresseBlockMeasAddRange_s(t *testing.T) {
	lc := NewCompressedBlock()
	checkStorageAddRange(t, lc)
}

func TestCompresseBlockMeasCheck(t *testing.T) {
	lc := NewCompressedBlock()
	checkStorage_singleId(t, lc, 0, 100, 5)
}

func TestCompressAddDifferentId(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("compressTime not panic")
		}
	}()

	cblock := NewCompressedBlock()
	cblock.Add(NewMeas(1, 1, 1, 1))
	cblock.Add(NewMeas(2, 1, 1, 1))
}
