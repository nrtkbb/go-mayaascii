package parser

import (
	"testing"

	"github.com/nrtkbb/go-mayaascii/cmd"
)

func TestMakeRename_Min(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`rename -uid "B0F0F886-4CD8-DA88-AC80-C1B83173300D";`)
	r := MakeRename(c.Parse())
	msg := `got Rename %s "%s", wont "%s"`
	if r.To == nil {
		t.Errorf(msg, "To", r.To, "(*string=0xFFFFFFFFF)")
	}
	if *r.To != "B0F0F886-4CD8-DA88-AC80-C1B83173300D" {
		t.Errorf(msg, "*To", *r.To, "B0F0F886-4CD8-DA88-AC80-C1B83173300D")
	}
	if r.From != nil {
		t.Errorf(msg, "From", r.From, nil)
	}
	if !r.UUID {
		t.Errorf(msg, "UUID", r.UUID, true)
	}
	if r.IgnoreShape {
		t.Errorf(msg, "IgnoreShape", r.IgnoreShape, false)
	}
}

func TestMakeRename_Max(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`rename -is "from" "to";`)
	r := MakeRename(c.Parse())
	msg := `got Rename %s "%s", wont "%s"`
	if r.To == nil {
		t.Errorf(msg, "To", r.To, "(*string=0xfffffffff)")
	}
	if *r.To != "to" {
		t.Errorf(msg, "*To", *r.To, "to")
	}
	if r.From == nil {
		t.Errorf(msg, "From", r.From, "(*string=0xfffffffff)")
	}
	if *r.From != "from" {
		t.Errorf(msg, "*From", *r.From, "from")
	}
	if r.UUID {
		t.Errorf(msg, "UUID", r.UUID, false)
	}
	if !r.IgnoreShape {
		t.Errorf(msg, "IgnoreShape", r.IgnoreShape, true)
	}
}
