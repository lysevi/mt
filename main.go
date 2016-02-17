package main

import (
	"fmt"
	"time"

	"github.com/pkg/profile"
)

func main() {
	defer profile.Start().Stop()
	fmt.Println("****************")
	storage := NewMemoryStorage(10000000)
	iterations := 250000
	tm := Time(1)
	var Val int64 = 1

	fmt.Println("add:")
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
	fmt.Println("elapsed: ", elapsed)
}
