package main

import (
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
