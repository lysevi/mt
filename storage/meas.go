package storage

import (
	"fmt"
)

type Id int32
type Flag uint64
type Time uint64

type Meas struct {
	Id     Id
	Tstamp Time
	Value  int64
	Flg    Flag
}

type MeasByTime []Meas

func (m MeasByTime) Len() int {
	return len(m)
}

func (m MeasByTime) Less(i, j int) bool {
	return m[i].Tstamp < m[j].Tstamp
}

func (m MeasByTime) Swap(i, j int) {
	t := m[i]
	m[i] = m[j]
	m[j] = t
}

func NewMeas(Id Id, Tstamp Time, Value int64, Flg Flag) Meas {
	m := Meas{}
	m.Id = Id
	m.Flg = Flg
	m.Tstamp = Tstamp
	m.Value = Value
	return m
}

func (m *Meas) String() string {
	return fmt.Sprintf("{i:%v, t:%v, f:%v, v:%v}", m.Id, m.Tstamp, m.Flg, m.Value)
}

func measEqual(m1, m2 Meas) bool {
	return m1.Id == m2.Id && m1.Tstamp == m2.Tstamp && m1.Value == m2.Value && m1.Flg == m2.Flg
}
