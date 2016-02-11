package mt

// paper http://www.vldb.org/pvldb/vol8/p1816-teller.pdf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	//	"unsafe"
)

const MAX_BLOCK_SIZE = 1024 * 1024
const MAX_BIT = 7

type CompressedBlock struct {
	StartTime  Time
	prev_delta int64
	prev_time  Time
	byteNum    uint64 //cur byte pos
	bitNum     uint8  //cur bit  pos
	data       [MAX_BLOCK_SIZE]byte
}

func NewCompressedBlock() *CompressedBlock {
	res := CompressedBlock{}
	res.StartTime = 0
	res.byteNum = 0
	res.bitNum = MAX_BIT
	return &res
}

/// if first insertion
func (c *CompressedBlock) compressFirstTime(t Time) {
	buf := new(bytes.Buffer)
	delta := t - c.StartTime
	binary.Write(buf, binary.LittleEndian, delta)

	write_size := uint64(buf.Len())
	byte_array := buf.Bytes()
	j := 0
	for i := c.byteNum; i < c.byteNum+write_size; i++ {
		c.data[i] = byte_array[j]
		j++
	}
	c.prev_delta = 0
	c.prev_time = t
	c.byteNum += write_size + 1
}

func (c *CompressedBlock) compressTime(t Time) {
	if t < c.StartTime {
		panic(fmt.Errorf("compressTime:"))
	}

	if c.byteNum == 0 {
		c.compressFirstTime(t)
	} else {
		delta := int64(t) - int64(c.prev_time)
		D := (int64)(delta)

		fmt.Println("D=", D, "t=", t, c.prev_time, c.prev_delta)

		if D == 0 {
			cur_byte := &c.data[c.byteNum]
			*cur_byte = setBit(*cur_byte, c.bitNum, 0)
			c.incBit()
		} else {
			if D > (-63) && D < 64 {
				c.d_63(D)
			} else {
				if D > (-255) && D < 256 {
					c.d_256(D)
				} else {
					if D > (-2047) && D < 2048 {
						c.d_2048(D)
					} else {
						c.d_bigger(D)
					}
				}
			}
		}
		c.prev_time = t
		c.prev_delta = delta
	}
}

// on D [-63;64]
func (c *CompressedBlock) d_63(D int64) {
	cur_byte := &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, 1)
	c.incBit()

	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, 0)
	c.incBit()

	//	fmt.Println("[-63;64] D=", D)
	//	fmt.Println("[-63;64]   >>>! ", *cur_byte, c.bitNum, c.byteNum)

	bite_value := byte(D)
	for i := uint8(0); i <= 5; i++ {
		cur_byte = &c.data[c.byteNum]
		time_bit := getBit(bite_value, i)
		*cur_byte = setBit(*cur_byte, c.bitNum, time_bit)
		//		fmt.Println("[-63;64] >  ", i, *cur_byte, c.bitNum, time_bit)
		c.incBit()
	}
}

// on D [-255;256]
func (c *CompressedBlock) d_256(D int64) {
	cur_byte := &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, 1)
	c.incBit()

	//fmt.Println("  >>> ", *cur_byte, c.bitNum, c.byteNum)
	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, 1)
	c.incBit()

	//	fmt.Println("  >>> ", *cur_byte, c.bitNum, c.byteNum)
	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, 0)
	c.incBit()

	//fmt.Println("  >>>! ", *cur_byte, c.bitNum, c.byteNum)
	bite_value := byte(D)
	for i := uint8(0); i <= 7; i++ {
		cur_byte = &c.data[c.byteNum]
		time_bit := getBit(bite_value, i)
		*cur_byte = setBit(*cur_byte, c.bitNum, time_bit)
		c.incBit()
	}
}

// on D [-2047;2048]
func (c *CompressedBlock) d_2048(D int64) {
	cur_byte := &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, 1)
	c.incBit()

	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, 1)
	c.incBit()

	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, 1)
	c.incBit()

	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, 0)
	c.incBit()

	//fmt.Println("  >>>! ", *cur_byte, c.bitNum, c.byteNum)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint32(D))
	byte_array := buf.Bytes()
	//	fmt.Println("*********** in ", byte_array)
	for _, bite_value := range byte_array[0:3] {
		//fmt.Println("*** ", bite_value)
		for i := uint8(0); i <= 7; i++ {
			cur_byte = &c.data[c.byteNum]
			time_bit := getBit(bite_value, i)
			*cur_byte = setBit(*cur_byte, c.bitNum, time_bit)
			//			fmt.Println(">  ", i, *cur_byte, c.bitNum, time_bit, c.byteNum)
			c.incBit()
		}
	}
}

// on D> 2048
func (c *CompressedBlock) d_bigger(D int64) {
	cur_byte := &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, 1)
	c.incBit()

	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, 1)
	c.incBit()

	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, 1)
	c.incBit()

	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, 1)
	c.incBit()

	//fmt.Println("  >>>! ", *cur_byte, c.bitNum, c.byteNum)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint32(D))
	byte_array := buf.Bytes()
	//fmt.Println("*********** in ", byte_array)
	for i := 0; i <= 3; i++ {
		bite_value := byte_array[i]
		//fmt.Println("*** ", bite_value)
		for i := uint8(0); i <= 7; i++ {
			cur_byte = &c.data[c.byteNum]
			time_bit := getBit(bite_value, i)
			*cur_byte = setBit(*cur_byte, c.bitNum, time_bit)
			//			fmt.Println(">  ", i, *cur_byte, c.bitNum, time_bit, c.byteNum)
			c.incBit()
		}
	}
}

func (c *CompressedBlock) readTime(prev_readed Time) Time {
	if c.byteNum == 0 {
		b := c.data[0:8]
		buf := bytes.NewBuffer(b)
		var readed_delta Time
		binary.Read(buf, binary.LittleEndian, &readed_delta)
		c.byteNum = 9
		c.bitNum = MAX_BIT
		return c.StartTime + readed_delta
	}
	cur_byte := &c.data[c.byteNum]
	res1 := getBit(*cur_byte, c.bitNum)
	//	fmt.Println("! ", res1)
	c.incBit()

	if res1 == 0 {
		//		fmt.Println("zero !>>>> ", res1)
		return prev_readed

	}

	res2 := getBit(*cur_byte, c.bitNum)
	c.incBit()
	//	fmt.Println("!>>>> ", res1, res2)
	if res1 == 1 && res2 == 0 {
		fmt.Println("-63 63")
		res := byte(0)
		for i := uint8(0); i <= 5; i++ {
			cur_byte = &c.data[c.byteNum]
			time_bit := getBit(*cur_byte, c.bitNum)
			c.incBit()
			res = setBit(res, i, time_bit)
		}
		return prev_readed + Time(res)
	}

	res3 := getBit(*cur_byte, c.bitNum)
	c.incBit()
	fmt.Println("res: ", res1, res2, res3)
	if res1 == 1 && res2 == 1 && res3 == 0 {
		fmt.Println("-255 255")
		res := byte(0)
		for i := uint8(0); i <= 7; i++ {
			cur_byte = &c.data[c.byteNum]
			time_bit := getBit(*cur_byte, c.bitNum)
			c.incBit()
			res = setBit(res, i, time_bit)
		}
		return prev_readed + Time(res)
	}

	res4 := getBit(*cur_byte, c.bitNum)
	c.incBit()
	fmt.Println("res: ", res1, res2, res3, res4)
	if res1 == 1 && res2 == 1 && res3 == 1 && res4 == 0 {
		fmt.Println("[-2047, 2048]")
		b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		for b_num := 0; b_num < 3; b_num++ {
			res := &b[b_num]
			for i := uint8(0); i <= 7; i++ {
				cur_byte = &c.data[c.byteNum]
				time_bit := getBit(*cur_byte, c.bitNum)
				c.incBit()
				*res = setBit(*res, i, time_bit)
			}
		}

		buf := bytes.NewBuffer(b)
		fmt.Println(buf.Bytes())
		var readed_delta uint32 = 0
		binary.Read(buf, binary.LittleEndian, &readed_delta)
		return prev_readed + Time(readed_delta)
	} else {
		fmt.Println("D>2048")
		b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		for b_num := 0; b_num < 8; b_num++ {
			res := &b[b_num]
			for i := uint8(0); i <= 7; i++ {
				cur_byte = &c.data[c.byteNum]
				time_bit := getBit(*cur_byte, c.bitNum)
				c.incBit()
				*res = setBit(*res, i, time_bit)
			}
		}

		buf := bytes.NewBuffer(b)
		fmt.Println(buf.Bytes())
		var readed_delta uint32 = 0
		binary.Read(buf, binary.LittleEndian, &readed_delta)
		return prev_readed + Time(readed_delta)
	}

	fmt.Println("END ")
	return 0
}
func (c *CompressedBlock) incByte() {
	c.byteNum++
}

func (c *CompressedBlock) incBit() {
	c.bitNum--
	if c.bitNum > 7 { // c.bitNum is uint8. 0-1==255
		c.bitNum = MAX_BIT
		c.incByte()
	}
}

func (c *CompressedBlock) addMeas(m Meas) {
	c.compressTime(m.Tstamp)
}
