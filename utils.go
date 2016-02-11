package mt

import (
	"unsafe"
)

func idFltr(ids []Id, idV Id) bool {
	if len(ids) == 0 {
		return true
	} else {
		for _, v := range ids {
			if v == idV {
				return true
			}
		}
	}
	return false
}

func flagFltr(fltrFlag, measFlag Flag) bool {
	if fltrFlag == Flag(0) {
		return true
	}
	return fltrFlag == measFlag
}

func inTimeInterval(from, to, tstamp Time) bool {
	// [from tstamp to]
	if from <= tstamp && to >= tstamp {
		return true
	} else {
		return false
	}
}

func getPtr(m []byte, num uint64, szOfElement uint64) *tstHeader {
	offset := num * szOfElement
	return (*tstHeader)(unsafe.Pointer(&m[offset]))
}
