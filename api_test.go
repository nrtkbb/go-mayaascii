package mayaascii

import (
	"strings"
	"testing"
)

// TODO: TestApi_File, TestApi_FileInfo, TestApi_Select

func TestApi_UnmarshalFocus(t *testing.T) {
	reader := strings.NewReader(`//Maya test scene
createNode transform -n "name1";
	rename -uid "CFAE1110-4845-2AC4-5BC0-CB8FB886A568";
	setAttr ".t" -type "double3" 0.1 1e-7 100;
// End of test scene`)

	focus := CommandTypes{
		CreateNodeCommand,
		SetAttrCommand,
	}

	mo, err := UnmarshalFocus(reader, focus)
	if err != nil {
		t.Error(err.Error())
	}

	if len(mo.Nodes) != 1 {
		t.Errorf("got len(mo.Nodes) %d, wont 1", len(mo.Nodes))
	}

	node, err := mo.GetNode("name1")
	if err != nil {
		t.Fatal(err.Error())
	}

	if node.Attrs == nil {
		t.Fatal("got node.Attrs was nil, wont not.")
	}

	if len(node.Attrs) != 1 {
		t.Errorf("got len(node.Attrs) %d, wont 1", len(node.Attrs))
	}
}