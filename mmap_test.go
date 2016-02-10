package mt

import (
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/lysevi/mt/mmap"
)

type tstHeader struct {
	first  int
	second int64
}

var testData = []byte("0000000000000000")
var testPath = filepath.Join(os.TempDir(), "testdata")

func openFile(flags int) *os.File {
	f, err := os.OpenFile(testPath, flags, 0644)
	if err != nil {
		panic(err.Error())
	}
	return f
}

func init() {
	f := openFile(os.O_RDWR | os.O_CREATE | os.O_TRUNC)
	f.Write(testData)
	f.Close()
}

var f *os.File

func openAndMap(t *testing.T) mmap.MMap {
	f := openFile(os.O_RDWR)

	mmaped, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		t.Errorf("error mapping: %s", err)
	}

	return mmaped
}

func TestMap(t *testing.T) {
	const magicOne int = 1171
	const magicTwo int64 = 1271
	{
		mmaped := openAndMap(t)
		hdr := (*tstHeader)(unsafe.Pointer(&mmaped[0]))
		hdr.first = magicOne
		hdr.second = magicTwo
		mmaped.Unmap()
		f.Close()
	}
	{
		mmaped := openAndMap(t)
		hdr := (*tstHeader)(unsafe.Pointer(&mmaped[0]))
		if hdr.first != magicOne {
			t.Errorf("hdr.first != 1171")
		}

		if hdr.second != magicTwo {
			t.Errorf("hdr.second = 1271")
		}
		mmaped.Unmap()
		f.Close()
	}
}
