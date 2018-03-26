package main

import (
	"testing"
)

func TestMakeRequires_Min(t *testing.T) {
	cb := &CmdBuilder{}
	cb.Append(`requires maya "2016";`)
	cmd := cb.Parse()
	r := MakeRequires(cmd)
	if r.PluginName != "maya" {
		t.Fatalf("got %v, wont %v", r.PluginName, "maya")
	}
	if r.Version != "2016" {
		t.Fatalf("got %v, wont %v", r.Version, "2016")
	}
	if len(r.NodeTypes) != 0 {
		t.Fatalf("got %v, wont %v", len(r.NodeTypes), 0)
	}
	if len(r.DataTypes) != 0 {
		t.Fatalf("got %v, wont %v", len(r.DataTypes), 0)
	}
}

func TestMakeRequires_Max(t *testing.T) {
	cb := &CmdBuilder{}
	cb.Append(`requires -nodeType "typeName1"
		-dataType "typeName2" "pluginName" "version";`)
	cmd := cb.Parse()
	r := MakeRequires(cmd)
	if r.PluginName != "pluginName" {
		t.Fatalf("got %v, wont %v", r.PluginName, "pluginName")
	}
	if r.Version != "version" {
		t.Fatalf("got %v, wont %v", r.Version, "version")
	}
	if len(r.NodeTypes) != 1 {
		t.Fatalf("got %v, wont %v", len(r.NodeTypes), 1)
	}
	if r.NodeTypes[0] != "typeName1" {
		t.Fatalf("got %v, wont %v", r.NodeTypes[0], "typeName1")
	}
	if len(r.DataTypes) != 1 {
		t.Fatalf("got %v, wont %v", len(r.DataTypes), 1)
	}
	if r.DataTypes[0] != "typeName2" {
		t.Fatalf("got %v, wont %v", r.DataTypes[0], "typeName2")
	}
}
