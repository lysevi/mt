package main

import (
	"fmt"
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
