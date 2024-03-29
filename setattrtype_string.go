// Code generated by "stringer -type=SetAttrType setattrtype.go"; DO NOT EDIT.

package mayaascii

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SetAttrTypeInvalid-0]
	_ = x[SetAttrTypeBool-1]
	_ = x[SetAttrTypeInt-2]
	_ = x[SetAttrTypeDouble-3]
	_ = x[SetAttrTypeShort2-4]
	_ = x[SetAttrTypeShort3-5]
	_ = x[SetAttrTypeLong2-6]
	_ = x[SetAttrTypeLong3-7]
	_ = x[SetAttrTypeInt32Array-8]
	_ = x[SetAttrTypeFloat2-9]
	_ = x[SetAttrTypeFloat3-10]
	_ = x[SetAttrTypeDouble2-11]
	_ = x[SetAttrTypeDouble3-12]
	_ = x[SetAttrTypeDoubleArray-13]
	_ = x[SetAttrTypeMatrix-14]
	_ = x[SetAttrTypeMatrixXform-15]
	_ = x[SetAttrTypePointArray-16]
	_ = x[SetAttrTypeVectorArray-17]
	_ = x[SetAttrTypeString-18]
	_ = x[SetAttrTypeStringArray-19]
	_ = x[SetAttrTypeSphere-20]
	_ = x[SetAttrTypeCone-21]
	_ = x[SetAttrTypeReflectanceRGB-22]
	_ = x[SetAttrTypeSpectrumRGB-23]
	_ = x[SetAttrTypeComponentList-24]
	_ = x[SetAttrTypeAttributeAlias-25]
	_ = x[SetAttrTypeNurbsCurve-26]
	_ = x[SetAttrTypeNurbsSurface-27]
	_ = x[SetAttrTypeNurbsTrimface-28]
	_ = x[SetAttrTypePolyFaces-29]
	_ = x[SetAttrTypeDataPolyComponent-30]
	_ = x[SetAttrTypeDataReferenceEdits-31]
	_ = x[SetAttrTypeMesh-32]
	_ = x[SetAttrTypeLattice-33]
}

const _SetAttrType_name = "SetAttrTypeInvalidSetAttrTypeBoolSetAttrTypeIntSetAttrTypeDoubleSetAttrTypeShort2SetAttrTypeShort3SetAttrTypeLong2SetAttrTypeLong3SetAttrTypeInt32ArraySetAttrTypeFloat2SetAttrTypeFloat3SetAttrTypeDouble2SetAttrTypeDouble3SetAttrTypeDoubleArraySetAttrTypeMatrixSetAttrTypeMatrixXformSetAttrTypePointArraySetAttrTypeVectorArraySetAttrTypeStringSetAttrTypeStringArraySetAttrTypeSphereSetAttrTypeConeSetAttrTypeReflectanceRGBSetAttrTypeSpectrumRGBSetAttrTypeComponentListSetAttrTypeAttributeAliasSetAttrTypeNurbsCurveSetAttrTypeNurbsSurfaceSetAttrTypeNurbsTrimfaceSetAttrTypePolyFacesSetAttrTypeDataPolyComponentSetAttrTypeDataReferenceEditsSetAttrTypeMeshSetAttrTypeLattice"

var _SetAttrType_index = [...]uint16{0, 18, 33, 47, 64, 81, 98, 114, 130, 151, 168, 185, 203, 221, 243, 260, 282, 303, 325, 342, 364, 381, 396, 421, 443, 467, 492, 513, 536, 560, 580, 608, 637, 652, 670}

func (i SetAttrType) String() string {
	if i < 0 || i >= SetAttrType(len(_SetAttrType_index)-1) {
		return "SetAttrType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SetAttrType_name[_SetAttrType_index[i]:_SetAttrType_index[i+1]]
}
