package main

import (
	"fmt"
	"time"

	//"github.com/pkg/profile"
)

func main() {
	//defer profile.Start().Stop()
	fmt.Println("****************")
	storage := NewStorage()
	iterations := 500000
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
	fmt.Println("elapsed: ", elapsed)
}
