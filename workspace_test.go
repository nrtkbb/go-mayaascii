package mayaascii

import (
	"testing"
)

func TestMakeWorkspace(t *testing.T) {
	cb := &CmdBuilder{}
	cb.Append(`workspace -fr "sourceImages" "sourceimages";`)
	cmd := cb.Parse()
	w := MakeWorkspace(cmd)
	if w.FileRule != "sourceImages" {
		t.Errorf("got %v, wont %v", w.FileRule, "sourceImages")
	}
	if w.Place != "sourceimages" {
		t.Errorf("got %v, wont %v", w.Place, "sourceimages")
	}
}
