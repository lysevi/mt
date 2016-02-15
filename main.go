package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("****************")
	cblock := NewCompressedBlock()

	iterations := 1000
	for i := 0; i < iterations; i++ {
		t := Time(time.Millisecond * time.Duration(i))
		m := NewMeas(1, t, int64(i), Flag(i%2))
		cblock.Add(m)
	}

	fmt.Println("used bytes: ", cblock.byteNum)
	fmt.Println("full size: ", 28*iterations)
}
