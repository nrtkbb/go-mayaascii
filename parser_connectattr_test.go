package mayaascii

import (
	"testing"
)

func TestMakeConnectAttr_Min(t *testing.T) {
	cb := &CmdBuilder{}
	cb.Append(`connectAttr "tmp_file11.msg" ":defaultTextureList1.tx" -na;`)
	c := cb.Parse()
	ca, err := MakeConnectAttr(c)
	if err != nil {
		t.Fatal(err)
	}
	msg := "got ConnectAttr %s %v, wont %v"
	if ca.SrcNode != "tmp_file11" {
		t.Errorf(msg, "SrcNode", ca.SrcNode, "tmp_file11")
	}
	if ca.SrcAttr != "msg" {
		t.Errorf(msg, "SrcAttr", ca.SrcAttr, "msg")
	}
	if ca.DstNode != ":defaultTextureList1" {
		t.Errorf(msg, "DstNode", ca.DstNode, ":defaultTextureList1")
	}
	if ca.DstAttr != "tx" {
		t.Errorf(msg, "DstAttr", ca.DstAttr, "tx")
	}
	if !ca.NextAvailable {
		t.Errorf(msg, "NextAvailable", ca.NextAvailable, true)
	}
	if ca.Force {
		t.Errorf(msg, "Force", ca.Force, false)
	}
	if ca.Lock != nil {
		t.Errorf(msg, "Lock", ca.Lock, nil)
	}
	if ca.ReferenceDest != nil {
		t.Errorf(msg, "ReferenceDest", ca.ReferenceDest, nil)
	}
}

func TestMakeConnectAttr_Max(t *testing.T) {
	cb := &CmdBuilder{}
	cb.Append(`connectAttr -f -l on -rd "test" "tmp_file11.msg" ":defaultTextureList1.tx" -na;`)
	c := cb.Parse()
	ca, err := MakeConnectAttr(c)
	if err != nil {
		t.Fatal(err)
	}
	msg := "got ConnectAttr %s %s, wont %s"
	if ca.SrcNode != "tmp_file11" {
		t.Errorf(msg, "SrcNode", ca.SrcNode, "tmp_file11")
	}
	if ca.SrcAttr != "msg" {
		t.Errorf(msg, "SrcAttr", ca.SrcAttr, "msg")
	}
	if ca.DstNode != ":defaultTextureList1" {
		t.Errorf(msg, "DstNode", ca.DstNode, ":defaultTextureList1")
	}
	if ca.DstAttr != "tx" {
		t.Errorf(msg, "DstAttr", ca.DstAttr, "tx")
	}
	if !ca.NextAvailable {
		t.Errorf(msg, "NextAvailable", ca.NextAvailable, true)
	}
	if !ca.Force {
		t.Errorf(msg, "Force", ca.Force, true)
	}
	if ca.Lock == nil {
		t.Errorf(msg, "Lock", ca.Lock, nil)
	}
	if !*ca.Lock {
		t.Errorf(msg, "Lock", *ca.Lock, true)
	}
	if ca.ReferenceDest == nil {
		t.Errorf(msg, "ReferenceDest", ca.ReferenceDest, nil)
	}
	if *ca.ReferenceDest != `"test"` {
		t.Errorf(msg, "ReferenceDest", *ca.ReferenceDest, `"test"`)
	}
}
