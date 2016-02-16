package main

import (
	"fmt"
)

var _ = fmt.Sprintf("")

type MemoryStorage struct {
	max_time Time
	cblocks  []*CompressedBlock
}

func NewMemoryStorage(sz int64) *MemoryStorage {
	res := &MemoryStorage{}
	res.max_time = 0
	return res
}

func (c *MemoryStorage) Add(m Meas) bool {
	if c.max_time < m.Tstamp {
		c.max_time = m.Tstamp
	}
	var freeBlock *CompressedBlock = nil
	for _, v := range c.cblocks {
		if v.id == m.Id && !v.IsFull() {
			freeBlock = v
			break
		}
	}
	success := false
	if freeBlock != nil {
		success = freeBlock.Add(m)
	} else {
		freeBlock = NewCompressedBlock()
		success = freeBlock.Add(m)
		c.cblocks = append(c.cblocks, freeBlock)
	}
	return success
}

func (c *MemoryStorage) Add_range(m []Meas) int64 {
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

func (c *MemoryStorage) Cap() int64 {
	return 1
}
func (c *MemoryStorage) IsFull() bool {
	return false
}

func (c *MemoryStorage) Close() {

}

func (c *MemoryStorage) ReadAll() []Meas {
	return c.Read([]Id{}, 0, c.max_time)
}
func (c *MemoryStorage) Read(ids []Id, from, to Time) []Meas {
	return c.ReadFltr(ids, 0, from, to)
}
func (c *MemoryStorage) ReadFltr(ids []Id, flg Flag, from, to Time) []Meas {
	res := []Meas{}

	for _, v := range c.cblocks {
		subres := v.ReadFltr(ids, flg, from, to)
		res = append(res, subres...)
	}
	return res
}
func (c *MemoryStorage) TimePoint(ids []Id, time Time) []Meas {
	return c.Read(ids, 0, time)
}
func (c *MemoryStorage) TimePointFltr(ids []Id, flg Flag, time Time) []Meas {
	return c.ReadFltr(ids, flg, 0, time)
}
