package mt

import (
	"testing"
)

func TestMeasInTimeInterval(t *testing.T) {
	if !inTimeInterval(0, 10, 1) {
		t.Error("!inTimeInterval(0,10,1)")
	}

	if inTimeInterval(0, 10, 11) {
		t.Error("inTimeInterval(0,10,11)")
	}
}

func TestMeasIdFltr(t *testing.T) {
	var ids []Id

	if !idFltr(ids, 1) {
		t.Error("!idFltr(ids,1)")
	}
	ids = append(ids, 1)

	if idFltr(ids, 2) {
		t.Error("idFltr(ids,2)")
	}

	if !idFltr(ids, 1) {
		t.Error("!idFltr(ids,1)")
	}
}

func TestMeasFlagFltr(t *testing.T) {
	if !flagFltr(0, 1) {
		t.Error("!flagFltr(0,1)")
	}

	if flagFltr(2, 1) {
		t.Error("flagFltr(2,1)")
	}

	if !flagFltr(0, 1) {
		t.Error("flagFltr(0,1)")
	}

	if !flagFltr(1, 1) {
		t.Error("flagFltr(1,1)")
	}
}
