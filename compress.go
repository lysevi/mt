package main

// paper http://www.vldb.org/pvldb/vol8/p1816-teller.pdf

import (
	"encoding/binary"
	"fmt"
)

var _ = fmt.Sprintf("")

const MAX_BLOCK_SIZE = 1024 * 1024
const MAX_BIT = 7

type CompressedBlock struct {
	StartTime  Time
	prev_delta int64
	prev_time  Time
	byteNum    uint64 //cur byte pos
	bitNum     uint8  //cur bit  pos
	data       [MAX_BLOCK_SIZE]uint8
}

func NewCompressedBlock() *CompressedBlock {
	res := CompressedBlock{}
	res.StartTime = 0
	res.byteNum = 0
	res.bitNum = MAX_BIT
	return &res
}

func (c *CompressedBlock) incByte() {
	c.byteNum++
	if c.byteNum >= MAX_BLOCK_SIZE {
		panic("out of bound")
	}
}

func (c *CompressedBlock) incBit() {
	c.bitNum--
	if c.bitNum > 7 { // c.bitNum is uint8. 0-1==255
		c.bitNum = MAX_BIT
		c.incByte()
	}
}

func (c CompressedBlock) String() string {
	res := "["
	for i := uint64(0); i < 20; i++ {
		cur_byte := c.data[i]
		res += fmt.Sprintf("%v: ", i)
		for j := (MAX_BIT); j >= 0; j-- {
			if j == 3 {
				res += " "
			}
			if checkBit(cur_byte, uint8(j)) {
				res += "1"
			} else {
				res += "0"
			}
		}
		res += "\n "
		if i == 7 {
			res += "\n "
		}
	}
	res += "]"
	return res
}

func (c CompressedBlock) delta_64(t Time) uint16 {
	res := uint16(t)
	res = setBit16(res, 8, 1)
	return res
}

func (c CompressedBlock) delta_256(t Time) uint16 {
	res := uint16(t)
	res = setBit16(res, 11, 1)
	res = setBit16(res, 10, 1)
	res = setBit16(res, 9, 0)
	return res
}

func (c CompressedBlock) delta_2048(t Time) uint16 {
	res := uint16(t)
	res = setBit16(res, 15, 1)
	res = setBit16(res, 14, 1)
	res = setBit16(res, 13, 1)
	res = setBit16(res, 12, 0)

	return res
}

func (c CompressedBlock) delta_big(t Time) uint64 {
	res := uint64(t)
	res = setBit64(res, 35, 1)
	res = setBit64(res, 34, 1)
	res = setBit64(res, 33, 1)
	res = setBit64(res, 32, 1)

	return res
}

func (c *CompressedBlock) write_64(D uint16) {
	bts := []byte{0, 0}
	binary.LittleEndian.PutUint16(bts, D)

	cur_byte := &c.data[c.byteNum]
	bvalue := getBit(bts[1], 0)
	*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
	c.incBit()

	for i := int8(7); i >= 0; i-- {
		bvalue = getBit(bts[0], uint8(i))
		cur_byte := &c.data[c.byteNum]
		*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
		c.incBit()
	}
}

func (c *CompressedBlock) write_256(D uint16) {
	bts := []byte{0, 0}
	binary.LittleEndian.PutUint16(bts, D)

	cur_byte := &c.data[c.byteNum]
	bvalue := getBit(bts[1], 3)
	*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
	c.incBit()

	bvalue = getBit(bts[1], 2)
	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
	c.incBit()

	bvalue = getBit(bts[1], 1)
	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
	c.incBit()

	bvalue = getBit(bts[1], 0)
	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
	c.incBit()

	for i := int8(7); i >= 0; i-- {
		bvalue = getBit(bts[0], uint8(i))
		cur_byte = &c.data[c.byteNum]
		*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
		c.incBit()
	}
}

func (c *CompressedBlock) write_2048(D uint16) {
	bts := []byte{0, 0}
	binary.LittleEndian.PutUint16(bts, D)

	for bn := range bts {
		b := bts[1-bn] //reverse iterations
		for i := int8(7); i >= 0; i-- {
			bvalue := getBit(b, uint8(i))
			cur_byte := &c.data[c.byteNum]
			*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
			c.incBit()
		}
	}
}

func (c *CompressedBlock) write_big(D uint64) {
	bts := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint64(bts, D)

	for i := 0; i < 4; i++ {
		cur_byte := &c.data[c.byteNum]
		*cur_byte = setBit(*cur_byte, c.bitNum, 1)
		c.incBit()
	}

	for bn := 3; bn >= 0; bn-- {
		b := bts[bn]
		for i := int8(7); i >= 0; i-- {
			bvalue := getBit(b, uint8(i))
			cur_byte := &c.data[c.byteNum]
			*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
			c.incBit()
		}
	}
}

func (c *CompressedBlock) readTime(prev_readed Time) Time {
	cur_byte := &c.data[c.byteNum]
	//	fmt.Println("cur_byte: ", *cur_byte)

	res1 := getBit(*cur_byte, c.bitNum)
	c.incBit()
	if res1 == 0 {
		return prev_readed

	}
	cur_byte = &c.data[c.byteNum]
	res2 := getBit(*cur_byte, c.bitNum)
	//	fmt.Println("pos: ", c.byteNum, c.bitNum, "ress:", res1, res2)
	c.incBit()

	if res1 == 1 && res2 == 0 {
		//		fmt.Println("R -63 63")
		res := byte(0)

		for i := int8(6); i >= 0; i-- {
			cur_byte = &c.data[c.byteNum]
			time_bit := getBit(*cur_byte, c.bitNum)
			c.incBit()
			res = setBit(res, uint8(i), time_bit)
		}
		return prev_readed + Time(res)
	}
	cur_byte = &c.data[c.byteNum]
	res3 := getBit(*cur_byte, c.bitNum)
	c.incBit()

	if res1 == 1 && res2 == 1 && res3 == 0 {
		//		fmt.Println("R -255 256")
		res := uint16(0)

		for i := int8(8); i >= 0; i-- {
			cur_byte = &c.data[c.byteNum]
			time_bit := getBit(*cur_byte, c.bitNum)
			c.incBit()
			res = setBit16(res, uint8(i), time_bit)
		}
		return prev_readed + Time(res)
	}
	cur_byte = &c.data[c.byteNum]
	res4 := getBit(*cur_byte, c.bitNum)
	c.incBit()

	if res1 == 1 && res2 == 1 && res3 == 1 && res4 == 0 {
		//		fmt.Println("R -2047 2048")
		res := uint32(0)

		for i := int8(11); i >= 0; i-- {
			cur_byte = &c.data[c.byteNum]
			time_bit := getBit(*cur_byte, c.bitNum)
			c.incBit()
			res = setBit32(res, uint8(i), time_bit)
		}
		return prev_readed + Time(res)
	}

	if res1 == 1 && res2 == 1 && res3 == 1 && res4 == 1 {
		//		fmt.Println("R big")
		res := uint32(0)

		for i := int8(31); i >= 0; i-- {
			cur_byte := &c.data[c.byteNum]
			time_bit := getBit(*cur_byte, c.bitNum)
			c.incBit()
			res = setBit32(res, uint8(i), time_bit)
		}
		return prev_readed + Time(res)
	}
	panic("read error!!!")
}

func (c *CompressedBlock) writeTime(t Time) {
	if t < c.StartTime {
		panic(fmt.Errorf("compressTime:"))
	}

	if c.byteNum == 0 {
		c.prev_time = c.StartTime
	}

	delta := int64(t) - int64(c.prev_time)
	D := (int64)(delta)

	if D == 0 {
		cur_byte := &c.data[c.byteNum]
		*cur_byte = setBit(*cur_byte, c.bitNum, 0)
		c.incBit()
	} else {
		if D > (-63) && D < 64 {
			c.write_64(c.delta_64(Time(D)))
		} else {
			if D > (-255) && D < 256 {
				c.write_256(c.delta_256(Time(D)))
			} else {
				if D > (-2047) && D < 2048 {
					c.write_2048(c.delta_2048(Time(D)))
				} else {
					c.write_big(c.delta_big(Time(D)))
				}
			}
		}
	}
	c.prev_time = t
	c.prev_delta = delta
}
