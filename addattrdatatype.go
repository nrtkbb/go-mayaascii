package mayaascii

import (
	"fmt"
	"strings"
)

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

func NewAddAttrDataType(name string) (AddAttrDataType, error) {
	switch name {
	case "int32Array":
		return AddAttrDataTypeInt32Array, nil
	case "double2":
		return AddAttrDataTypeDouble2, nil
	case "double3":
		return AddAttrDataTypeDouble3, nil
	case "doubleArray":
		return AddAttrDataTypeDoubleArray, nil
	case "float2":
		return AddAttrDataTypeFloat2, nil
	case "float3":
		return AddAttrDataTypeFloat3, nil
	case "floatArray":
		return AddAttrDataTypeFloatArray, nil
	case "lattice":
		return AddAttrDataTypeLattice, nil
	case "long2":
		return AddAttrDataTypeLong2, nil
	case "long3":
		return AddAttrDataTypeLong3, nil
	case "matrix":
		return AddAttrDataTypeMatrix, nil
	case "mesh":
		return AddAttrDataTypeMesh, nil
	case "nurbsCurve":
		return AddAttrDataTypeNurbsCurve, nil
	case "nurbsSurface":
		return AddAttrDataTypeNurbsSurface, nil
	case "pointArray":
		return AddAttrDataTypePointArray, nil
	case "reflectanceRGB":
		return AddAttrDataTypeReflectanceRGB, nil
	case "short2":
		return AddAttrDataTypeShort2, nil
	case "short3":
		return AddAttrDataTypeShort3, nil
	case "spectrumRGB":
		return AddAttrDataTypeSpectrumRGB, nil
	case "string":
		return AddAttrDataTypeString, nil
	case "stringArray":
		return AddAttrDataTypeStringArray, nil
	case "vectorArray":
		return AddAttrDataTypeVectorArray, nil
	}
	return AddAttrDataTypeInvalid, fmt.Errorf(
		"%s is not AddAttrDataType name", name)
}
