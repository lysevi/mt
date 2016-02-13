package main

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

func setBit(v uint8, bitNum uint8, bitValue uint8) byte {
	if bitValue == 1 {
		return v | (1 << bitNum)
	} else {
		v &^= (1 << bitNum)
		return v
	}
}

func getBit(v byte, bitNum uint8) uint8 {
	return ((v >> bitNum) & 1)
}

func checkBit(v byte, bitNum uint8) bool {
	return (getBit(v, bitNum)) == 1
}

func setBit16(v uint16, bitNum uint8, bitValue uint8) uint16 {
	if bitValue == 1 {
		return v | (1 << bitNum)
	} else {
		v &^= (1 << bitNum)
		return v
	}
}

func getBit16(v uint16, bitNum uint8) uint16 {
	return ((v >> bitNum) & 1)
}

func checkBit16(v uint16, bitNum uint8) bool {
	return (getBit16(v, bitNum)) == 1
}

func setBit32(v uint32, bitNum uint8, bitValue uint8) uint32 {
	if bitValue == 1 {
		return v | (1 << bitNum)
	} else {
		v &^= (1 << bitNum)
		return v
	}
}

func getBit32(v uint32, bitNum uint8) uint32 {
	return ((v >> bitNum) & 1)
}

func checkBit32(v uint32, bitNum uint8) bool {
	return (getBit32(v, bitNum)) == 1
}
