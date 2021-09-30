package mayaascii

import (
	"testing"
)

func TestMakeWorkspace(t *testing.T) {
	cb := &CmdBuilder{}
	cb.Append(`workspace -fr "sourceImages" "sourceimages";`)
	c := cb.Parse()
	w := ParseWorkspace(c)
	if w.FileRule != "sourceImages" {
		t.Errorf("got %v, wont %v", w.FileRule, "sourceImages")
	}
	if w.Place != "sourceimages" {
		t.Errorf("got %v, wont %v", w.Place, "sourceimages")
	}
}
