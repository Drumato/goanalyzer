package a_test

import (
	"testing"
	"undefined_unit_test/a"
)

func TestF2(t *testing.T) {
	if a.F2("Drum") != "f2Drum" {
		t.FailNow()
	}
}