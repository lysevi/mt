package storage

import (
	"testing"
)

func BenchmarkCompressedBlockTimeDeltas64(b *testing.B) {
	cblock := NewCompressedBlock()
	for i := 0; i < b.N; i++ {
		cblock.delta_64(63)
	}
}

func BenchmarkCompressedBlockTimeDeltas256(b *testing.B) {
	cblock := NewCompressedBlock()
	for i := 0; i < b.N; i++ {
		cblock.delta_256(255)
	}
}

func BenchmarkCompressedBlockTimeDeltas2048(b *testing.B) {
	cblock := NewCompressedBlock()
	for i := 0; i < b.N; i++ {
		cblock.delta_2048(2048)
	}
}

func BenchmarkCompressedBlockTimeDeltas_big(b *testing.B) {
	cblock := NewCompressedBlock()
	for i := 0; i < b.N; i++ {
		cblock.delta_big(4294967295)
	}
}

func BenchmarkCompressedBlockDeltaTimeWrite(b *testing.B) {
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

func BenchmarkCompressedBlockDeltaTimeRW(b *testing.B) {
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
		cblock.bitNum = maxBit

		rs := readStatus{}
		for i := 0; i < iterations; i++ {
			cblock.readTime(0, &rs)
			cblock.readTime(0, &rs)
			cblock.readTime(0, &rs)
		}
	}
}

func BenchmarkCompressedBlockTimeRW(b *testing.B) {
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

		cblock.bitNum = maxBit
		cblock.byteNum = 0

		readed_t := Time(cblock.StartTime)
		rs := readStatus{}
		for _, tm = range times {
			readed_t = cblock.readTime(readed_t, &rs)
		}
	}
}

func BenchmarkCompressedBlockValuesRW(b *testing.B) {
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

		cblock.bitNum = maxBit
		cblock.byteNum = 0

		readed_v := cblock.startValue
		rs := readStatus{}
		for _, _ = range values[1:] {
			readed_v = cblock.readValue(readed_v, &rs)
		}
	}
}
func BenchmarkCompressedBlockValuesFlags(b *testing.B) {
	for tnum := 0; tnum < b.N; tnum++ {
		cblock := NewCompressedBlock()

		flags := []Flag{}
		cblock.writeFlag(0)
		cblock.firstValue = false
		for i := Flag(1); i < 10; i++ {
			cblock.writeFlag(i)
			cblock.writeFlag(i)
			flags = append(flags, i)
		}

		cblock.bitNum = maxBit
		cblock.byteNum = 0
		readed_flag := cblock.prevFlag

		rs := readStatus{}
		for _, _ = range flags {
			readed_flag = cblock.readFlag(readed_flag, &rs)
			readed_flag = cblock.readFlag(readed_flag, &rs)
		}
	}
}

func BenchmarkCompressedBlockMeasWrite(b *testing.B) {
	for tnum := 0; tnum < b.N; tnum++ {
		cblock := NewCompressedBlock()
		iterations := 100
		for i := 0; i < iterations; i++ {
			m := NewMeas(1, Time(i), int64(i), Flag(i))
			cblock.Add(m)
		}
		_ = cblock.ReadAll()

	}
}
