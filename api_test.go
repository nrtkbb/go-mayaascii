package mayaascii

import (
	"strings"
	"testing"
)

type stringTestData struct {
	title string
	value string
	wont  string
}

func stringTester(d stringTestData, t *testing.T) {
	if d.value != d.wont {
		t.Errorf("got %s was \"%s\", wont \"%s\"",
			d.title, d.value, d.wont)
	}
}

type intTestData struct {
	title string
	value int
	wont  int
}

func intTester(d intTestData, t *testing.T) {
	if d.value != d.wont {
		t.Errorf("got %s was %d, wont %d",
			d.title, d.value, d.wont)
	}
}

type boolTestData struct {
	title string
	value bool
	wont  bool
}

func boolTester(d boolTestData, t *testing.T) {
	if d.value != d.wont {
		t.Errorf("got %s was %v, wont %v",
			d.title, d.value, d.wont)
	}
}

func TestApi_File(t *testing.T) {
	reader := strings.NewReader(`//Maya test scene
file -rdi 1 -ns "test" -rfn "testRN" -typ "mayaAscii" "c:/test_data/test01.ma";
file -r -ns "test" -dr 1 -rfn "testRN" -typ "mayaAscii" "c:/test_data/test01.ma";
//End of test scene`)

	focus := CommandTypes{
		FileCommand,
	}

	mo, err := UnmarshalFocus(reader, focus)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(mo.Files) != 2 {
		t.Fatalf("got len(mo.Files) %d, wont 2", len(mo.Files))
	}

	ff := mo.Files[0] // first file = ff
	intTester(intTestData{"ff.GetReferenceDepthInfo()", ff.GetReferenceDepthInfo(), 1}, t)
	for _, d := range []stringTestData{
		{"ff.GetNamespace()", ff.GetNamespace(), "test"},
		{"ff.GetReferenceNode()", ff.GetReferenceNode(), "testRN"},
		{"ff.GetType()", ff.GetType(), "mayaAscii"},
		{"ff.Path()", ff.GetPath(), "c:/test_data/test01.ma"},
	} {
		stringTester(d, t)
	}

	sf := mo.Files[1]
	for _, d := range []boolTestData{
		{"sf.IsDeferReference()", sf.IsDeferReference(), true},
		{"sf.IsReference()", sf.IsReference(), true},
	} {
		boolTester(d, t)
	}

	for _, d := range []stringTestData{
		{"sf.GetNamespace()", sf.GetNamespace(), "test"},
		{"sf.GetReferenceNode()", sf.GetReferenceNode(), "testRN"},
		{"sf.GetType()", sf.GetType(), "mayaAscii"},
		{"sf.Path()", sf.GetPath(), "c:/test_data/test01.ma"},
	} {
		stringTester(d, t)
	}
}

func TestApi_FileInfo(t *testing.T) {
	reader := strings.NewReader(`//Maya test scene
fileInfo "application" "maya";
fileInfo "product" "Maya 2022";
fileInfo "version" "2022";
fileInfo "cutIdentifier" "202106180615-26a94e7f8c";
fileInfo "osv" "Windows 10 Pro v2009 (Build: 19042)";
fileInfo "UUID" "1225D207-450E-1399-04FC-89B0FA6B1B7F";
//End of test scene`)

	focus := CommandTypes{
		FileInfoCommand,
	}

	mo, err := UnmarshalFocus(reader, focus)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(mo.FileInfos) != 6 {
		t.Fatalf("got len(mo.FileInfos) %d, wont 6", len(mo.FileInfos))
	}

	for _, d := range []stringTestData{
		{"mo.FileInfos[0].GetName()", mo.FileInfos[0].GetName(), "application"},
		{"mo.FileInfos[1].GetName()", mo.FileInfos[1].GetName(), "product"},
		{"mo.FileInfos[2].GetName()", mo.FileInfos[2].GetName(), "version"},
		{"mo.FileInfos[3].GetName()", mo.FileInfos[3].GetName(), "cutIdentifier"},
		{"mo.FileInfos[4].GetName()", mo.FileInfos[4].GetName(), "osv"},
		{"mo.FileInfos[5].GetName()", mo.FileInfos[5].GetName(), "UUID"},
		{"mo.FileInfos[0].GetValue()", mo.FileInfos[0].GetValue(), "maya"},
		{"mo.FileInfos[1].GetValue()", mo.FileInfos[1].GetValue(), "Maya 2022"},
		{"mo.FileInfos[2].GetValue()", mo.FileInfos[2].GetValue(), "2022"},
		{"mo.FileInfos[3].GetValue()", mo.FileInfos[3].GetValue(), "202106180615-26a94e7f8c"},
		{"mo.FileInfos[4].GetValue()", mo.FileInfos[4].GetValue(), "Windows 10 Pro v2009 (Build: 19042)"},
		{"mo.FileInfos[5].GetValue()", mo.FileInfos[5].GetValue(), "1225D207-450E-1399-04FC-89B0FA6B1B7F"},
	} {
		stringTester(d, t)
	}
}

func TestApi_Select(t *testing.T) {
	reader := strings.NewReader(`//Maya test scene
select -ne :time1;
	setAttr ".o" 1;
	setAttr ".unw" 1;
select -ne :hardwareRenderGlobals;
	setAttr ".ctrs" 256;
	setAttr ".btrs" 512;
// End of test scene`)

	focus := CommandTypes{
		SelectCommand,
		SetAttrCommand,
	}

	mo, err := UnmarshalFocus(reader, focus)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(mo.Selects) != 2 {
		t.Fatalf("got len(mo.Select) %d, wont 2", len(mo.Selects))
	}

	for _, d := range []stringTestData{
		{"mo.Selects[0].GetName()", mo.Selects[0].GetName(), ":time1"},
		{"mo.Selects[1].GetName()", mo.Selects[1].GetName(), ":hardwareRenderGlobals"},
	} {
		stringTester(d, t)
	}

	for _, d := range []intTestData{
		{"len(mo.Selects[0].Attrs)", len(mo.Selects[0].Attrs), 2},
		{"len(mo.Selects[1].Attrs)", len(mo.Selects[1].Attrs), 2},
	} {
		intTester(d, t)
	}
}

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
		t.Fatal(err.Error())
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
