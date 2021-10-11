package mayaascii

import (
	"strings"
	"testing"
)

func getTestMa() string {
	return `//Maya ASCII 2019 scene
//Name: test.ma
//Last modified: Thu, Oct 1, 2019 01:00:00 AM
//Codeset: UTF-8
requires maya "201iff07";
requires -nodeType "nearestPointOnMesh" "nearestPointOnMesh" "4.0";
currentUnit -l centimeter -a degree -t film;
fileInfo "application" "maya";
fileInfo "product" "Maya 2018";
fileInfo "version" "2018";
fileInfo "cutIdentifier" "xxxx";
fileInfo "osv" "Microsoft Windows 8 Home Premium Edition, 64-bit  (Build 9200)\n";
createNode transform -s -n "persp";
	rename -uid "CFAE1109-4845-2AC4-5BC0-CB8FB886A568";
	setAttr ".v" no;
	setAttr ".t" -type "double3" 28 21 28 ;
	setAttr ".r" -type "double3" -27.938352729602379 44.999999999999972 -5.172681101354183e-14 ;
createNode camera -s -n "perspShape" -p "persp";
	rename -uid "9FA883FE-404A-E503-71B1-8796E0AAEE09";
	setAttr -k off ".v" no;
	setAttr ".fl" 34.999999999999993;
	setAttr ".coi" 44.82186966202994;
	setAttr ".imn" -type "string" "persp";
	setAttr ".den" -type "string" "persp_depth";
	setAttr ".man" -type "string" "persp_mask";
	setAttr ".hc" -type "string" "viewSet -p %camera";
createNode transform -s -n "top";
	rename -uid "2B8E5E49-4563-34B7-B04E-73860F4AD5F9";
	setAttr ".v" no;
	setAttr ".t" -type "double3" 0 1000.1 0 ;
	setAttr ".r" -type "double3" -89.999999999999986 0 0 ;
createNode camera -s -n "topShape" -p "top";
	rename -uid "1F6627D2-4A29-20C2-597E-3CB4B7E72DA6";
	setAttr -k off ".v" no;
	setAttr ".rnd" no;
	setAttr ".coi" 1000.1;
	setAttr ".ow" 30;
	setAttr ".imn" -type "string" "top";
	setAttr ".den" -type "string" "top_depth";
	setAttr ".man" -type "string" "top_mask";
	setAttr ".hc" -type "string" "viewSet -t %camera";
	setAttr ".o" yes;
createNode transform -s -n "front";
	rename -uid "B372B829-4FBA-BAD1-05F6-A18F24F904FB";
	setAttr ".v" no;
	setAttr ".t" -type "double3" 0 0 1000.1 ;
createNode camera -s -n "frontShape" -p "front";
	rename -uid "02F66D57-41B2-05FC-4792-43A3935BB379";
	setAttr -k off ".v" no;
	setAttr ".rnd" no;
	setAttr ".coi" 1000.1;
	setAttr ".ow" 30;
	setAttr ".imn" -type "string" "front";
	setAttr ".den" -type "string" "front_depth";
	setAttr ".man" -type "string" "front_mask";
	setAttr ".hc" -type "string" "viewSet -f %camera";
	setAttr ".o" yes;
createNode transform -s -n "side";
	rename -uid "04F463E5-4453-1D93-FF79-93B09FCD8732";
	setAttr ".v" no;
	setAttr ".t" -type "double3" 1000.1 0 0 ;
	setAttr ".r" -type "double3" 0 89.999999999999986 0 ;
createNode camera -s -n "sideShape" -p "side";
	rename -uid "AA9781D0-43FF-A334-53AA-0D966EF4CBA2";
	setAttr -k off ".v" no;
	setAttr ".rnd" no;
	setAttr ".coi" 1000.1;
	setAttr ".ow" 30;
	setAttr ".imn" -type "string" "side";
	setAttr ".den" -type "string" "side_depth";
	setAttr ".man" -type "string" "side_mask";
	setAttr ".hc" -type "string" "viewSet -s %camera";
	setAttr ".o" yes;
createNode transform -n "group1";
	rename -uid "E0ED7F6A-4729-596E-DA41-C0A33F50AAA9";
createNode transform -n "pCube1" -p "group1";
	rename -uid "4CB23388-47E0-1310-8177-A78A230BDD29";
createNode mesh -n "pCubeShape1" -p "|group1|pCube1";
	rename -uid "1029B43B-4E30-EB32-2C75-2687B0A1E81B";
	setAttr -k off ".v";
	setAttr ".vir" yes;
	setAttr ".vif" yes;
	setAttr ".uvst[0].uvsn" -type "string" "map1";
	setAttr ".cuvs" -type "string" "map1";
	setAttr ".dcc" -type "string" "Ambient+Diffuse";
	setAttr ".covm[0]"  0 1 1;
	setAttr ".cdvm[0]"  0 1 1;
createNode transform -n "group2";
	rename -uid "61B503D1-4EC4-B442-C619-21AB8861DCB9";
createNode transform -n "pCube1" -p "group2";
	rename -uid "413F3759-4779-7F3C-92AB-58B0C8EB1D3D";
createNode mesh -n "pCubeShape1" -p "|group2|pCube1";
	rename -uid "63986C9E-4D3E-783F-23B1-FB970A15C370";
	setAttr -k off ".v";
	setAttr ".vir" yes;
	setAttr ".vif" yes;
	setAttr ".uvst[0].uvsn" -type "string" "map1";
	setAttr -s 14 ".uvst[0].uvsp[0:13]" -type "float2" 0.375 0 0.625 0 0.375
		 0.25 0.625 0.25 0.375 0.5 0.625 0.5 0.375 0.75 0.625 0.75 0.375 1 0.625 1 0.875 0
		 0.875 0.25 0.125 0 0.125 0.25;
	setAttr ".cuvs" -type "string" "map1";
	setAttr ".dcc" -type "string" "Ambient+Diffuse";
	setAttr ".covm[0]"  0 1 1;
	setAttr ".cdvm[0]"  0 1 1;
	setAttr -s 8 ".vt[0:7]"  -0.5 -0.5 0.5 0.5 -0.5 0.5 -0.5 0.5 0.5 0.5 0.5 0.5
		 -0.5 0.5 -0.5 0.5 0.5 -0.5 -0.5 -0.5 -0.5 0.5 -0.5 -0.5;
	setAttr -s 12 ".ed[0:11]"  0 1 0 2 3 0 4 5 0 6 7 0 0 2 0 1 3 0 2 4 0
		 3 5 0 4 6 0 5 7 0 6 0 0 7 1 0;
	setAttr -s 6 -ch 24 ".fc[0:5]" -type "polyFaces" 
		f 4 0 5 -2 -5
		mu 0 4 0 1 3 2
		f 4 1 7 -3 -7
		mu 0 4 2 3 5 4
		f 4 2 9 -4 -9
		mu 0 4 4 5 7 6
		f 4 3 11 -1 -11
		mu 0 4 6 7 9 8
		f 4 -12 -10 -8 -6
		mu 0 4 1 10 11 3
		f 4 10 4 6 8
		mu 0 4 12 0 2 13;
	setAttr ".cd" -type "dataPolyComponent" Index_Data Edge 0 ;
	setAttr ".cvd" -type "dataPolyComponent" Index_Data Vertex 0 ;
	setAttr ".pd[0]" -type "dataPolyComponent" Index_Data UV 0 ;
	setAttr ".hfd" -type "dataPolyComponent" Index_Data Face 0 ;
createNode lightLinker -s -n "lightLinker1";
	rename -uid "3C3DFFBA-4F59-FEFE-138D-DDABD5AC5AE0";
	setAttr -s 2 ".lnk";
	setAttr -s 2 ".slnk";
createNode shapeEditorManager -n "shapeEditorManager";
	rename -uid "1A7DF032-4A33-A0D0-E8CD-BEB8CD090FCD";
createNode poseInterpolatorManager -n "poseInterpolatorManager";
	rename -uid "78BDC8AD-4D9C-1F85-5C36-C4A92D829280";
createNode displayLayerManager -n "layerManager";
	rename -uid "AAC57732-457E-1066-CB14-F4817796B7D7";
createNode displayLayer -n "defaultLayer";
	rename -uid "68054D8E-4300-36D1-49A0-E4A8A1D5071F";
createNode renderLayerManager -n "renderLayerManager";
	rename -uid "9DF0D069-49FC-2905-2607-CF8481821B5F";
createNode renderLayer -n "defaultRenderLayer";
	rename -uid "F943CCCE-4BCC-D170-402E-2F91E75DD736";
	setAttr ".g" yes;
createNode polyCube -n "polyCube1";
	rename -uid "66200A8D-46DC-7804-2ED9-E983C0496492";
	setAttr ".cuv" 4;
createNode nearestPointOnMesh -n "nearestPointOnMesh1";
	rename -uid "B2C0D7BA-4FA2-A74E-5D6A-8DAA23AAF72E";
createNode script -n "sceneConfigurationScriptNode";
	rename -uid "4F61132E-4CDD-B5DD-E33E-AD9341041F6D";
	setAttr ".b" -type "string" "playbackOptions -min 1 -max 120 -ast 1 -aet 200 ";
	setAttr ".st" 6;
select -ne :time1;
	setAttr ".o" 1;
	setAttr ".unw" 1;
select -ne :hardwareRenderingGlobals;
	setAttr ".otfna" -type "stringArray" 22 "NURBS Curves" "NURBS Surfaces" "Polygons" "Subdiv Surface" "Particles" "Particle Instance" "Fluids" "Strokes" "Image Planes" "UI" "Lights" "Cameras" "Locators" "Joints" "IK Handles" "Deformers" "Motion Trails" "Components" "Hair Systems" "Follicles" "Misc. UI" "Ornaments"  ;
	setAttr ".otfva" -type "Int32Array" 22 0 1 1 1 1 1
		 1 1 1 0 0 0 0 0 0 0 0 0
		 0 0 0 0 ;
	setAttr ".fprt" yes;
select -ne :renderPartition;
	setAttr -s 2 ".st";
select -ne :renderGlobalsList1;
select -ne :defaultShaderList1;
	setAttr -s 4 ".s";
select -ne :postProcessList1;
	setAttr -s 2 ".p";
select -ne :defaultRenderingList1;
select -ne :initialShadingGroup;
	setAttr -s 2 ".dsm";
	setAttr ".ro" yes;
select -ne :initialParticleSE;
	setAttr ".ro" yes;
select -ne :defaultResolution;
	setAttr ".pa" 1;
select -ne :hardwareRenderGlobals;
	setAttr ".ctrs" 256;
	setAttr ".btrs" 512;
select -ne :ikSystem;
	setAttr -s 4 ".sol";
/*block comment1*/
/*
	block
	comment2
*/
connectAttr "polyCube1.out" "|group1|pCube1|pCubeShape1.i";
relationship "link" ":lightLinker1" ":initialShadingGroup.message" ":defaultLightSet.message";
relationship "link" ":lightLinker1" ":initialParticleSE.message" ":defaultLightSet.message";
relationship "shadowLink" ":lightLinker1" ":initialShadingGroup.message" ":defaultLightSet.message";
relationship "shadowLink" ":lightLinker1" ":initialParticleSE.message" ":defaultLightSet.message";
connectAttr "layerManager.dli[0]" "defaultLayer.id";
connectAttr "renderLayerManager.rlmi[0]" "defaultRenderLayer.rlid";
connectAttr "defaultRenderLayer.msg" ":defaultRenderingList1.r" -na;
connectAttr "|group1|pCube1|pCubeShape1.iog" ":initialShadingGroup.dsm" -na;
connectAttr "|group2|pCube1|pCubeShape1.iog" ":initialShadingGroup.dsm" -na;
// End of test.ma`
}

func TestRequires(t *testing.T) {
	reader := strings.NewReader(getTestMa())

	mo, err := Unmarshal(reader)
	if err != nil {
		t.Error(err.Error())
	}

	if mo == nil {
		t.Error("got nil, wont *Object")
	}

	if len(mo.Requires) != 2 {
		t.Errorf("got len(mo.RequiresCmd) %d, wont 2", len(mo.Requires))
	}

	require := mo.Requires[1]
	if require.GetPluginName() != "nearestPointOnMesh" {
		t.Errorf("got %s, wont nearestPointOnMesh", require.GetPluginName())
	}

	if require.GetVersion() != "4.0" {
		t.Errorf("got %s, wont 4.0", require.GetVersion())
	}

	if len(require.GetDataTypes()) != 0 {
		t.Errorf("got %d, wont 0", len(require.GetDataTypes()))
	}

	if len(require.GetNodeTypes()) != 1 {
		t.Fatalf("got %d, wont 1", len(require.GetNodeTypes()))
	}

	nodeType := require.GetNodeTypes()[0]
	if nodeType != "nearestPointOnMesh" {
		t.Errorf("got %s, wont nearestPointOnMesh", nodeType)
	}

	if len(require.Nodes) != 1 {
		t.Errorf("got %d, wont 1", len(require.Nodes))
	}

	nearestPointOnMeshNode := require.Nodes[0]
	if nearestPointOnMeshNode.GetName()!= "nearestPointOnMesh1" {
		t.Errorf("got %s, wont nearestPointOnMesh1", nearestPointOnMeshNode.GetName())
	}

	if nearestPointOnMeshNode.renameCmd == nil {
		t.Errorf("got nil, wont *RenameCmd")
	}

	if uuid, err := nearestPointOnMeshNode.GetUUID(); err != nil &&
		uuid == "B2C0D7BA-4FA2-A74E-5D6A-8DAA23AAF72E" {
		t.Errorf("got %s, wont B2C0D7BA-4FA2-A74E-5D6A-8DAA23AAF72E", uuid)
	}
}

func TestLineComment(t *testing.T) {
	reader := strings.NewReader(getTestMa())

	mo, err := Unmarshal(reader)
	if err != nil {
		t.Error(err.Error())
	}

	if mo == nil {
		t.Error("got nil, wont *Object")
	}

	lineComments := mo.LineComments
	if len(lineComments) != 5 {
		t.Fatalf("got mo.LineComments len %d, wont 5", len(lineComments))
	}

	firstComment := lineComments[0]
	if firstComment.GetComment() != "Maya ASCII 2019 scene" {
		t.Errorf("got \"%s\", wont \"Maya ASCII 2019 scene\"", firstComment.GetComment())
	}
}

func TestBlockComment(t *testing.T) {
	reader := strings.NewReader(getTestMa())

	mo, err := Unmarshal(reader)
	if err != nil {
		t.Error(err.Error())
	}

	blockComments := mo.BlockComments
	if len(blockComments) != 2 {
		t.Fatalf("got mo.BlockComments len %d, wont 2.", len(blockComments))
	}

	if blockComments[0].GetComment() != "block comment1" {
		t.Errorf("got mo.BlockComments[0].Comment \"%s\", wont \"block comment1\".",
			blockComments[0].GetComment())
	}

	comment2 := `
	block
	comment2
`
	if blockComments[1].GetComment() != comment2 {
		t.Errorf("got mo.BlockComments[1].Comment \"%s\", wont \"block comment2\".",
			blockComments[1].GetComment())
	}
}

func TestNodes(t *testing.T) {
	reader := strings.NewReader(getTestMa())

	mo, err := Unmarshal(reader)
	if err != nil {
		t.Error(err.Error())
	}

	if mo == nil {
		t.Error("got nil, wont *Object")
	}

	if len(mo.Nodes) != 22 {
		t.Errorf("got %d, wont 22", len(mo.Nodes))
	}

	perspShape, err := mo.GetNode("perspShape")
	if err != nil {
		t.Errorf(err.Error())
	}

	if perspShape.GetName() != "perspShape" {
		t.Errorf("got %s, wont perspShape", perspShape.GetName())
	}

	if perspShape.GetType() != "camera" {
		t.Errorf("got %s, wont camera", perspShape.GetType())
	}

	if len((*perspShape).Attrs) != 7 {
		t.Errorf("got %d, wont 7", len(perspShape.Attrs))
	}

	for _, attr := range (*perspShape).Attrs {
		switch attr.GetName() {
		case ".v":
			if attr.GetAttrType() != SetAttrTypeBool {
				t.Errorf("got %v, wont cmd.SetAttrTypeBool", attr.GetAttrType())
			}
			if len(attr.GetAttrValue()) != 1 {
				t.Errorf("got %d, wont 1", len(attr.GetAttrValue()))
			}

			attrBools, err := ToAttrBool(attr.GetAttrValue())
			if err != nil {
				t.Error("AttrBool convert NG! " + err.Error())
			}

			if len(attrBools) != 1 {
				t.Errorf("got %d, wont 1", len(attrBools))
			}

			attrBool := attrBools[0]
			if attrBool.Bool() != false {
				t.Errorf("got %v, wont false", attrBool.Bool())
			}

		}
	}

	if uuid, err := perspShape.GetUUID(); err != nil &&
		uuid == "CFAE1109-4845-2AC4-5BC0-CB8FB886A568" {
		t.Errorf("got %v, wont CFAE1109-4845-2AC4-5BC0-CB8FB886A568", uuid)
	}
}
