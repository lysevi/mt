package mt

import (
	"testing"
)

func checkWriterAdd(t *testing.T, writer MeasWriter) {
	cap_start := writer.Cap()
	m := Meas{}
	if !writer.Add(m) {
		t.Error("!writer.Add(m)")
	}

	cap_start -= 1
	if writer.Cap() != cap_start {
		t.Error("writer.Cap()!=cap_start ", writer.Cap(), cap_start)
	}

	var i int64 = 0
	for ; i < cap_start; i++ {
		if !writer.Add(m) {
			t.Errorf("add %d value error", i)
		}
	}

	if i == 0 {
		t.Error("test logic error")
	}

	if !writer.IsFull() {
		t.Error("not full?!")
	}
}

func checkWriterAddRange(t *testing.T, writer MeasWriter) {
	size := writer.Cap()
	twice := size * 2

	m := make([]Meas, twice, twice)

	writed := writer.Add_range(m)
	if writed != size {
		t.Errorf("writed!=size, %d %d", writed, size)
	}
}

func checkStorageAddRange(t *testing.T, storage MeasWriter) {
	size := storage.Cap()
	twice := size * 2

	m := make([]Meas, twice, twice)

	writed := storage.Add_range(m)
	if writed <= size {
		t.Errorf("writed<=size, %d %d", writed, size)
	}
}

func checkStorage(t *testing.T, storage MeasStorage, from, to, step Time) {
	checkAll := func(res []Meas, msg string) {
		i := from
		for _, m := range res {
			if m.Id != Id(i) || m.Flg != Flag(i) || m.Tstamp != Time(i) {
				t.Errorf("msg: ", m)
			}
			i += step
		}
	}

	m := Meas{}
	total_count := 0
	for i := from; i < to; i += step {
		m.Id = Id(i)
		m.Flg = Flag(i)
		m.Tstamp = Time(i)
		storage.Add(m)
		total_count++
	}

	all := storage.ReadAll()
	if len(all) != total_count {
		t.Error("len(all)!=total_count", len(all), total_count)
	}

	checkAll(all, "readAll error: ")

	var ids []Id
	all = storage.Read(ids, from, to)
	if len(all) != total_count {
		t.Error("len(all)!=total_count", len(all), total_count)
	}

	checkAll(all, "read error: ")

	ids = append(ids, Id(from+step))
	fltr_res := storage.ReadFltr(ids, 0, from, to)
	if len(fltr_res) != 1 {
		t.Error("len(fltr_res)!=1", len(fltr_res))
	} else {
		if fltr_res[0].Id != ids[0] {
			t.Error("ReadFltr: ", fltr_res[0])
		}
	}

	fltr_res = storage.ReadFltr(ids, Flag(to+1), from, to)
	if len(fltr_res) != 0 {
		t.Error("len(fltr_res)!=0", len(fltr_res))
	}

	var empty_id []Id
	all = storage.TimePoint(empty_id, to)
	if len(all) != total_count {
		t.Error("timepoint: len(all)!=total_count", len(all), total_count)
	}

	checkAll(all, "TimePoint error: ")

	var emptyIDs []Id
	fltr_res = storage.TimePointFltr(emptyIDs, 0, to)
	if len(fltr_res) != total_count {
		t.Error("len(fltr_res)!=total_count", len(fltr_res))
	}

	checkAll(all, "TimePointFltr error: ")

	magicFlag := Flag(from + step)
	fltr_res = storage.TimePointFltr(emptyIDs, magicFlag, to)
	if len(fltr_res) != 1 {
		t.Error("len(fltr_res)!=1", len(fltr_res))
	}

	if fltr_res[0].Flg != magicFlag {
		t.Error("TimePointFltr: ", fltr_res)
	}
}
