package mayaascii

import "testing"

func TestMakeAddAttr_int(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`addAttr -shortName ms -longName mass -defaultValue 1.0 -minValue 0.001 -maxValue 10000;`)
	addAttr, err := ParseAddAttr(c.Parse())
	if err != nil {
		t.Fatal(err)
	}

	msg := `got AddAttrCmd %s %v, wont %v`
	if addAttr.ShortName == nil {
		t.Fatalf(msg, "shotName", nil, "ms")
	}
	if *addAttr.ShortName != "ms" {
		t.Errorf(msg, "shotName", *addAttr.ShortName, "ms")
	}
	if addAttr.LongName == nil {
		t.Fatalf(msg, "longName", nil, "mass")
	}
	if *addAttr.LongName != "mass" {
		t.Errorf(msg, "longName", *addAttr.LongName, "mass")
	}
	if addAttr.DefaultValue == nil {
		t.Fatalf(msg, "defaultValue", nil, 1.0)
	}
	if *addAttr.DefaultValue != 1.0 {
		t.Errorf(msg, "defautlValue", *addAttr.DefaultValue, 1.0)
	}
	if addAttr.MinValue == nil {
		t.Fatalf(msg, "minValue", nil, 0.001)
	}
	if *addAttr.MinValue != 0.001 {
		t.Errorf(msg, "minValue", *addAttr.MinValue, 0.001)
	}
	if addAttr.MaxValue == nil {
		t.Fatalf(msg, "maxValue", nil, 10000)
	}
	if *addAttr.MaxValue != 10000 {
		t.Errorf(msg, "maxValue", 10000)
	}
}

func TestMakeAddAttr_bool(t *testing.T) {
	c := &CmdBuilder{}
	c.Append(`addAttr -ci true -sn "liw" -ln "lockInfluenceWeights" -min 0 -max 1 -at "bool";`)
	addAttr, err := ParseAddAttr(c.Parse())
	if err != nil {
		t.Fatal(err)
	}

	msg := `got AddAttrCmd %s %v, wont %v`
	if addAttr.CachedInternally == nil {
		t.Fatalf(msg, "cacheInternally", nil, true)
	}
	if *addAttr.CachedInternally != true {
		t.Errorf(msg, "cacheInternally", *addAttr.CachedInternally, true)
	}
	if addAttr.ShortName == nil {
		t.Fatalf(msg, "shortName", nil, "liw")
	}
	if *addAttr.ShortName != "liw" {
		t.Errorf(msg, "shortName", *addAttr.ShortName, "liw")
	}
	if addAttr.LongName == nil {
		t.Fatalf(msg, "longName", nil, "lockInfluenceWeights")
	}
	if *addAttr.LongName != "lockInfluenceWeights" {
		t.Errorf(msg, "longName", *addAttr.LongName, "lockInfluenceWeights")
	}
	if addAttr.MinValue == nil {
		t.Fatalf(msg, "minValue", nil, 0)
	}
	if *addAttr.MinValue != 0 {
		t.Errorf(msg, "minValue", *addAttr.MinValue, 0)
	}
	if addAttr.MaxValue == nil {
		t.Fatalf(msg, "maxValue", nil, 1)
	}
	if *addAttr.MaxValue != 1 {
		t.Errorf(msg, "maxValue", *addAttr.MaxValue, 1)
	}
	if addAttr.AttributeType == nil {
		t.Fatalf(msg, "attributeType", nil, "bool")
	}
	if *addAttr.AttributeType != AddAttrAttributeTypeBool {
		t.Errorf(msg, "attributeType", *addAttr.AttributeType, AddAttrAttributeTypeBool)
	}
}
