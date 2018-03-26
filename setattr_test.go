package main

import (
	"testing"
)

//func TestArrayArray(t *testing.T) {
//	var a [][2]int
//	a = append(a, [2]int{1})
//	a = append(a, [2]int{2})
//	a[0][1] = 3
//	log.Print(a)
//	log.Print(a[0][1])
//	for i, v := range a {
//		log.Print(i, v)
//	}
//}

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
	i, ok := sa.Attr.(*[]int);
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
		t.Errorf(msg, "Attr", sa.Attr, []int{1, 2})
	}
	c.Append(`setAttr ".attrName" -type "short2" 3 4;`)
	_, err = MakeSetAttr(c.Parse(), sa)
	if err != nil {
		t.Fatal(err)
	}
	if len(*s2) != 2 {
		t.Errorf(msg, "len(Attr)", len(*s2), 2)
	}
	if (*s2)[0][0] != 1 || (*s2)[0][1] != 2 ||
		(*s2)[1][0] != 3 || (*s2)[1][1] != 4 {
		t.Errorf(msg, "Attr", sa.Attr, []AttrShort2{
			{1, 2},
			{3, 4},
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
		t.Errorf(msg, "Attr", sa.Attr, []int{1, 2})
	}
	if len(*l2) != 1 {
		t.Errorf(msg, "len(Attr)", len(*l2), 1)
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
		t.Errorf(msg, "Attr", sa.Attr, []int{1, 2, 3})
	}
	if len(*l3) != 1 {
		t.Errorf(msg, "len(Attr)", len(*l3), 1)
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
		t.Errorf(msg, "Attr", sa.Attr, []float64{1, 2.2})
	}
	if len(*f2) != 1 {
		t.Errorf(msg, "len(Attr)", len(*f2), 1)
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
}
