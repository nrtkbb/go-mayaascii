package mayaascii

import "strings"

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
