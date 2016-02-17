package main

import (
	"fmt"
	"testing"
	"time"
)

var _ = fmt.Sprintf("")

const OneSecond = time.Duration(1) * time.Second

func TestIntegrationCompressedBlock(t *testing.T) {
	cblock := NewCompressedBlock()
	iterations := 5000
	tm := Time(1)
	meases := []Meas{}
	for i := 0; i < iterations; i++ {
		m := NewMeas(1, tm, int64(i), Flag(0x002202))
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

func TestIntegrationMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage(10000000)
	iterations := 250000
	tm := Time(1)
	var Val int64 = 1

	startTime := time.Now()
	for i := 0; i < iterations; i++ {
		storage.Add(NewMeas(1, tm+Time(10), Val, Flag(0x002202)))
		storage.Add(NewMeas(1, tm+Time(250), Val, Flag(0x1)))
		storage.Add(NewMeas(1, tm+Time(1000), Val, Flag(0x002202)))
		storage.Add(NewMeas(1, tm+Time(4000), Val, Flag(0x002202)))
		tm = tm + Time(4000)
		Val *= 2
	}
	elapsed := time.Since(startTime)

	if elapsed > (OneSecond) {
		t.Error("so slow: ", elapsed)
	}
}

func TestIntegrationStorage(t *testing.T) {
	storage := NewStorage()
	iterations := 250000
	tm := Time(1)
	var Val int64 = 1

	startTime := time.Now()
	for i := 0; i < iterations; i++ {
		storage.Add(NewMeas(1, tm+Time(10), Val, Flag(0x002202)))
		storage.Add(NewMeas(2, tm+Time(250), Val, Flag(0x1)))
		storage.Add(NewMeas(3, tm+Time(1000), Val, Flag(0x002202)))
		storage.Add(NewMeas(4, tm+Time(4000), Val, Flag(0x002202)))
		tm = tm + Time(4000)
		Val *= 2
	}
	storage.Close()
	elapsed := time.Since(startTime)
	if elapsed > (OneSecond) {
		t.Error("so slow: ", elapsed)
	}
}
