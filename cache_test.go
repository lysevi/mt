package mt

import (
	"testing"
)

func TestCacheAdd(t *testing.T) {
	lc := NewLinearCache(100)
	checkWriterAdd(t, lc)
}

func TestCacheAddRange(t *testing.T) {
	lc := NewLinearCache(200)
	checkWriterAddRange(t, lc)
}

func TestCacheCheck(t *testing.T) {
	lc := NewLinearCache(1000)
	checkStorage(t, lc, 0, 100, 5)
}
