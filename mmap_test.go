package main

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

var testData = []byte("0000000000000000000000000000000000000000000000000000000000000000")
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

func getPtr(m []byte, num uint64, szOfElement uint64) *tstHeader {
	offset := num * szOfElement
	return (*tstHeader)(unsafe.Pointer(&m[offset]))
}

func TestMap(t *testing.T) {
	const magicOne int = 1171
	const magicTwo int64 = 1271
	const szOfElement uint64 = uint64(unsafe.Sizeof(tstHeader{}))
	{
		mmaped := openAndMap(t)
		for i := uint64(0); i < 2; i++ {
			hdr := getPtr(mmaped, i, szOfElement)
			hdr.first = magicOne
			hdr.second = magicTwo
		}

		mmaped.Unmap()
		f.Close()
	}
	{
		mmaped := openAndMap(t)
		for i := uint64(0); i < 2; i++ {
			hdr := getPtr(mmaped, i, szOfElement)

			if hdr.first != magicOne {
				t.Error("hdr.first != 1171", hdr.first)
			}

			if hdr.second != magicTwo {
				t.Error("hdr.second = 1271", hdr.second)
			}

		}
		mmaped.Unmap()
		f.Close()
	}
}
