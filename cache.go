package main

import (
	_ "fmt"
)

type LinearCache struct {
	MeasStorage
	meases []Meas
	pos    int64
	sz     int64
}

func NewLinearCache(sz int64) *LinearCache {
	res := &LinearCache{}
	res.meases = make([]Meas, sz, sz)
	res.pos = 0
	res.sz = sz
	return res
}

func (c *LinearCache) Add(m Meas) bool {
	if c.pos < c.sz {
		c.meases[c.pos] = m
		c.pos++
		return true
	}
	return false
}

func (c *LinearCache) Add_range(m []Meas) int64 {
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

func (c *LinearCache) Cap() int64 {
	return c.sz - c.pos
}
func (c *LinearCache) IsFull() bool {
	return c.sz == c.pos
}

func (c *LinearCache) Close() {

}

func (c LinearCache) ReadAll() []Meas {
	return c.meases[:c.pos]
}
func (c LinearCache) Read(ids []Id, from, to Time) []Meas {
	return c.ReadFltr(ids, 0, from, to)
}
func (c LinearCache) ReadFltr(ids []Id, flg Flag, from, to Time) []Meas {
	res := make([]Meas, 0, 0)
	for i := int64(0); i < c.pos; i++ {
		v := &c.meases[i]
		if idFltr(ids, v.Id) && inTimeInterval(from, to, v.Tstamp) && flagFltr(flg, v.Flg) {
			res = append(res, *v)
		}
	}
	return res
}
func (c LinearCache) TimePoint(ids []Id, time Time) []Meas {
	return c.Read(ids, 0, time)
}
func (c LinearCache) TimePointFltr(ids []Id, flg Flag, time Time) []Meas {
	return c.ReadFltr(ids, flg, 0, time)
}
