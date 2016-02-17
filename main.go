package main

import (
	"fmt"
	"time"

	"github.com/lysevi/mt/storage"
	//"github.com/pkg/profile"
)

func main() {
	//defer profile.Start().Stop()
	fmt.Println("****************")
	stor := storage.NewStorage()
	iterations := 250000
	tm := storage.Time(1)
	var Val int64 = 1

	startTime := time.Now()
	for i := 0; i < iterations; i++ {
		stor.Add(storage.NewMeas(1, tm+storage.Time(10), Val, storage.Flag(0x002202)))
		stor.Add(storage.NewMeas(2, tm+storage.Time(250), Val, storage.Flag(0x1)))
		stor.Add(storage.NewMeas(3, tm+storage.Time(1000), Val, storage.Flag(0x002202)))
		stor.Add(storage.NewMeas(4, tm+storage.Time(4000), Val, storage.Flag(0x002202)))
		tm = tm + storage.Time(4000)
		Val *= 2
	}
	stor.WaitSync()
	elapsed := time.Since(startTime)
	fmt.Println("elapsed: ", elapsed)
	stor.Close()
}
