package main

import (
	"fmt"
	"testing"
)

func TestMeas2String(t *testing.T) {
	m := NewMeas(0, 1, 2, 3)
	s := fmt.Sprintf("%v", m)

	if len(s) == 0 {
		t.Error("meas.String==0")
	}
}
