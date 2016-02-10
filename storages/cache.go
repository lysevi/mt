package storages

import (
	"github.com/lysevi/mt/common"
)

type LinearCache struct {
	common.MeasStorage
}

func NewLinearCache() *LinearCache {
	return &LinearCache{}
}

func (c *LinearCache) Add(m common.Meas) bool {
	return false
}

func (c *LinearCache) Add_range(m []common.Meas) int64 {
	return 0
}

func (c *LinearCache) Cap() int64 {
	return 0
}
func (c *LinearCache) IsFull() bool {
	return false
}

func (c *LinearCache) Close() {

}

func (c LinearCache) ReadAll() []common.Meas {
	return make([]common.Meas, 0, 0)
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
