package main

import (
	"sync"
	"testing"
)

func TestMemoryStorageAddSingle(t *testing.T) {
	lc := NewMemoryStorage(200)
	lc.Add(NewMeas(11, 3, 3, 2))
	if len(lc.cblocks) != 1 {
		t.Error("cblock len error: ", len(lc.cblocks))
	}
}

func TestMemoryStorageAddRange(t *testing.T) {
	lc := NewMemoryStorage(200)
	checkWriterAddRange(t, lc)
}

func TestMemoryStorageAddRange_s(t *testing.T) {
	lc := NewMemoryStorage(200)
	checkStorageAddRange(t, lc)
}

func TestMemoryStorageCheck(t *testing.T) {
	lc := NewMemoryStorage(1000)
	checkStorage(t, lc, 0, 100, 5)
}

func TestMemoryStorageThreads(t *testing.T) {
	lc := NewMemoryStorage(1000)

	f := func(id Id, count int, storage *MemoryStorage, wg *sync.WaitGroup) {

		defer wg.Done()
		var t Time = 1
		for i := 0; i < count; i++ {
			m := NewMeas(id, t, int64(i), Flag(i))
			storage.Add(m)
		}
	}
	f_range := func(id Id, count int, storage *MemoryStorage, wg *sync.WaitGroup) {

		defer wg.Done()
		var t Time = 1
		meases := []Meas{}
		for i := 0; i < count; i++ {
			m := NewMeas(id, t, int64(i), Flag(i))
			meases = append(meases, m)
		}
		storage.Add_range(meases)

	}
	const write_count = 1000
	wg := sync.WaitGroup{}
	wg.Add(1)
	go f(1, write_count, lc, &wg)
	wg.Add(1)
	go f_range(6, write_count, lc, &wg)

	wg.Wait()

	all := lc.ReadAll()

	if len(lc.cblocks) != 2 {
		t.Error("len(lc.cblocks)!=2", len(lc.cblocks))
	}

	if len(all) != 2*write_count {
		t.Error(len(all))
	}
}
