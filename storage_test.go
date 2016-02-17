package main

import (
	"fmt"
	"testing"
)

var _ = fmt.Sprintf("")

func TestStorageAddSingle(t *testing.T) {
	//	storage := NewStorage()

}

func TestStorageAddRange(t *testing.T) {
	lc := NewStorage()
	checkWriterAddRange(t, lc)
	lc.Close()
}

func TestStorageAddRange_s(t *testing.T) {
	lc := NewStorage()
	checkStorageAddRange(t, lc)
	lc.Close()
}

func TestStorageCheck(t *testing.T) {
	lc := NewStorage()
	checkStorage(t, lc, 0, 100, 5)
	lc.Close()
}

func TestMemoryStorageCacheSync(t *testing.T) {
	lc := NewStorage()
	writes_count := CACHE_DEFAULT_SIZE * 2
	for i := 0; i < writes_count; i++ {
		m := NewMeas(1, Time(i), int64(i), Flag(i))
		lc.Add(m)
	}

	for {
		if lc.sync_complete {
			break
		}
	}
	lc.Close()
}
