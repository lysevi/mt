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
	prev_delta Time
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
	c.prev_delta = delta
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
		delta := t - c.prev_time
		D := (int64)(delta - c.prev_delta)

		//fmt.Println("D=", D)

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

	//	fmt.Println("D=", D)
	//	fmt.Println("  >>>! ", *cur_byte, c.bitNum, c.byteNum)

	bite_value := byte(D)
	for i := uint8(0); i <= 5; i++ {
		cur_byte = &c.data[c.byteNum]
		time_bit := getBit(bite_value, i)
		*cur_byte = setBit(*cur_byte, c.bitNum, time_bit)
		//		fmt.Println(">  ", i, *cur_byte, c.bitNum, time_bit)
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

	fmt.Println("  >>> ", *cur_byte, c.bitNum, c.byteNum)
	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, 0)
	c.incBit()

	//fmt.Println("  >>>! ", *cur_byte, c.bitNum, c.byteNum)
	bite_value := byte(D)
	for i := uint8(0); i <= 5; i++ {
		cur_byte = &c.data[c.byteNum]
		time_bit := getBit(bite_value, i)
		*cur_byte = setBit(*cur_byte, c.bitNum, time_bit)
		//fmt.Println(">  ", i, *cur_byte, c.bitNum, time_bit, c.byteNum)
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
	binary.Write(buf, binary.LittleEndian, uint16(D))
	byte_array := buf.Bytes()
	//fmt.Println(byte_array)
	for _, bite_value := range byte_array {
		//fmt.Println("*** ", bite_value)
		for i := uint8(0); i <= 5; i++ {
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
	//	fmt.Println(byte_array)
	for i := 0; i <= 3; i++ {
		bite_value := byte_array[i]
		//fmt.Println("*** ", bite_value)
		for i := uint8(0); i <= 5; i++ {
			cur_byte = &c.data[c.byteNum]
			time_bit := getBit(bite_value, i)
			*cur_byte = setBit(*cur_byte, c.bitNum, time_bit)
			//			fmt.Println(">  ", i, *cur_byte, c.bitNum, time_bit, c.byteNum)
			c.incBit()
		}
	}
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
