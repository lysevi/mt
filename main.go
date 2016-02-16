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
	iterations := 100000
	tm := Time(1)
	var m Meas
	fmt.Println("add:")
	startTime := time.Now()
	for i := 0; i < iterations; i++ {
		m = NewMeas(1, tm, int64(i), Flag(0x002202))
		tm += Time(10)
		storage.Add(m)
	}
	endTime := time.Now()
	fmt.Println("elapsed: ", endTime.Sub(startTime))
}
