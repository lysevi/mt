package main

import (
	"testing"
)

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
