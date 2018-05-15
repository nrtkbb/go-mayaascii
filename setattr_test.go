package main

import (
	"testing"
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
	c := &CmdBuilder{}
	c.Append(`setAttr -s 4 ".attrName";`)
	beforeSetAttr, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if beforeSetAttr.AttrType != TypeInvalid {
		t.Errorf(msg, "AttrType", beforeSetAttr.AttrType, TypeInvalid)
	}
	if *beforeSetAttr.Size != uint(4) {
		t.Errorf(msg, "Size", *beforeSetAttr.Size, uint(4))
	}
}

func TestMakeSetAttr_int(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" 1 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeInt {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeInt)
	}
	i, ok := sa.Attr.(*[]int)
	if len(*i) != 2 {
		t.Errorf(msg, "len(Attr)", len(*i), 2)
	}
	if !ok || (*i)[0] != 1 || (*i)[1] != 2 {
		t.Errorf(msg, "Attr", sa.Attr, []int{1, 2})
	}
	c.Clear()
	c.Append(`setAttr ".attrName" 3 4;`)
	sa, err = MakeSetAttr(c.Parse(), sa)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != TypeInt {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeInt)
	}
	if len(*i) != 4 {
		t.Errorf(msg, "len(Attr)", len(*i), 4)
	}
	if (*i)[0] != 1 || (*i)[1] != 2 ||
		(*i)[2] != 3 || (*i)[3] != 4 {
		t.Errorf(msg, "Attr", sa.Attr, []int{1, 2, 3, 4})
	}
}

func TestMakeSetAttr_int_toDouble_toInt(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" 1 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeInt {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeInt)
	}
	i, ok := sa.Attr.(*[]int)
	if len(*i) != 2 {
		t.Errorf(msg, "len(Attr)", len(*i), 2)
	}
	if !ok || (*i)[0] != 1 || (*i)[1] != 2 {
		t.Errorf(msg, "Attr", sa.Attr, []int{1, 2})
	}
	c.Clear()
	c.Append(`setAttr ".attrName" 3.3 4e+020 5e-020;`)
	sa, err = MakeSetAttr(c.Parse(), sa)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != TypeDouble {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDouble)
	}
	f, ok := sa.Attr.(*[]float64)
	if len(*f) != 5 {
		t.Errorf(msg, "len(Attr)", len(*f), 5)
	}
	if (*f)[0] != 1 || (*f)[1] != 2 ||
		(*f)[2] != 3.3 || (*f)[3] != 4E+020 || (*f)[4] != 5E-020 {
		t.Errorf(msg, "Attr", sa.Attr, []float64{
			1, 2, 3.3, 4E+020, 5E-020})
	}
	c.Clear()
	c.Append(`setAttr ".attrName" 5 6;`)
	sa, err = MakeSetAttr(c.Parse(), sa)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != TypeDouble {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDouble)
	}
	f, ok = sa.Attr.(*[]float64)
	if len(*f) != 7 {
		t.Errorf(msg, "len(Attr)", len(*f), 7)
	}
	if (*f)[0] != 1 || (*f)[1] != 2 ||
		(*f)[2] != 3.3 || (*f)[3] != 4E+020 || (*f)[4] != 5E-020 ||
		(*f)[5] != 5 || (*f)[6] != 6 {
		t.Errorf(msg, "Attr", sa.Attr, []float64{
			1, 2, 3.3, 4E+20, 5E-20, 5, 6})
	}
}

func TestMakeSetAttr_doubleWithExponent(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" 1e+020 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeDouble {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDouble)
	}
	f, ok := sa.Attr.(*[]float64)
	if len(*f) != 2 {
		t.Errorf(msg, "len(Attr)", len(*f), 2)
	}
	if !ok || (*f)[0] != 1E+020 || (*f)[1] != 2 {
		t.Errorf(msg, "Attr", sa.Attr, []float64{1E+020, 2})
	}
}

func TestMakeSetAttr_double(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" 1.1 2.2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeDouble {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDouble)
	}
	f, ok := sa.Attr.(*[]float64)
	if len(*f) != 2 {
		t.Errorf(msg, "len(Attr)", len(*f), 2)
	}
	if !ok || (*f)[0] != 1.1 || (*f)[1] != 2.2 {
		t.Errorf(msg, "Attr", sa.Attr, []float64{1.1, 2.2})
	}
	c.Clear()
	c.Append(`setAttr ".attrName" 3.3 4.4;`)
	sa, err = MakeSetAttr(c.Parse(), sa)
	if sa.AttrType != TypeDouble {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDouble)
	}
	if len(*f) != 4 {
		t.Errorf(msg, "len(Attr)", len(*f), 4)
	}
	if (*f)[0] != 1.1 || (*f)[1] != 2.2 || (*f)[2] != 3.3 || (*f)[3] != 4.4 {
		t.Errorf(msg, "Attr", sa.Attr, []float64{1.1, 2.2, 3.3, 4.4})
	}
}

func testBool(t *testing.T, boolString string, wont bool) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" ` + boolString + `;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeBool {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeBool)
	}
	if b, ok := sa.Attr.(*bool); !ok || *b != wont {
		t.Errorf(msg, "Attr", sa.Attr, wont)
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
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "short2" 1 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeShort2 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeShort2)
	}
	s2, ok := sa.Attr.(*[]AttrShort2)
	if len(*s2) != 1 {
		t.Errorf(msg, "len(Attr)", len(*s2), 1)
	}
	if !ok || (*s2)[0][0] != 1 || (*s2)[0][1] != 2 {
		t.Errorf(msg, "Attr", sa.Attr, []AttrShort2{
			{1, 2},
		})
	}
}

func TestMakeSetAttr_short2_add(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "short2" 1 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	c.Append(`setAttr ".attrName" -type "short2" 3 4;`)
	sa, err = MakeSetAttr(c.Parse(), sa)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	s2, ok := sa.Attr.(*[]AttrShort2)
	if len(*s2) != 2 {
		t.Errorf(msg, "len(Attr)", len(*s2), 2)
	}
	if !ok ||
		(*s2)[0][0] != 1 || (*s2)[0][1] != 2 ||
		(*s2)[1][0] != 3 || (*s2)[1][1] != 4 {
		t.Errorf(msg, "Attr", sa.Attr, []AttrShort2{
			{1, 2},
			{3, 4},
		})
	}
}

func TestMakeSetAttr_short2_size(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr -s 2 ".attrName" -type "short2" 1 2 1 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeShort2 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeShort2)
	}
	s2, ok := sa.Attr.(*[]AttrShort2)
	if len(*s2) != 2 {
		t.Errorf(msg, "len(Attr)", len(*s2), 2)
	}
	if !ok ||
		(*s2)[0][0] != 1 ||
		(*s2)[0][1] != 2 ||
		(*s2)[1][0] != 1 ||
		(*s2)[1][1] != 2 {
		t.Errorf(msg, "Attr", sa.Attr, []AttrShort2{
			{1, 2},
			{1, 2},
		})
	}
}

func TestMakeSetAttr_short2_sizeOver(t *testing.T) {
	c := &CmdBuilder{}
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
	if sa.AttrType != TypeShort2 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeShort2)
	}
	s2, ok := sa.Attr.(*[]AttrShort2)
	if len(*s2) != 4 {
		t.Errorf(msg, "len(Attr)", len(*s2), 4)
	}
	if !ok ||
		(*s2)[0][0] != 1 ||
		(*s2)[0][1] != 2 ||
		(*s2)[1][0] != 1 ||
		(*s2)[1][1] != 2 ||
		(*s2)[2][0] != 1 ||
		(*s2)[2][1] != 2 ||
		(*s2)[3][0] != 1 ||
		(*s2)[3][1] != 2 {
		t.Errorf(msg, "Attr", sa.Attr, []AttrShort2{
			{1, 2},
			{1, 2},
			{1, 2},
			{1, 2},
		})
	}
}

func TestMakeSetAttr_long2(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "long2" 1 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeLong2 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeLong2)
	}
	l2, ok := sa.Attr.(*[]AttrLong2)
	if !ok || (*l2)[0][0] != 1 || (*l2)[0][1] != 2 {
		t.Errorf(msg, "Attr", sa.Attr, [][2]int{{1, 2}})
	}
	if len(*l2) != 1 {
		t.Errorf(msg, "len(Attr)", len(*l2), 1)
	}
	c.Clear()
	c.Append(`setAttr -s 2 ".attrName" -type "long2" 1 2 1 2;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != TypeLong2 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeLong2)
	}
	l2, ok = sa.Attr.(*[]AttrLong2)
	if !ok ||
		(*l2)[0][0] != 1 ||
		(*l2)[0][1] != 2 ||
		(*l2)[1][0] != 1 ||
		(*l2)[1][1] != 2 {
		t.Errorf(msg, "Attr", sa.Attr, [][2]int{
			{1, 2},
			{1, 2},
		})
	}
	if len(*l2) != 2 {
		t.Errorf(msg, "len(Attr)", len(*l2), 2)
	}
}

func TestMakeSetAttr_short3(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "short3" 1 2 3;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeShort3 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeShort3)
	}
	s3, ok := sa.Attr.(*[]AttrShort3)
	if !ok || (*s3)[0][0] != 1 || (*s3)[0][1] != 2 || (*s3)[0][2] != 3 {
		t.Errorf(msg, "Attr", sa.Attr, []int{1, 2, 3})
	}
	if len(*s3) != 1 {
		t.Errorf(msg, "len(Attr)", len(*s3), 1)
	}
	c.Append(`setAttr -s 2 ".attrName" -type "short3" 1 2 3 1 2 3;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != TypeShort3 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeShort3)
	}
	s3, ok = sa.Attr.(*[]AttrShort3)
	if !ok ||
		(*s3)[0][0] != 1 ||
		(*s3)[0][1] != 2 ||
		(*s3)[0][2] != 3 ||
		(*s3)[1][0] != 1 ||
		(*s3)[1][1] != 2 ||
		(*s3)[1][2] != 3 {
		t.Errorf(msg, "Attr", sa.Attr, [][3]int{
			{1, 2, 3},
			{1, 2, 3},
		})
	}
	if len(*s3) != 2 {
		t.Errorf(msg, "len(Attr)", len(*s3), 2)
	}
}

func TestMakeSetAttr_long3(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "long3" 1 2 3;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeLong3 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeLong3)
	}
	l3, ok := sa.Attr.(*[]AttrLong3)
	if !ok || (*l3)[0][0] != 1 || (*l3)[0][1] != 2 || (*l3)[0][2] != 3 {
		t.Errorf(msg, "Attr", sa.Attr, [][3]int{
			{1, 2, 3},
		})
	}
	if len(*l3) != 1 {
		t.Errorf(msg, "len(Attr)", len(*l3), 1)
	}
	c.Clear()
	c.Append(`setAttr -s 2 ".attrName" -type "long3" 1 2 3 1 2 3;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != TypeLong3 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeLong3)
	}
	l3, ok = sa.Attr.(*[]AttrLong3)
	if !ok ||
		(*l3)[0][0] != 1 ||
		(*l3)[0][1] != 2 ||
		(*l3)[0][2] != 3 ||
		(*l3)[1][0] != 1 ||
		(*l3)[1][1] != 2 ||
		(*l3)[1][2] != 3 {
		t.Errorf(msg, "Attr", sa.Attr, [][3]int{
			{1, 2, 3},
			{1, 2, 3},
		})
	}
	if len(*l3) != 2 {
		t.Errorf(msg, "len(Attr)", len(*l3), 2)
	}
}

func TestMakeSetAttr_Int32Array(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "Int32Array" 2 1 2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeInt32Array {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeInt32Array)
	}
	i32a, ok := sa.Attr.(*AttrInt32Array)
	if !ok || (*i32a)[0] != 1 || (*i32a)[1] != 2 {
		t.Errorf(msg, "Attr", []int{(*i32a)[0], (*i32a)[1]}, []int{1, 2})
	}
	if len(*i32a) != 2 {
		t.Errorf(msg, "len(Attr)", len(*i32a), 2)
	}
}

func TestMakeSetAttr_float2(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "float2" 1 2.2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeFloat2 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeFloat2)
	}
	f2, ok := sa.Attr.(*[]AttrFloat2)
	if !ok || (*f2)[0][0] != 1.0 || (*f2)[0][1] != 2.2 {
		t.Errorf(msg, "Attr", sa.Attr, [][2]float64{
			{1, 2.2},
		})
	}
	if len(*f2) != 1 {
		t.Errorf(msg, "len(Attr)", len(*f2), 1)
	}
	c.Clear()
	c.Append(`setAttr -s 2 ".attrName" -type "float2" 1 2.2 1 2.2;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != TypeFloat2 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeFloat2)
	}
	f2, ok = sa.Attr.(*[]AttrFloat2)
	if !ok ||
		(*f2)[0][0] != 1.0 ||
		(*f2)[0][1] != 2.2 ||
		(*f2)[1][0] != 1.0 ||
		(*f2)[1][1] != 2.2 {
		t.Errorf(msg, "Attr", sa.Attr, [][2]float64{
			{1, 2.2,},
			{1, 2.2,},
		})
	}
	if len(*f2) != 2 {
		t.Errorf(msg, "len(Attr)", len(*f2), 2)
	}
}

func TestMakeSetAttr_float3(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "float3" 1 2.2 3.3;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeFloat3 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeFloat3)
	}
	f3, ok := sa.Attr.(*[]AttrFloat3)
	if !ok || (*f3)[0][0] != 1.0 || (*f3)[0][1] != 2.2 || (*f3)[0][2] != 3.3 {
		t.Errorf(msg, "Attr", sa.Attr, []float64{1, 2.2, 3.3})
	}
	if len(*f3) != 1 {
		t.Errorf(msg, "len(Attr)", len(*f3), 1)
	}
	c.Clear()
	c.Append(`setAttr -s 2 ".attrName" -type "float3" 1 2.2 3.3 1 2.2 3.3;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != TypeFloat3 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeFloat3)
	}
	f3, ok = sa.Attr.(*[]AttrFloat3)
	if !ok ||
		(*f3)[0][0] != 1.0 ||
		(*f3)[0][1] != 2.2 ||
		(*f3)[0][2] != 3.3 ||
		(*f3)[1][0] != 1.0 ||
		(*f3)[1][1] != 2.2 ||
		(*f3)[1][2] != 3.3 {
		t.Errorf(msg, "Attr", sa.Attr, [][3]float64{
			{1, 2.2, 3.3,},
			{1, 2.2, 3.3,},
		})
	}
	if len(*f3) != 2 {
		t.Errorf(msg, "len(Attr)", len(*f3), 2)
	}
}

func TestMakeSetAttr_double2(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "double2" 1 2.2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeDouble2 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDouble2)
	}
	d2, ok := sa.Attr.(*[]AttrDouble2)
	if !ok || (*d2)[0][0] != 1.0 || (*d2)[0][1] != 2.2 {
		t.Errorf(msg, "Attr", sa.Attr, []float64{1, 2.2})
	}
	if len(*d2) != 1 {
		t.Errorf(msg, "len(Attr)", len(*d2), 1)
	}
	c.Clear()
	c.Append(`setAttr -s 2 ".attrName" -type "double2" 1 2.2 1 2.2;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != TypeDouble2 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDouble2)
	}
	d2, ok = sa.Attr.(*[]AttrDouble2)
	if !ok ||
		(*d2)[0][0] != 1.0 ||
		(*d2)[0][1] != 2.2 ||
		(*d2)[1][0] != 1.0 ||
		(*d2)[1][1] != 2.2 {
		t.Errorf(msg, "Attr", sa.Attr, [][2]float64{
			{1, 2.2,},
			{1, 2.2,},
		})
	}
	if len(*d2) != 2 {
		t.Errorf(msg, "len(Attr)", len(*d2), 2)
	}
}

func TestMakeSetAttr_double3(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "double3" 1 2.2 3.3;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeDouble3 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDouble3)
	}
	d3, ok := sa.Attr.(*[]AttrDouble3)
	if !ok || (*d3)[0][0] != 1.0 || (*d3)[0][1] != 2.2 || (*d3)[0][2] != 3.3 {
		t.Errorf(msg, "Attr", sa.Attr, []float64{1, 2.2, 3.3})
	}
	if len(*d3) != 1 {
		t.Errorf(msg, "len(Attr)", len(*d3), 1)
	}
	c.Clear()
	c.Append(`setAttr -s 2 ".attrName" -type "double3" 1 2.2 3.3 1 2.2 3.3;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != TypeDouble3 {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDouble3)
	}
	d3, ok = sa.Attr.(*[]AttrDouble3)
	if !ok ||
		(*d3)[0][0] != 1.0 ||
		(*d3)[0][1] != 2.2 ||
		(*d3)[0][2] != 3.3 ||
		(*d3)[1][0] != 1.0 ||
		(*d3)[1][1] != 2.2 ||
		(*d3)[1][2] != 3.3 {
		t.Errorf(msg, "Attr", sa.Attr, [][3]float64{
			{1, 2.2, 3.3,},
			{1, 2.2, 3.3,},
		})
	}
	if len(*d3) != 2 {
		t.Errorf(msg, "len(Attr)", len(*d3), 2)
	}
}

func TestMakeSetAttr_doubleArray(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "doubleArray" 2 1.1 2.2;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeDoubleArray {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDoubleArray)
	}
	da, ok := sa.Attr.(*AttrDoubleArray)
	if !ok || (*da)[0] != 1.1 || (*da)[1] != 2.2 {
		t.Errorf(msg, "Attr", sa.Attr, []float64{1, 2.2, 3.3})
	}
	if len(*da) != 2 {
		t.Errorf(msg, "len(Attr)", len(*da), 2)
	}
	c.Clear()
	c.Append(`setAttr ".attrName" -type "doubleArray" 0;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != TypeDoubleArray {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDoubleArray)
	}
	da, ok = sa.Attr.(*AttrDoubleArray)
	if !ok {
		t.Errorf(msg, "Attr", sa.Attr, []float64{})
	}
	if len(*da) != 0 {
		t.Errorf(msg, "len(Attr)", len(*da), 0)
	}
}

func TestMakeSetAttr_matrix(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".ix" -type "matrix" 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeMatrix {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeMatrix)
	}
	mt, ok := sa.Attr.(*AttrMatrix)
	wontMt := [16]float64{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}
	allOk := true
	for i, v := range *mt {
		if v != wontMt[i] {
			allOk = false
			break
		}
	}
	if !ok || !allOk {
		t.Errorf(msg, "Attr", sa.Attr, wontMt)
	}
	if len(*mt) != 16 {
		t.Errorf(msg, "len(Attr)", len(*mt), 16)
	}
}

func TestMakeSetAttr_matrix_xform(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".ix" -type "matrix" "xform" 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
0 0 0 0 0 0 0 0 0 0 1 0 0 0 1 1 1 1 yes;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeMatrixXform {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeMatrixXform)
	}
	mtx, ok := sa.Attr.(*AttrMatrixXform)
	if !ok ||
		mtx.Scale.X != 1 || mtx.Scale.Y != 1 || mtx.Scale.Z != 1 ||
		mtx.Rotate.X != 0 || mtx.Rotate.Y != 0 || mtx.Rotate.Z != 0 ||
		mtx.RotateOrder != RotateOrderXYZ ||
		mtx.Translate.X != 0 || mtx.Translate.Y != 0 || mtx.Translate.Z != 0 ||
		mtx.Shear.XY != 0 || mtx.Shear.XZ != 0 || mtx.Shear.YZ != 0 ||
		mtx.ScalePivot.X != 0 || mtx.ScalePivot.Y != 0 || mtx.ScalePivot.Z != 0 ||
		mtx.ScaleTranslate.X != 0 || mtx.ScaleTranslate.Y != 0 || mtx.ScaleTranslate.Z != 0 ||
		mtx.RotatePivot.X != 0 || mtx.RotatePivot.Y != 0 || mtx.RotatePivot.Z != 0 ||
		mtx.RotateTranslation.X != 0 || mtx.RotateTranslation.Y != 0 || mtx.RotateTranslation.Z != 0 ||
		mtx.RotateOrient.W != 0 || mtx.RotateOrient.X != 0 || mtx.RotateOrient.Y != 0 || mtx.RotateOrient.Z != 1 ||
		mtx.JointOrient.W != 0 || mtx.JointOrient.X != 0 || mtx.JointOrient.Y != 0 || mtx.JointOrient.Z != 1 ||
		mtx.InverseParentScale.X != 1 || mtx.InverseParentScale.Y != 1 || mtx.InverseParentScale.Z != 1 ||
		mtx.CompensateForParentScale == false {
		t.Errorf(msg, "Attr", sa.Attr, nil)
	}
}

func TestMakeSetAttr_pointArray(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".attrName" -type "pointArray" 1 1.1 2.2 3.3 4.4;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypePointArray {
		t.Errorf(msg, "AttrType", sa.AttrType, TypePointArray)
	}
	pa, ok := sa.Attr.(*AttrPointArray)
	if !ok || (*pa)[0].X != 1.1 || (*pa)[0].Y != 2.2 || (*pa)[0].Z != 3.3 || (*pa)[0].W != 4.4 {
		t.Errorf(msg, "Attr", sa.Attr, AttrPointArray{
			{1.1, 2.2, 3.3, 4.4,},
		})
	}
	if len(*pa) != 1 {
		t.Errorf(msg, "len(Attr)", len(*pa), 1)
	}
	c.Clear()
	c.Append(`setAttr ".attrName" -type "pointArray" 2 1.1 2.2 3.3 4.4 1.1 2.2 3.3 4.4;`)
	sa, err = MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sa.AttrType != TypePointArray {
		t.Errorf(msg, "AttrType", sa.AttrType, TypePointArray)
	}
	pa, ok = sa.Attr.(*AttrPointArray)
	if !ok ||
		(*pa)[0].X != 1.1 ||
		(*pa)[0].Y != 2.2 ||
		(*pa)[0].Z != 3.3 ||
		(*pa)[0].W != 4.4 ||
		(*pa)[1].X != 1.1 ||
		(*pa)[1].Y != 2.2 ||
		(*pa)[1].Z != 3.3 ||
		(*pa)[1].W != 4.4 {
		t.Errorf(msg, "Attr", sa.Attr, AttrPointArray{
			{1.1, 2.2, 3.3, 4.4,},
			{1.1, 2.2, 3.3, 4.4,},
		})
	}
	if len(*pa) != 2 {
		t.Errorf(msg, "len(Attr)", len(*pa), 2)
	}
}

func TestMakeSetAttr_polyFaces(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr -s 2 ".attrName" -type "polyFaces" f 3 1 2 3 mc 1 3 0 1 2 f 3 2 3 4 mc 2 3 2 3 4;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypePolyFaces {
		t.Errorf(msg, "AttrType", sa.AttrType, TypePolyFaces)
	}
	pfs, ok := sa.Attr.(*[]AttrPolyFaces)
	if !ok ||
		(*pfs)[0].FaceEdge[0] != 1 ||
		(*pfs)[0].FaceEdge[1] != 2 ||
		(*pfs)[0].FaceEdge[2] != 3 ||
		(*pfs)[0].MultiColor[0].ColorIndex != 1 ||
		(*pfs)[0].MultiColor[0].ColorIDs[0] != 0 ||
		(*pfs)[0].MultiColor[0].ColorIDs[1] != 1 ||
		(*pfs)[0].MultiColor[0].ColorIDs[2] != 2 ||
		(*pfs)[1].FaceEdge[0] != 2 ||
		(*pfs)[1].FaceEdge[1] != 3 ||
		(*pfs)[1].FaceEdge[2] != 4 ||
		(*pfs)[1].MultiColor[0].ColorIndex != 2 ||
		(*pfs)[1].MultiColor[0].ColorIDs[0] != 2 ||
		(*pfs)[1].MultiColor[0].ColorIDs[1] != 3 ||
		(*pfs)[1].MultiColor[0].ColorIDs[2] != 4 {
		t.Errorf(msg, "Attr", sa.Attr, []AttrPolyFaces{
			{
				FaceEdge: []int{1, 2, 3},
				MultiColor: []AttrMultiColor{
					{
						ColorIndex: 1,
						ColorIDs:   []int{1, 2, 3},
					},
				},
			},
			{
				FaceEdge: []int{2, 3, 4},
				MultiColor: []AttrMultiColor{
					{
						ColorIndex: 2,
						ColorIDs:   []int{2, 3, 4},
					},
				},
			},
		})
	}
	if len(*pfs) != 2 {
		t.Errorf(msg, "len(Attr)", len(*pfs), 2)
	}
}

func TestMakeSetAttr_polyFacesMax(t *testing.T) {
	c := &CmdBuilder{}
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
	if sa.AttrType != TypePolyFaces {
		t.Errorf(msg, "AttrType", sa.AttrType, TypePolyFaces)
	}
	pfs, ok := sa.Attr.(*[]AttrPolyFaces)
	if !ok ||
		(*pfs)[0].FaceEdge[0] != 1 ||
		(*pfs)[0].FaceEdge[1] != 2 ||
		(*pfs)[0].FaceEdge[2] != 3 ||
		(*pfs)[0].HoleEdge[0] != 5 ||
		(*pfs)[0].HoleEdge[1] != 6 ||
		(*pfs)[0].HoleEdge[2] != 7 ||
		(*pfs)[0].FaceUV[0].UVSet != 0 ||
		(*pfs)[0].FaceUV[0].FaceUV[0] != 0 ||
		(*pfs)[0].FaceUV[0].FaceUV[1] != 1 ||
		(*pfs)[0].FaceUV[0].FaceUV[2] != 3 ||
		(*pfs)[0].FaceUV[1].UVSet != 1 ||
		(*pfs)[0].FaceUV[1].FaceUV[0] != 0 ||
		(*pfs)[0].FaceUV[1].FaceUV[1] != 1 ||
		(*pfs)[0].FaceUV[1].FaceUV[2] != 3 ||
		(*pfs)[0].MultiColor[0].ColorIndex != 1 ||
		(*pfs)[0].MultiColor[0].ColorIDs[0] != 0 ||
		(*pfs)[0].MultiColor[0].ColorIDs[1] != 1 ||
		(*pfs)[0].MultiColor[0].ColorIDs[2] != 2 ||
		(*pfs)[1].FaceEdge[0] != 2 ||
		(*pfs)[1].FaceEdge[1] != 3 ||
		(*pfs)[1].FaceEdge[2] != 4 ||
		(*pfs)[1].FaceUV[0].UVSet != 0 ||
		(*pfs)[1].FaceUV[0].FaceUV[0] != 2 ||
		(*pfs)[1].FaceUV[0].FaceUV[1] != 3 ||
		(*pfs)[1].FaceUV[0].FaceUV[2] != 4 ||
		(*pfs)[1].FaceUV[1].UVSet != 1 ||
		(*pfs)[1].FaceUV[1].FaceUV[0] != 2 ||
		(*pfs)[1].FaceUV[1].FaceUV[1] != 3 ||
		(*pfs)[1].FaceUV[1].FaceUV[2] != 4 ||
		(*pfs)[1].MultiColor[0].ColorIndex != 2 ||
		(*pfs)[1].MultiColor[0].ColorIDs[0] != 2 ||
		(*pfs)[1].MultiColor[0].ColorIDs[1] != 3 ||
		(*pfs)[1].MultiColor[0].ColorIDs[2] != 4 {
		t.Errorf(msg, "Attr", sa.Attr, []AttrPolyFaces{
			{
				FaceEdge: []int{1, 2, 3},
				HoleEdge: []int{5, 6, 7},
				FaceUV: []AttrFaceUV{
					{
						UVSet:  0,
						FaceUV: []int{0, 1, 3},
					},
					{
						UVSet:  1,
						FaceUV: []int{0, 1, 3},
					},
				},
				MultiColor: []AttrMultiColor{
					{
						ColorIndex: 0,
						ColorIDs:   []int{0, 1, 2},
					},
				},
			},
			{
				FaceEdge: []int{2, 3, 4},
				FaceUV: []AttrFaceUV{
					{
						UVSet:  0,
						FaceUV: []int{2, 3, 4},
					},
					{
						UVSet:  1,
						FaceUV: []int{2, 3, 4},
					},
				},
				MultiColor: []AttrMultiColor{
					{
						ColorIndex: 2,
						ColorIDs:   []int{2, 3, 4},
					},
				},
			},
		})
	}
	if len(*pfs) != 2 {
		t.Errorf(msg, "len(Attr)", len(*pfs), 2)
	}
}

func TestMakeDataPolyComponent(t *testing.T) {
	c := &CmdBuilder{}
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
	if sa.AttrType != TypeDataPolyComponent {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDataPolyComponent)
	}
	dpc, ok := sa.Attr.(*[]AttrDataPolyComponent)
	if !ok ||
		(*dpc)[0].PolyComponentType != DPCedge ||
		(*dpc)[0].IndexValue[0] != 9.0600013732910156 ||
		(*dpc)[0].IndexValue[1] != 4.409998893737793 ||
		(*dpc)[0].IndexValue[2] != 9.0600013732910156 ||
		(*dpc)[0].IndexValue[3] != 9.2000007629394531 ||
		(*dpc)[0].IndexValue[4] != 9.0600013732910156 ||
		(*dpc)[0].IndexValue[5] != 9.0600013732910156 ||
		(*dpc)[0].IndexValue[9] != 9.0600013732910156 ||
		(*dpc)[0].IndexValue[10] != 9.0600013732910156 ||
		(*dpc)[0].IndexValue[11] != 9.0600013732910156 ||
		(*dpc)[0].IndexValue[12] != 9.0600013732910156 ||
		(*dpc)[0].IndexValue[13] != 4.409998893737793 ||
		(*dpc)[0].IndexValue[15] != 4.6099758148193359 ||
		(*dpc)[0].IndexValue[18] != 4.409998893737793 ||
		(*dpc)[0].IndexValue[19] != 4.6099758148193359 ||
		(*dpc)[0].IndexValue[20] != 4.409998893737793 ||
		(*dpc)[0].IndexValue[21] != 9.2000007629394531 ||
		(*dpc)[0].IndexValue[23] != 9.2000007629394531 ||
		(*dpc)[0].IndexValue[26] != 9.2000007629394531 ||
		(*dpc)[0].IndexValue[27] != 9.2000007629394531 ||
		(*dpc)[0].IndexValue[28] != 9.2000007629394531 ||
		(*dpc)[0].IndexValue[30] != 4.6099758148193359 ||
		(*dpc)[0].IndexValue[32] != 4.6099758148193359 ||
		(*dpc)[0].IndexValue[34] != 4.6099758148193359 ||
		(*dpc)[0].IndexValue[35] != 4.6099758148193359 {
		t.Errorf(msg, "Attr", dpc, []AttrDataPolyComponent{
			{
				PolyComponentType: DPCedge,
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
	if len(*dpc) != 1 {
		t.Errorf(msg, "len(Attr)", len(*dpc), 1)
	}
	if len((*dpc)[0].IndexValue) != 24 {
		t.Errorf(msg, "len(Attr.IndexValue)", len((*dpc)[0].IndexValue), 24)
	}
}

func sameDPC(t *testing.T, ok bool, dpc *[]AttrDataPolyComponent, dpcType AttrDPCType) {
	if !ok || (*dpc)[0].PolyComponentType != dpcType {
		msg := `got SetAttr %s %s, wont %s`
		t.Errorf(msg, "Attr", dpc, []AttrDataPolyComponent{
			{
				PolyComponentType: dpcType,
				IndexValue:        map[int]float64{},
			},
		})
	}
}

func TestMakeDataPolyComponentVertex(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".cvd" -type "dataPolyComponent" Index_Data Vertex 0;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeDataPolyComponent {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDataPolyComponent)
	}
	dpc, ok := sa.Attr.(*[]AttrDataPolyComponent)
	sameDPC(t, ok, dpc, DPCvertex)
	if len(*dpc) != 1 {
		t.Errorf(msg, "len(Attr)", len(*dpc), 1)
	}
	if len((*dpc)[0].IndexValue) != 0 {
		t.Errorf(msg, "len(Attr.IndexValue)", len((*dpc)[0].IndexValue), 0)
	}
}

func TestMakeDataPolyComponentUV(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".pd[0]" -type "dataPolyComponent" Index_Data UV 0 ;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeDataPolyComponent {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDataPolyComponent)
	}
	dpc, ok := sa.Attr.(*[]AttrDataPolyComponent)
	sameDPC(t, ok, dpc, DPCuv)
	if len(*dpc) != 1 {
		t.Errorf(msg, "len(Attr)", len(*dpc), 1)
	}
	if len((*dpc)[0].IndexValue) != 0 {
		t.Errorf(msg, "len(Attr.IndexValue)", len((*dpc)[0].IndexValue), 0)
	}
}

func TestMakeDataPolyComponentFace(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".hfd" -type "dataPolyComponent" Index_Data Face 0 ;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeDataPolyComponent {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDataPolyComponent)
	}
	dpc, ok := sa.Attr.(*[]AttrDataPolyComponent)
	sameDPC(t, ok, dpc, DPCface)
	if len(*dpc) != 1 {
		t.Errorf(msg, "len(Attr)", len(*dpc), 1)
	}
	if len((*dpc)[0].IndexValue) != 0 {
		t.Errorf(msg, "len(Attr.IndexValue)", len((*dpc)[0].IndexValue), 0)
	}
}

func TestMakeAttributeAlias(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".aal" -type "attributeAlias" {"detonationFrame","borderConnections[0]","incandescence"
		,"borderConnections[1]","color","borderConnections[2]","nucleusSolver","publishedNodeInfo[0]"
		} ;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeAttributeAlias {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeAttributeAlias)
	}
	aaa, ok := sa.Attr.(*[]AttrAttributeAlias)
	if !ok ||
		(*aaa)[0].NewAlias != "detonationFrame" ||
		(*aaa)[0].CurrentName != "borderConnections[0]" ||
		(*aaa)[1].NewAlias != "incandescence" ||
		(*aaa)[1].CurrentName != "borderConnections[1]" ||
		(*aaa)[2].NewAlias != "color" ||
		(*aaa)[2].CurrentName != "borderConnections[2]" ||
		(*aaa)[3].NewAlias != "nucleusSolver" ||
		(*aaa)[3].CurrentName != "publishedNodeInfo[0]" {
		msg := `got SetAttr %s %s, wont %s`
		t.Errorf(msg, "Attr", aaa, &[]AttrAttributeAlias{
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
	if len(*aaa) != 4 {
		t.Errorf(msg, "len(Attr)", len(*aaa), 4)
	}
}

func TestMakeComponentList(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".ics" -type "componentList" 2 "vtx[130]" "vtx[147]";`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeComponentList {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeComponentList)
	}
	cl, ok := sa.Attr.(*AttrComponentList)
	if !ok ||
		(*cl)[0] != "vtx[130]" ||
		(*cl)[1] != "vtx[147]" {
		msg := `got SetAttr %s %s, wont %s`
		t.Errorf(msg, "Attr", cl, &AttrComponentList{
			"vtx[130]",
			"vtx[147]",
		})
	}
	if len(*cl) != 2 {
		t.Errorf(msg, "len(Attr)", len(*cl), 2)
	}
}

func TestMakeCone(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".cone" -type "cone" 45.0 5.0;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeCone {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeCone)
	}
	cone, ok := sa.Attr.(*[]AttrCone)
	if !ok ||
		(*cone)[0].ConeAngle != 45.0 ||
		(*cone)[0].ConeCap != 5.0 {
		msg := `got SetAttr %s %s, wont %s`
		t.Errorf(msg, "Attr", sa.Attr, &[]AttrCone{
			{
				ConeAngle: 45.0,
				ConeCap:   5.0,
			},
		})
	}
	if len(*cone) != 1 {
		t.Errorf(msg, "len(Attr)", len(*cone), 1)
	}
}

func TestMakeDoubleArray(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`setAttr ".dd" -type "doubleArray" 7 -1 1 0 0 0.5 1 -0.11000000000000004 ;`)
	sa, err := MakeSetAttr(c.Parse(), nil)
	if err != nil {
		t.Fatal(err)
	}
	msg := `got SetAttr %s %s, wont %s`
	if sa.AttrType != TypeDoubleArray {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeDoubleArray)
	}
	da, ok := sa.Attr.(*AttrDoubleArray)
	if !ok ||
		(*da)[0] != -1 ||
		(*da)[1] != 1 ||
		(*da)[2] != 0 ||
		(*da)[3] != 0 ||
		(*da)[4] != 0.5 ||
		(*da)[5] != 1 ||
		(*da)[6] != -0.11000000000000004 {
		msg := `got SetAttr %s %s, wont %s`
		t.Errorf(msg, "Attr", sa.Attr, &AttrDoubleArray{
			-1, 1, 0, 0, 0.5, 1, -0.11000000000000004,
		})
	}
	if len(*da) != 7 {
		t.Errorf(msg, "len(Attr)", len(*da), 7)
	}
}

func TestMakeLattice(t *testing.T) {
	c := &CmdBuilder{}
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
	if sa.AttrType != TypeLattice {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeLattice)
	}
	la, ok := sa.Attr.(*[]AttrLattice)
	if !ok ||
		(*la)[0].DivisionS != 2 ||
		(*la)[0].DivisionT != 2 ||
		(*la)[0].DivisionU != 2 ||
		(*la)[0].Points[0].S != -0.5 ||
		(*la)[0].Points[0].T != -0.5 ||
		(*la)[0].Points[0].U != -0.5 ||
		(*la)[0].Points[1].S != 0.5 ||
		(*la)[0].Points[1].T != -0.5 ||
		(*la)[0].Points[1].U != -0.5 ||
		(*la)[0].Points[2].S != -0.5 ||
		(*la)[0].Points[2].T != 0.5 ||
		(*la)[0].Points[2].U != -0.5 ||
		(*la)[0].Points[3].S != 0.5 ||
		(*la)[0].Points[3].T != 0.5 ||
		(*la)[0].Points[3].U != -0.5 ||
		(*la)[0].Points[4].S != -0.5 ||
		(*la)[0].Points[4].T != -0.5 ||
		(*la)[0].Points[4].U != 0.5 ||
		(*la)[0].Points[5].S != 0.5 ||
		(*la)[0].Points[5].T != -0.5 ||
		(*la)[0].Points[5].U != 0.5 ||
		(*la)[0].Points[6].S != -0.5 ||
		(*la)[0].Points[6].T != 0.5 ||
		(*la)[0].Points[6].U != 0.5 ||
		(*la)[0].Points[7].S != 0.5 ||
		(*la)[0].Points[7].T != 0.5 ||
		(*la)[0].Points[7].U != 0.5 {
		msg := `got SetAttr %s %s, wont %s`
		t.Errorf(msg, "Attr", sa.Attr, &[]AttrLattice{
			{
				DivisionS: 2,
				DivisionT: 2,
				DivisionU: 2,
				Points: []AttrLaticePoint{
					{-0.5, -0.5, -0.5},
					{0.5, -0.5, -0.5},
					{-0.5, 0.5, -0.5},
					{0.5, 0.5, -0.5},
					{-0.5, -0.5, 0.5},
					{0.5, -0.5, 0.5},
					{-0.5, 0.5, 0.5},
					{0.5, 0.5, 0.5},
				},
			},
		})
	}
	if len(*la) != 1 {
		t.Errorf(msg, "len(Attr)", len(*la), 1)
	}
}

func TestMakeNurbsCurve(t *testing.T) {
	c := &CmdBuilder{}
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
	if sa.AttrType != TypeNurbsCurve {
		t.Errorf(msg, "AttrType", sa.AttrType, TypeNurbsCurve)
	}
	nc, ok := sa.Attr.(*[]AttrNurbsCurve)
	if !ok ||
		(*nc)[0].Degree != 3 ||
		(*nc)[0].Spans != 1 ||
		(*nc)[0].Form != AttrFormOpen ||
		(*nc)[0].IsRational != false ||
		(*nc)[0].Dimension != 3 ||
		(*nc)[0].KnotValues[0] != 0 ||
		(*nc)[0].KnotValues[1] != 0 ||
		(*nc)[0].KnotValues[2] != 0 ||
		(*nc)[0].KnotValues[3] != 1 ||
		(*nc)[0].KnotValues[4] != 1 ||
		(*nc)[0].KnotValues[5] != 1 ||
		(*nc)[0].CvValues[0].X != 0 ||
		(*nc)[0].CvValues[0].Y != 0 ||
		*(*nc)[0].CvValues[0].Z != 0 ||
		(*nc)[0].CvValues[0].W != nil ||
		(*nc)[0].CvValues[1].X != 0.33333333333333326 ||
		(*nc)[0].CvValues[1].Y != 0 ||
		*(*nc)[0].CvValues[1].Z != -0.33333333333333326 ||
		(*nc)[0].CvValues[1].W != nil ||
		(*nc)[0].CvValues[2].X != 0.66666666666666663 ||
		(*nc)[0].CvValues[2].Y != 0 ||
		*(*nc)[0].CvValues[2].Z != -0.66666666666666663 ||
		(*nc)[0].CvValues[2].W != nil ||
		(*nc)[0].CvValues[3].X != 1 ||
		(*nc)[0].CvValues[3].Y != 0 ||
		*(*nc)[0].CvValues[3].Z != -1 ||
		(*nc)[0].CvValues[3].W != nil {
		msg := `got SetAttr %s %s, wont %s`
		zero := 0.0
		minus03 := -0.33333333333333326
		minus06 := -0.66666666666666663
		minus1 := -1.0
		t.Errorf(msg, "Attr", sa.Attr, &[]AttrNurbsCurve{
			{
				Degree:     3,
				Spans:      1,
				Form:       AttrFormOpen,
				IsRational: false,
				Dimension:  3,
				KnotValues: []float64{
					0, 0, 0, 1, 1, 1,
				},
				CvValues: []AttrCvValue{
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
	if len(*nc) != 1 {
		t.Errorf(msg, "len(Attr)", len(*nc), 1)
	}
}
