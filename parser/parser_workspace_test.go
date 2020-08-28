package parser

import (
	"testing"

	"github.com/nrtkbb/go-mayaascii/cmd"
)

func TestMakeWorkspace(t *testing.T) {
	cb := &cmd.CmdBuilder{}
	cb.Append(`workspace -fr "sourceImages" "sourceimages";`)
	c := cb.Parse()
	w := MakeWorkspace(c)
	if w.FileRule != "sourceImages" {
		t.Errorf("got %v, wont %v", w.FileRule, "sourceImages")
	}
	if w.Place != "sourceimages" {
		t.Errorf("got %v, wont %v", w.Place, "sourceimages")
	}
}
