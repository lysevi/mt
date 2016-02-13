package main

import (
	"testing"
)

func BenchmarkTimeDeltas64(b *testing.B) {
	cblock := NewCompressedBlock()
	for i := 0; i < b.N; i++ {
		cblock.delta_64(1)
		cblock.delta_64(64)
		cblock.delta_64(63)
	}
}

func BenchmarkTimeDeltas256(b *testing.B) {
	cblock := NewCompressedBlock()
	for i := 0; i < b.N; i++ {
		cblock.delta_256(256)
		cblock.delta_256(255)
		cblock.delta_256(65)
	}
}

func BenchmarkTimeDeltas2048(b *testing.B) {
	cblock := NewCompressedBlock()
	for i := 0; i < b.N; i++ {
		cblock.delta_2048(2048)
		cblock.delta_2048(257)
		cblock.delta_2048(4095)
	}
}

func BenchmarkTimeDeltas_big(b *testing.B) {
	cblock := NewCompressedBlock()
	for i := 0; i < b.N; i++ {
		cblock.delta_big(2049)
		cblock.delta_big(65535)
		cblock.delta_big(4095)
		cblock.delta_big(4294967295)
	}
}

func BenchmarkTimeWrite(b *testing.B) {
	cblock := NewCompressedBlock()
	for i := 0; i < b.N; i++ {
		cblock.write_64(cblock.delta_64(1))
		cblock.write_64(cblock.delta_64(64))
		cblock.write_64(cblock.delta_64(64))

		cblock.write_256(cblock.delta_256(256))
		cblock.write_256(cblock.delta_256(255))
		cblock.write_256(cblock.delta_256(65))

		cblock.write_2048(cblock.delta_2048(2048))
		cblock.write_2048(cblock.delta_2048(257))
		cblock.write_2048(cblock.delta_2048(4095))

		cblock.write_big(cblock.delta_big(2049))
		cblock.write_big(cblock.delta_big(65535))
		cblock.write_big(cblock.delta_big(4095))
	}
}

func BenchmarkTimeRW(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cblock := NewCompressedBlock()
		iterations := 1
		for i := 0; i < iterations; i++ {
			cblock.write_64(cblock.delta_64(1))
			cblock.write_64(cblock.delta_64(64))
			cblock.write_64(cblock.delta_64(64))

			cblock.write_256(cblock.delta_256(256))
			cblock.write_256(cblock.delta_256(255))
			cblock.write_256(cblock.delta_256(65))

			cblock.write_2048(cblock.delta_2048(2048))
			cblock.write_2048(cblock.delta_2048(257))
			cblock.write_2048(cblock.delta_2048(4095))

			cblock.write_big(cblock.delta_big(2049))
			cblock.write_big(cblock.delta_big(65535))
			cblock.write_big(cblock.delta_big(4095))
		}

		cblock.byteNum = 0
		cblock.bitNum = MAX_BIT

		for i := 0; i < iterations; i++ {
			cblock.readTime(0)
			cblock.readTime(0)
			cblock.readTime(0)

			cblock.readTime(0)
			cblock.readTime(0)
			cblock.readTime(0)

			cblock.readTime(0)
			cblock.readTime(0)
			cblock.readTime(0)

			cblock.readTime(0)
			cblock.readTime(0)
			cblock.readTime(0)

		}
	}
}
