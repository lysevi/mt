package main

// paper http://www.vldb.org/pvldb/vol8/p1816-teller.pdf

import (
	_ "bytes"
	_ "encoding/binary"
	_ "fmt"
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

func (c CompressedBlock) compressTime(t Time) {
	panic("not implemented")
}

func (c CompressedBlock) delta_64(t Time) uint16 {
	return 0
}

func (c CompressedBlock) delta_256(t Time) uint16 {
	return 0
}

func (c CompressedBlock) delta_2048(t Time) uint16 {
	return 0
}

func (c CompressedBlock) delta_big(t Time) uint64 {
	return 0
}
