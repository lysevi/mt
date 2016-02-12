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
