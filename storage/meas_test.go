package storage

import (
	"fmt"
	"sort"
	"testing"
)

func TestMeas2String(t *testing.T) {
	m := NewMeas(0, 1, 2, 3)
	s := fmt.Sprintf("%v", m.String())

	if len(s) == 0 {
		t.Error("meas.String==0")
	}
}

func TestMeasEqual(t *testing.T) {
	m1 := NewMeas(0, 1, 2, 3)
	m2 := NewMeas(0, 1, 2, 2)

	if !measEqual(m1, m1) {
		t.Error("not equal", m1)
	}
	if measEqual(m1, m2) {
		t.Error(" equal", m1, m2)
	}
}

func TestMeasSort(t *testing.T) {
	var m MeasByTime

	m = append(m, NewMeas(0, 0, 2, 3))
	m = append(m, NewMeas(0, 1, 2, 3))
	m = append(m, NewMeas(0, 2, 2, 3))
	sort.Sort(m)

	for i, v := range m {
		if Time(i) != v.Tstamp {
			t.Error("sort by time error ", m)
			break
		}
	}
}
