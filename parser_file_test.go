package mayaascii

import (
	"testing"
)

func TestMakeFile_ReferenceDepthInfo(t *testing.T) {
	c := &CmdBuilder{}
	fileLine := `file -rdi 1 -ns "baseA" -rfn "baseARN" -op "v=0;" -typ "mayaAscii" "C:/baseA.ma";`
	c.Append(fileLine)
	f := ParseFile(c.Parse())
	msg := `got FileCmd %v "%v", wont "%v"`
	if f.ReferenceDepthInfo != 1 {
		t.Errorf(msg, "ReferenceDepthInfo", f.ReferenceDepthInfo, 1)
	}
	if f.Namespace != "baseA" {
		t.Errorf(msg, "NodeName", f.Namespace, "baseA")
	}
	if f.ReferenceNode != "baseARN" {
		t.Errorf(msg, "ReferenceNode", f.ReferenceNode, "baseARN")
	}
	if f.Options != "v=0;" {
		t.Errorf(msg, "Options", f.Options, "v=0;")
	}
	if f.Type != "mayaAscii" {
		t.Errorf(msg, "Type", f.Type, "mayaAscii")
	}
	if f.Path != "C:/baseA.ma" {
		t.Errorf(msg, "Path", f.Path, "C:/baseA.ma")
	}
	if f.String() != fileLine {
		t.Errorf(msg, "f.String()", f.String(), fileLine)
	}
}

func TestMakeFile_Reference(t *testing.T) {
	c := &CmdBuilder{}
	fileLine := `file -r -ns "baseA" -dr 1 -rfn "baseARN" -op "v=0;" -typ "mayaAscii" "C:/baseA.ma";`
	c.Append(fileLine)
	f := ParseFile(c.Parse())
	msg := `got FileCmd %v "%v", wont "%v"`
	if f.Reference != true {
		t.Errorf(msg, "Reference", f.Reference, true)
	}
	if f.ReferenceDepthInfo != 0 {
		t.Errorf(msg, "ReferenceDepthInfo", f.ReferenceDepthInfo, 0)
	}
	if f.Namespace != "baseA" {
		t.Errorf(msg, "NodeName", f.Namespace, "baseA")
	}
	if f.ReferenceNode != "baseARN" {
		t.Errorf(msg, "ReferenceNode", f.ReferenceNode, "baseARN")
	}
	if f.Options != "v=0;" {
		t.Errorf(msg, "Options", f.Options, "v=0;")
	}
	if f.Type != "mayaAscii" {
		t.Errorf(msg, "Type", f.Type, "mayaAscii")
	}
	if f.Path != "C:/baseA.ma" {
		t.Errorf(msg, "Path", f.Path, "C:/baseA.ma")
	}
	if f.String() != fileLine {
		t.Errorf(msg, "f.String()", f.String(), fileLine)
	}
}
