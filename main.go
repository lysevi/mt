package mt

import (
	"fmt"
)

func main() {

	{ // !
		fmt.Println("****************")
		//		var t2 = Time(321)
		cblock := NewCompressedBlock()
		cblock.StartTime = 0

		cblock.compressTime(1)
		cblock.compressTime(65)
		cblock.compressTime(129)
		cblock.compressTime(193)
		cblock.compressTime(257)
		cblock.compressTime(321)
		fmt.Println("+++++++++++++++++++++++")
		cblock.compressTime(385)
		fmt.Println("-----------------------")

		fmt.Println("cblock: ", cblock.String())
		fmt.Println("pos: ", cblock.bitNum, cblock.byteNum)

		cblock.bitNum = 0
		cblock.byteNum = 0

		tm := cblock.readTime(0)
		if tm != 1 {
			fmt.Println("tm!=t1", tm, 1)
		}

		tm = cblock.readTime(tm)
		if tm != 65 {
			fmt.Println("tm!=t1", tm, 65)
		}
		tm = cblock.readTime(tm)
		if tm != 129 {
			fmt.Println("tm!=t1", tm, 129)
		}

		tm = cblock.readTime(tm)
		if tm != 193 {
			fmt.Println("tm!=t1", tm, 193)
		}
		tm = cblock.readTime(tm)
		if tm != 257 {
			fmt.Println("tm!=t1", tm, 257)
		}

		tm = cblock.readTime(tm)
		if tm != 321 {
			fmt.Println("tm!=t1", tm, 321)
		}
		fmt.Println("pos: ", cblock.bitNum, cblock.byteNum)

		fmt.Println("+++++++++++++++++++++++")
		tm = cblock.readTime(tm)
		fmt.Println("-----------------------")
		if tm != 385 {
			fmt.Println("tm!=t2", tm, 385)
		}

		fmt.Println("cblock: ", cblock.String())
	}

}
