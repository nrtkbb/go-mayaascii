package parser

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/nrtkbb/go-mayaascii/cmd"
)

func MakeWorkspace(c *cmd.Cmd) *cmd.Workspace {
	w := cmd.Workspace{Cmd: c}
	for i := 1; i < len(w.Token); i++ {
		switch w.Token[i] {
		case "-fr":
			w.FileRule = strings.Trim(w.Token[i+1], "\"")
			w.Place = strings.Trim(w.Token[i+2], "\"")
			i += 2
		default:
			panic("this option can not parse yet " + w.Token[i])
		}
	}
	return &w
}

func MakeRequires(c *cmd.Cmd) *cmd.Requires {
	// max Token = [requires, -nodeType, "typeName1", -dataType, "typeName2", "pluginName", "version"]
	// min Token = [requires, "pluginName", "version"]
	r := cmd.Requires{Cmd: c}
	r.PluginName = strings.Trim(r.Token[len(r.Token)-2], "\"")
	r.Version = strings.Trim(r.Token[len(r.Token)-1], "\"")
	if len(r.Token) > 3 {
		for i := 1; i < len(r.Token)-2; i += 2 {
			if r.Token[i] == "-nodeType" {
				r.NodeTypes = append(r.NodeTypes, strings.Trim(r.Token[i+1], "\""))
				continue
			}
			if r.Token[i] == "-dataType" {
				r.DataTypes = append(r.DataTypes, strings.Trim(r.Token[i+1], "\""))
				continue
			}
		}
	}
	return &r
}

func MakeConnectAttr(c *cmd.Cmd) (*cmd.ConnectAttr, error) {
	ca := &cmd.ConnectAttr{Cmd: c}
	for i := 1; i < len(ca.Token); i++ {
		switch ca.Token[i] {
		case "-f":
			ca.Force = true
		case "-l":
			var lock bool
			lock, err := isOnYesOrOffNo(ca.Token[i+1])
			if err != nil {
				log.Print(err)
				return nil, err
			}
			ca.Lock = &lock
			i++
		case "-na":
			ca.NextAvailable = true
		case "-rd":
			ca.ReferenceDest = &ca.Token[i+1]
			i++
		default:
			// trim "nodeName.attrName" -> {node: nodeName, attr: .attrName}
			dotIndex := strings.Index(ca.Token[i], ".")
			node := ca.Token[i][1:dotIndex]
			attr := ca.Token[i][dotIndex+1 : len(ca.Token[i])-1]
			if ca.SrcNode == "" {
				ca.SrcNode = node
				ca.SrcAttr = attr
			} else {
				ca.DstNode = node
				ca.DstAttr = attr
			}
		}
	}
	return ca, nil
}

func MakeCreateNode(c *cmd.Cmd) *cmd.CreateNode {
	n := &cmd.CreateNode{Cmd: c}
	n.NodeType = c.Token[1]
	for i := 2; i < len(n.Token); i++ {
		switch n.Token[i] {
		case "-n":
			i++
			n.NodeName = n.Token[i][1 : len(c.Token[i])-1]
		case "-p":
			i++
			p := n.Token[i][1 : len(c.Token[i])-1]
			n.Parent = &p
		case "-s":
			n.Shared = true
		case "-ss":
			n.SkipSelect = true
		}
	}
	return n
}

func MakeRename(c *cmd.Cmd) *cmd.Rename {
	r := &cmd.Rename{Cmd: c}
	for i := 1; i < len(r.Token); i++ {
		switch r.Token[i] {
		case "-uid":
			r.UUID = true
		case "-is":
			r.IgnoreShape = true
		default:
			if r.To == nil {
				to := strings.Trim(r.Token[i], "\"")
				r.To = &to
			} else {
				r.From = r.To
				to := strings.Trim(r.Token[i], "\"")
				r.To = &to
			}
		}
	}
	return r
}

func MakeSelect(c *cmd.Cmd) *cmd.Select {
	s := &cmd.Select{Cmd: c}
	for i := 1; i < len(c.Token); i++ {
		switch s.Token[i] {
		case "-add":
			s.Add = true
		case "-af":
			s.AddFirst = true
		case "-all":
			s.All = true
		case "-ado":
			s.AllDagObjects = true
		case "-adn":
			s.AllDependencyNodes = true
		case "-cl":
			s.Clear = true
		case "-cc":
			s.ContainerCentric = true
		case "-d":
			s.Deselect = true
		case "-hi":
			s.Hierarchy = true
		case "-ne":
			s.NoExpand = true
		case "-r":
			s.Replace = true
		case "-sym":
			s.Symmetry = true
		case "-sys":
			s.SymmetrySide = true
		case "-tgl":
			s.Toggle = true
		case "-vis":
			s.Visible = true
		default:
			s.Names = append(s.Names, s.Token[i])
		}
	}
	return s
}

func getAttrNameFromSetAttr(token *[]string) (int, string) {
	for i := 1; i < len(*token); i++ {
		t := (*token)[i]
		if t[0] == '"' &&
			t[1] == '.' &&
			t[len(t)-1] == '"' {
			// ".attr" -> .attr
			return i, t[1 : len(t)-1]
		}
	}
	return -1, ""
}

func isSameAttr(name1, name2 string) bool {
	if name1 == name2 {
		return true
	}
	if name2[len(name2)-1] != ']' {
		return false
	}
	openIdx2 := strings.LastIndex(name2, "[")
	if openIdx2 == -1 {
		return false
	}
	if name2[:openIdx2] == name1 {
		// .attrName[1] == .attrName
		// .attrName[0:499] == .attrName
		// .attrName[0].subName[0:499] == .attrName[0].subName
		return true
	}
	if name1[len(name1)-1] != ']' {
		return false
	}
	openIdx1 := strings.LastIndex(name1, "[")
	if openIdx1 == -1 {
		return false
	}
	if name1[:openIdx1] == name2[:openIdx2] {
		// .attrName[1] == .attrName[0]
		// .attrName[500:999] == .attrName[0:499]
		// .attrName[0].subName[500:999] == .attrName[0].subName[0:499]
		return true
	}
	return false
}

func fixSizeOver(end int, token *[]string) int {
	lenToken := len(*token)
	if lenToken < end {
		return lenToken
	}
	return end
}

var PrimitiveTypes = map[cmd.AttrType]struct{}{
	cmd.TypeInvalid: struct{}{},
	cmd.TypeBool:    struct{}{},
	cmd.TypeInt:     struct{}{},
	cmd.TypeDouble:  struct{}{},
}

func MakeSetAttr(c *cmd.Cmd, beforeSetAttr *cmd.SetAttr) (*cmd.SetAttr, error) {
	attrNameIdx, attrName := getAttrNameFromSetAttr(&c.Token)
	sa := &cmd.SetAttr{Cmd: c}
	sa.AttrName = attrName
	if beforeSetAttr != nil && isSameAttr(beforeSetAttr.AttrName, attrName) {
		sa.AlteredValue = beforeSetAttr.AlteredValue
		sa.Caching = beforeSetAttr.Caching
		sa.CapacityHint = beforeSetAttr.CapacityHint
		sa.ChannelBox = beforeSetAttr.ChannelBox
		sa.Clamp = beforeSetAttr.Clamp
		sa.Keyable = beforeSetAttr.Keyable
		sa.Lock = beforeSetAttr.Lock
		sa.Size = beforeSetAttr.Size
		sa.AttrType = beforeSetAttr.AttrType
		sa.Attr = beforeSetAttr.Attr
	}
	for i := 1; i < len(sa.Token); i++ {
		if i == attrNameIdx {
			continue
		}
		v := sa.Token[i]
		switch v {
		case "-av":
			sa.AlteredValue = true
		case "-ca":
			ca := true
			sa.Caching = &ca
		case "-ch":
			i++
			ch, err := strconv.Atoi(sa.Token[i])
			if err != nil {
				return nil, err
			}
			uch := uint(ch)
			sa.CapacityHint = &uch
		case "-cb":
			i++
			cb, err := isOnYesOrOffNo(sa.Token[i])
			if err != nil {
				return nil, err
			}
			sa.ChannelBox = &cb
		case "-c":
			sa.Clamp = true
		case "-k":
			i++
			k, err := isOnYesOrOffNo(sa.Token[i])
			if err != nil {
				return nil, err
			}
			sa.Keyable = &k
		case "-l":
			i++
			l, err := isOnYesOrOffNo(sa.Token[i])
			if err != nil {
				return nil, err
			}
			sa.Lock = &l
		case "-s":
			i++
			s, err := strconv.Atoi(sa.Token[i])
			if err != nil {
				return nil, err
			}
			us := uint(s)
			sa.Size = &us
		case "-type":
			i++
			attrType, err := MakeAttrType(&sa.Token, i)
			if err != nil {
				return nil, err
			}
			sa.AttrType = attrType
		default:
			if _, ok := PrimitiveTypes[sa.AttrType]; !ok {
				a, count, err := MakeAttr(&sa.Token, i, sa.Size, sa.AttrType)
				if err != nil {
					log.Println(sa.Token)
					log.Println(i)
					log.Println(sa.Size)
					log.Println(sa.AttrType)
					return nil, err
				}
				for _, aa := range a {
					sa.Attr = append(sa.Attr, aa)
				}
				if count == -1 {
					break
				}
				i += count
				break
			}
			b, err := isOnYesOrOffNo(v)
			if err == nil {
				sa.AttrType = cmd.TypeBool
				ab := cmd.AttrBool(b)
				abs := []cmd.Attr{&ab}
				sa.Attr = abs
				return sa, nil
			}

			var isInt bool
			for iii, token := range sa.Token[i:] {
				if strings.Contains(token, ".") ||
					strings.Contains(token, "e") {
					break
				}
				if iii+i == len(sa.Token)-1 {
					isInt = true
				}
			}
			if isInt && sa.AttrType != cmd.TypeDouble {
				intArray, err := ParseInts(sa.Token[i:]...)
				if err != nil {
					return nil, err
				}
				for _, ia := range intArray {
					ai := cmd.AttrInt(ia)
					sa.Attr = append(sa.Attr, &ai)
				}
				sa.AttrType = cmd.TypeInt
				return sa, nil
			}

			floatArray, err := ParseFloats(sa.Token[i:]...)
			if err != nil {
				return nil, err
			}
			if sa.AttrType == cmd.TypeInt {
				if 0 < len(sa.Attr) {
					first := sa.Attr[0]
					_, ok := first.(*cmd.AttrInt)
					if !ok {
						return nil, errors.New(
							fmt.Sprintf("invalid pattern %v and %v",
								sa.Attr, floatArray))
					}
					floatArrayAttr := make([]cmd.Attr, len(sa.Attr))
					for i, ia := range sa.Attr {
						ai, _ := ia.(*cmd.AttrInt)
						af := cmd.AttrFloat(float64(ai.Int()))
						floatArrayAttr[i] = &af
					}
					sa.Attr = floatArrayAttr
				}
			}
			sa.AttrType = cmd.TypeDouble
			for _, f := range floatArray {
				af := cmd.AttrFloat(f)
				sa.Attr = append(sa.Attr, &af)
			}
			return sa, nil
		}
	}
	return sa, nil
}

func isOnYesOrOffNo(t string) (bool, error) {
	if t == "on" || t == "yes" || t == "true" {
		return true, nil
	}
	if t == "off" || t == "no" || t == "false" {
		return false, nil
	}
	return false,
		errors.New(
			fmt.Sprintf("this string is not bool word. \"%s\"", t))
}

func ParseInts(token ...string) ([]int, error) {
	var result []int
	for _, t := range token {
		i, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			value, err := isOnYesOrOffNo(t)
			if err != nil {
				panic(err)
			}
			if value {
				result = append(result, 1)
			} else {
				result = append(result, 0)
			}
		} else {
			result = append(result, int(i))
		}
	}
	return result, nil
}

func ParseFloats(token ...string) ([]float64, error) {
	result := make([]float64, len(token))
	for i := 0; i < len(token); i++ {
		f, err := strconv.ParseFloat(token[i], 64)
		if err != nil {
			return nil, err
		}
		result[i] = f
	}
	return result, nil
}

func MakeShort2Long2(
	token *[]string,
	start int,
	size *uint,
	at *cmd.AttrType) ([]cmd.Attr, int, error) {
	var end int
	if size != nil {
		end = fixSizeOver(start+(2*int(*size)), token)
	} else {
		end = start + 2
	}
	v, err := ParseInts((*token)[start:end]...)
	if err != nil {
		return nil, 0, err
	}
	if at != nil && *at == cmd.TypeShort2 {
		s2 := make([]cmd.AttrShort2, (end-start)/2)
		for i := 0; i < len(s2); i++ {
			s2[i][0] = v[i*2]
			s2[i][1] = v[i*2+1]
		}
		a := make([]cmd.Attr, len(s2))
		for i, s := range s2 {
			a[i] = &s
		}
		return a, end - start, nil
	} else {
		l2 := make([]cmd.AttrLong2, (end-start)/2)
		for i := 0; i < len(l2); i++ {
			l2[i][0] = v[i*2]
			l2[i][1] = v[i*2+1]
		}
		a := make([]cmd.Attr, len(l2))
		for i, l := range l2 {
			a[i] = &l
		}
		return a, end - start, nil
	}
}

func MakeShort3Long3(token *[]string, start int, size *uint, at *cmd.AttrType) ([]cmd.Attr, int, error) {
	var end int
	if size != nil {
		end = fixSizeOver(start+(3*int(*size)), token)
	} else {
		end = start + 3
	}
	v, err := ParseInts((*token)[start:end]...)
	if err != nil {
		return nil, 0, err
	}
	if at != nil && *at == cmd.TypeShort3 {
		s3 := make([]cmd.AttrShort3, (end-start)/3)
		for i := 0; i < len(s3); i++ {
			s3[i][0] = v[i*3]
			s3[i][1] = v[i*3+1]
			s3[i][2] = v[i*3+2]
		}
		a := make([]cmd.Attr, len(s3))
		for i, s := range s3 {
			a[i] = &s
		}
		return a, end - start, nil
	} else {
		l3 := make([]cmd.AttrLong3, (end-start)/3)
		for i := 0; i < len(l3); i++ {
			l3[i][0] = v[i*3]
			l3[i][1] = v[i*3+1]
			l3[i][2] = v[i*3+2]
		}
		a := make([]cmd.Attr, len(l3))
		for i, l := range l3 {
			a[i] = &l
		}
		return a, end - start, nil
	}
}

func MakeInt32Array(token *[]string, start int) ([]cmd.Attr, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	a := make([]cmd.Attr, 1)
	if numberOfArray != 0 {
		result, err := ParseInts((*token)[start+1 : start+1+numberOfArray]...)
		if err != nil {
			return nil, 0, err
		}
		ia := cmd.AttrInt32Array(result)
		a[0] = &ia
	} else {
		ia := cmd.AttrInt32Array{}
		a[0] = &ia
	}
	return a, 1 + numberOfArray, nil
}

func MakeFloat2Double2(token *[]string, start int, size *uint, at *cmd.AttrType) ([]cmd.Attr, int, error) {
	var end int
	if size != nil {
		end = fixSizeOver(start+(2*int(*size)), token)
	} else {
		end = start + 2
	}
	v, err := ParseFloats((*token)[start:end]...)
	if err != nil {
		return nil, 0, err
	}
	if at != nil && *at == cmd.TypeFloat2 {
		f2 := make([]cmd.AttrFloat2, (end-start)/2)
		for i := 0; i < len(f2); i++ {
			f2[i][0] = v[i*2]
			f2[i][1] = v[i*2+1]
		}
		a := make([]cmd.Attr, len(f2))
		for i, f := range f2 {
			a[i] = &f
		}
		return a, end - start, nil
	} else {
		d2 := make([]cmd.AttrDouble2, (end-start)/2)
		for i := 0; i < len(d2); i++ {
			d2[i][0] = v[i*2]
			d2[i][1] = v[i*2+1]
		}
		a := make([]cmd.Attr, len(d2))
		for i, d := range d2 {
			a[i] = &d
		}
		return a, end - start, nil
	}
}

func MakeFloat3Double3(token *[]string, start int, size *uint, at *cmd.AttrType) ([]cmd.Attr, int, error) {
	var end int
	if size != nil {
		end = fixSizeOver(start+(3*int(*size)), token)
	} else {
		end = start + 3
	}
	v, err := ParseFloats((*token)[start:end]...)
	if err != nil {
		return nil, 0, err
	}
	if at != nil && *at == cmd.TypeFloat3 {
		f3 := make([]cmd.AttrFloat3, (end-start)/3)
		for i := 0; i < len(f3); i++ {
			f3[i][0] = v[i*3]
			f3[i][1] = v[i*3+1]
			f3[i][2] = v[i*3+2]
		}
		a := make([]cmd.Attr, len(f3))
		for i, f := range f3 {
			a[i] = &f
		}
		return a, end - start, nil
	} else {
		d3 := make([]cmd.AttrDouble3, (end-start)/3)
		for i := 0; i < len(d3); i++ {
			d3[i][0] = v[i*3]
			d3[i][1] = v[i*3+1]
			d3[i][2] = v[i*3+2]
		}
		a := make([]cmd.Attr, len(d3))
		for i, d := range d3 {
			a[i] = &d
		}
		return a, end - start, nil
	}
}

func MakeDoubleArray(token *[]string, start int) ([]cmd.Attr, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	a := make([]cmd.Attr, 1)
	if numberOfArray != 0 {
		f, err := ParseFloats((*token)[start+1 : start+1+numberOfArray]...)
		if err != nil {
			return nil, 0, err
		}
		da := cmd.AttrDoubleArray(f)
		a[0] = &da
	} else {
		da := cmd.AttrDoubleArray{}
		a[0] = &da
	}
	return a, 1 + numberOfArray, nil
}

func MakeMatrix(token *[]string, start int) ([]cmd.Attr, int, error) {
	mat4x4, err := ParseFloats((*token)[start : start+16]...)
	if err != nil {
		return nil, 0, err
	}
	a := make([]cmd.Attr, 1)
	a[0] = &cmd.AttrMatrix{
		mat4x4[0], mat4x4[1], mat4x4[2], mat4x4[3],
		mat4x4[4], mat4x4[5], mat4x4[6], mat4x4[7],
		mat4x4[8], mat4x4[9], mat4x4[10], mat4x4[11],
		mat4x4[12], mat4x4[13], mat4x4[14], mat4x4[15],
	}
	return a, 16, nil
}

func MakeMatrixXform(token *[]string, start int) ([]cmd.Attr, int, error) {
	// type:
	// string double double double
	// double double double
	// integer
	// double double double
	// double double double
	// double double double
	// double double double
	// double double double
	// double double double
	// double double double double
	// double double double double
	// double double double
	// boolean
	// mean:
	// xform scaleX scaleY scaleZ
	// rotateX rotateY rotateZ
	// rotationOrder (0=XYZ, 1=YZX, 2=ZXY, 3=XZY, 4=YXZ, 5=ZYX)
	// translateX translateY translateZ
	// shearXY shearXZ shearYZ
	// scalePivotX scalePivotY scalePivotZ
	// scaleTranslationX scaleTranslationY scaleTranslationZ
	// rotatePivotX rotatePivotY rotatePivotZ
	// rotateTranslationX rotateTranslationY rotateTranslationZ
	// rotateOrientW rotateOrientX rotateOrientY rotateOrientZ
	// jointOrientW jointOrientX jointOrientY jointOrientZ
	// inverseParentScaleX inverseParentScaleY inverseParentScaleZ
	// compensateForParentScale
	floats, err := ParseFloats((*token)[start+1 : start+37]...)
	if err != nil {
		return nil, 0, err
	}
	rotateOrder, err := cmd.ConvertAttrRotateOrder(int(floats[6]))
	if err != nil {
		return nil, 0, err
	}
	componentForParentScale, err := isOnYesOrOffNo((*token)[start+37])
	if err != nil {
		return nil, 0, err
	}
	mx := cmd.AttrMatrixXform{
		Scale:                    cmd.AttrVector{X: floats[0], Y: floats[1], Z: floats[2]},
		Rotate:                   cmd.AttrVector{X: floats[3], Y: floats[4], Z: floats[5]},
		RotateOrder:              rotateOrder,
		Translate:                cmd.AttrVector{X: floats[7], Y: floats[8], Z: floats[9]},
		Shear:                    cmd.AttrShear{XY: floats[10], XZ: floats[11], YZ: floats[12]},
		ScalePivot:               cmd.AttrVector{X: floats[13], Y: floats[14], Z: floats[15]},
		ScaleTranslate:           cmd.AttrVector{X: floats[16], Y: floats[17], Z: floats[18]},
		RotatePivot:              cmd.AttrVector{X: floats[19], Y: floats[20], Z: floats[21]},
		RotateTranslation:        cmd.AttrVector{X: floats[22], Y: floats[23], Z: floats[24]},
		RotateOrient:             cmd.AttrOrient{W: floats[25], X: floats[26], Y: floats[27], Z: floats[28]},
		JointOrient:              cmd.AttrOrient{W: floats[29], X: floats[30], Y: floats[31], Z: floats[32]},
		InverseParentScale:       cmd.AttrVector{X: floats[33], Y: floats[34], Z: floats[35]},
		CompensateForParentScale: componentForParentScale,
	}
	a := make([]cmd.Attr, 1)
	a[0] = &mx
	return a, 38, nil
}

func MakePointArray(token *[]string, start int) ([]cmd.Attr, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	a := make([]cmd.Attr, 1)
	if numberOfArray != 0 {
		f, err := ParseFloats((*token)[start+1 : start+1+(numberOfArray*4)]...)
		if err != nil {
			return nil, 0, err
		}
		pa := make([]cmd.AttrPoint, numberOfArray)
		for i := 0; i < numberOfArray; i++ {
			pa[i].X = f[i*4]
			pa[i].Y = f[i*4+1]
			pa[i].Z = f[i*4+2]
			pa[i].W = f[i*4+3]
		}
		paa := cmd.AttrPointArray(pa)
		a[0] = &paa
	} else {
		paa := cmd.AttrPointArray{}
		a[0] = &paa
	}
	return a, 1 + (numberOfArray * 4), nil
}

func MakeVectorArray(token *[]string, start int) ([]cmd.Attr, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	a := make([]cmd.Attr, 1)
	if numberOfArray != 0 {
		f, err := ParseFloats((*token)[start+1 : start+1+numberOfArray]...)
		if err != nil {
			return nil, 0, err
		}
		va := make([]cmd.AttrVector, numberOfArray)
		for i := 0; i < numberOfArray; i += 3 {
			va[i] = cmd.AttrVector{
				X: f[i],
				Y: f[i+1],
				Z: f[i+2],
			}
		}
		vaa := cmd.AttrVectorArray(va)
		a[0] = &vaa
	} else {
		vaa := cmd.AttrVectorArray{}
		a[0] = &vaa
	}
	return a, 1 + (numberOfArray * 3), nil
}

func MakeString(token *[]string, start int) ([]cmd.Attr, int, error) {
	s := cmd.AttrString((*token)[start][1 : len((*token)[start])-1])
	a := []cmd.Attr{&s}
	return a, 1, nil
}

func MakeStringArray(token *[]string, start int) ([]cmd.Attr, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	sa := cmd.AttrStringArray((*token)[start+1 : start+1+numberOfArray])
	for i, s := range sa {
		sa[i] = s[1 : len(s)-1]
	}
	a := []cmd.Attr{&sa}
	return a, 1 + numberOfArray, nil
}

func MakeSphere(token *[]string, start int) ([]cmd.Attr, int, error) {
	s, err := strconv.ParseFloat((*token)[start], 64)
	if err != nil {
		return nil, 0, err
	}
	sp := cmd.AttrSphere(s)
	a := []cmd.Attr{&sp}
	return a, 1, nil
}

func MakeCone(token *[]string, start int) ([]cmd.Attr, int, error) {
	f, err := ParseFloats((*token)[start : start+2]...)
	if err != nil {
		return nil, 0, err
	}
	c := cmd.AttrCone{
		ConeAngle: f[0],
		ConeCap:   f[1],
	}
	a := []cmd.Attr{&c}
	return a, 2, nil
}

func MakeReflectanceRGB(token *[]string, start int) ([]cmd.Attr, int, error) {
	f, err := ParseFloats((*token)[start : start+3]...)
	if err != nil {
		return nil, 0, err
	}
	a := make([]cmd.Attr, 1)
	a[0] = &cmd.AttrReflectanceRGB{
		RedReflect:   f[0],
		GreenReflect: f[1],
		BlueReflect:  f[2],
	}
	return a, 3, nil
}

func MakeSpectrumRGB(token *[]string, start int) ([]cmd.Attr, int, error) {
	f, err := ParseFloats((*token)[start : start+3]...)
	if err != nil {
		return nil, 0, err
	}
	a := make([]cmd.Attr, 1)
	a[0] = &cmd.AttrSpectrumRGB{
		RedSpectrum:   f[0],
		GreenSpectrum: f[1],
		BlueSpectrum:  f[2],
	}
	return a, 3, nil
}

func MakeComponentList(token *[]string, start int) ([]cmd.Attr, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	var cl cmd.AttrComponentList
	for _, c := range (*token)[start+1 : start+1+numberOfArray] {
		cl = append(cl, strings.Trim(c, "\""))
	}
	a := []cmd.Attr{&cl}
	return a, 1 + numberOfArray, nil
}

func MakeAttributeAlias(token *[]string, start int) ([]cmd.Attr, int, error) {
	if (*token)[start] != "{" {
		return nil, 0, errors.New("there was no necessary token")
	}
	var aaa []cmd.AttrAttributeAlias
	for i := start + 1; i < len(*token); i += 2 {
		if (*token)[i] == "}" {
			break
		}
		aaa = append(aaa, cmd.AttrAttributeAlias{
			NewAlias:    (*token)[i],
			CurrentName: (*token)[i+1],
		})
	}
	a := make([]cmd.Attr, len(aaa))
	for i := range aaa {
		a[i] = &aaa[i]
	}
	return a, 2 + (len(aaa) * 2), nil
}

func MakeNurbsCurve(token *[]string, start int) ([]cmd.Attr, int, error) {
	i1, err := ParseInts((*token)[start : start+3]...)
	if err != nil {
		return nil, 0, err
	}
	degree := i1[0]
	spans := i1[1]
	form, err := cmd.ConvertAttrFormType(i1[2]) // open(0), closed(1), periodic(2)
	if err != nil {
		return nil, 0, err
	}
	isRational, err := isOnYesOrOffNo((*token)[start+3])
	if err != nil {
		return nil, 0, err
	}
	i2, err := ParseInts((*token)[start+4 : start+6]...)
	if err != nil {
		return nil, 0, err
	}
	dimension := i2[0]
	knotCount := i2[1]
	kv, err := ParseFloats((*token)[start+6 : start+6+knotCount]...)
	if err != nil {
		return nil, 0, err
	}
	cvCount, err := strconv.Atoi((*token)[start+6+knotCount])
	if err != nil {
		return nil, 0, err
	}
	divideCv := 2
	if isRational {
		divideCv += 1
	}
	if dimension == 3 {
		divideCv += 1
	}
	cv, err := ParseFloats((*token)[start+7+knotCount : start+7+knotCount+(cvCount*divideCv)]...)
	if err != nil {
		return nil, 0, err
	}
	cvValues := make([]cmd.AttrCvValue, len(cv)/divideCv)
	for i := 0; i < cvCount; i++ {
		cvValues[i].X = cv[i*divideCv]
		cvValues[i].Y = cv[i*divideCv+1]
		if dimension == 3 {
			cvValues[i].Z = &cv[i*divideCv+2]
		}
		if isRational {
			if dimension == 3 {
				cvValues[i].W = &cv[i*divideCv+3]
			} else {
				cvValues[i].W = &cv[i*divideCv+2]
			}
		}
	}
	a := make([]cmd.Attr, 1)
	a[0] = &cmd.AttrNurbsCurve{
		Degree:     degree,
		Spans:      spans,
		Form:       form,
		IsRational: isRational,
		Dimension:  dimension,
		KnotValues: kv,
		CvValues:   cvValues,
	}
	count := 7 + knotCount + (cvCount * divideCv)
	return a, count, nil
}

func MakeNurbsSurface(token *[]string, start int) ([]cmd.Attr, int, error) {
	i1, err := ParseInts((*token)[start : start+4]...)
	if err != nil {
		return nil, 0, err
	}
	uDegree := i1[0]
	vDegree := i1[1]
	uForm, err := cmd.ConvertAttrFormType(i1[2])
	if err != nil {
		return nil, 0, err
	}
	vForm, err := cmd.ConvertAttrFormType(i1[3])
	if err != nil {
		return nil, 0, err
	}
	isRational, err := isOnYesOrOffNo((*token)[start+4])
	uKnotCount, err := strconv.Atoi((*token)[start+5])
	if err != nil {
		return nil, 0, err
	}
	uKnotValues, err := ParseFloats((*token)[start+6 : start+6+uKnotCount]...)
	if err != nil {
		return nil, 0, err
	}
	vKnotCount, err := strconv.Atoi((*token)[start+6+uKnotCount])
	if err != nil {
		return nil, 0, err
	}
	vKnotValues, err := ParseFloats(
		(*token)[start+7+uKnotCount : start+7+uKnotCount+vKnotCount]...)
	if err != nil {
		return nil, 0, err
	}
	var isTrim *bool
	if (*token)[start+7+uKnotCount+vKnotCount] == "\"TRIM\"" {
		v := true
		isTrim = &v
	} else if (*token)[start+7+uKnotCount+vKnotCount] == "\"NOTRIM\"" {
		v := false
		isTrim = &v
	}
	cvStart := start + 7 + uKnotCount + vKnotCount
	if isTrim != nil {
		cvStart++
	}
	cvCount, err := strconv.Atoi((*token)[cvStart])
	if err != nil {
		return nil, 0, err
	}
	divideCv := 3
	if isRational {
		divideCv++
	}
	cv, err := ParseFloats((*token)[cvStart+1 : cvStart+1+(cvCount*divideCv)]...)
	if err != nil {
		return nil, 0, err
	}
	cvValue := make([]cmd.AttrCvValue, cvCount)
	for i := 0; i < cvCount; i++ {
		cvValue[i].X = cv[i*divideCv]
		cvValue[i].Y = cv[i*divideCv+1]
		cvValue[i].Z = &cv[i*divideCv+2]
		if isRational {
			cvValue[i].W = &cv[i*divideCv+3]
		}
	}
	a := make([]cmd.Attr, 1)
	a[0] = &cmd.AttrNurbsSurface{
		UDegree:     uDegree,
		VDegree:     vDegree,
		UForm:       uForm,
		VForm:       vForm,
		IsRational:  isRational,
		UKnotValues: uKnotValues,
		VKnotValues: vKnotValues,
		IsTrim:      isTrim,
		CvValues:    cvValue,
	}
	count := (cvStart + (cvCount * divideCv)) - start
	return a, count, nil
}

func MakeNurbsTrimface(token *[]string, start int) ([]cmd.Attr, int, error) {
	// TODO: Waiting for Autodesk
	a := []cmd.Attr{&cmd.AttrNurbsTrimface{}}
	return a, -1, nil
}

func MakeCountInt(token *[]string, start int) ([]int, error) {
	count, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, err
	}
	result, err := ParseInts((*token)[start+1 : start+1+count]...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func MakePolyFace(token *[]string, start int, size *uint) ([]cmd.Attr, int, error) {
	switchNumber := start
	var pfs []cmd.AttrPolyFaces
	var fCount uint
	for _, v := range (*token)[start:] {
		if v == "f" {
			fCount++
		}
	}
	if size != nil {
		s := *size
		if fCount < s {
			s = fCount
		}
		pfs = make([]cmd.AttrPolyFaces, s)
	}
	i := -1
	loop := true
	for loop && len(*token) > switchNumber {
		switch (*token)[switchNumber] {
		case "f":
			fe, err := MakeCountInt(token, switchNumber+1)
			if err != nil {
				log.Println("error case f")
				log.Printf("token number is %d", switchNumber+1)
				return nil, 0, err
			}
			i++
			if i >= len(pfs) {
				pf := cmd.AttrPolyFaces{}
				pfs = append(pfs, pf)
			}
			pfs[i].FaceEdge = fe
			switchNumber += 2 + len(fe)
		case "h":
			he, err := MakeCountInt(token, switchNumber+1)
			if err != nil {
				return nil, 0, err
			}
			pfs[i].HoleEdge = he
			switchNumber += 2 + len(he)
		case "fc":
			fc, err := MakeCountInt(token, switchNumber+1)
			if err != nil {
				return nil, 0, err
			}
			pfs[i].FaceColor = fc
			switchNumber += 2 + len(fc)
		case "mc":
			colorIndex, err := strconv.ParseInt((*token)[switchNumber+1], 10, 64)
			if err != nil {
				return nil, 0, err
			}
			colorIDs, err := MakeCountInt(token, switchNumber+2)
			if err != nil {
				return nil, 0, err
			}
			mc := cmd.AttrMultiColor{
				ColorIndex: int(colorIndex),
				ColorIDs:   colorIDs,
			}
			pfs[i].MultiColor = append(pfs[i].MultiColor, mc)
			switchNumber += 3 + len(colorIDs)
		case "mu":
			var fuv cmd.AttrFaceUV
			uvSet, err := strconv.Atoi((*token)[switchNumber+1])
			if err != nil {
				return nil, 0, err
			}
			fuv.UVSet = uvSet
			uv, err := MakeCountInt(token, switchNumber+2)
			if err != nil {
				return nil, 0, err
			}
			fuv.FaceUV = uv
			pfs[i].FaceUV = append(pfs[i].FaceUV, fuv)
			switchNumber += 3 + len(uv)
		default:
			loop = false
			break
		}
	}
	a := make([]cmd.Attr, len(pfs))
	for i := range pfs {
		a[i] = &pfs[i]
	}
	return a, switchNumber - start, nil
}

func MakeDataPolyComponent(token *[]string, start int) ([]cmd.Attr, int, error) {
	if "Index_Data" != (*token)[start] {
		return nil, 0, errors.New(
			"since the Index_Data did not exist, " +
				"this token is an unknown dataPolyComponent")
	}
	var dpc cmd.AttrDataPolyComponent
	switch (*token)[start+1] {
	case "Edge":
		dpc.PolyComponentType = cmd.DPCedge
	case "Face":
		dpc.PolyComponentType = cmd.DPCface
	case "Vertex":
		dpc.PolyComponentType = cmd.DPCvertex
	case "UV":
		dpc.PolyComponentType = cmd.DPCuv
	default:
		return nil, 0, errors.New(
			"it is an unknown dataPolyComponent " +
				"that is neither Edge, Face, Vertex, UV")
	}
	count, err := strconv.Atoi((*token)[start+2])
	if err != nil {
		return nil, 0, err
	}
	dpc.IndexValue = map[int]float64{}
	for i := 0; i < count; i++ {
		index, err := strconv.Atoi((*token)[start+3+(i*2)])
		if err != nil {
			return nil, 0, err
		}
		value, err := strconv.ParseFloat((*token)[start+4+(i*2)], 64)
		if err != nil {
			return nil, 0, err
		}
		dpc.IndexValue[index] = value
	}
	a := []cmd.Attr{&dpc}
	return a, 3 + (count * 2), nil
}

func MakeMesh(_ *[]string, _ int) ([]cmd.Attr, int, error) {
	// Not Implement
	var a []cmd.Attr
	return a, -1, nil
}

func MakeLattice(token *[]string, start int) ([]cmd.Attr, int, error) {
	c, err := ParseInts((*token)[start : start+4]...)
	if err != nil {
		return nil, 0, err
	}
	la := cmd.AttrLattice{
		DivisionS: c[0],
		DivisionT: c[1],
		DivisionU: c[2],
	}
	la.Points = make([]cmd.AttrLatticePoint, c[3])
	for i := 0; i < c[3]*3; i += 3 {
		p, err := ParseFloats((*token)[start+4+i : start+4+i+3]...)
		if err != nil {
			return nil, 0, err
		}
		la.Points[i/3].S = p[0]
		la.Points[i/3].T = p[1]
		la.Points[i/3].U = p[2]
	}
	a := []cmd.Attr{&la}
	return a, 4 + (c[3] * 3), nil
}

func MakeAttr(token *[]string, start int, size *uint, attrType cmd.AttrType) ([]cmd.Attr, int, error) {
	switch attrType {
	case cmd.TypeShort2, cmd.TypeLong2:
		return MakeShort2Long2(token, start, size, &attrType)
	case cmd.TypeShort3, cmd.TypeLong3:
		return MakeShort3Long3(token, start, size, &attrType)
	case cmd.TypeInt32Array:
		return MakeInt32Array(token, start)
	case cmd.TypeFloat2, cmd.TypeDouble2:
		return MakeFloat2Double2(token, start, size, &attrType)
	case cmd.TypeFloat3, cmd.TypeDouble3:
		return MakeFloat3Double3(token, start, size, &attrType)
	case cmd.TypeDoubleArray:
		return MakeDoubleArray(token, start)
	case cmd.TypeMatrix:
		return MakeMatrix(token, start)
	case cmd.TypeMatrixXform:
		return MakeMatrixXform(token, start)
	case cmd.TypePointArray:
		return MakePointArray(token, start)
	case cmd.TypeVectorArray:
		return MakeVectorArray(token, start)
	case cmd.TypeString:
		return MakeString(token, start)
	case cmd.TypeStringArray:
		return MakeStringArray(token, start)
	case cmd.TypeSphere:
		return MakeSphere(token, start)
	case cmd.TypeCone:
		return MakeCone(token, start)
	case cmd.TypeReflectanceRGB:
		return MakeReflectanceRGB(token, start)
	case cmd.TypeSpectrumRGB:
		return MakeSpectrumRGB(token, start)
	case cmd.TypeComponentList:
		return MakeComponentList(token, start)
	case cmd.TypeAttributeAlias:
		return MakeAttributeAlias(token, start)
	case cmd.TypeNurbsCurve:
		return MakeNurbsCurve(token, start)
	case cmd.TypeNurbsSurface:
		return MakeNurbsSurface(token, start)
	case cmd.TypeNurbsTrimface:
		return MakeNurbsTrimface(token, start)
	case cmd.TypePolyFaces:
		return MakePolyFace(token, start, size)
	case cmd.TypeDataPolyComponent:
		return MakeDataPolyComponent(token, start)
	case cmd.TypeMesh:
		return MakeMesh(token, start)
	case cmd.TypeLattice:
		return MakeLattice(token, start)
	}
	return nil, 0, nil
}

func MakeAttrType(token *[]string, start int) (cmd.AttrType, error) {
	typeString := (*token)[start][1 : len((*token)[start])-1]
	switch typeString {
	case "short2":
		return cmd.TypeShort2, nil
	case "short3":
		return cmd.TypeShort3, nil
	case "long2":
		return cmd.TypeLong2, nil
	case "long3":
		return cmd.TypeLong3, nil
	case "Int32Array":
		return cmd.TypeInt32Array, nil
	case "float2":
		return cmd.TypeFloat2, nil
	case "double2":
		return cmd.TypeDouble2, nil
	case "float3":
		return cmd.TypeFloat3, nil
	case "double3":
		return cmd.TypeDouble3, nil
	case "doubleArray":
		return cmd.TypeDoubleArray, nil
	case "matrix":
		typeString2 := (*token)[start+1]
		if typeString2 == "\"xform\"" {
			return cmd.TypeMatrixXform, nil
		} else {
			return cmd.TypeMatrix, nil
		}
	case "pointArray":
		return cmd.TypePointArray, nil
	case "vectorArray":
		return cmd.TypeVectorArray, nil
	case "string":
		return cmd.TypeString, nil
	case "stringArray":
		return cmd.TypeStringArray, nil
	case "sphere":
		return cmd.TypeSphere, nil
	case "cone":
		return cmd.TypeCone, nil
	case "reflectanceRGB":
		return cmd.TypeReflectanceRGB, nil
	case "spectrumRGB":
		return cmd.TypeSpectrumRGB, nil
	case "componentList":
		return cmd.TypeComponentList, nil
	case "attributeAlias":
		return cmd.TypeAttributeAlias, nil
	case "nurbsCurve":
		return cmd.TypeNurbsCurve, nil
	case "nurbsSurface":
		return cmd.TypeNurbsSurface, nil
	case "nurbsTrimface":
		return cmd.TypeNurbsTrimface, nil
	case "polyFaces":
		return cmd.TypePolyFaces, nil
	case "dataPolyComponent":
		return cmd.TypeDataPolyComponent, nil
	case "mesh":
		return cmd.TypeMesh, nil
	case "lattice":
		return cmd.TypeLattice, nil
	}
	return cmd.TypeInvalid, errors.New("Invalid type " + typeString)
}

func MakeAddAttr(c *cmd.Cmd) *cmd.AddAttr {
	aa := &cmd.AddAttr{Cmd: c}
	// TODO: Do finish!
	return aa
}
