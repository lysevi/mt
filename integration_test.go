package main

import (
	"math"
	"testing"
)

func TestIntegrationCompressedBlock(t *testing.T) {
	cblock := NewCompressedBlock()
	iterations := 5000
	tm := Time(1)
	meases := []Meas{}
	for i := 0; i < iterations; i++ {
		m := NewMeas(1, tm, int64(math.Sin(float64(i))), Flag(0x002202))
		tm += Time(1000)
		meases = append(meases, m)
	}

	cblock.Add_range(meases)

	uncompressedSize := uint64(28 * iterations)
	if cblock.byteNum >= uncompressedSize {
		t.Error("compression not work: ", cblock.byteNum, uncompressedSize)
	}
	//	fmt.Println("used bytes: ", cblock.byteNum)
	//	fmt.Println("full size: ", 28*iterations)
	//	fmt.Println("count: ", len(meases))

	readed := cblock.ReadAll()
	for i, _ := range meases {
		if !measEqual(readed[i], meases[i]) {
			t.Error("i: ", i, readed[i].String(), meases[i].String())
		}
	}
}
