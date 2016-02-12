package main

// paper http://www.vldb.org/pvldb/vol8/p1816-teller.pdf

import (
	_ "bytes"
	"encoding/binary"
	"fmt"
	//	"unsafe"
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
	bts := []byte{0, 1}
	binary.LittleEndian.PutUint16(bts, uint16(t))
	bts[1] = 1
	subres := binary.LittleEndian.Uint16(bts)
	fmt.Println(">>> ", subres)
	return subres
}

func (c CompressedBlock) delta_256(t Time) uint16 {
	bts := []byte{0, 0}
	binary.LittleEndian.PutUint16(bts, uint16(t))
	bts[1] = setBit(bts[1], 3, 1)
	bts[1] = setBit(bts[1], 2, 1)
	bts[1] = setBit(bts[1], 1, 0)
	subres := binary.LittleEndian.Uint16(bts)
	return subres
}

func (c CompressedBlock) delta_2048(t Time) uint16 {
	bts := []byte{0, 0}
	binary.LittleEndian.PutUint16(bts, uint16(t))
	const num uint8 = 1
	bts[num] = setBit(bts[num], 7, 1)
	bts[num] = setBit(bts[num], 6, 1)
	bts[num] = setBit(bts[num], 5, 1)
	bts[num] = setBit(bts[num], 4, 0)
	subres := binary.LittleEndian.Uint16(bts)
	return uint16(subres)
}

func (c CompressedBlock) delta_big(t Time) uint64 {
	bts := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(bts, uint32(t))
	bts[4] = 15
	subres := binary.LittleEndian.Uint64(bts)
	return subres
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

}

func (c *CompressedBlock) write_2048(D uint16) {

}

func (c *CompressedBlock) write_big(D uint64) {

}
