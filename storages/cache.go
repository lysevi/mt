package storages

import (
	"github.com/lysevi/mt/common"
)

type LinearCache struct {
	common.MeasStorage
	meases []common.Meas
	pos    int64
	sz     int64
}

func NewLinearCache(sz int64) *LinearCache {
	res := &LinearCache{}
	res.meases = make([]common.Meas, sz, sz)
	res.pos = 0
	res.sz = sz
	return res
}

func (c *LinearCache) Add(m common.Meas) bool {
	if c.pos < c.sz {
		c.meases[c.pos] = m
		c.pos++
		return true
	}
	return false
}

func (c *LinearCache) Add_range(m []common.Meas) int64 {
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

func (c LinearCache) ReadAll() []common.Meas {
	return c.meases[:c.pos]
}
func (c LinearCache) Read(ids []common.Id, from, to common.Time) []common.Meas {
	return make([]common.Meas, 0, 0)
}
func (c LinearCache) ReadFltr(ids []common.Id, flg common.Flag, from, to common.Time) []common.Meas {
	return make([]common.Meas, 0, 0)
}
func (c LinearCache) TimePoint(ids []common.Id, time common.Time) []common.Meas {
	return make([]common.Meas, 0, 0)
}
func (c LinearCache) TimePointFltr(ids []common.Id, flg common.Flag, time common.Time) []common.Meas {
	return make([]common.Meas, 0, 0)
}
