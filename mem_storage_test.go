package main

import (
	"fmt"
	"sync"
	"testing"
)

var _ = fmt.Sprintf("")

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
	f_reader := func(count int, storage *MemoryStorage, wg *sync.WaitGroup, stop *bool) {
		for {
			_ = storage.ReadAll()
			if *stop {
				break
			}
		}
	}
	var stop bool = false
	const write_count = 100
	wg := sync.WaitGroup{}
	go f_reader(100, lc, &wg, &stop)
	wg.Add(1)
	go f(1, write_count, lc, &wg)
	wg.Add(1)
	go f_range(6, write_count, lc, &wg)

	wg.Wait()
	stop = true
	all := lc.ReadAll()

	if len(lc.cblocks) != 2 {
		t.Error("len(lc.cblocks)!=2", len(lc.cblocks))
	}

	if len(all) != 2*write_count {
		t.Error(len(all))
	}
}

func TestMemoryStorageArchive(t *testing.T) {
	lc := NewMemoryStorage(1000)
	i := 1
	for {
		m := NewMeas(1, Time(i), int64(i), Flag(i))
		lc.Add(m)
		if len(lc.archive) != 0 {
			break
		}
		i++
	}
	all := lc.ReadAll()
	if len(all) != i {
		t.Error(len(all), i)
	}
}
