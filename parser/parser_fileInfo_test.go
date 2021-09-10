package parser

import (
	"testing"

	"github.com/nrtkbb/go-mayaascii/cmd"
)

func TestMakeFileInfo(t *testing.T) {
	c := &cmd.CmdBuilder{}
	fileInfoLine := `fileInfo "fileInfoName" "fileInfoValue";`
	c.Append(fileInfoLine)
	fi := MakeFileInfo(c.Parse())
	msg := `got FileInfo %v "%v", wont "%v"`
	if fi.Name != "fileInfoName" {
		t.Errorf(msg, "Name", fi.Name, "fileInfoName")
	}
	if fi.Value != "fileInfoValue" {
		t.Errorf(msg, "Value", fi.Value, "fileInfoValue")
	}
	if fi.String() != fileInfoLine {
		t.Errorf(msg, "String()", fi.String(), fileInfoLine)
	}
}
