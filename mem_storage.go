package main

type MemoryStorage struct {
	max_time Time
}

func NewMemoryStorage(sz int64) *MemoryStorage {
	res := &MemoryStorage{}
	return res
}

func (c *MemoryStorage) Add(m Meas) bool {
	return false
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
	return 0
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
	res := make([]Meas, 0, 0)
	return res
}
func (c *MemoryStorage) TimePoint(ids []Id, time Time) []Meas {
	return c.Read(ids, 0, time)
}
func (c *MemoryStorage) TimePointFltr(ids []Id, flg Flag, time Time) []Meas {
	return c.ReadFltr(ids, flg, 0, time)
}
