package parser

import (
	"testing"

	"github.com/nrtkbb/go-mayaascii/cmd"
)

func TestIsSameAttr(t *testing.T) {
	msg := `got isSameAttr "%s" vs "%s" is %v, wont %v`
	attrName := ".attrName"
	if !isSameAttr(attrName, attrName) {
		t.Errorf(
			msg,
			attrName, attrName,
			isSameAttr(attrName, attrName),
			true)
	}
	attrNameIdx := ".attrName[0]"
	if !isSameAttr(attrName, attrNameIdx) {
		t.Errorf(
			msg,
			attrName, attrNameIdx,
			isSameAttr(attrName, attrNameIdx),
			true)
	}
	if !isSameAttr(attrNameIdx, attrNameIdx) {
		t.Errorf(
			msg,
			attrNameIdx, attrNameIdx,
			isSameAttr(attrNameIdx, attrNameIdx),
			true)
	}
	attrSubName := ".attrName[0].subName"
	if isSameAttr(attrName, attrSubName) {
		t.Errorf(
			msg,
			attrName, attrSubName,
			isSameAttr(attrName, attrSubName),
			false)
	}
	attrSubNameIdx := ".attrName[0].subName[0]"
	if !isSameAttr(attrSubName, attrSubNameIdx) {
		t.Errorf(
			msg,
			attrSubName, attrSubNameIdx,
			isSameAttr(attrSubName, attrSubNameIdx),
			true)
	}
}

func TestMakeSetAttr_Size(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr -s 4 ".attrName";`)
	beforeSetAttr, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if beforeSetAttr.AttrType != cmd.TypeInvalid {
		t.Errorf(msg, "AttrType", beforeSetAttr.AttrType, cmd.TypeInvalid)
	}
	if *beforeSetAttr.Size != uint(4) {
		t.Errorf(msg, "Size", *beforeSetAttr.Size, uint(4))
	}
}

func TestMakeSetAttr_int(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" 1 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeInt {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeInt)
	}
	ret, err := cmd.ToAttrInt(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0].Int() != 1 || ret[1].Int() != 2 {
		t.Errorf(msg, "Attr", sa.Attr, []int{1, 2})
	}
	if len(sa.Attr) != 2 {
		t.Errorf(msg, "len(Attr)", len(sa.Attr), 2)
	}
	c.Clear()
	c.Append(`setAttr ".attrName" 3 4;`)
	sa, err = MakeSetAttr(c.Parse(), sa)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != cmd.TypeInt {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeInt)
	}
	ret, err = cmd.ToAttrInt(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if len(sa.Attr) != 4 {
		t.Errorf(msg, "len(Attr)", len(sa.Attr), 4)
	}
	if ret[0].Int() != 1 || ret[1].Int() != 2 ||
		ret[2].Int() != 3 || ret[3].Int() != 4 {
		t.Errorf(msg, "Attr", sa.Attr, []int{1, 2, 3, 4})
	}
}

func TestMakeSetAttr_int_toDouble_toInt(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" 1 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeInt {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeInt)
	}
	if len(sa.Attr) != 2 {
		t.Errorf(msg, "len(Attr)", len(sa.Attr), 2)
	}
	ret, err := cmd.ToAttrInt(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0].Int() != 1 || ret[1].Int() != 2 {
		var ai0 cmd.AttrInt = 1
		var ai1 cmd.AttrInt = 2
		t.Errorf(msg, "Attr", sa.Attr, [2]*cmd.AttrInt{&ai0, &ai1})
	}
	c.Clear()
	c.Append(`setAttr ".attrName" 3.3 4e+020 5e-020;`)
	sa, err = MakeSetAttr(c.Parse(), sa)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != cmd.TypeDouble {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDouble)
	}
	af, err := cmd.ToAttrFloat(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if len(af) != 5 {
		t.Errorf(msg, "len(Attr)", len(af), 5)
	}
	if af[0].Float() != 1 ||
		af[1].Float() != 2 ||
		af[2].Float() != 3.3 ||
		af[3].Float() != 4E+020 ||
		af[4].Float() != 5E-020 {
		var af0 cmd.AttrFloat = 1
		var af1 cmd.AttrFloat = 2
		var af2 cmd.AttrFloat = 3.3
		var af3 cmd.AttrFloat = 4E+020
		var af4 cmd.AttrFloat = 5E-020
		t.Errorf(msg, "Attr", sa.Attr, []*cmd.AttrFloat{
			&af0, &af1, &af2, &af3, &af4})
	}
	c.Clear()
	c.Append(`setAttr ".attrName" 5 6;`)
	sa, err = MakeSetAttr(c.Parse(), sa)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != cmd.TypeDouble {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDouble)
	}
	af, err = cmd.ToAttrFloat(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if len(af) != 7 {
		t.Errorf(msg, "len(Attr)", len(af), 7)
	}
	if af[0].Float() != 1 ||
		af[1].Float() != 2 ||
		af[2].Float() != 3.3 ||
		af[3].Float() != 4E+020 ||
		af[4].Float() != 5E-020 ||
		af[5].Float() != 5 ||
		af[6].Float() != 6 {
		var af0 cmd.AttrFloat = 1
		var af1 cmd.AttrFloat = 2
		var af2 cmd.AttrFloat = 3.3
		var af3 cmd.AttrFloat = 4E+020
		var af4 cmd.AttrFloat = 5E-020
		var af5 cmd.AttrFloat = 5
		var af6 cmd.AttrFloat = 6
		t.Errorf(msg, "Attr", sa.Attr, []*cmd.AttrFloat{
			&af0, &af1, &af2, &af3, &af4, &af5, &af6})
	}
}

func TestMakeSetAttr_string(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "string" "//network/folder/texture.jpg";`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeString {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeString)
	}
	if len(sa.Attr) != 1 {
		t.Errorf(msg, "len(Attr)", len(sa.Attr), 1)
	}
	ret, err := cmd.ToAttrString(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0].String() != "//network/folder/texture.jpg" {
		t.Errorf(msg, "Attr", ret[0].String(),
			"//network/folder/texture.jpg")
	}
}

func TestMakeSetAttr_stringArray(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "stringArray" 2 "//network/folder/texture.jpg" "//network/folder/texture.jpg";`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeStringArray {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeStringArray)
	}
	if len(sa.Attr) != 1 {
		t.Errorf(msg, "len(Attr)", len(sa.Attr), 1)
	}
	ret, err := cmd.ToAttrStringArray(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if (*ret[0])[0] != "//network/folder/texture.jpg" &&
		(*ret[0])[1] != "//network/folder/texture.jpg" {
		t.Errorf(msg, "Attr", *ret[0], cmd.AttrStringArray{
			"//network/folder/texture.jpg",
			"//network/folder/texture.jpg",
		})
	}
}

func TestMakeSetAttr_doubleWithExponent(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" 1e+020 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeDouble {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDouble)
	}
	if len(sa.Attr) != 2 {
		t.Errorf(msg, "len(Attr)", len(sa.Attr), 2)
	}
	ret, err := cmd.ToAttrFloat(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0].Float() != 1E+020 ||
		ret[1].Float() != 2 {
		t.Errorf(msg, "Attr", sa.Attr, []float64{1E+020, 2})
	}
}

func TestMakeSetAttr_double(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" 1.1 2.2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeDouble {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDouble)
	}
	if len(sa.Attr) != 2 {
		t.Errorf(msg, "len(Attr)", len(sa.Attr), 2)
	}
	ret, err := cmd.ToAttrFloat(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0].Float() != 1.1 ||
		ret[1].Float() != 2.2 {
		var af0 cmd.AttrFloat = 1.1
		var af1 cmd.AttrFloat = 2.2
		t.Errorf(msg, "Attr", sa.Attr, []*cmd.AttrFloat{
			&af0, &af1,
		})
	}
	c.Clear()
	c.Append(`setAttr ".attrName" 3.3 4.4;`)
	sa, err = MakeSetAttr(c.Parse(), sa)
	if sa.AttrType != cmd.TypeDouble {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDouble)
	}
	ret, err = cmd.ToAttrFloat(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) != 4 {
		t.Errorf(msg, "len(Attr)", len(ret), 4)
	}
	if ret[0].Float() != 1.1 ||
		ret[1].Float() != 2.2 ||
		ret[2].Float() != 3.3 ||
		ret[3].Float() != 4.4 {
		var af0 cmd.AttrFloat = 1.1
		var af1 cmd.AttrFloat = 2.2
		var af2 cmd.AttrFloat = 3.3
		var af3 cmd.AttrFloat = 4.4
		t.Errorf(msg, "Attr", sa.Attr, []*cmd.AttrFloat{
			&af0, &af1, &af2, &af3,
		})
	}
}

func testBool(t *testing.T, boolString string, wont bool) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" ` + boolString + `;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeBool {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeBool)
	}
	if len(sa.Attr) != 1 {
		t.Errorf(msg, "len(Attr)", len(sa.Attr), 1)
	}
	if b, ok := sa.Attr[0].(*cmd.AttrBool); !ok || b.Bool() != wont {
		t.Errorf(msg, "Attr", sa.Attr[0], wont)
	}
}

func TestMakeSetAttr_boolNo(t *testing.T) {
	testBool(t, "no", false)
	testBool(t, "false", false)
	testBool(t, "off", false)
	testBool(t, "yes", true)
	testBool(t, "true", true)
	testBool(t, "on", true)
}

func TestMakeSetAttr_short2(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "short2" 1 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeShort2 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeShort2)
	}
	if len(sa.Attr) != 1 {
		t.Errorf(msg, "len(Attr)", len(sa.Attr), 1)
	}
	s2, ok := sa.Attr[0].(*cmd.AttrShort2)
	if !ok || s2[0] != 1 || s2[1] != 2 {
		t.Errorf(msg, "Attr", sa.Attr, []cmd.AttrShort2{
			{1, 2},
		})
	}
}

func TestMakeSetAttr_short2_add(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "short2" 1 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	c.Append(`setAttr ".attrName" -type "short2" 3 4;`)
	sa, err = MakeSetAttr(c.Parse(), sa)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if len(sa.Attr) != 2 {
		t.Errorf(msg, "len(Attr)", len(sa.Attr), 2)
	}
	ret, err := cmd.ToAttrShort2(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) != 2 {
		t.Errorf(msg, "len(Attr)", len(ret), 2)
	}
	if ret[0][0] != 1 || ret[0][1] != 2 ||
		ret[1][0] != 3 || ret[1][1] != 4 {
		t.Errorf(msg, "Attr", sa.Attr, []cmd.AttrShort2{
			{1, 2},
			{3, 4},
		})
	}
}

func TestMakeSetAttr_short2_size(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr -s 2 ".attrName" -type "short2" 1 2 1 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeShort2 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeShort2)
	}
	ret, err := cmd.ToAttrShort2(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) != 2 {
		t.Errorf(msg, "len(Attr)", len(ret), 2)
	}
	if ret[0][0] != 1 || ret[0][1] != 2 ||
		ret[1][0] != 1 || ret[1][1] != 2 {
		t.Errorf(msg, "Attr", sa.Attr, []cmd.AttrShort2{
			{1, 2},
			{1, 2},
		})
	}
}

func TestMakeSetAttr_short2_sizeOver(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr -s 4 ".attrName";`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	c.Append(`setAttr ".attrName" -type "short2" 1 2 1 2;`)
	sa, err = MakeSetAttr(c.Parse(), sa)
	c.Append(`setAttr ".attrName" -type "short2" 1 2 1 2;`)
	sa, err = MakeSetAttr(c.Parse(), sa)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeShort2 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeShort2)
	}
	ret, err := cmd.ToAttrShort2(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) != 4 {
		t.Errorf(msg, "len(Attr)", len(ret), 4)
	}
	if ret[0][0] != 1 ||
		ret[0][1] != 2 ||
		ret[1][0] != 1 ||
		ret[1][1] != 2 ||
		ret[2][0] != 1 ||
		ret[2][1] != 2 ||
		ret[3][0] != 1 ||
		ret[3][1] != 2 {
		t.Errorf(msg, "Attr", sa.Attr, []cmd.AttrShort2{
			{1, 2},
			{1, 2},
			{1, 2},
			{1, 2},
		})
	}
}

func TestMakeSetAttr_long2(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "long2" 1 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeLong2 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeLong2)
	}
	ret, err := cmd.ToAttrLong2(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0][0] != 1 || ret[0][1] != 2 {
		var l20 cmd.AttrLong2 = [2]int{1, 2}
		t.Errorf(msg, "Attr", sa.Attr, []*cmd.AttrLong2{
			&l20,
		})
	}
	if len(ret) != 1 {
		t.Errorf(msg, "len(Attr)", len(ret), 1)
	}
	c.Clear()
	c.Append(`setAttr -s 2 ".attrName" -type "long2" 1 2 1 2;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != cmd.TypeLong2 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeLong2)
	}
	ret, err = cmd.ToAttrLong2(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0][0] != 1 ||
		ret[0][1] != 2 ||
		ret[1][0] != 1 ||
		ret[1][1] != 2 {
		var l20 cmd.AttrLong2 = [2]int{1, 2}
		var l21 cmd.AttrLong2 = [2]int{1, 2}
		t.Errorf(msg, "Attr", sa.Attr, []*cmd.AttrLong2{
			&l20, &l21,
		})
	}
	if len(ret) != 2 {
		t.Errorf(msg, "len(Attr)", len(ret), 2)
	}
}

func TestMakeSetAttr_short3(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "short3" 1 2 3;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeShort3 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeShort3)
	}
	ret, err := cmd.ToAttrShort3(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0][0] != 1 || ret[0][1] != 2 || ret[0][2] != 3 {
		t.Errorf(msg, "Attr", sa.Attr, []*cmd.AttrShort3{
			{1, 2, 3},
		})
	}
	if len(ret) != 1 {
		t.Errorf(msg, "len(Attr)", len(ret), 1)
	}
	c.Append(`setAttr -s 2 ".attrName" -type "short3" 1 2 3 1 2 3;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != cmd.TypeShort3 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeShort3)
	}
	ret, err = cmd.ToAttrShort3(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0][0] != 1 ||
		ret[0][1] != 2 ||
		ret[0][2] != 3 ||
		ret[1][0] != 1 ||
		ret[1][1] != 2 ||
		ret[1][2] != 3 {
		t.Errorf(msg, "Attr", sa.Attr, []*cmd.AttrShort3{
			{1, 2, 3},
			{1, 2, 3},
		})
	}
	if len(ret) != 2 {
		t.Errorf(msg, "len(Attr)", len(ret), 2)
	}
}

func TestMakeSetAttr_long3(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "long3" 1 2 3;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeLong3 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeLong3)
	}
	ret, err := cmd.ToAttrLong3(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0][0] != 1 || ret[0][1] != 2 || ret[0][2] != 3 {
		t.Errorf(msg, "Attr", sa.Attr, []*cmd.AttrLong3{
			{1, 2, 3},
		})
	}
	if len(ret) != 1 {
		t.Errorf(msg, "len(Attr)", len(ret), 1)
	}
	c.Clear()
	c.Append(`setAttr -s 2 ".attrName" -type "long3" 1 2 3 1 2 3;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != cmd.TypeLong3 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeLong3)
	}
	ret, err = cmd.ToAttrLong3(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0][0] != 1 ||
		ret[0][1] != 2 ||
		ret[0][2] != 3 ||
		ret[1][0] != 1 ||
		ret[1][1] != 2 ||
		ret[1][2] != 3 {
		t.Errorf(msg, "Attr", sa.Attr, [][3]int{
			{1, 2, 3},
			{1, 2, 3},
		})
	}
	if len(ret) != 2 {
		t.Errorf(msg, "len(Attr)", len(ret), 2)
	}
}

func TestMakeSetAttr_Int32Array(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "Int32Array" 2 1 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeInt32Array {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeInt32Array)
	}
	ret, err := cmd.ToAttrInt32Array(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if (*ret[0])[0] != 1 || (*ret[0])[1] != 2 {
		t.Errorf(msg, "Attr", sa.Attr, [][]int{{1, 2}})
	}
	if len(*ret[0]) != 2 {
		t.Errorf(msg, "len(Attr)", len(*ret[0]), 2)
	}
}

func TestMakeSetAttr_float2(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "float2" 1 2.2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeFloat2 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeFloat2)
	}
	ret, err := cmd.ToAttrFloat2(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0][0] != 1.0 || ret[0][1] != 2.2 {
		t.Errorf(msg, "Attr", sa.Attr, []*cmd.AttrFloat2{
			{1, 2.2},
		})
	}
	if len(ret) != 1 {
		t.Errorf(msg, "len(Attr)", len(ret), 1)
	}
	c.Clear()
	c.Append(`setAttr -s 2 ".attrName" -type "float2" 1 2.2 1 2.2;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != cmd.TypeFloat2 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeFloat2)
	}
	ret, err = cmd.ToAttrFloat2(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0][0] != 1.0 ||
		ret[0][1] != 2.2 ||
		ret[1][0] != 1.0 ||
		ret[1][1] != 2.2 {
		t.Errorf(msg, "Attr", sa.Attr, []*cmd.AttrFloat2{
			{1, 2.2},
			{1, 2.2},
		})
	}
	if len(ret) != 2 {
		t.Errorf(msg, "len(Attr)", len(ret), 2)
	}
}

func TestMakeSetAttr_float3(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "float3" 1 2.2 3.3;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeFloat3 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeFloat3)
	}
	ret, err := cmd.ToAttrFloat3(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0][0] != 1.0 || ret[0][1] != 2.2 || ret[0][2] != 3.3 {
		t.Errorf(msg, "Attr", sa.Attr, []*cmd.AttrFloat3{
			{1, 2.2, 3.3},
		})
	}
	if len(ret) != 1 {
		t.Errorf(msg, "len(Attr)", len(ret), 1)
	}
	c.Clear()
	c.Append(`setAttr -s 2 ".attrName" -type "float3" 1 2.2 3.3 1 2.2 3.3;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != cmd.TypeFloat3 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeFloat3)
	}
	ret, err = cmd.ToAttrFloat3(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0][0] != 1.0 ||
		ret[0][1] != 2.2 ||
		ret[0][2] != 3.3 ||
		ret[1][0] != 1.0 ||
		ret[1][1] != 2.2 ||
		ret[1][2] != 3.3 {
		t.Errorf(msg, "Attr", sa.Attr, []*cmd.AttrFloat3{
			{1, 2.2, 3.3},
			{1, 2.2, 3.3},
		})
	}
	if len(ret) != 2 {
		t.Errorf(msg, "len(Attr)", len(ret), 2)
	}
}

func TestMakeSetAttr_double2(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "double2" 1 2.2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeDouble2 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDouble2)
	}
	ret, err := cmd.ToAttrDouble2(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0][0] != 1.0 || ret[0][1] != 2.2 {
		t.Errorf(msg, "Attr", sa.Attr, []*cmd.AttrDouble2{
			{1, 2.2},
		})
	}
	if len(ret) != 1 {
		t.Errorf(msg, "len(Attr)", len(ret), 1)
	}
	c.Clear()
	c.Append(`setAttr -s 2 ".attrName" -type "double2" 1 2.2 1 2.2;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != cmd.TypeDouble2 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDouble2)
	}
	ret, err = cmd.ToAttrDouble2(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0][0] != 1.0 ||
		ret[0][1] != 2.2 ||
		ret[1][0] != 1.0 ||
		ret[1][1] != 2.2 {
		t.Errorf(msg, "Attr", sa.Attr, [][2]float64{
			{1, 2.2},
			{1, 2.2},
		})
	}
	if len(ret) != 2 {
		t.Errorf(msg, "len(Attr)", len(ret), 2)
	}
}

func TestMakeSetAttr_double3(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "double3" 1 2.2 3.3;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeDouble3 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDouble3)
	}
	ret, err := cmd.ToAttrDouble3(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0][0] != 1.0 || ret[0][1] != 2.2 || ret[0][2] != 3.3 {
		t.Errorf(msg, "Attr", sa.Attr, []float64{1, 2.2, 3.3})
	}
	if len(ret) != 1 {
		t.Errorf(msg, "len(Attr)", len(ret), 1)
	}
	c.Clear()
	c.Append(`setAttr -s 2 ".attrName" -type "double3" 1 2.2 3.3 1 2.2 3.3;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != cmd.TypeDouble3 {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDouble3)
	}
	ret, err = cmd.ToAttrDouble3(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0][0] != 1.0 ||
		ret[0][1] != 2.2 ||
		ret[0][2] != 3.3 ||
		ret[1][0] != 1.0 ||
		ret[1][1] != 2.2 ||
		ret[1][2] != 3.3 {
		t.Errorf(msg, "Attr", sa.Attr, [][3]float64{
			{1, 2.2, 3.3},
			{1, 2.2, 3.3},
		})
	}
	if len(ret) != 2 {
		t.Errorf(msg, "len(Attr)", len(ret), 2)
	}
}

func TestMakeSetAttr_doubleArray(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "doubleArray" 2 1.1 2.2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeDoubleArray {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDoubleArray)
	}
	ret, err := cmd.ToAttrDoubleArray(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if (*ret[0])[0] != 1.1 || (*ret[0])[1] != 2.2 {
		t.Errorf(msg, "Attr", sa.Attr, []float64{1, 2.2, 3.3})
	}
	if len(*ret[0]) != 2 {
		t.Errorf(msg, "len(Attr)", len(*ret[0]), 2)
	}
	c.Clear()
	c.Append(`setAttr ".attrName" -type "doubleArray" 0;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != cmd.TypeDoubleArray {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDoubleArray)
	}
	ret, err = cmd.ToAttrDoubleArray(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if len(*ret[0]) != 0 {
		t.Errorf(msg, "len(Attr)", len(*ret[0]), 0)
	}
}

func TestMakeSetAttr_matrix(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".ix" -type "matrix" 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeMatrix {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeMatrix)
	}
	ret, err := cmd.ToAttrMatrix(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	wontMt := [16]float64{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}
	allOk := true
	for i, v := range *ret[0] {
		if v != wontMt[i] {
			allOk = false
			break
		}
	}
	if !allOk {
		t.Errorf(msg, "Attr", sa.Attr, wontMt)
	}
	if len(*ret[0]) != 16 {
		t.Errorf(msg, "len(Attr)", len(*ret[0]), 16)
	}
}

func TestMakeSetAttr_matrix_xform(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".ix" -type "matrix" "xform" 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
0 0 0 0 0 0 0 0 0 0 1 0 0 0 1 1 1 1 yes;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeMatrixXform {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeMatrixXform)
	}
	ret, err := cmd.ToAttrMatrixXform(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if ret[0].Scale.X != 1 || ret[0].Scale.Y != 1 || ret[0].Scale.Z != 1 ||
		ret[0].Rotate.X != 0 || ret[0].Rotate.Y != 0 || ret[0].Rotate.Z != 0 ||
		ret[0].RotateOrder != cmd.RotateOrderXYZ ||
		ret[0].Translate.X != 0 || ret[0].Translate.Y != 0 || ret[0].Translate.Z != 0 ||
		ret[0].Shear.XY != 0 || ret[0].Shear.XZ != 0 || ret[0].Shear.YZ != 0 ||
		ret[0].ScalePivot.X != 0 || ret[0].ScalePivot.Y != 0 || ret[0].ScalePivot.Z != 0 ||
		ret[0].ScaleTranslate.X != 0 || ret[0].ScaleTranslate.Y != 0 || ret[0].ScaleTranslate.Z != 0 ||
		ret[0].RotatePivot.X != 0 || ret[0].RotatePivot.Y != 0 || ret[0].RotatePivot.Z != 0 ||
		ret[0].RotateTranslation.X != 0 || ret[0].RotateTranslation.Y != 0 || ret[0].RotateTranslation.Z != 0 ||
		ret[0].RotateOrient.W != 0 || ret[0].RotateOrient.X != 0 || ret[0].RotateOrient.Y != 0 || ret[0].RotateOrient.Z != 1 ||
		ret[0].JointOrient.W != 0 || ret[0].JointOrient.X != 0 || ret[0].JointOrient.Y != 0 || ret[0].JointOrient.Z != 1 ||
		ret[0].InverseParentScale.X != 1 || ret[0].InverseParentScale.Y != 1 || ret[0].InverseParentScale.Z != 1 ||
		ret[0].CompensateForParentScale == false {
		t.Errorf(msg, "Attr", sa.Attr, nil)
	}
}

func TestMakeSetAttr_pointArray(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "pointArray" 1 1.1 2.2 3.3 4.4;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypePointArray {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypePointArray)
	}
	ret, err := cmd.ToAttrPointArray(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	pa := *ret[0]
	if pa[0].X != 1.1 || pa[0].Y != 2.2 || pa[0].Z != 3.3 || pa[0].W != 4.4 {
		t.Errorf(msg, "Attr", sa.Attr, cmd.AttrPointArray{
			{X: 1.1, Y: 2.2, Z: 3.3, W: 4.4},
		})
	}
	if len(*ret[0]) != 1 {
		t.Errorf(msg, "len(Attr)", len(*ret[0]), 1)
	}
	c.Clear()
	c.Append(`setAttr ".attrName" -type "pointArray" 2 1.1 2.2 3.3 4.4 1.1 2.2 3.3 4.4;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != cmd.TypePointArray {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypePointArray)
	}
	ret, err = cmd.ToAttrPointArray(sa.Attr)
	pa = *ret[0]
	if pa[0].X != 1.1 ||
		pa[0].Y != 2.2 ||
		pa[0].Z != 3.3 ||
		pa[0].W != 4.4 ||
		pa[1].X != 1.1 ||
		pa[1].Y != 2.2 ||
		pa[1].Z != 3.3 ||
		pa[1].W != 4.4 {
		t.Errorf(msg, "Attr", sa.Attr, cmd.AttrPointArray{
			{X: 1.1, Y: 2.2, Z: 3.3, W: 4.4},
			{X: 1.1, Y: 2.2, Z: 3.3, W: 4.4},
		})
	}
	if len(pa) != 2 {
		t.Errorf(msg, "len(Attr)", len(pa), 2)
	}
}

func TestMakeSetAttr_polyFaces(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr -s 2 ".attrName" -type "polyFaces"
f 3 1 2 3
mc 1 3 0 1 2
f 3 2 3 4
mc 2 3 2 3 4;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypePolyFaces {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypePolyFaces)
	}
	ret, err := cmd.ToAttrPolyFaces(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if (*ret[0]).FaceEdge[0] != 1 ||
		(*ret[0]).FaceEdge[1] != 2 ||
		(*ret[0]).FaceEdge[2] != 3 ||
		(*ret[0]).MultiColor[0].ColorIndex != 1 ||
		(*ret[0]).MultiColor[0].ColorIDs[0] != 0 ||
		(*ret[0]).MultiColor[0].ColorIDs[1] != 1 ||
		(*ret[0]).MultiColor[0].ColorIDs[2] != 2 ||
		(*ret[1]).FaceEdge[0] != 2 ||
		(*ret[1]).FaceEdge[1] != 3 ||
		(*ret[1]).FaceEdge[2] != 4 ||
		(*ret[1]).MultiColor[0].ColorIndex != 2 ||
		(*ret[1]).MultiColor[0].ColorIDs[0] != 2 ||
		(*ret[1]).MultiColor[0].ColorIDs[1] != 3 ||
		(*ret[1]).MultiColor[0].ColorIDs[2] != 4 {
		t.Errorf(msg, "Attr", sa.Attr, []cmd.AttrPolyFaces{
			{
				FaceEdge: []int{1, 2, 3},
				MultiColor: []cmd.AttrMultiColor{
					{
						ColorIndex: 1,
						ColorIDs:   []int{1, 2, 3},
					},
				},
			},
			{
				FaceEdge: []int{2, 3, 4},
				MultiColor: []cmd.AttrMultiColor{
					{
						ColorIndex: 2,
						ColorIDs:   []int{2, 3, 4},
					},
				},
			},
		})
	}
	if len(ret) != 2 {
		t.Errorf(msg, "len(Attr)", len(ret), 2)
	}
}

func TestMakeSetAttr_polyFacesMax(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr -s 2 ".attrName" -type "polyFaces"
	f 3 1 2 3
	h 3 5 6 7
	mu 0 3 0 1 3
	mu 1 3 0 1 3
	mc 1 3 0 1 2
	f 3 2 3 4
	mu 0 3 2 3 4
	mu 1 3 2 3 4
	mc 2 3 2 3 4;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypePolyFaces {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypePolyFaces)
	}
	ret, err := cmd.ToAttrPolyFaces(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if (*ret[0]).FaceEdge[0] != 1 ||
		(*ret[0]).FaceEdge[1] != 2 ||
		(*ret[0]).FaceEdge[2] != 3 ||
		(*ret[0]).HoleEdge[0] != 5 ||
		(*ret[0]).HoleEdge[1] != 6 ||
		(*ret[0]).HoleEdge[2] != 7 ||
		(*ret[0]).FaceUV[0].UVSet != 0 ||
		(*ret[0]).FaceUV[0].FaceUV[0] != 0 ||
		(*ret[0]).FaceUV[0].FaceUV[1] != 1 ||
		(*ret[0]).FaceUV[0].FaceUV[2] != 3 ||
		(*ret[0]).FaceUV[1].UVSet != 1 ||
		(*ret[0]).FaceUV[1].FaceUV[0] != 0 ||
		(*ret[0]).FaceUV[1].FaceUV[1] != 1 ||
		(*ret[0]).FaceUV[1].FaceUV[2] != 3 ||
		(*ret[0]).MultiColor[0].ColorIndex != 1 ||
		(*ret[0]).MultiColor[0].ColorIDs[0] != 0 ||
		(*ret[0]).MultiColor[0].ColorIDs[1] != 1 ||
		(*ret[0]).MultiColor[0].ColorIDs[2] != 2 ||
		(*ret[1]).FaceEdge[0] != 2 ||
		(*ret[1]).FaceEdge[1] != 3 ||
		(*ret[1]).FaceEdge[2] != 4 ||
		(*ret[1]).FaceUV[0].UVSet != 0 ||
		(*ret[1]).FaceUV[0].FaceUV[0] != 2 ||
		(*ret[1]).FaceUV[0].FaceUV[1] != 3 ||
		(*ret[1]).FaceUV[0].FaceUV[2] != 4 ||
		(*ret[1]).FaceUV[1].UVSet != 1 ||
		(*ret[1]).FaceUV[1].FaceUV[0] != 2 ||
		(*ret[1]).FaceUV[1].FaceUV[1] != 3 ||
		(*ret[1]).FaceUV[1].FaceUV[2] != 4 ||
		(*ret[1]).MultiColor[0].ColorIndex != 2 ||
		(*ret[1]).MultiColor[0].ColorIDs[0] != 2 ||
		(*ret[1]).MultiColor[0].ColorIDs[1] != 3 ||
		(*ret[1]).MultiColor[0].ColorIDs[2] != 4 {
		t.Errorf(msg, "Attr", sa.Attr, []cmd.AttrPolyFaces{
			{
				FaceEdge: []int{1, 2, 3},
				HoleEdge: []int{5, 6, 7},
				FaceUV: []cmd.AttrFaceUV{
					{
						UVSet:  0,
						FaceUV: []int{0, 1, 3},
					},
					{
						UVSet:  1,
						FaceUV: []int{0, 1, 3},
					},
				},
				MultiColor: []cmd.AttrMultiColor{
					{
						ColorIndex: 0,
						ColorIDs:   []int{0, 1, 2},
					},
				},
			},
			{
				FaceEdge: []int{2, 3, 4},
				FaceUV: []cmd.AttrFaceUV{
					{
						UVSet:  0,
						FaceUV: []int{2, 3, 4},
					},
					{
						UVSet:  1,
						FaceUV: []int{2, 3, 4},
					},
				},
				MultiColor: []cmd.AttrMultiColor{
					{
						ColorIndex: 2,
						ColorIDs:   []int{2, 3, 4},
					},
				},
			},
		})
	}
	if len(ret) != 2 {
		t.Errorf(msg, "len(Attr)", len(ret), 2)
	}
}

func TestMakeDataPolyComponent(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".cd" -type "dataPolyComponent" Index_Data Edge 24
		0 9.0600013732910156
		1 4.409998893737793
		2 9.0600013732910156
		3 9.2000007629394531
		4 9.0600013732910156
		5 9.0600013732910156
		9 9.0600013732910156
		10 9.0600013732910156
		11 9.0600013732910156
		12 9.0600013732910156
		13 4.409998893737793
		15 4.6099758148193359
		18 4.409998893737793
		19 4.6099758148193359
		20 4.409998893737793
		21 9.2000007629394531
		23 9.2000007629394531
		26 9.2000007629394531
		27 9.2000007629394531
		28 9.2000007629394531
		30 4.6099758148193359
		32 4.6099758148193359
		34 4.6099758148193359
		35 4.6099758148193359 ;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeDataPolyComponent {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDataPolyComponent)
	}
	ret, err := cmd.ToAttrDataPolyComponent(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if (*ret[0]).PolyComponentType != cmd.DPCedge ||
		(*ret[0]).IndexValue[0] != 9.0600013732910156 ||
		(*ret[0]).IndexValue[1] != 4.409998893737793 ||
		(*ret[0]).IndexValue[2] != 9.0600013732910156 ||
		(*ret[0]).IndexValue[3] != 9.2000007629394531 ||
		(*ret[0]).IndexValue[4] != 9.0600013732910156 ||
		(*ret[0]).IndexValue[5] != 9.0600013732910156 ||
		(*ret[0]).IndexValue[9] != 9.0600013732910156 ||
		(*ret[0]).IndexValue[10] != 9.0600013732910156 ||
		(*ret[0]).IndexValue[11] != 9.0600013732910156 ||
		(*ret[0]).IndexValue[12] != 9.0600013732910156 ||
		(*ret[0]).IndexValue[13] != 4.409998893737793 ||
		(*ret[0]).IndexValue[15] != 4.6099758148193359 ||
		(*ret[0]).IndexValue[18] != 4.409998893737793 ||
		(*ret[0]).IndexValue[19] != 4.6099758148193359 ||
		(*ret[0]).IndexValue[20] != 4.409998893737793 ||
		(*ret[0]).IndexValue[21] != 9.2000007629394531 ||
		(*ret[0]).IndexValue[23] != 9.2000007629394531 ||
		(*ret[0]).IndexValue[26] != 9.2000007629394531 ||
		(*ret[0]).IndexValue[27] != 9.2000007629394531 ||
		(*ret[0]).IndexValue[28] != 9.2000007629394531 ||
		(*ret[0]).IndexValue[30] != 4.6099758148193359 ||
		(*ret[0]).IndexValue[32] != 4.6099758148193359 ||
		(*ret[0]).IndexValue[34] != 4.6099758148193359 ||
		(*ret[0]).IndexValue[35] != 4.6099758148193359 {
		t.Errorf(msg, "Attr", sa.Attr, []cmd.AttrDataPolyComponent{
			{
				PolyComponentType: cmd.DPCedge,
				IndexValue: map[int]float64{
					0:  9.0600013732910156,
					1:  4.409998893737793,
					2:  9.0600013732910156,
					3:  9.2000007629394531,
					4:  9.0600013732910156,
					5:  9.0600013732910156,
					9:  9.0600013732910156,
					10: 9.0600013732910156,
					11: 9.0600013732910156,
					12: 9.0600013732910156,
					13: 4.409998893737793,
					15: 4.6099758148193359,
					18: 4.409998893737793,
					19: 4.6099758148193359,
					20: 4.409998893737793,
					21: 9.2000007629394531,
					23: 9.2000007629394531,
					26: 9.2000007629394531,
					27: 9.2000007629394531,
					28: 9.2000007629394531,
					30: 4.6099758148193359,
					32: 4.6099758148193359,
					34: 4.6099758148193359,
					35: 4.6099758148193359,
				},
			},
		})
	}
	if len(ret) != 1 {
		t.Errorf(msg, "len(Attr)", len(ret), 1)
	}
	if len((*ret[0]).IndexValue) != 24 {
		t.Errorf(msg, "len(Attr.IndexValue)", len((*ret[0]).IndexValue), 24)
	}
}

func sameDPC(t *testing.T, dpc []*cmd.AttrDataPolyComponent, dpcType cmd.AttrDPCType) {
	d := *dpc[0]
	if d.PolyComponentType != dpcType {
		msg := `got SetAttr %s %s, wont %s`
		t.Errorf(msg, "Attr", dpc, []cmd.AttrDataPolyComponent{
			{
				PolyComponentType: dpcType,
				IndexValue:        map[int]float64{},
			},
		})
	}
}

func TestMakeDataPolyComponentVertex(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".cvd" -type "dataPolyComponent" Index_Data Vertex 0;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeDataPolyComponent {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDataPolyComponent)
	}
	ret, err := cmd.ToAttrDataPolyComponent(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	sameDPC(t, ret, cmd.DPCvertex)
	if len(ret) != 1 {
		t.Errorf(msg, "len(Attr)", len(ret), 1)
	}
	if len((*ret[0]).IndexValue) != 0 {
		t.Errorf(msg, "len(Attr.IndexValue)", len((*ret[0]).IndexValue), 0)
	}
}

func TestMakeDataPolyComponentUV(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".pd[0]" -type "dataPolyComponent" Index_Data UV 0 ;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeDataPolyComponent {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDataPolyComponent)
	}
	ret, err := cmd.ToAttrDataPolyComponent(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	sameDPC(t, ret, cmd.DPCuv)
	if len(ret) != 1 {
		t.Errorf(msg, "len(Attr)", len(ret), 1)
	}
	if len((*ret[0]).IndexValue) != 0 {
		t.Errorf(msg, "len(Attr.IndexValue)", len((*ret[0]).IndexValue), 0)
	}
}

func TestMakeDataPolyComponentFace(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".hfd" -type "dataPolyComponent" Index_Data Face 0 ;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeDataPolyComponent {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDataPolyComponent)
	}
	ret, err := cmd.ToAttrDataPolyComponent(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	sameDPC(t, ret, cmd.DPCface)
	if len(ret) != 1 {
		t.Errorf(msg, "len(Attr)", len(ret), 1)
	}
	if len((*ret[0]).IndexValue) != 0 {
		t.Errorf(msg, "len(Attr.IndexValue)", len((*ret[0]).IndexValue), 0)
	}
}

func TestMakeDataReferenceEdits(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".ed" -type "dataReferenceEdits"
	"namespaceRN"
	"namespace:childNameSpaceRN" 11
    0 "nodeNameA" "nodeNameB" "-s -r "
	1 |namespace:topNode "extraAttr" "ea" " -ci 1 -nn \"ea\" -at \"double\""
	2 "|namespace:topNode" "ea" " -k 1 0"
	3 "namespace:nodeNameA.attrNameA" "namespace:nodeNameB.attrNameB" ""
	4 "|namespace:topNode" "dexAttr" ""
	5 0 "namespace:childNameSpaceRN" "|namespace:topNode.attrNameA" "|namespace:topNode.attrNameB"
	"namespace:childNameSpaceRN.placeHolderList[1]" "namespace:childNameSpaceRN.placeHolderList[2]" ""
	5 3 "namespace:childNameSpaceRN" "|namespace:topNode|namespace:childNode.attrNameA"
	"namespace:childNameSpaceRN.placeHolderList[3]" ""
	5 4 "namespace:childNameSpaceRN" "|namespace:topNode|namespace:childNode.attrNameB"
	"namespace:childNameSpaceRN.placeHolderList[4]" ""
    7 "fcurve" "|namespace:nodeName_attrName_X" 1
	"add 396 -4131.291016 18 18 1 0 0 423 -4131.291016 18 18 1 0 0" 0
    8 "|namespace:topNode" "attrNameA"
    9 "|namespace:topNode" "attrNameA";`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %v, wont %v`
	if sa.AttrType != cmd.TypeDataReferenceEdits {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDataReferenceEdits)
	}
	ret, err := cmd.ToAttrDataReferenceEdits(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) != 1 {
		t.Errorf(msg, "len(ret)", len(ret), 1)
	}
	re := ret[0]
	if re.TopReferenceNode != "namespaceRN" {
		t.Errorf(msg, "TopReferenceNode", re.TopReferenceNode, "namespaceRN")
	}
	if len(re.ReferenceEdits) != 1 {
		t.Errorf(msg, "len(re.ReferenceEdits)", len(re.ReferenceEdits), 1)
	}
	referenceEdit := re.ReferenceEdits[0]
	if referenceEdit.ReferenceNode != "namespace:childNameSpaceRN" {
		t.Errorf(msg, "referenceEdit.ReferenceNode",
			referenceEdit.ReferenceNode, "namespace:childNameSpaceRN")
	}
	if referenceEdit.CommandNum != 11 {
		t.Errorf(msg, "referenceEdit.CommandNum", referenceEdit.CommandNum, 11)
	}

	// Parent
	if len(referenceEdit.Parents) != 1 {
		t.Errorf(msg, "len(referenceEdit.Parents)", len(referenceEdit.Parents), 1)
	}
	parent := referenceEdit.Parents[0]
	if parent.NodeA != "nodeNameA" {
		t.Errorf(msg, "parent.NodeA", parent.NodeA, "nodeNameA")
	}
	if parent.NodeB != "nodeNameB" {
		t.Errorf(msg, "parent.NodeB", parent.NodeB, "nodeNameB")
	}
	if parent.Arguments != "-s -r " {
		t.Errorf(msg, "parent.Arguments", parent.Arguments, "-s -r ")
	}

	// AddAttr
	if len(referenceEdit.AddAttrs) != 1 {
		t.Errorf(msg, "len(referenceEdit.AddAttrs)", len(referenceEdit.AddAttrs), 1)
	}
	addAttr := referenceEdit.AddAttrs[0]
	if addAttr.Node != "|namespace:topNode" {
		t.Errorf(msg, "addAttr.Node", addAttr.Node, "|namespace:topNode")
	}
	if addAttr.LongAttr != "extraAttr" {
		t.Errorf(msg, "addAttr.LongAttr", addAttr.LongAttr, "extraAttr")
	}
	if addAttr.ShortAttr != "ea" {
		t.Errorf(msg, "addAttr.ShortAttr", addAttr.ShortAttr, "ea")
	}
	if addAttr.Arguments != " -ci 1 -nn \\\"ea\\\" -at \\\"double\\\"" {
		t.Errorf(msg, "addAttr.Arguments", addAttr.Arguments,
			" -ci 1 -nn \\\"ea\\\" -at \\\"double\\\"")
	}

	// SetAttr
	if len(referenceEdit.SetAttrs) != 1 {
		t.Errorf(msg, "len(referenceEdit.SetAttrs)", len(referenceEdit.SetAttrs), 1)
	}
	setAttr := referenceEdit.SetAttrs[0]
	if setAttr.Node != "|namespace:topNode" {
		t.Errorf(msg, "setAttr.Node", setAttr.Node, "|namespace:topNode")
	}
	if setAttr.Attr != "ea" {
		t.Errorf(msg, "setAttr.Attr", setAttr.Attr, "ea")
	}
	if setAttr.Arguments != " -k 1 0" {
		t.Errorf(msg, "setAttr.Arguments", setAttr.Arguments, " -k 1 0")
	}

	// DisconnectAttr
	if len(referenceEdit.DisconnectAttrs) != 1 {
		t.Errorf(msg, "len(referenceEdit.DisconnectAttrs)",
			len(referenceEdit.DisconnectAttrs), 1)
	}
	disconnectAttr := referenceEdit.DisconnectAttrs[0]
	if disconnectAttr.SourcePlug != "namespace:nodeNameA.attrNameA" {
		t.Errorf(msg, "disconnectAttr.SourcePlug", disconnectAttr.SourcePlug,
			"namespace:nodeNameA.attrNameA")
	}
	if disconnectAttr.DistPlug != "namespace:nodeNameB.attrNameB" {
		t.Errorf(msg, "disconnectAttr.DistPlug", disconnectAttr.DistPlug,
			"namespace:nodeNameB.attrNameB")
	}
	if disconnectAttr.Arguments != "" {
		t.Errorf(msg, "disconnectAttr.Arguments", disconnectAttr.Arguments, "")
	}

	// DeleteAttr
	if len(referenceEdit.DeleteAttrs) != 1 {
		t.Errorf(msg, "referenceEdit.DeleteAttrs", len(referenceEdit.DeleteAttrs), 1)
	}
	deleteAttr := referenceEdit.DeleteAttrs[0]
	if deleteAttr.Node != "|namespace:topNode" {
		t.Errorf(msg, "deleteAttr.Node", deleteAttr.Node, "|namespace:topNode")
	}
	if deleteAttr.Attr != "dexAttr" {
		t.Errorf(msg, "deleteAttr.Attr", deleteAttr.Attr, "dexAttr")
	}
	if deleteAttr.Arguments != "" {
		t.Errorf(msg, "deleteAttr.Arguments", deleteAttr.Arguments, "")
	}

	// ConnectAttr
	if len(referenceEdit.ConnectAttrs) != 3 {
		t.Errorf(msg, "len(referenceEdit.ConnectAttrs)", len(referenceEdit.ConnectAttrs), 3)
	}
	connectAttr1 := referenceEdit.ConnectAttrs[0]
	if connectAttr1.MagicNumber != 0 {
		t.Errorf(msg, "connectAttr1.MagicNumber", connectAttr1.MagicNumber, 0)
	}
	if connectAttr1.ReferenceNode != "namespace:childNameSpaceRN" {
		t.Errorf(msg, "connectAttr1.ReferenceNode", connectAttr1.ReferenceNode,
			"namespace:childNameSpaceRN")
	}
	if connectAttr1.SourcePlug != "|namespace:topNode.attrNameA" {
		t.Errorf(msg, "connectAttr1.SourcePlug", connectAttr1.SourcePlug,
			"|namespace:topNode.attrNameA")
	}
	if connectAttr1.DistPlug != "|namespace:topNode.attrNameB" {
		t.Errorf(msg, "connectAttr1.DistPlug", connectAttr1.DistPlug,
			"|namespace:topNode.attrNameB")
	}
	if *connectAttr1.SourcePHL != "namespace:childNameSpaceRN.placeHolderList[1]" {
		t.Errorf(msg, "*connectAttr1.SourcePHL", *connectAttr1.SourcePHL,
			"namespace:childNameSpaceRN.placeHolderList[1]")
	}
	if *connectAttr1.DistPHL != "namespace:childNameSpaceRN.placeHolderList[2]" {
		t.Errorf(msg, "*connectAttr1.DistPHL", *connectAttr1.DistPHL,
			"namespace:childNameSpaceRN.placeHolderList[2]")
	}
	if connectAttr1.Arguments != "" {
		t.Errorf(msg, "connectAttr1.Arguments", connectAttr1.Arguments, "")
	}

	connectAttr2 := referenceEdit.ConnectAttrs[1]
	if connectAttr2.MagicNumber != 3 {
		t.Errorf(msg, "connectAttr2.MagicNumber", connectAttr2.MagicNumber, 3)
	}
	if connectAttr2.ReferenceNode != "namespace:childNameSpaceRN" {
		t.Errorf(msg, "connectAttr2.ReferenceNode", connectAttr2.ReferenceNode,
			"namespace:childNameSpaceRN")
	}
	if connectAttr2.SourcePlug != "|namespace:topNode|namespace:childNode.attrNameA" {
		t.Errorf(msg, "connectAttr2.SourcePlug", connectAttr2.SourcePlug,
			"|namespace:topNode|namespace:childNode.attrNameA")
	}
	if connectAttr2.DistPlug != "namespace:childNameSpaceRN.placeHolderList[3]" {
		t.Errorf(msg, "connectAttr2.DistPlug", connectAttr2.DistPlug,
			"namespace:childNameSpaceRN.placeHolderList[3]")
	}
	if connectAttr2.SourcePHL != nil {
		t.Errorf(msg, "connectAttr2.SourcePHL", connectAttr2.SourcePHL, nil)
	}
	if connectAttr2.DistPHL != nil {
		t.Errorf(msg, "connectAttr2.DistPHL", connectAttr2.DistPHL, nil)
	}
	if connectAttr2.Arguments != "" {
		t.Errorf(msg, "connectAttr2.Arguments", connectAttr2.Arguments, "")
	}

	connectAttr3 := referenceEdit.ConnectAttrs[2]
	if connectAttr3.MagicNumber != 4 {
		t.Errorf(msg, "connectAttr3.MagicNumber", connectAttr3.MagicNumber, 4)
	}
	if connectAttr3.ReferenceNode != "namespace:childNameSpaceRN" {
		t.Errorf(msg, "connectAttr2.ReferenceNode", connectAttr3.ReferenceNode,
			"namespace:childNameSpaceRN")
	}
	if connectAttr3.SourcePlug != "|namespace:topNode|namespace:childNode.attrNameB" {
		t.Errorf(msg, "connectAttr3.SourcePlug", connectAttr3.SourcePlug,
			"|namespace:topNode|namespace:childNode.attrNameB")
	}
	if connectAttr3.DistPlug != "namespace:childNameSpaceRN.placeHolderList[4]" {
		t.Errorf(msg, "connectAttr3.DistPlug", connectAttr3.DistPlug,
			"namespace:childNameSpaceRN.placeHolderList[4]")
	}
	if connectAttr3.SourcePHL != nil {
		t.Errorf(msg, "connectAttr3.SourcePHL", connectAttr3.SourcePHL, nil)
	}
	if connectAttr3.DistPHL != nil {
		t.Errorf(msg, "connectAttr3.DistPHL", connectAttr3.DistPHL, nil)
	}
	if connectAttr3.Arguments != "" {
		t.Errorf(msg, "connectAttr3.Arguments", connectAttr3.Arguments, "")
	}

	// Relationship
	if len(referenceEdit.Relationships) != 1 {
		t.Errorf(msg, "len(referenceEdit.Relationship)", len(referenceEdit.Relationships), 1)
	}
	relationship := referenceEdit.Relationships[0]
	if relationship.Type != "fcurve" {
		t.Errorf(msg, "relationship.Type", relationship.Type, "fcurve")
	}
	if relationship.CommandNum != 1 {
		t.Errorf(msg, "relationship.CommandNum", relationship.CommandNum, 1)
	}
	if relationship.NodeName != "|namespace:nodeName_attrName_X" {
		t.Errorf(msg, "relationship.NodeName", relationship.NodeName,
			"|namespace:nodeName_attrName_X")
	}
	if len(relationship.Commands) != 1 {
		t.Errorf(msg, "len(relationship.Commands)", len(relationship.Commands), 1)
	}
	relationshipCommand := relationship.Commands[0]
	if relationshipCommand != "add 396 -4131.291016 18 18 1 0 0 423 -4131.291016 18 18 1 0 0" {
		t.Errorf(msg, "relationshipCommand", relationshipCommand,
			"add 396 -4131.291016 18 18 1 0 0 423 -4131.291016 18 18 1 0 0")
	}

	// Lock
	if len(referenceEdit.Locks) != 1 {
		t.Errorf(msg, "len(referenceEdit.Locks)", len(referenceEdit.Locks), 1)
	}
	lock := referenceEdit.Locks[0]
	if lock.Node != "|namespace:topNode" {
		t.Errorf(msg, "lock.Node", lock.Node, "|namespace:topNode")
	}
	if lock.Attr != "attrNameA" {
		t.Errorf(msg, "lock.Attr", lock.Attr, "attrNameA")
	}

	// Unlock
	if len(referenceEdit.Unlocks) != 1 {
		t.Errorf(msg, "len(referenceEdit.Unlocks)", len(referenceEdit.Unlocks), 1)
	}
	unlock := referenceEdit.Unlocks[0]
	if unlock.Node != "|namespace:topNode" {
		t.Errorf(msg, "unlock.Node", unlock.Node, "|namespace:topNode")
	}
	if unlock.Attr != "attrNameA" {
		t.Errorf(msg, "unlock.Attr", unlock.Attr, "attrNameA")
	}
}

func TestMakeAttributeAlias(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".aal" -type "attributeAlias" {"detonationFrame","borderConnections[0]","incandescence"
		,"borderConnections[1]","color","borderConnections[2]","nucleusSolver","publishedNodeInfo[0]"
		} ;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeAttributeAlias {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeAttributeAlias)
	}
	ret, err := cmd.ToAttrAttributeAlias(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if (*ret[0]).NewAlias != "detonationFrame" ||
		(*ret[0]).CurrentName != "borderConnections[0]" ||
		(*ret[1]).NewAlias != "incandescence" ||
		(*ret[1]).CurrentName != "borderConnections[1]" ||
		(*ret[2]).NewAlias != "color" ||
		(*ret[2]).CurrentName != "borderConnections[2]" ||
		(*ret[3]).NewAlias != "nucleusSolver" ||
		(*ret[3]).CurrentName != "publishedNodeInfo[0]" {
		msg := `got SetAttr %s %s, wont %s`
		t.Errorf(msg, "Attr", ret, &[]cmd.AttrAttributeAlias{
			{
				NewAlias:    "detonationFrame",
				CurrentName: "borderConnections[0]",
			},
			{
				NewAlias:    "incandescence",
				CurrentName: "borderConnections[1]",
			},
			{
				NewAlias:    "color",
				CurrentName: "borderConnections[2]",
			},
			{
				NewAlias:    "nucleusSolver",
				CurrentName: "publishedNodeInfo[0]",
			},
		})
	}
	if len(ret) != 4 {
		t.Errorf(msg, "len(Attr)", len(ret), 4)
	}
}

func TestMakeComponentList(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".ics" -type "componentList" 2 "vtx[130]" "vtx[147]";`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeComponentList {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeComponentList)
	}
	ret, err := cmd.ToAttrComponentList(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if (*ret[0])[0] != "vtx[130]" ||
		(*ret[0])[1] != "vtx[147]" {
		msg := `got SetAttr %s %s, wont %s`
		t.Errorf(msg, "Attr", ret, &cmd.AttrComponentList{
			"vtx[130]",
			"vtx[147]",
		})
	}
	if len(*ret[0]) != 2 {
		t.Errorf(msg, "len(Attr)", len(*ret[0]), 2)
	}
}

func TestMakeCone(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".cone" -type "cone" 45.0 5.0;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeCone {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeCone)
	}
	ret, err := cmd.ToAttrCone(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if (*ret[0]).ConeAngle != 45.0 ||
		(*ret[0]).ConeCap != 5.0 {
		msg := `got SetAttr %s %s, wont %s`
		t.Errorf(msg, "Attr", sa.Attr, &[]cmd.AttrCone{
			{
				ConeAngle: 45.0,
				ConeCap:   5.0,
			},
		})
	}
	if len(ret) != 1 {
		t.Errorf(msg, "len(Attr)", len(ret), 1)
	}
}

func TestMakeDoubleArray(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".dd" -type "doubleArray" 7 -1 1 0 0 0.5 1 -0.11000000000000004 ;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeDoubleArray {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeDoubleArray)
	}
	ret, err := cmd.ToAttrDoubleArray(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if (*ret[0])[0] != -1 ||
		(*ret[0])[1] != 1 ||
		(*ret[0])[2] != 0 ||
		(*ret[0])[3] != 0 ||
		(*ret[0])[4] != 0.5 ||
		(*ret[0])[5] != 1 ||
		(*ret[0])[6] != -0.11000000000000004 {
		msg := `got SetAttr %s %s, wont %s`
		t.Errorf(msg, "Attr", sa.Attr, &cmd.AttrDoubleArray{
			-1, 1, 0, 0, 0.5, 1, -0.11000000000000004,
		})
	}
	if len(*ret[0]) != 7 {
		t.Errorf(msg, "len(Attr)", len(*ret[0]), 7)
	}
}

func TestMakeLattice(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".cc" -type "lattice" 2 2 2 8
	-0.5 -0.5 -0.5
	0.5 -0.5 -0.5
	-0.5 0.5 -0.5
	0.5 0.5 -0.5
	-0.5 -0.5 0.5
	0.5 -0.5 0.5
	-0.5 0.5 0.5
	0.5 0.5 0.5 ;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeLattice {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeLattice)
	}
	ret, err := cmd.ToAttrLattice(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if (*ret[0]).DivisionS != 2 ||
		(*ret[0]).DivisionT != 2 ||
		(*ret[0]).DivisionU != 2 ||
		(*ret[0]).Points[0].S != -0.5 ||
		(*ret[0]).Points[0].T != -0.5 ||
		(*ret[0]).Points[0].U != -0.5 ||
		(*ret[0]).Points[1].S != 0.5 ||
		(*ret[0]).Points[1].T != -0.5 ||
		(*ret[0]).Points[1].U != -0.5 ||
		(*ret[0]).Points[2].S != -0.5 ||
		(*ret[0]).Points[2].T != 0.5 ||
		(*ret[0]).Points[2].U != -0.5 ||
		(*ret[0]).Points[3].S != 0.5 ||
		(*ret[0]).Points[3].T != 0.5 ||
		(*ret[0]).Points[3].U != -0.5 ||
		(*ret[0]).Points[4].S != -0.5 ||
		(*ret[0]).Points[4].T != -0.5 ||
		(*ret[0]).Points[4].U != 0.5 ||
		(*ret[0]).Points[5].S != 0.5 ||
		(*ret[0]).Points[5].T != -0.5 ||
		(*ret[0]).Points[5].U != 0.5 ||
		(*ret[0]).Points[6].S != -0.5 ||
		(*ret[0]).Points[6].T != 0.5 ||
		(*ret[0]).Points[6].U != 0.5 ||
		(*ret[0]).Points[7].S != 0.5 ||
		(*ret[0]).Points[7].T != 0.5 ||
		(*ret[0]).Points[7].U != 0.5 {
		msg := `got SetAttr %s %s, wont %s`
		t.Errorf(msg, "Attr", sa.Attr, &[]cmd.AttrLattice{
			{
				DivisionS: 2,
				DivisionT: 2,
				DivisionU: 2,
				Points: []cmd.AttrLatticePoint{
					{S: -0.5, T: -0.5, U: -0.5},
					{S: 0.5, T: -0.5, U: -0.5},
					{S: -0.5, T: 0.5, U: -0.5},
					{S: 0.5, T: 0.5, U: -0.5},
					{S: -0.5, T: -0.5, U: 0.5},
					{S: 0.5, T: -0.5, U: 0.5},
					{S: -0.5, T: 0.5, U: 0.5},
					{S: 0.5, T: 0.5, U: 0.5},
				},
			},
		})
	}
	if len(ret) != 1 {
		t.Errorf(msg, "len(Attr)", len(ret), 1)
	}
}

func TestMakeNurbsCurve(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".cc" -type "nurbsCurve"
	3 1 0 no 3
	6 0 0 0 1 1 1
	4
	0 0 0
	0.33333333333333326 0 -0.33333333333333326
	0.66666666666666663 0 -0.66666666666666663
	1 0 -1
	;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeNurbsCurve {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeNurbsCurve)
	}
	ret, err := cmd.ToAttrNurbsCurve(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if (*ret[0]).Degree != 3 ||
		(*ret[0]).Spans != 1 ||
		(*ret[0]).Form != cmd.AttrFormOpen ||
		(*ret[0]).IsRational != false ||
		(*ret[0]).Dimension != 3 ||
		(*ret[0]).KnotValues[0] != 0 ||
		(*ret[0]).KnotValues[1] != 0 ||
		(*ret[0]).KnotValues[2] != 0 ||
		(*ret[0]).KnotValues[3] != 1 ||
		(*ret[0]).KnotValues[4] != 1 ||
		(*ret[0]).KnotValues[5] != 1 ||
		(*ret[0]).CvValues[0].X != 0 ||
		(*ret[0]).CvValues[0].Y != 0 ||
		*(*ret[0]).CvValues[0].Z != 0 ||
		(*ret[0]).CvValues[0].W != nil ||
		(*ret[0]).CvValues[1].X != 0.33333333333333326 ||
		(*ret[0]).CvValues[1].Y != 0 ||
		*(*ret[0]).CvValues[1].Z != -0.33333333333333326 ||
		(*ret[0]).CvValues[1].W != nil ||
		(*ret[0]).CvValues[2].X != 0.66666666666666663 ||
		(*ret[0]).CvValues[2].Y != 0 ||
		*(*ret[0]).CvValues[2].Z != -0.66666666666666663 ||
		(*ret[0]).CvValues[2].W != nil ||
		(*ret[0]).CvValues[3].X != 1 ||
		(*ret[0]).CvValues[3].Y != 0 ||
		*(*ret[0]).CvValues[3].Z != -1 ||
		(*ret[0]).CvValues[3].W != nil {
		msg := `got SetAttr %s %s, wont %s`
		zero := 0.0
		minus03 := -0.33333333333333326
		minus06 := -0.66666666666666663
		minus1 := -1.0
		t.Errorf(msg, "Attr", sa.Attr, &[]cmd.AttrNurbsCurve{
			{
				Degree:     3,
				Spans:      1,
				Form:       cmd.AttrFormOpen,
				IsRational: false,
				Dimension:  3,
				KnotValues: []float64{
					0, 0, 0, 1, 1, 1,
				},
				CvValues: []cmd.AttrCvValue{
					{
						X: 0,
						Y: 0,
						Z: &zero,
						W: nil,
					},
					{
						X: 0.33333333333333326,
						Y: 0,
						Z: &minus03,
						W: nil,
					},
					{
						X: 0.66666666666666663,
						Y: 0,
						Z: &minus06,
						W: nil,
					},
					{
						X: 1,
						Y: 0,
						Z: &minus1,
						W: nil,
					},
				},
			},
		})
	}
	if len(ret) != 1 {
		t.Errorf(msg, "len(Attr)", len(ret), 1)
	}
}

func TestMakeNurbsSurface(t *testing.T) {
	c := &cmd.CmdBuilder{}
	c.Append(`setAttr ".cc" -type "nurbsSurface"
		3 3 0 0 no
		6 0 0 0 1 1 1
		6 0 0 0 1 1 1
		16
		-0.5 -3.061616997868383e-17 0.5
		-0.5 -1.0205389992894611e-17 0.16666666666666669
		-0.5 1.0205389992894608e-17 -0.16666666666666663
		-0.5 3.061616997868383e-17 -0.5
		-0.16666666666666669 -3.061616997868383e-17 0.5
		-0.16666666666666669 -1.0205389992894611e-17 0.16666666666666669
		-0.16666666666666669 1.0205389992894608e-17 -0.16666666666666663
		-0.16666666666666669 3.061616997868383e-17 -0.5
		0.16666666666666663 -3.061616997868383e-17 0.5
		0.16666666666666663 -1.0205389992894611e-17 0.16666666666666669
		0.16666666666666663 1.0205389992894608e-17 -0.16666666666666663
		0.16666666666666663 3.061616997868383e-17 -0.5
		0.5 -3.061616997868383e-17 0.5
		0.5 -1.0205389992894611e-17 0.16666666666666669
		0.5 1.0205389992894608e-17 -0.16666666666666663
		0.5 3.061616997868383e-17 -0.5;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != cmd.TypeNurbsSurface {
		t.Errorf(msg, "AttrType", sa.AttrType, cmd.TypeNurbsSurface)
	}
	ret, err := cmd.ToAttrNurbsSurface(sa.Attr)
	if err != nil {
		t.Fatal(err)
	}
	if (*ret[0]).UDegree != 3 ||
		(*ret[0]).VDegree != 3 ||
		(*ret[0]).UForm != cmd.AttrFormOpen ||
		(*ret[0]).VForm != cmd.AttrFormOpen ||
		(*ret[0]).IsRational != false ||
		len((*ret[0]).UKnotValues) != 6 ||
		(*ret[0]).UKnotValues[0] != 0 ||
		(*ret[0]).UKnotValues[1] != 0 ||
		(*ret[0]).UKnotValues[2] != 0 ||
		(*ret[0]).UKnotValues[3] != 1 ||
		(*ret[0]).UKnotValues[4] != 1 ||
		(*ret[0]).UKnotValues[5] != 1 ||
		len((*ret[0]).VKnotValues) != 6 ||
		(*ret[0]).VKnotValues[0] != 0 ||
		(*ret[0]).VKnotValues[1] != 0 ||
		(*ret[0]).VKnotValues[2] != 0 ||
		(*ret[0]).VKnotValues[3] != 1 ||
		(*ret[0]).VKnotValues[4] != 1 ||
		(*ret[0]).VKnotValues[5] != 1 ||
		(*ret[0]).IsTrim != nil ||
		len((*ret[0]).CvValues) != 16 ||
		(*ret[0]).CvValues[0].X != -0.5 ||
		(*ret[0]).CvValues[0].Y != -3.061616997868383e-17 ||
		(*ret[0]).CvValues[0].Z == nil ||
		*(*ret[0]).CvValues[0].Z != 0.5 ||
		(*ret[0]).CvValues[0].W != nil ||
		(*ret[0]).CvValues[1].X != -0.5 ||
		(*ret[0]).CvValues[1].Y != -1.0205389992894611e-17 ||
		(*ret[0]).CvValues[1].Z == nil ||
		*(*ret[0]).CvValues[1].Z != 0.16666666666666669 ||
		(*ret[0]).CvValues[1].W != nil ||
		(*ret[0]).CvValues[2].X != -0.5 ||
		(*ret[0]).CvValues[2].Y != 1.0205389992894608e-17 ||
		(*ret[0]).CvValues[2].Z == nil ||
		*(*ret[0]).CvValues[2].Z != -0.16666666666666663 ||
		(*ret[0]).CvValues[2].W != nil ||
		(*ret[0]).CvValues[3].X != -0.5 ||
		(*ret[0]).CvValues[3].Y != 3.061616997868383e-17 ||
		(*ret[0]).CvValues[3].Z == nil ||
		*(*ret[0]).CvValues[3].Z != -0.5 ||
		(*ret[0]).CvValues[3].W != nil ||
		(*ret[0]).CvValues[4].X != -0.16666666666666669 ||
		(*ret[0]).CvValues[4].Y != -3.061616997868383e-17 ||
		(*ret[0]).CvValues[4].Z == nil ||
		*(*ret[0]).CvValues[4].Z != 0.5 ||
		(*ret[0]).CvValues[4].W != nil ||
		(*ret[0]).CvValues[5].X != -0.16666666666666669 ||
		(*ret[0]).CvValues[5].Y != -1.0205389992894611e-17 ||
		(*ret[0]).CvValues[5].Z == nil ||
		*(*ret[0]).CvValues[5].Z != 0.16666666666666669 ||
		(*ret[0]).CvValues[5].W != nil ||
		(*ret[0]).CvValues[6].X != -0.16666666666666669 ||
		(*ret[0]).CvValues[6].Y != 1.0205389992894608e-17 ||
		(*ret[0]).CvValues[6].Z == nil ||
		*(*ret[0]).CvValues[6].Z != -0.16666666666666663 ||
		(*ret[0]).CvValues[6].W != nil ||
		(*ret[0]).CvValues[7].X != -0.16666666666666669 ||
		(*ret[0]).CvValues[7].Y != 3.061616997868383e-17 ||
		(*ret[0]).CvValues[7].Z == nil ||
		*(*ret[0]).CvValues[7].Z != -0.5 ||
		(*ret[0]).CvValues[7].W != nil ||
		(*ret[0]).CvValues[8].X != 0.16666666666666663 ||
		(*ret[0]).CvValues[8].Y != -3.061616997868383e-17 ||
		(*ret[0]).CvValues[8].Z == nil ||
		*(*ret[0]).CvValues[8].Z != 0.5 ||
		(*ret[0]).CvValues[8].W != nil ||
		(*ret[0]).CvValues[9].X != 0.16666666666666663 ||
		(*ret[0]).CvValues[9].Y != -1.0205389992894611e-17 ||
		(*ret[0]).CvValues[9].Z == nil ||
		*(*ret[0]).CvValues[9].Z != 0.16666666666666669 ||
		(*ret[0]).CvValues[9].W != nil ||
		(*ret[0]).CvValues[10].X != 0.16666666666666663 ||
		(*ret[0]).CvValues[10].Y != 1.0205389992894608e-17 ||
		(*ret[0]).CvValues[10].Z == nil ||
		*(*ret[0]).CvValues[10].Z != -0.16666666666666663 ||
		(*ret[0]).CvValues[10].W != nil ||
		(*ret[0]).CvValues[11].X != 0.16666666666666663 ||
		(*ret[0]).CvValues[11].Y != 3.061616997868383e-17 ||
		(*ret[0]).CvValues[11].Z == nil ||
		*(*ret[0]).CvValues[11].Z != -0.5 ||
		(*ret[0]).CvValues[11].W != nil ||
		(*ret[0]).CvValues[12].X != 0.5 ||
		(*ret[0]).CvValues[12].Y != -3.061616997868383e-17 ||
		(*ret[0]).CvValues[12].Z == nil ||
		*(*ret[0]).CvValues[12].Z != 0.5 ||
		(*ret[0]).CvValues[12].W != nil ||
		(*ret[0]).CvValues[13].X != 0.5 ||
		(*ret[0]).CvValues[13].Y != -1.0205389992894611e-17 ||
		(*ret[0]).CvValues[13].Z == nil ||
		*(*ret[0]).CvValues[13].Z != 0.16666666666666669 ||
		(*ret[0]).CvValues[13].W != nil ||
		(*ret[0]).CvValues[14].X != 0.5 ||
		(*ret[0]).CvValues[14].Y != 1.0205389992894608e-17 ||
		(*ret[0]).CvValues[14].Z == nil ||
		*(*ret[0]).CvValues[14].Z != -0.16666666666666663 ||
		(*ret[0]).CvValues[14].W != nil ||
		(*ret[0]).CvValues[15].X != 0.5 ||
		(*ret[0]).CvValues[15].Y != 3.061616997868383e-17 ||
		(*ret[0]).CvValues[15].Z == nil ||
		*(*ret[0]).CvValues[15].Z != -0.5 ||
		(*ret[0]).CvValues[15].W != nil {
		msg := `got SetAttr %s %s, wont %s`
		plus05 := 0.5
		plus016 := 0.16666666666666669
		minus016 := -0.16666666666666663
		minus05 := -0.5
		t.Errorf(msg, "Attr", sa.Attr, &[]cmd.AttrNurbsSurface{
			{
				UDegree:    3,
				VDegree:    3,
				UForm:      cmd.AttrFormOpen,
				VForm:      cmd.AttrFormOpen,
				IsRational: false,
				UKnotValues: []float64{
					0, 0, 0, 1, 1, 1,
				},
				VKnotValues: []float64{
					0, 0, 0, 1, 1, 1,
				},
				IsTrim: nil,
				CvValues: []cmd.AttrCvValue{
					{X: -0.5, Y: -3.061616997868383e-17, Z: &plus05, W: nil},
					{X: -0.5, Y: -1.0205389992894611e-17, Z: &plus016, W: nil},
					{X: -0.5, Y: 1.0205389992894611e-17, Z: &minus016, W: nil},
					{X: -0.5, Y: 3.061616997868383e-17, Z: &minus05, W: nil},
					{X: -0.16666666666666669, Y: -3.061616997868383e-17, Z: &plus05, W: nil},
					{X: -0.16666666666666669, Y: -1.0205389992894611e-17, Z: &plus016, W: nil},
					{X: -0.16666666666666669, Y: 1.0205389992894611e-17, Z: &minus016, W: nil},
					{X: -0.16666666666666669, Y: 3.061616997868383e-17, Z: &minus05, W: nil},
					{X: 0.16666666666666663, Y: -3.061616997868383e-17, Z: &plus05, W: nil},
					{X: 0.16666666666666663, Y: -1.0205389992894611e-17, Z: &plus016, W: nil},
					{X: 0.16666666666666663, Y: 1.0205389992894611e-17, Z: &minus016, W: nil},
					{X: 0.16666666666666663, Y: 3.061616997868383e-17, Z: &minus05, W: nil},
					{X: 0.5, Y: -3.061616997868383e-17, Z: &plus05, W: nil},
					{X: 0.5, Y: -1.0205389992894611e-17, Z: &plus016, W: nil},
					{X: 0.5, Y: 1.0205389992894611e-17, Z: &minus016, W: nil},
					{X: 0.5, Y: 3.061616997868383e-17, Z: &minus05, W: nil},
				},
			},
		})
	}
	if len(ret) != 1 {
		t.Errorf(msg, "len(Attr)", len(ret), 1)
	}
}

func TestMakeNurbsTrimface(t *testing.T) {

}
