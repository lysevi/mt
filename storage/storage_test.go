package storage

import (
	"fmt"
	"testing"
)

var _ = fmt.Sprintf("")

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

func TestStorageCacheSync(t *testing.T) {
	lc := NewStorage()
	writes_count := defaultCacheSize * 3
	writes := 0
	for i := 0; i < writes_count; i++ {
		m := NewMeas(1, Time(i), int64(i), Flag(i))
		lc.Add(m)
		writes++
		if i == writes_count/2 {
			all := lc.ReadAll()
			if len(all) != writes {
				t.Error("storage readall error: ", len(all), writes)
			}
		}
	}

	all := lc.ReadAll()
	if len(all) != writes {
		t.Error("storage readall error: ", len(all), writes)
	}
	lc.Close()
}
