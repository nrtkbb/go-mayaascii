package mayaascii

import "strings"

// go:generate stringer -type=AddAttrDataType addattrdatatype.go
type AddAttrDataType int

const (
	AddAttrDataTypeInvalid AddAttrDataType = iota
	AddAttrDataTypeInt32Array
	AddAttrDataTypeDouble2
	AddAttrDataTypeDouble3
	AddAttrDataTypeDoubleArray
	AddAttrDataTypeFloat2
	AddAttrDataTypeFloat3
	AddAttrDataTypeFloatArray
	AddAttrDataTypeLattice
	AddAttrDataTypeLong2
	AddAttrDataTypeLong3
	AddAttrDataTypeMatrix
	AddAttrDataTypeMesh
	AddAttrDataTypeNurbsCurve
	AddAttrDataTypeNurbsSurface
	AddAttrDataTypePointArray
	AddAttrDataTypeReflectanceRGB
	AddAttrDataTypeShort2
	AddAttrDataTypeShort3
	AddAttrDataTypeSpectrumRGB
	AddAttrDataTypeString
	AddAttrDataTypeStringArray
	AddAttrDataTypeVectorArray
)

func (i AddAttrDataType) Name() string {
	s := i.String()
	s = s[15:]                            // remove AddAttrDataType prefix.
	return strings.ToLower(s[:1]) + s[1:] // ToLower head one string.
}
