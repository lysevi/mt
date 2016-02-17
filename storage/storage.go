package storage

import (
	"fmt"
	"sort"
	"sync"
)

var _ = fmt.Sprintf("")

type Storage struct {
	cache *LinearCache
	mstor *MemoryStorage

	wg            sync.WaitGroup
	stop          chan interface{}
	cache_sync    chan *LinearCache
	sync_complete bool
	lock          sync.Mutex
}

const CACHE_DEFAULT_SIZE = 100000

func NewStorage() *Storage {
	res := &Storage{}
	res.cache = NewLinearCache(CACHE_DEFAULT_SIZE)
	res.mstor = NewMemoryStorage(0)
	res.stop = make(chan interface{})
	res.cache_sync = make(chan *LinearCache)
	res.sync_complete = true
	res.wg.Add(1)
	go res.cacheSync()
	return res
}

func (c *Storage) cacheSync() {
	for {
		var ch *LinearCache = nil
		select {
		case ch = <-c.cache_sync:
			c.sync_complete = false

			all := ch.ReadAll()
			id2meases := make(map[Id]MeasByTime)
			for _, v := range all {
				items, ok := id2meases[v.Id]
				if ok {
					id2meases[v.Id] = append(items, v)
				} else {
					id2meases[v.Id] = MeasByTime{v}
				}
			}

			for _, val := range id2meases {
				sorted_vals := val[:]
				sort.Sort(sorted_vals)
				c.mstor.Add_range(sorted_vals)
			}

			c.sync_complete = true

		case <-c.stop:
			c.wg.Done()
			break
		}

	}
}

func (c *Storage) Add(m Meas) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.cache.Add(m) {
		old_cache := c.cache
		c.cache = NewLinearCache(CACHE_DEFAULT_SIZE)
		c.cache_sync <- old_cache
		c.cache.Add(m)
	}

	return true
}

func (c *Storage) Add_range(m []Meas) int64 {
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

func (c *Storage) Cap() int64 {
	return c.mstor.Cap()
}
func (c *Storage) IsFull() bool {
	return c.mstor.IsFull()
}
func (c *Storage) WaitSync() {
	for {
		if c.sync_complete {
			break
		}
	}
}
func (c *Storage) Close() {
	c.WaitSync()
	c.stop <- 1
	c.wg.Wait()
}

func (c *Storage) ReadAll() []Meas {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.WaitSync()
	from_cache := c.cache.ReadAll()
	from_stor := c.mstor.ReadAll()
	return append(from_cache, from_stor...)
}
func (c *Storage) Read(ids []Id, from, to Time) []Meas {
	return append(c.cache.ReadFltr(ids, 0, from, to), c.mstor.ReadFltr(ids, 0, from, to)...)
}
func (c *Storage) ReadFltr(ids []Id, flg Flag, from, to Time) []Meas {
	c.lock.Lock()

	defer c.lock.Unlock()

	c.WaitSync()

	from_cache := c.cache.ReadFltr(ids, flg, from, to)
	from_stor := c.mstor.ReadFltr(ids, flg, from, to)
	return append(from_cache, from_stor...)
}

func (c *Storage) TimePoint(ids []Id, time Time) []Meas {
	return c.Read(ids, 0, time)
}

func (c *Storage) TimePointFltr(ids []Id, flg Flag, time Time) []Meas {
	return c.ReadFltr(ids, flg, 0, time)
}
