package main

import (
	"fmt"
	"math"
	"testing"
	_ "time"
)

var _ = fmt.Sprintf("")

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
	if cblock.byteNum >= uncompressedSize/5 {
		t.Error("compression not work: ", cblock.byteNum, uncompressedSize)
	}
	//	used := cblock.byteNum
	//	unpacked := 28 * iterations
	//	fmt.Println("used bytes: ", used)
	//	fmt.Println("full size: ", unpacked)
	//	fmt.Printf("compression: %v \n", float32(unpacked)/float32(used))
	//	fmt.Println("count: ", len(meases))

	readed := cblock.ReadAll()
	for i, _ := range meases {
		if !measEqual(readed[i], meases[i]) {
			t.Error("i: ", i, readed[i].String(), meases[i].String())
		}
	}
}

//func TestIntegrationMemoryStorage(t *testing.T) {
//	storage := NewMemoryStorage(10000000)
//	iterations := 1000000
//	tm := Time(1)
//	meases := []Meas{}
//	for i := 0; i < iterations; i++ {
//		m := NewMeas(1, tm, int64(i), Flag(0x002202))
//		tm += Time(1000)
//		meases = append(meases, m)
//	}
//	fmt.Println("add:")
//	startTime := time.Now()
//	storage.Add_range(meases)
//	endTime := time.Now()
//	fmt.Println("elapsed: ", endTime.Sub(startTime))
//}
