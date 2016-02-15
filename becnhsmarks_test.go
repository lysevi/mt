package main

import (
	"testing"
)

func BenchmarkTimeDeltas64(b *testing.B) {
	cblock := NewCompressedBlock()
	for i := 0; i < b.N; i++ {
		cblock.delta_64(63)
	}
}

func BenchmarkTimeDeltas256(b *testing.B) {
	cblock := NewCompressedBlock()
	for i := 0; i < b.N; i++ {
		cblock.delta_256(255)
	}
}

func BenchmarkTimeDeltas2048(b *testing.B) {
	cblock := NewCompressedBlock()
	for i := 0; i < b.N; i++ {
		cblock.delta_2048(2048)
	}
}

func BenchmarkTimeDeltas_big(b *testing.B) {
	cblock := NewCompressedBlock()
	for i := 0; i < b.N; i++ {
		cblock.delta_big(4294967295)
	}
}

func BenchmarkDeltaTimeWrite(b *testing.B) {
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

func BenchmarkDeltaTimeRW(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cblock := NewCompressedBlock()
		iterations := 1
		for i := 0; i < iterations; i++ {
			cblock.write_64(cblock.delta_64(64))

			cblock.write_256(cblock.delta_256(255))

			cblock.write_2048(cblock.delta_2048(4095))

			cblock.write_big(cblock.delta_big(65535))
		}

		cblock.byteNum = 0
		cblock.bitNum = MAX_BIT

		for i := 0; i < iterations; i++ {
			cblock.readTime(0)
			cblock.readTime(0)
			cblock.readTime(0)
		}
	}
}

func BenchmarkTimeRW(b *testing.B) {
	for tnum := 0; tnum < b.N; tnum++ {
		cblock := NewCompressedBlock()
		iterations := 1000
		tm := Time(1)
		times := []Time{}
		for i := 0; i < iterations; i++ {
			cblock.writeTime(tm)
			times = append(times, tm)
			tm *= 2
		}

		cblock.bitNum = MAX_BIT
		cblock.byteNum = 0

		readed_t := Time(cblock.StartTime)
		for _, tm = range times {
			readed_t = cblock.readTime(readed_t)
		}
	}
}

func BenchmarkValuesRW(b *testing.B) {
	for tnum := 0; tnum < b.N; tnum++ {
		cblock := NewCompressedBlock()
		values := []uint64{}
		delta := uint64(1)
		for i := uint64(0); i < 50; i++ {
			v := i * delta
			cblock.writeValue(v)
			values = append(values, v)
			delta *= 2
		}

		cblock.bitNum = MAX_BIT
		cblock.byteNum = 0

		readed_v := cblock.startValue

		for _, _ = range values[1:] {
			readed_v = cblock.readValue(readed_v)
		}
	}
}
