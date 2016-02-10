package storages

import (
	"testing"
)

func TestCacheAdd(t *testing.T) {
	lc := NewLinearCache()
	checkWriterAdd(t, lc)
}

func TestCacheAddRange(t *testing.T) {
	lc := NewLinearCache()
	checkWriterAddRange(t, lc)
}

func TestCacheCheck(t *testing.T) {
	lc := NewLinearCache()
	checkStorage(t, lc, 0, 100, 5)
}
