package a

import "testing"

func TestF2(t *testing.T) {
	if f2("Drum") != "f2Drum" {
		t.FailNow()
	}
}