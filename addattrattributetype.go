package mayaascii

import (
	"fmt"
	"strings"
)

// go:generate stringer -type=AddAttrAttributeType addattrattributetype.go
type AddAttrAttributeType int

const (
	AddAttrAttributeTypeInvalid AddAttrAttributeType = iota
	AddAttrAttributeTypeBool
	AddAttrAttributeTypeByte
	AddAttrAttributeTypeChar
	AddAttrAttributeTypeCompound
	AddAttrAttributeTypeDouble
	AddAttrAttributeTypeDouble2
	AddAttrAttributeTypeDouble3
	AddAttrAttributeTypeDoubleAngle
	AddAttrAttributeTypeDoubleLinear
	AddAttrAttributeTypeEnum
	AddAttrAttributeTypeFloat
	AddAttrAttributeTypeFloat2
	AddAttrAttributeTypeFloat3
	AddAttrAttributeTypeFltMatrix
	AddAttrAttributeTypeLong
	AddAttrAttributeTypeLong2
	AddAttrAttributeTypeLong3
	AddAttrAttributeTypeMessage
	AddAttrAttributeTypeReflectance
	AddAttrAttributeTypeShort
	AddAttrAttributeTypeShort2
	AddAttrAttributeTypeShort3
	AddAttrAttributeTypeSpectrum
	AddAttrAttributeTypeTime
)

func (i AddAttrAttributeType) Name() string {
	s := i.String()
	s = s[20:]                            // remove AddAttrAttributeType prefix.
	return strings.ToLower(s[:1]) + s[1:] // ToLower head one string.
}

func NewAddAttrAttributeType(name string) (AddAttrAttributeType, error) {
	switch name {
	case "bool":
		return AddAttrAttributeTypeBool, nil
	case "byte":
		return AddAttrAttributeTypeByte, nil
	case "char":
		return AddAttrAttributeTypeChar, nil
	case "compound":
		return AddAttrAttributeTypeCompound, nil
	case "double":
		return AddAttrAttributeTypeDouble, nil
	case "double2":
		return AddAttrAttributeTypeDouble2, nil
	case "double3":
		return AddAttrAttributeTypeDouble3, nil
	case "doubleAngle":
		return AddAttrAttributeTypeDoubleAngle, nil
	case "doubleLinear":
		return AddAttrAttributeTypeDoubleLinear, nil
	case "enum":
		return AddAttrAttributeTypeEnum, nil
	case "float":
		return AddAttrAttributeTypeFloat, nil
	case "float2":
		return AddAttrAttributeTypeFloat2, nil
	case "float3":
		return AddAttrAttributeTypeFloat3, nil
	case "fltMatrix":
		return AddAttrAttributeTypeFltMatrix, nil
	case "long":
		return AddAttrAttributeTypeLong, nil
	case "long2":
		return AddAttrAttributeTypeLong2, nil
	case "long3":
		return AddAttrAttributeTypeLong3, nil
	case "message":
		return AddAttrAttributeTypeMessage, nil
	case "reflectance":
		return AddAttrAttributeTypeReflectance, nil
	case "short":
		return AddAttrAttributeTypeShort, nil
	case "short2":
		return AddAttrAttributeTypeShort2, nil
	case "short3":
		return AddAttrAttributeTypeShort3, nil
	case "spectrum":
		return AddAttrAttributeTypeSpectrum, nil
	case "time":
		return AddAttrAttributeTypeTime, nil
	}
	return AddAttrAttributeTypeInvalid, fmt.Errorf(
		"%s is not AddAttrAttributeType name", name)
}
