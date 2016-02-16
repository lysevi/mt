package main

import (
	"fmt"
	"math"
	"time"

	//"github.com/pkg/profile"
)

func main() {
	//defer profile.Start().Stop()
	fmt.Println("****************")
	storage := NewMemoryStorage(10000000)
	iterations := 10000
	tm := Time(1)
	meases := []Meas{}
	for i := 0; i < iterations; i++ {
		m := NewMeas(1, tm, int64(math.Sin(float64(i))), Flag(0x002202))
		tm += Time(1000)
		meases = append(meases, m)
	}
	fmt.Println("add:")
	startTime := time.Now()
	storage.Add_range(meases)
	endTime := time.Now()
	fmt.Println("elapsed: ", endTime.Sub(startTime))
}
