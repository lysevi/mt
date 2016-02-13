package main

import (
	"fmt"
)

func main() {
	fmt.Println("****************")

	cblock := NewCompressedBlock()
	iterations := 2
	for i := 0; i < iterations; i++ {
		cblock.write_64(cblock.delta_64(64))
		cblock.write_256(cblock.delta_256(256))

	}

	cblock.byteNum = 0
	cblock.bitNum = MAX_BIT
	fmt.Println(">>>> 64: ", cblock.readTime(0))
	fmt.Println(" +> ", cblock.byteNum, cblock.bitNum)
	fmt.Println(">>>> 256: ", cblock.readTime(0))
	fmt.Println(" +> ", cblock.byteNum, cblock.bitNum)
	fmt.Println("***************")
	for i := 0; i < iterations; i++ {
		fmt.Println("!+> ", cblock.byteNum, cblock.bitNum)
		t_1 := cblock.readTime(0)
		if t_1 != 64 {
			fmt.Print("d64 read error i:", i, t_1)
			fmt.Print(cblock.String())
			return
		}

		t_256 := cblock.readTime(0)

		if t_256 != 256 {
			fmt.Print("d256 read error i:", i, t_256)
			fmt.Print(cblock.String())
			return
		}
	}
}
