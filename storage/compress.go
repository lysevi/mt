package storage

// paper http://www.vldb.org/pvldb/vol8/p1816-teller.pdf

import (
	"fmt"
)

var _ = fmt.Sprintf("")

const (
	maxBlockSize          = (1024 * 1024)
	maxBit                = 7
	compressedMeasMaxSize = 33 + 65 + 64 // time max size + flag + value
)

type CompressedBlock struct {
	id         Id
	StartTime  Time
	prev_delta int64
	prev_time  Time
	max_time   Time
	firstValue bool
	startValue uint64
	prevLead   uint8
	prevTail   uint8
	prevValue  uint64

	startFlag uint64
	prevFlag  Flag
	byteNum   uint64 //cur byte pos
	bitNum    uint8  //cur bit  pos
	data      [maxBlockSize]uint8
}

type readStatus struct {
	byteNum    uint64 //cur byte pos
	bitNum     uint8  //cur bit  pos
	prevLead   uint8
	prevTail   uint8
	prevValue  uint64
	startValue uint64
}

func NewCompressedBlock() *CompressedBlock {
	res := CompressedBlock{}
	res.StartTime = 0
	res.byteNum = 0
	res.bitNum = maxBit
	res.firstValue = true
	res.id = -1
	return &res
}

func newReadStatus() readStatus {
	res := readStatus{}
	res.bitNum = maxBit
	res.byteNum = 0
	return res
}

func (c *readStatus) incByte() {
	c.byteNum++
	if c.byteNum >= maxBlockSize {
		panic("out of bound")
	}
}

func (c *readStatus) incBit() {
	c.bitNum--
	if c.bitNum > 7 { // c.bitNum is uint8. 0-1==255
		c.bitNum = maxBit
		c.incByte()
	}
}

func (c *CompressedBlock) incByte() {
	c.byteNum++
	if c.byteNum >= maxBlockSize {
		panic("out of bound")
	}
}

func (c *CompressedBlock) incBit() {
	c.bitNum--
	if c.bitNum > 7 { // c.bitNum is uint8. 0-1==255
		c.bitNum = maxBit
		c.incByte()
	}
}

func (c CompressedBlock) String() string {
	res := fmt.Sprintf("cblock byte: %v bit: %v[", c.byteNum, c.bitNum)
	for i := uint64(0); i <= c.byteNum; i++ {
		cur_byte := c.data[i]
		res += fmt.Sprintf("%v: ", i)
		for j := (maxBit); j >= 0; j-- {
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
	return uint16(t) | uint16(256)
}

func (c CompressedBlock) delta_256(t Time) uint16 {
	return uint16(t) | uint16(3072)
}

func (c CompressedBlock) delta_2048(t Time) uint16 {
	return uint16(t) | uint16(57344)
}

func (c CompressedBlock) delta_big(t Time) uint64 {
	return uint64(t) | uint64(64424509440)
}

func (c *CompressedBlock) write_64(D uint16) {
	cur_byte := &c.data[c.byteNum]
	bvalue := getBit16(D, 8)
	*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
	c.incBit()

	if c.byteNum == maxBit {
		cur_byte := &c.data[c.byteNum]
		*cur_byte = (*cur_byte) | byte(D)
		c.byteNum++
	} else {
		step_h := maxBit - c.bitNum
		step_l := c.bitNum + 1
		high := byte(D) >> step_h
		low := byte(D) << (step_l)

		cur_byte := &c.data[c.byteNum]
		*cur_byte = (*cur_byte) | high
		c.byteNum++

		cur_byte = &c.data[c.byteNum]
		*cur_byte = (*cur_byte) | low
		c.bitNum = maxBit - step_h
	}

}

func (c *CompressedBlock) write_256(D uint16) {
	cur_byte := &c.data[c.byteNum]
	bvalue := getBit16(D, 11)
	*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
	c.incBit()

	bvalue = getBit16(D, 10)
	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
	c.incBit()

	bvalue = getBit16(D, 9)
	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
	c.incBit()

	bvalue = getBit16(D, 8)
	cur_byte = &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
	c.incBit()

	for i := int8(7); i >= 0; i-- {
		bvalue = getBit16(D, uint8(i))
		cur_byte = &c.data[c.byteNum]
		*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
		c.incBit()
	}
}

func (c *CompressedBlock) write_2048(D uint16) {
	for bn := 15; bn >= 0; bn-- {
		bvalue := getBit16(D, uint8(bn))
		cur_byte := &c.data[c.byteNum]
		*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
		c.incBit()
	}
}

func (c *CompressedBlock) write_big(D uint64) {
	for i := 0; i < 4; i++ {
		cur_byte := &c.data[c.byteNum]
		*cur_byte = setBit(*cur_byte, c.bitNum, 1)
		c.incBit()
	}

	for bn := 31; bn >= 0; bn-- {
		bvalue := getBit64(D, uint8(bn))
		cur_byte := &c.data[c.byteNum]
		*cur_byte = setBit(*cur_byte, c.bitNum, bvalue)
		c.incBit()
	}
}

func (c *CompressedBlock) readTime(prev_readed Time, rs *readStatus) Time {
	cur_byte := &c.data[rs.byteNum]
	//	fmt.Println("cur_byte: ", *cur_byte)

	res1 := getBit(*cur_byte, rs.bitNum)
	rs.incBit()
	if res1 == 0 {
		return prev_readed

	}
	cur_byte = &c.data[rs.byteNum]
	res2 := getBit(*cur_byte, rs.bitNum)
	//	fmt.Println("pos: ", c.byteNum, c.bitNum, "ress:", res1, res2)
	rs.incBit()

	if res1 == 1 && res2 == 0 {
		//		fmt.Println("R -63 63")
		res := byte(0)

		for i := int8(6); i >= 0; i-- {
			cur_byte = &c.data[rs.byteNum]
			time_bit := getBit(*cur_byte, rs.bitNum)
			rs.incBit()
			res = setBit(res, uint8(i), time_bit)
		}
		return prev_readed + Time(res)
	}
	cur_byte = &c.data[rs.byteNum]
	res3 := getBit(*cur_byte, rs.bitNum)
	rs.incBit()

	if res1 == 1 && res2 == 1 && res3 == 0 {
		//		fmt.Println("R -255 256")
		res := uint16(0)

		for i := int8(8); i >= 0; i-- {
			cur_byte = &c.data[rs.byteNum]
			time_bit := getBit(*cur_byte, rs.bitNum)
			rs.incBit()
			res = setBit16(res, uint8(i), time_bit)
		}
		return prev_readed + Time(res)
	}
	cur_byte = &c.data[rs.byteNum]
	res4 := getBit(*cur_byte, rs.bitNum)
	rs.incBit()

	if res1 == 1 && res2 == 1 && res3 == 1 && res4 == 0 {
		//		fmt.Println("R -2047 2048")
		res := uint32(0)

		for i := int8(11); i >= 0; i-- {
			cur_byte = &c.data[rs.byteNum]
			time_bit := getBit(*cur_byte, rs.bitNum)
			rs.incBit()
			res = setBit32(res, uint8(i), time_bit)
		}
		return prev_readed + Time(res)
	}

	if res1 == 1 && res2 == 1 && res3 == 1 && res4 == 1 {
		//		fmt.Println("R big")
		res := uint32(0)

		for i := int8(31); i >= 0; i-- {
			cur_byte := &c.data[rs.byteNum]
			time_bit := getBit(*cur_byte, rs.bitNum)
			rs.incBit()
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
	//fmt.Println(D)
	if D == 0 {
		cur_byte := &c.data[c.byteNum]
		*cur_byte = setBit(*cur_byte, c.bitNum, 0)
		c.incBit()
	} else {
		if D < 127 {
			c.write_64(uint16(D) | uint16(256))
			//c.write_64(c.delta_64(Time(D)))
		} else {
			if D < 511 {
				c.write_256(uint16(D) | uint16(3072))
				//c.write_256(c.delta_256(Time(D)))
			} else {
				if D < 4095 {
					c.write_2048(uint16(D) | uint16(57344))
					//c.write_2048(c.delta_2048(Time(D)))
				} else {
					c.write_big(uint64(D) | uint64(64424509440))
					//c.write_big(c.delta_big(Time(D)))
				}
			}
		}
	}
	c.prev_time = t
	c.prev_delta = delta
}

func (c *CompressedBlock) leadingZeros(xor uint64) uint8 {
	leading_zeros := uint8(0)
	for i := 63; i >= 0; i-- {
		if checkBit64(xor, uint8(i)) {
			break
		}
		leading_zeros++
	}
	return leading_zeros
}

func (c *CompressedBlock) tailngZeros(xor uint64) uint8 {
	tailng_zeros := uint8(0)
	for i := uint8(0); i <= 63; i++ {
		if checkBit64(xor, i) {
			break
		}
		tailng_zeros++
	}
	return tailng_zeros
}

func (c *CompressedBlock) compressValue(prev, cur uint64, prevLead, prevTail uint8) {
	xor := prev ^ cur
	if xor == 0 {
		cur_byte := &c.data[c.byteNum]
		*cur_byte = setBit(*cur_byte, c.bitNum, 0)
		c.incBit()
		return
	}

	cur_byte := &c.data[c.byteNum]
	*cur_byte = setBit(*cur_byte, c.bitNum, 1) //1000....000
	c.incBit()
	tailZeros := c.tailngZeros(xor)
	leadingZeros := c.leadingZeros(xor)

	if prevLead == leadingZeros && prevTail == tailZeros {
		cur_byte := &c.data[c.byteNum]
		*cur_byte = setBit(*cur_byte, c.bitNum, 0)
		c.incBit()

	} else {
		cur_byte := &c.data[c.byteNum]
		*cur_byte = setBit(*cur_byte, c.bitNum, 1)
		c.incBit()

		for i := 5; i >= 0; i-- {
			b := getBit(leadingZeros, uint8(i))
			cur_byte := &c.data[c.byteNum]
			*cur_byte = setBit(*cur_byte, c.bitNum, b)
			c.incBit()
		}

		for i := 5; i >= 0; i-- {
			b := getBit(tailZeros, uint8(i))
			cur_byte := &c.data[c.byteNum]
			*cur_byte = setBit(*cur_byte, c.bitNum, b)
			c.incBit()
		}
	}

	for i := int8(63 - leadingZeros); i >= int8(tailZeros); i-- {
		b := getBit64(xor, uint8(i))
		cur_byte := &c.data[c.byteNum]
		*cur_byte = setBit(*cur_byte, c.bitNum, b)
		c.incBit()
	}

	c.prevValue = cur
	c.prevLead = leadingZeros
	c.prevTail = tailZeros
}

func (c *CompressedBlock) writeValue(value uint64) {
	if c.firstValue {
		c.startValue = value
		c.prevValue = value
		c.firstValue = false
		return
	}

	c.compressValue(c.prevValue, value, c.prevLead, c.prevTail)
}

func (c *CompressedBlock) readValue(prev uint64, rs *readStatus) uint64 {
	cur_byte := &c.data[rs.byteNum]
	res0 := getBit(*cur_byte, rs.bitNum)
	rs.incBit()

	if res0 == 0 {
		//		fmt.Println("res0 == 0")
		return rs.prevValue
	}

	cur_byte = &c.data[rs.byteNum]
	res1 := getBit(*cur_byte, rs.bitNum)
	rs.incBit()

	if res1 == 1 {
		//		fmt.Println("res0 ==1")
		leading := uint8(0)
		for i := 5; i >= 0; i-- {
			cur_byte = &c.data[rs.byteNum]
			b := getBit(*cur_byte, rs.bitNum)
			rs.incBit()
			leading = setBit(leading, uint8(i), b)
		}

		tail := uint8(0)
		for i := 5; i >= 0; i-- {
			cur_byte = &c.data[rs.byteNum]
			b := getBit(*cur_byte, rs.bitNum)
			rs.incBit()
			tail = setBit(tail, uint8(i), b)
		}
		result := uint64(0)
		//		fmt.Println("lead/tail", leading, tail)
		for i := int8(63 - leading); i >= int8(tail); i-- {

			cur_byte = &c.data[rs.byteNum]
			b := getBit(*cur_byte, rs.bitNum)
			rs.incBit()
			result = setBit64(result, uint8(i), b)
		}
		//		fmt.Println("xor: ", result, "prev:", prev)
		rs.prevLead = leading
		rs.prevTail = tail
		return result ^ prev
	} else {
		result := uint64(0)
		leading := rs.prevLead
		tail := rs.prevTail
		//		fmt.Println("lead/tail", leading, tail)
		for i := int8(63 - leading); i >= int8(tail); i-- {

			cur_byte = &c.data[rs.byteNum]
			b := getBit(*cur_byte, rs.bitNum)
			rs.incBit()
			result = setBit64(result, uint8(i), b)
		}
		//		fmt.Println("xor: ", result)
		return result ^ prev
	}
}

func (c *CompressedBlock) writeFlag(f Flag) {
	if c.firstValue {
		c.prevFlag = f
		return
	}

	if c.prevFlag == f {
		cur_byte := &c.data[c.byteNum]
		*cur_byte = setBit(*cur_byte, c.bitNum, 0)
		c.incBit()
		return
	} else {
		cur_byte := &c.data[c.byteNum]
		*cur_byte = setBit(*cur_byte, c.bitNum, 1)
		c.incBit()

		for i := int8(63); i >= int8(0); i-- {
			b := getBit64(uint64(f), uint8(i))
			cur_byte = &c.data[c.byteNum]
			*cur_byte = setBit(*cur_byte, c.bitNum, b)
			c.incBit()
		}
	}
}

func (c *CompressedBlock) readFlag(prev Flag, rs *readStatus) Flag {
	cur_byte := &c.data[rs.byteNum]
	b := getBit(*cur_byte, rs.bitNum)
	rs.incBit()
	if b == 0 {
		return prev
	} else {
		result := uint64(0)
		for i := int8(63); i >= int8(0); i-- {
			cur_byte = &c.data[rs.byteNum]
			b := getBit(*cur_byte, rs.bitNum)
			rs.incBit()
			result = setBit64(uint64(result), uint8(i), b)
		}
		return Flag(result)
	}
}

func (c *CompressedBlock) Add(m Meas) bool {

	if !c.firstValue && m.Id != c.id {
		panic("m.Id!=c.id")
	}

	if m.Tstamp > c.max_time {
		c.max_time = m.Tstamp
	}

	if c.firstValue {
		c.id = m.Id
		c.StartTime = m.Tstamp
		c.writeFlag(Flag(m.Flg))
		c.writeValue(uint64(m.Value))
	} else {
		c.writeTime(m.Tstamp)
		c.writeFlag(m.Flg)
		c.writeValue(uint64(m.Value))
		c.firstValue = false
	}
	return true
}

func (c *CompressedBlock) Add_range(m []Meas) int64 {
	var res int64 = 0
	for _, v := range m {
		add_res := c.Add(v)
		if !add_res {
			break
		}
		res++
	}
	return res
}

func (c *CompressedBlock) Cap() int64 {
	in_bytes := int64((maxBlockSize - c.byteNum) / compressedMeasMaxSize)
	if !c.firstValue {
		in_bytes--
	}
	return in_bytes
}

func (c *CompressedBlock) IsFull() bool {
	return (maxBlockSize - c.byteNum) < compressedMeasMaxSize
}

func (c *CompressedBlock) Close() {}

func (c *CompressedBlock) ReadAll() []Meas {
	return c.Read([]Id{}, 0, c.max_time)
}

func (c *CompressedBlock) Read(ids []Id, from, to Time) []Meas {
	return c.ReadFltr(ids, 0, from, to)
}

func (c *CompressedBlock) ReadFltr(ids []Id, flg Flag, from, to Time) []Meas {
	if len(ids) != 0 && !idFltr(ids, c.id) {
		return []Meas{}
	}

	rs := newReadStatus()
	rs.startValue = c.startValue
	rs.prevValue = c.prevValue
	rs.prevLead = c.prevLead
	rs.prevTail = c.prevTail

	prev_time := c.StartTime
	prev_value := c.startValue
	prev_flag := c.prevFlag

	m := NewMeas(c.id, prev_time, int64(prev_value), prev_flag)
	result := []Meas{}
	if inTimeInterval(from, to, m.Tstamp) && flagFltr(flg, m.Flg) && !c.firstValue {
		result = append(result, m)
	}

	if c.byteNum == 0 && c.bitNum == maxBit {
		return result
	}

	for {
		prev_time = c.readTime(prev_time, &rs)
		prev_flag = c.readFlag(prev_flag, &rs)
		prev_value = c.readValue(prev_value, &rs)
		if inTimeInterval(from, to, prev_time) && flagFltr(flg, prev_flag) {
			m = NewMeas(c.id, prev_time, int64(prev_value), prev_flag)

			result = append(result, m)
		}
		if c.byteNum == rs.byteNum && c.bitNum >= rs.bitNum {
			break
		}
	}
	return result
}

func (c *CompressedBlock) TimePoint(ids []Id, time Time) []Meas {
	return c.Read(ids, 0, time)
}
func (c *CompressedBlock) TimePointFltr(ids []Id, flg Flag, time Time) []Meas {
	return c.ReadFltr(ids, flg, 0, time)
}
