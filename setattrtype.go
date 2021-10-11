package mayaascii

import "strings"

// go:generate stringer -type=SetAttrType setattrtype.go
type SetAttrType int

const (
	// SetAttrTypeInvalid parse error
	SetAttrTypeInvalid SetAttrType = iota

	SetAttrTypeBool
	SetAttrTypeInt
	SetAttrTypeDouble

	// SetAttrTypeShort2 type: short short
	// mean: value1 value2
	// setAttr node.short2Attr -type short2 1 2;
	SetAttrTypeShort2

	// SetAttrTypeShort3 type: short short short
	// mean: value1 value2 value3
	// setAttr node.short3Attr -type short3 1 2 3;
	SetAttrTypeShort3

	// SetAttrTypeLong2 type: long long
	// mean: value1 value2
	// setAttr node.long2Attr -type long2 1000000 2000000;
	SetAttrTypeLong2

	// SetAttrTypeLong3 type: long long long
	// mean: value1 value2 value3
	// setAttr node.long3Attr -type long3 1000000 2000000 3000000;
	SetAttrTypeLong3

	// SetAttrTypeInt32Array type: int [int]
	// mean: numberOfArrayValues {arrayValue}
	// setAttr node.int32ArrayAttr -type Int32Array 2 12 75;
	SetAttrTypeInt32Array

	// SetAttrTypeFloat2 type: float float
	// mean: value1 value2
	// setAttr node.float2Attr -type float2 1.1 2.2;
	SetAttrTypeFloat2

	// SetAttrTypeFloat3 type: float float float
	// mean: value1 value2 value3
	// setAttr node.float3Attr -type float3 1.1 2.2 3.3;
	SetAttrTypeFloat3

	// SetAttrTypeDouble2 type: double double
	// mean: value1 value2
	// setAttr node.double2Attr -type double2 1.1 2.2;
	SetAttrTypeDouble2

	// SetAttrTypeDouble3 type: double double double
	// mean: value1 value2 value3
	// setAttr node.double3Attr -type double3 1.1 2.2 3.3;
	SetAttrTypeDouble3

	// SetAttrTypeDoubleArray type: int {double}
	// mean: numberOfArrayValues {arrayValue}
	// setAttr node.doubleArrayAttr -type doubleArray 2 3.14159 2.782;
	SetAttrTypeDoubleArray

	// SetAttrTypeMatrix type: double double double double double double double double double double double double double double double double
	// mean: row1col1 row1col2 row1col3 row1col4 row2col1 row2col2 row2col3 row2col4 row3col1 row3col2 row3col3 row3col4 row4col1 row4col2 row4col3 row4col4
	// setAttr ".ix" -type "matrix" 5 0 0 0 0 0 0 0 0 0 5 0 0 0 0 1;
	SetAttrTypeMatrix

	// SetAttrTypeMatrixXform type: string double double double
	//       double double double
	//       integer
	//       double double double
	//       double double double
	//       double double double
	//       double double double
	//       double double double
	//       double double double
	//       double double double double
	//       double double double double
	//       double double double
	//       boolean
	// mean: xform scaleX scaleY scaleZ
	//       rotateX rotateY rotateZ
	//       rotationOrder (0=XYZ, 1=YZX, 2=ZXY, 3=XZY, 4=YXZ, 5=ZYX)
	//       translateX translateY translateZ
	//       shearXY shearXZ shearYZ
	//       scalePivotX scalePivotY scalePivotZ
	//       scaleTranslationX scaleTranslationY scaleTranslationZ
	//       rotatePivotX rotatePivotY rotatePivotZ
	//       rotateTranslationX rotateTranslationY rotateTranslationZ
	//       rotateOrientW rotateOrientX rotateOrientY rotateOrientZ
	//       jointOrientW jointOrientX jointOrientY jointOrientZ
	//       inverseParentScaleX inverseParentScaleY inverseParentScaleZ
	//       compensateForParentScale
	// setAttr ".xm[0]" -type "matrix" "xform" 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0 1 1 1 1 yes;
	SetAttrTypeMatrixXform

	// SetAttrTypePointArray type: int {double double double double}
	// mean: numberOfArrayValues {xValue yValue zValue wValue}
	// setAttr node.pointArrayAttr -type pointArray 2 1 1 1 1 2 2 2 1;
	SetAttrTypePointArray

	// SetAttrTypeVectorArray type: int {double double double}
	// mean: numberOfArrayValues {xValue yValue zValue}
	// setAttr node.vectorArrayAttr -type vectorArray 2 1 1 1 2 2 2;
	SetAttrTypeVectorArray

	// SetAttrTypeString type: string
	// mean: characterStringValue
	// setAttr node.stringAttr -type "string" "blarg";
	SetAttrTypeString

	// SetAttrTypeStringArray type: int {string}
	// mean: numberOfArrayValues {arrayValue}
	// setAttr node.stringArrayAttr -type stringArray 3 "a" "b" "c";
	SetAttrTypeStringArray

	// SetAttrTypeSphere type: 倍精度浮動小数点数
	// mean: sphereRadius
	// setAttr node.sphereAttr -type sphere 5.0;
	SetAttrTypeSphere

	// SetAttrTypeCone type: double double
	// mean: coneAngle coneCap
	// setAttr node.coneAttr -type cone 45.0 5.0;
	SetAttrTypeCone

	// SetAttrTypeReflectanceRGB type: double double double
	// mean: redReflect greenReflect blueReflect
	// setAttr node.reflectanceRGBAttr -type reflectanceRGB 0.5 0.5 0.1;
	SetAttrTypeReflectanceRGB

	// SetAttrTypeSpectrumRGB type: double double double
	// mean: redSpectrum greenSpectrum blueSpectrum
	// setAttr node.spectrumRGBAttr -type spectrumRGB 0.5 0.5 0.1;
	SetAttrTypeSpectrumRGB

	// SetAttrTypeComponentList type: int {string}
	// mean: numberOfComponents {componentName}
	// setAttr node.componentListAttr -type componentList 3 cv[1] cv[12] cv[3];
	SetAttrTypeComponentList

	// SetAttrTypeAttributeAlias type: string string
	// mean: newAlias currentName
	// setAttr node.attrAliasAttr -type attributeAlias {"GoUp","translateY", "GoLeft","translateX"};
	SetAttrTypeAttributeAlias

	// SetAttrTypeNurbsCurve type: int int int bool int int {double}
	//       int {double double double}
	// mean: degree spans form isRational dimension knotCount {knotValue}
	//       cvCount {xCVValue yCVValue [zCVValue] [wCVValue]}
	// degree はカーブの次数(1 ～ 7)
	// spans はスパン数
	// form は open (0)、closed (1)、periodic (2)
	// dimension は 2 または 3 (カーブの次元に依存)
	// isRational カーブの CV に有理コンポーネントが含まれる場合に true
	// knotCount はノット リストのサイズ
	// knotValue はノット リストの単一エントリ
	// cvCount はカーブの CV 数
	// xCVValue、yCVValue、[zCVValue] [wCVValue] は単一の CV
	// zCVValue が存在するのは次元が 3 の場合のみ
	// wCVValue が存在するのは isRational が true の場合のみ
	//
	// setAttr node.curveAttr -type nurbsCurve 3 1 0 no 3 6 0 0 0 1 1 1
	// 4 -2 3 0 -2 1 0 -2 -1 0 -2 -3 0;
	SetAttrTypeNurbsCurve

	// SetAttrTypeNurbsSurface type: int int int int bool
	//       int {double}
	//       int {double}
	//       [string] int {double double double}
	// mean: uDegree vDegree uForm vForm isRational
	//       uKnotCount {uKnotValue}
	//       vKnotCount {vKnotValue} ["TRIM"|"NOTRIM"]
	//       cvCount {xCVValue yCVValue zCVValue [wCVValue]}
	// uDegree は U 方向のサーフェスの次数(範囲 1 ～ 7)
	// vDegree は V 方向のサーフェスの次数(範囲 1 ～ 7)
	// uForm は U 方向での open (0)、closed (1)、periodic (2)
	// vForm は V 方向での open (0)、closed (1)、periodic (2)
	// isRational はサーフェスの CV に有理コンポーネントが含まれるに true
	// uKnotCount は U ノット リストのサイズ
	// uKnotValue は U ノット リストの単一エントリ
	// vKnotCount は V ノット リストのサイズ
	// vKnotValue は V ノット リストの単一エントリ
	// "TRIM" を指定する場合は、トリム情報が必要
	// "NOTRIM" を指定すると、サーフェスはトリムされない
	// cvCount はサーフェスの CV 数
	// xCVValue、yCVValue、[zCVValue] [wCVValue] は単一の CV
	// zCVValue が存在するのは次元が 3 の場合のみ
	// wCVValue が存在するのは isRational が true の場合のみ
	//
	// setAttr node.surfaceAttr -type nurbsSurface 3 3 0 0 no
	// 6 0 0 0 1 1 1
	// 6 0 0 0 1 1 1
	// 16 -2 3 0 -2 1 0 -2 -1 0 -2 -3 0
	// -1 3 0 -1 1 0 -1 -1 0 -1 -3 0
	// 1 3 0 1 1 0 1 -1 0 1 -3 0
	// 3 3 0 3 1 0 3 -1 0 3 -3 0;
	SetAttrTypeNurbsSurface

	// SetAttrTypeNurbsTrimface type: bool int {int {int {int int int} int {int int}}}
	// mean: flipNormal boundaryCount {boundaryType tedgeCountOnBoundary
	//       {splineCountOnEdge {edgeTolerance isEdgeReversed geometricContinuity}
	//       {splineCountOnPedge {isMonotone pedgeTolerance}}}
	//
	// ↑こちらは online help に記載されてる内容だが、間違い。
	// ↓こちらはAutodeskから回答いただいた内容。
	//
	// BSPR-30157 - Doc: setAttr -type nurbsTrimface description shows "int"
	// but actually double and bool type value can be input
	//
	// 正しくは、下記のような値タイプになります。
	//
	// type: bool int {int {int {double bool bool} int {bool double}}}
	// mean: flipNormal boundaryCount {boundaryType tedgeCountOnBoundary
	//       {splineCountOnEdge {edgeTolerance isEdgeReversed geometricContinuity}
	//       {splineCountOnPedge {isMonotone pedgeTolerance}}}
	//
	// flipNormal は true の場合にサーフェスを反転させる -> Bool
	// boundaryCount: 境界の数 -> Int
	// boundaryType: -> Int
	// tedgeCountOnBoundary : 境界のエッジ数 -> Int
	// splineCountOnEdge : エッジのスプライン数 -> Int
	// edgeTolerance : 3D エッジを構築する際に使用する許容値 -> Double
	// isEdgeReversed : true の場合、エッジは逆向きになる -> Bool
	// geometricContinuity : true の場合、エッジは接線連続性を持つ -> Bool
	// splineCountOnPedge : 2D エッジのスプライン数 -> Int
	// isMonotone : true の場合、曲率は単調になる -> Bool
	// pedgeTolerance : 2D エッジの許容値 -> Double
	//
	SetAttrTypeNurbsTrimface

	// SetAttrTypePolyFaces type: {"f" int {int}}
	//       {"h" int {int}}
	//       {"mf" int {int}}
	//       {"mh" int {int}}
	//       {"mu" int int {int}}
	//       {"mc" int int {int}}
	//       {"fc" int {int}}
	// mean: {"f" faceEdgeCount {edgeIdValue}}
	//       {"h" holeEdgeCount {edgeIdValue}}
	//       {"mf" faceUVCount {uvIdValue}}
	//       {"mh" holeUVCount {uvIdValue}}
	//       {"mu" uvSet faceUVCount {uvIdValue}}
	//       {"mc" colorIndex colorIdCount {colorIdValue}}
	//       {"fc" faceColorCount {colorIndexValue}}
	// このデータ型(polyFace)は、setAttrs で頂点位置配列、
	// エッジ接続性配列(および対応する開始/終了頂点の記述)、
	// テクスチャ座標配列、カラー配列を書き出した後に
	// ファイルの読み取りや書き出しで使用するためのものです。
	// このデータ型は以前の型で
	// 作成された ID を使用してすべてのデータを参照します。
	//
	// "f" はフェースを構成するエッジの ID を指定 -
	// フェースでエッジが反転する場合は負の値
	// "h" は穴を構成するエッジの ID を指定 -
	// フェースでエッジが反転する場合は負の値
	// "mf" はフェースのテクスチャ座標(UV)の ID を指定
	// このデータ型はバージョン 3.0 で廃止されており。代わりに "mu" が使用されています。
	// "mh" は穴のテクスチャ座標(UV)を指定
	// このデータ型はバージョン 3.0 で廃止されており。代わりに "mu" が使用されています。
	// "mu" 最初の引数は UV セットです。これはゼロから始まる
	// 整数値です。2 番目の引数は有効な UV 値を持つフェース上の
	// 頂点の数です。最後の値はフェースの
	// テクスチャ座標(UV)の UV ID です。 これらのインデックスは
	// "mf" や "mh" を指定する際に使用するものです。
	// "mu" は複数指定することもできます(固有の UV セットごとに 1 つ)。
	// "fc" はフェースのカラー インデックス値を指定します。
	//
	//
	// `mc` (multi-color) is a replacement the old `fc` flag for color maps.
	// The first argument to `mc` is the color map (index) to use.
	//
	// setAttr ".fc[0]" -type "polyFaces"
	// f 4 0 2 -4 -2
	// mu 0 4 0 1 3 2
	// mc 0 4 0 1 3 2;
	//
	// Looking at the code the first value is the colour index, the second
	// one is the number of colour IDs to follow, then the rest are the list
	// of those colour IDs. In this case it’s colour index 0 with 4 colour
	// IDs of 0, 1, 3, and 2.
	//
	//
	// setAttr node.polyFaceAttr -type polyFaces "f" 3 1 2 3 "fc" 3 4 4 6;
	SetAttrTypePolyFaces

	// SetAttrTypeDataPolyComponent From the code
	// _dataPolyComponent_ takes data of the form
	// Index_Data Edge|Face|Vertex|UV
	// COUNT_OF_INDEX_VALUES {Index Value}
	SetAttrTypeDataPolyComponent

	// SetAttrTypeDataReferenceEdits type: "string"
	//       "string" int
	//       0 "string" "string" "string"
	//       1 string "string" "string" "string"
	//       2 "string" "string" "string"
	//       3 "string" "string" "string"
	//       4 "string" "string" ""
	//       5 0 "string" "string" "string" "string" "string" ""
	//       5 3 "string" "string" "string" ""
	//       5 4 "string" "string" "string" ""
	//       7 "string" "string" int ["string"...]
	//       0
	//       8 "string" "string"
	//       9 "string" "string"
	// mean: "topReferenceNodeName"
	//       "topReferenceName:childReferenceNodeName" {int:number of edit}
	//       {0:parent} "{nodeNameA}" "{nodeNameB}" "{argument of parent}"
	//       {1:addAttr} {nodeName} "{longAttrName}" "{shortAttrName}" "{arguments of addAttr}"
	//       {2:setAttr} "{nodeName}" "{attrName}" "{arguments of setAttr}"
	//       {3:disconnectAttr} "{sourcePlug}" "{distPlug}" ""
	//       {4:deleteAttr} "{nodeName}" "{attrName}" ""
	//       {5:connectAttr1} 0 "{namespace}:{referenceNodeName}" "{sourcePlug}" "{distPlug}"
	//       "{sourcePlaceHolder}" "{distPlaceHolder}" ""
	//       {5:connectAttr2} [34] "{namespace}:{referenceNodeName}" "{sourcePlug}" "{distPlug}" ""
	//       {7:relationship} "{fcurve}" "{namespace}:{node_name}" {int:number of command} "{command}" 0
	//       {8:lock} "{namespace}:{nodeName}" "{attrName}"
	//       {9:unLock} "{namespace}:{nodeName}" "{attrName}"
	//
	// example:
	//       setAttr ".ed" -type "dataReferenceEdits"
	//           "namespaceRN"
	//           "namespace:childNameSpaceRN" 11
	//           0 "nodeNameA" "nodeNameB" "-s -r "
	//           1 |namespace:topNode "extraAttr" "ea" " -ci 1 -nn \"ea\" -at \"double\""
	//           2 "|namespace:topNode" "ea" " -k 1 0"
	//           3 "namespace:nodeNameA.attrNameA" "namespace:nodeNameB.attrNameB" ""
	//           4 "|namespace:topNode" "dexAttr" ""
	//           5 0 "namespace:childNameSpaceRN" "|namespace:topNode.attrNameA" "|namespace:topNode.attrNameB"
	//           "namespace:childNameSpaceRN.placeHolderList[1]" "namespace:childNameSpaceRN.placeHolderList[2]" ""
	//           5 3 "namespace:childNameSpaceRN" "|namespace:topNode|namespace:childNode.attrNameA"
	//           "namespace:childNameSpaceRN.placeHolderList[3]"
	//           5 4 "namespace:childNameSpaceRN" "|namespace:topNode|namespace:childNode.attrNameB"
	//           "namespace:childNameSpaceRN.placeHolderList[4]"
	//           7 "fcurve" "|namespace:nodeName_attrName_X" 1
	//           "add 396 -4131.291016 18 18 1 0 0 423 -4131.291016 18 18 1 0 0" 0
	//           8 "|namespace:topNode" "attrNameA"
	//           9 "|namespace:topNode" "attrNameA"
	//           ;
	SetAttrTypeDataReferenceEdits

	// SetAttrTypeMesh type: {string [int {double double double}]}
	//       {string [int {double double double}]}
	//       [{string [int {double double}]}]
	//       {string [int {double double string}]}
	// mean: "v" [vertexCount {vertexX vertexY vertexZ}]
	//       "vn" [normalCount {normalX normalY normalZ}]
	//       ["vt" [uvCount {uValue vValue}]]
	//       "e" [edgeCount {startVertex endVertex "smooth"|"hard"}]
	// "v" はポリゴン メッシュの頂点を指定
	// "vn" は各頂点の法線を指定
	// "vt" はオプションで、各頂点の U,V テクスチャ座標を指定
	// "e" は頂点間のエッジの接続情報を指定
	//
	// setAttr node.meshAttr -type mesh "v" 3 0 0 0 0 1 0 0 0 1
	// "vn" 3 1 0 0 1 0 0 1 0 0
	// "vt" 3 0 0 0 1 1 0
	// "e" 3 0 1 "hard" 1 2 "hard" 2 0 "hard";
	SetAttrTypeMesh

	// SetAttrTypeLattice type: int int int int {double double double}
	// mean: sDivisionCount tDivisionCount uDivisionCount
	//       pointCount {pointX pointY pointZ}
	// sDivisionCount は水平方向のラティス分割数
	// tDivisionCount は垂直方向のラティス分割数
	// uDivisionCount は深度のラティス分割数
	// pointCount はラティス ポイントの総数
	// pointX、pointY、pointZ は単一のラティス ポイントこのリストは
	// S、T、U の順に異なる値を使用して指定されるため
	// 最初の 2 つのエントリは(S=0,T=0,U=0) (S=1,T=0,U=0) となる
	//
	// setAttr node.latticeAttr -type lattice 2 5 2 20
	// -2 -2 -2 2 -2 -2 -2 -1 -2 2 -1 -2 -2 0 -2
	// 2 0 -2 -2 1 -2 2 1 -2 -2 2 -2 2 2 -2
	// -2 -2 2 2 -2 2 -2 -1 2 2 -1 2 -2 0 2
	// 2 0 2 -2 1 2 2 1 2 -2 2 2 2 2 2;
	SetAttrTypeLattice
)

func (sa SetAttrType) Name() string {
	s := sa.String()
	s = s[11:] // remove SetAttrType prefix.
	return strings.ToLower(s[:1]) + s[1:] // ToLower head one string.
}

