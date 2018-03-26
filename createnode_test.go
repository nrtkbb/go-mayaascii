package main

import "testing"

func TestMakeCreateNode_Min(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`createNode transform -n "nodeName";`)
	cn := MakeCreateNode(c.Parse())
	msg := `got CreateNode %s "%s", wont "%s"`
	if cn.NodeType != "transform" {
		t.Errorf(msg, "NodeType", cn.NodeName, "transform")
	}
	if cn.NodeName != "nodeName" {
		t.Errorf(msg, "NodeName", cn.NodeName, "nodeName")
	}
	if cn.Parent != nil {
		t.Errorf(msg, "Parent", cn.Parent, nil)
	}
	if cn.Shared {
		t.Errorf(msg, "Shared", cn.Shared, false)
	}
	if cn.SkipSelect {
		t.Errorf(msg, "SkipSelect", cn.SkipSelect, false)
	}
}

func TestMakeCreateNode_Max(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`createNode camera -s -n "ns:grp|ns:cam|ns:camShape" -p "ns:grp|ns:cam";`)
	cn := MakeCreateNode(c.Parse())
	msg := `got CreateNode %s "%s", wont "%s"`
	if cn.NodeType != "camera" {
		t.Errorf(msg, "NodeType", cn.NodeName, "camera")
	}
	if cn.NodeName != "ns:grp|ns:cam|ns:camShape" {
		t.Errorf(msg, "NodeName", cn.NodeName, "ns:grp|ns:cam|ns:camShape")
	}
	if cn.Parent == nil {
		t.Errorf(msg, "Parent", *cn.Parent, "ns:grp|ns:cam")
	}
	if *cn.Parent != "ns:grp|ns:cam" {
		t.Errorf(msg, "Parent", *cn.Parent, "ns:grp|ns:cam")
	}
	if !cn.Shared {
		t.Errorf(msg, "Shared", cn.Shared, true)
	}
	if cn.SkipSelect {
		t.Errorf(msg, "SkipSelect", cn.SkipSelect, false)
	}
}
