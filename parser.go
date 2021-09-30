package mayaascii

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func ParseLineComment(c *Cmd) *LineCommentCmd {
	lc := LineCommentCmd{Cmd: c}
	lc.Comment = lc.Token[1]
	return &lc
}

func ParseFile(c *Cmd) *FileCmd {
	// [file, -rdi, 1, -ns, "ns", -rfn, "nsRN", -op, "v=0;", -typ, "mayaAscii", "path/to/file.ma"]
	// [file, -r, -ns, "namespace", -dr, 1, -rfn, "nsRN", -op, "v=0;", -typ, "mayaAscii", "path/to/file.ma"]
	f := FileCmd{Cmd: c}
	f.Path = strings.Trim(f.Token[len(f.Token)-1], "\"")
	var err error
	for i := 1; i < len(f.Token)-1; i++ {
		switch f.Token[i] {
		case "-rdi":
			f.ReferenceDepthInfo, err = strconv.Atoi(f.Token[i+1])
			if err != nil {
				log.Print(err)
				return &f
			}
			i++
		case "-ns":
			f.Namespace = strings.Trim(f.Token[i+1], "\"")
			i++
		case "-rfn":
			f.ReferenceNode = strings.Trim(f.Token[i+1], "\"")
			i++
		case "-op":
			f.Options = strings.Trim(f.Token[i+1], "\"")
			i++
		case "-typ":
			f.Type = strings.Trim(f.Token[i+1], "\"")
			i++
		case "-r":
			f.Reference = true
		case "-dr":
			var dr int
			dr, err = strconv.Atoi(f.Token[i+1])
			if dr == 1 {
				f.DeferReference = true
			}
			i++
		}
	}
	return &f
}

func ParseFileInfo(c *Cmd) *FileInfoCmd {
	fi := &FileInfoCmd{Cmd: c}
	fi.Name = strings.Trim(fi.Token[1], "\"")
	fi.Value = strings.Trim(fi.Token[2], "\"")
	return fi
}

func ParseWorkspace(c *Cmd) *WorkspaceCmd {
	w := WorkspaceCmd{Cmd: c}
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

func ParseRequires(c *Cmd) *RequiresCmd {
	// max Token = [requires, -nodeType, "typeName1", -dataType, "typeName2", "pluginName", "version"]
	// min Token = [requires, "pluginName", "version"]
	r := RequiresCmd{Cmd: c}
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

func ParseConnectAttr(c *Cmd) (*ConnectAttrCmd, error) {
	ca := &ConnectAttrCmd{Cmd: c}
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

func ParseCreateNode(c *Cmd) *CreateNodeCmd {
	n := &CreateNodeCmd{Cmd: c}
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

func ParseRename(c *Cmd) *RenameCmd {
	r := &RenameCmd{Cmd: c}
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

func ParseSelect(c *Cmd) *SelectCmd {
	s := &SelectCmd{Cmd: c}
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
		if name2[:openIdx2] == ".phl" {
			// setAttr ".phl[1476]" -type "matrix" 0.35087719298245618 0 0 0 0 0.35087719298245618 0 0
			//     0 0 0.35087719298245618 0 0 0 0 1;
			// setAttr ".phl[1504]" 0;
			return false
		}
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

var PrimitiveTypes = map[AttrType]struct{}{
	TypeInvalid: struct{}{},
	TypeBool:    struct{}{},
	TypeInt:     struct{}{},
	TypeDouble:  struct{}{},
}

func ParseSetAttr(c *Cmd, beforeSetAttr *SetAttrCmd) (*SetAttrCmd, error) {
	attrNameIdx, attrName := getAttrNameFromSetAttr(&c.Token)
	sa := &SetAttrCmd{Cmd: c}
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
			attrType, err := ParseAttrType(&sa.Token, i)
			if err != nil {
				return nil, err
			}
			sa.AttrType = attrType
		default:
			if _, ok := PrimitiveTypes[sa.AttrType]; !ok {
				a, count, err := ParseAttr(&sa.Token, i, sa.Size, sa.AttrType)
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
				sa.AttrType = TypeBool
				ab := AttrBool(b)
				abs := []AttrValue{&ab}
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
			if isInt && sa.AttrType != TypeDouble {
				intArray, err := ParseInts(sa.Token[i:]...)
				if err != nil {
					return nil, err
				}
				for _, ia := range intArray {
					ai := AttrInt(ia)
					sa.Attr = append(sa.Attr, &ai)
				}
				sa.AttrType = TypeInt
				return sa, nil
			}

			floatArray, err := ParseFloats(sa.Token[i:]...)
			if err != nil {
				return nil, err
			}
			if sa.AttrType == TypeInt {
				if 0 < len(sa.Attr) {
					first := sa.Attr[0]
					_, ok := first.(*AttrInt)
					if !ok {
						return nil, errors.New(
							fmt.Sprintf("invalid pattern %v and %v",
								sa.Attr, floatArray))
					}
					floatArrayAttr := make([]AttrValue, len(sa.Attr))
					for i, ia := range sa.Attr {
						ai, _ := ia.(*AttrInt)
						af := AttrFloat(float64(ai.Int()))
						floatArrayAttr[i] = &af
					}
					sa.Attr = floatArrayAttr
				}
			}
			sa.AttrType = TypeDouble
			for _, f := range floatArray {
				af := AttrFloat(f)
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

func ParseShort2Long2(
	token *[]string,
	start int,
	size *uint,
	at *AttrType) ([]AttrValue, int, error) {
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
	if at != nil && *at == TypeShort2 {
		s2 := make([]AttrShort2, (end-start)/2)
		for i := 0; i < len(s2); i++ {
			s2[i][0] = v[i*2]
			s2[i][1] = v[i*2+1]
		}
		a := make([]AttrValue, len(s2))
		for i, s := range s2 {
			a[i] = &s
		}
		return a, end - start, nil
	} else {
		l2 := make([]AttrLong2, (end-start)/2)
		for i := 0; i < len(l2); i++ {
			l2[i][0] = v[i*2]
			l2[i][1] = v[i*2+1]
		}
		a := make([]AttrValue, len(l2))
		for i, l := range l2 {
			a[i] = &l
		}
		return a, end - start, nil
	}
}

func ParseShort3Long3(token *[]string, start int, size *uint, at *AttrType) ([]AttrValue, int, error) {
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
	if at != nil && *at == TypeShort3 {
		s3 := make([]AttrShort3, (end-start)/3)
		for i := 0; i < len(s3); i++ {
			s3[i][0] = v[i*3]
			s3[i][1] = v[i*3+1]
			s3[i][2] = v[i*3+2]
		}
		a := make([]AttrValue, len(s3))
		for i, s := range s3 {
			a[i] = &s
		}
		return a, end - start, nil
	} else {
		l3 := make([]AttrLong3, (end-start)/3)
		for i := 0; i < len(l3); i++ {
			l3[i][0] = v[i*3]
			l3[i][1] = v[i*3+1]
			l3[i][2] = v[i*3+2]
		}
		a := make([]AttrValue, len(l3))
		for i, l := range l3 {
			a[i] = &l
		}
		return a, end - start, nil
	}
}

func ParseInt32Array(token *[]string, start int) ([]AttrValue, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	a := make([]AttrValue, 1)
	if numberOfArray != 0 {
		result, err := ParseInts((*token)[start+1 : start+1+numberOfArray]...)
		if err != nil {
			return nil, 0, err
		}
		ia := AttrInt32Array(result)
		a[0] = &ia
	} else {
		ia := AttrInt32Array{}
		a[0] = &ia
	}
	return a, 1 + numberOfArray, nil
}

func ParseFloat2Double2(token *[]string, start int, size *uint, at *AttrType) ([]AttrValue, int, error) {
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
	if at != nil && *at == TypeFloat2 {
		f2 := make([]AttrFloat2, (end-start)/2)
		for i := 0; i < len(f2); i++ {
			f2[i][0] = v[i*2]
			f2[i][1] = v[i*2+1]
		}
		a := make([]AttrValue, len(f2))
		for i, f := range f2 {
			a[i] = &f
		}
		return a, end - start, nil
	} else {
		d2 := make([]AttrDouble2, (end-start)/2)
		for i := 0; i < len(d2); i++ {
			d2[i][0] = v[i*2]
			d2[i][1] = v[i*2+1]
		}
		a := make([]AttrValue, len(d2))
		for i, d := range d2 {
			a[i] = &d
		}
		return a, end - start, nil
	}
}

func ParseFloat3Double3(token *[]string, start int, size *uint, at *AttrType) ([]AttrValue, int, error) {
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
	if at != nil && *at == TypeFloat3 {
		f3 := make([]AttrFloat3, (end-start)/3)
		for i := 0; i < len(f3); i++ {
			f3[i][0] = v[i*3]
			f3[i][1] = v[i*3+1]
			f3[i][2] = v[i*3+2]
		}
		a := make([]AttrValue, len(f3))
		for i, f := range f3 {
			a[i] = &f
		}
		return a, end - start, nil
	} else {
		d3 := make([]AttrDouble3, (end-start)/3)
		for i := 0; i < len(d3); i++ {
			d3[i][0] = v[i*3]
			d3[i][1] = v[i*3+1]
			d3[i][2] = v[i*3+2]
		}
		a := make([]AttrValue, len(d3))
		for i, d := range d3 {
			a[i] = &d
		}
		return a, end - start, nil
	}
}

func ParseDoubleArray(token *[]string, start int) ([]AttrValue, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	a := make([]AttrValue, 1)
	if numberOfArray != 0 {
		f, err := ParseFloats((*token)[start+1 : start+1+numberOfArray]...)
		if err != nil {
			return nil, 0, err
		}
		da := AttrDoubleArray(f)
		a[0] = &da
	} else {
		da := AttrDoubleArray{}
		a[0] = &da
	}
	return a, 1 + numberOfArray, nil
}

func ParseMatrix(token *[]string, start int) ([]AttrValue, int, error) {
	mat4x4, err := ParseFloats((*token)[start : start+16]...)
	if err != nil {
		return nil, 0, err
	}
	a := make([]AttrValue, 1)
	a[0] = &AttrMatrix{
		mat4x4[0], mat4x4[1], mat4x4[2], mat4x4[3],
		mat4x4[4], mat4x4[5], mat4x4[6], mat4x4[7],
		mat4x4[8], mat4x4[9], mat4x4[10], mat4x4[11],
		mat4x4[12], mat4x4[13], mat4x4[14], mat4x4[15],
	}
	return a, 16, nil
}

func ParseMatrixXform(token *[]string, start int) ([]AttrValue, int, error) {
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
	rotateOrder, err := ConvertAttrRotateOrder(int(floats[6]))
	if err != nil {
		return nil, 0, err
	}
	componentForParentScale, err := isOnYesOrOffNo((*token)[start+37])
	if err != nil {
		return nil, 0, err
	}
	mx := AttrMatrixXform{
		Scale:                    AttrVector{X: floats[0], Y: floats[1], Z: floats[2]},
		Rotate:                   AttrVector{X: floats[3], Y: floats[4], Z: floats[5]},
		RotateOrder:              rotateOrder,
		Translate:                AttrVector{X: floats[7], Y: floats[8], Z: floats[9]},
		Shear:                    AttrShear{XY: floats[10], XZ: floats[11], YZ: floats[12]},
		ScalePivot:               AttrVector{X: floats[13], Y: floats[14], Z: floats[15]},
		ScaleTranslate:           AttrVector{X: floats[16], Y: floats[17], Z: floats[18]},
		RotatePivot:              AttrVector{X: floats[19], Y: floats[20], Z: floats[21]},
		RotateTranslation:        AttrVector{X: floats[22], Y: floats[23], Z: floats[24]},
		RotateOrient:             AttrOrient{W: floats[25], X: floats[26], Y: floats[27], Z: floats[28]},
		JointOrient:              AttrOrient{W: floats[29], X: floats[30], Y: floats[31], Z: floats[32]},
		InverseParentScale:       AttrVector{X: floats[33], Y: floats[34], Z: floats[35]},
		CompensateForParentScale: componentForParentScale,
	}
	a := make([]AttrValue, 1)
	a[0] = &mx
	return a, 38, nil
}

func ParsePointArray(token *[]string, start int) ([]AttrValue, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	a := make([]AttrValue, 1)
	if numberOfArray != 0 {
		f, err := ParseFloats((*token)[start+1 : start+1+(numberOfArray*4)]...)
		if err != nil {
			return nil, 0, err
		}
		pa := make([]AttrPoint, numberOfArray)
		for i := 0; i < numberOfArray; i++ {
			pa[i].X = f[i*4]
			pa[i].Y = f[i*4+1]
			pa[i].Z = f[i*4+2]
			pa[i].W = f[i*4+3]
		}
		paa := AttrPointArray(pa)
		a[0] = &paa
	} else {
		paa := AttrPointArray{}
		a[0] = &paa
	}
	return a, 1 + (numberOfArray * 4), nil
}

func ParseVectorArray(token *[]string, start int) ([]AttrValue, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	a := make([]AttrValue, 1)
	if numberOfArray != 0 {
		f, err := ParseFloats((*token)[start+1 : start+1+numberOfArray]...)
		if err != nil {
			return nil, 0, err
		}
		va := make([]AttrVector, numberOfArray)
		for i := 0; i < numberOfArray; i += 3 {
			//log.Printf("numberOfArray: %d, size: %d, indeces:[%d, %d, %d]", numberOfArray, len(f), i, i+1, i+2)
			// Ornatrix ClumpNode は 3 で割り切れない数の vectorArray を扱う...
			add1 := i+1
			if add1 >= numberOfArray {
				add1 = i
			}
			add2 := i+2
			if add2 >= numberOfArray {
				add2 = add1
			}
			va[i] = AttrVector{
				X: f[i],
				Y: f[add1],
				Z: f[add2],
			}
		}
		vaa := AttrVectorArray(va)
		a[0] = &vaa
	} else {
		vaa := AttrVectorArray{}
		a[0] = &vaa
	}
	return a, 1 + (numberOfArray * 3), nil
}

func ParseString(token *[]string, start int) ([]AttrValue, int, error) {
	s := AttrString((*token)[start][1 : len((*token)[start])-1])
	a := []AttrValue{&s}
	return a, 1, nil
}

func ParseStringArray(token *[]string, start int) ([]AttrValue, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	sa := AttrStringArray((*token)[start+1 : start+1+numberOfArray])
	for i, s := range sa {
		sa[i] = s[1 : len(s)-1]
	}
	a := []AttrValue{&sa}
	return a, 1 + numberOfArray, nil
}

func ParseSphere(token *[]string, start int) ([]AttrValue, int, error) {
	s, err := strconv.ParseFloat((*token)[start], 64)
	if err != nil {
		return nil, 0, err
	}
	sp := AttrSphere(s)
	a := []AttrValue{&sp}
	return a, 1, nil
}

func ParseCone(token *[]string, start int) ([]AttrValue, int, error) {
	f, err := ParseFloats((*token)[start : start+2]...)
	if err != nil {
		return nil, 0, err
	}
	c := AttrCone{
		ConeAngle: f[0],
		ConeCap:   f[1],
	}
	a := []AttrValue{&c}
	return a, 2, nil
}

func ParseReflectanceRGB(token *[]string, start int) ([]AttrValue, int, error) {
	f, err := ParseFloats((*token)[start : start+3]...)
	if err != nil {
		return nil, 0, err
	}
	a := make([]AttrValue, 1)
	a[0] = &AttrReflectanceRGB{
		RedReflect:   f[0],
		GreenReflect: f[1],
		BlueReflect:  f[2],
	}
	return a, 3, nil
}

func ParseSpectrumRGB(token *[]string, start int) ([]AttrValue, int, error) {
	f, err := ParseFloats((*token)[start : start+3]...)
	if err != nil {
		return nil, 0, err
	}
	a := make([]AttrValue, 1)
	a[0] = &AttrSpectrumRGB{
		RedSpectrum:   f[0],
		GreenSpectrum: f[1],
		BlueSpectrum:  f[2],
	}
	return a, 3, nil
}

func ParseComponentList(token *[]string, start int) ([]AttrValue, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	var cl AttrComponentList
	for _, c := range (*token)[start+1 : start+1+numberOfArray] {
		cl = append(cl, strings.Trim(c, "\""))
	}
	a := []AttrValue{&cl}
	return a, 1 + numberOfArray, nil
}

func ParseAttributeAlias(token *[]string, start int) ([]AttrValue, int, error) {
	if (*token)[start] != "{" {
		return nil, 0, errors.New("there was no necessary token")
	}
	var aaa []AttrAttributeAlias
	for i := start + 1; i < len(*token); i += 2 {
		if (*token)[i] == "}" {
			break
		}
		aaa = append(aaa, AttrAttributeAlias{
			NewAlias:    (*token)[i],
			CurrentName: (*token)[i+1],
		})
	}
	a := make([]AttrValue, len(aaa))
	for i := range aaa {
		a[i] = &aaa[i]
	}
	return a, 2 + (len(aaa) * 2), nil
}

func ParseNurbsCurve(token *[]string, start int) ([]AttrValue, int, error) {
	i1, err := ParseInts((*token)[start : start+3]...)
	if err != nil {
		return nil, 0, err
	}
	degree := i1[0]
	spans := i1[1]
	form, err := ConvertAttrFormType(i1[2]) // open(0), closed(1), periodic(2)
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
	cvValues := make([]AttrCvValue, len(cv)/divideCv)
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
	a := make([]AttrValue, 1)
	a[0] = &AttrNurbsCurve{
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

func ParseNurbsSurface(token *[]string, start int) ([]AttrValue, int, error) {
	i1, err := ParseInts((*token)[start : start+4]...)
	if err != nil {
		return nil, 0, err
	}
	uDegree := i1[0]
	vDegree := i1[1]
	uForm, err := ConvertAttrFormType(i1[2])
	if err != nil {
		return nil, 0, err
	}
	vForm, err := ConvertAttrFormType(i1[3])
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
	cvValue := make([]AttrCvValue, cvCount)
	for i := 0; i < cvCount; i++ {
		cvValue[i].X = cv[i*divideCv]
		cvValue[i].Y = cv[i*divideCv+1]
		cvValue[i].Z = &cv[i*divideCv+2]
		if isRational {
			cvValue[i].W = &cv[i*divideCv+3]
		}
	}
	a := make([]AttrValue, 1)
	a[0] = &AttrNurbsSurface{
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

func ParseNurbsTrimface(token *[]string, start int) ([]AttrValue, int, error) {
	// TODO: Waiting for Autodesk
	a := []AttrValue{&AttrNurbsTrimface{}}
	return a, -1, nil
}

func ParseCountInt(token *[]string, start int) ([]int, error) {
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

func ParsePolyFace(token *[]string, start int, size *uint) ([]AttrValue, int, error) {
	switchNumber := start
	var pfs []AttrPolyFaces
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
		pfs = make([]AttrPolyFaces, s)
	}
	i := -1
	loop := true
	for loop && len(*token) > switchNumber {
		switch (*token)[switchNumber] {
		case "f":
			fe, err := ParseCountInt(token, switchNumber+1)
			if err != nil {
				log.Println("error case f")
				log.Printf("token number is %d", switchNumber+1)
				return nil, 0, err
			}
			i++
			if i >= len(pfs) {
				pf := AttrPolyFaces{}
				pfs = append(pfs, pf)
			}
			pfs[i].FaceEdge = fe
			switchNumber += 2 + len(fe)
		case "h":
			he, err := ParseCountInt(token, switchNumber+1)
			if err != nil {
				return nil, 0, err
			}
			pfs[i].HoleEdge = he
			switchNumber += 2 + len(he)
		case "fc":
			fc, err := ParseCountInt(token, switchNumber+1)
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
			colorIDs, err := ParseCountInt(token, switchNumber+2)
			if err != nil {
				return nil, 0, err
			}
			mc := AttrMultiColor{
				ColorIndex: int(colorIndex),
				ColorIDs:   colorIDs,
			}
			pfs[i].MultiColor = append(pfs[i].MultiColor, mc)
			switchNumber += 3 + len(colorIDs)
		case "mu":
			var fuv AttrFaceUV
			uvSet, err := strconv.Atoi((*token)[switchNumber+1])
			if err != nil {
				return nil, 0, err
			}
			fuv.UVSet = uvSet
			uv, err := ParseCountInt(token, switchNumber+2)
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
	a := make([]AttrValue, len(pfs))
	for i := range pfs {
		a[i] = &pfs[i]
	}
	return a, switchNumber - start, nil
}

func ParseDataPolyComponent(token *[]string, start int) ([]AttrValue, int, error) {
	if "Index_Data" != (*token)[start] {
		return nil, 0, errors.New(
			"since the Index_Data did not exist, " +
				"this token is an unknown dataPolyComponent")
	}
	var dpc AttrDataPolyComponent
	switch (*token)[start+1] {
	case "Edge":
		dpc.PolyComponentType = DPCedge
	case "Face":
		dpc.PolyComponentType = DPCface
	case "Vertex":
		dpc.PolyComponentType = DPCvertex
	case "UV":
		dpc.PolyComponentType = DPCuv
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
	a := []AttrValue{&dpc}
	return a, 3 + (count * 2), nil
}

func ParseMesh(_ *[]string, _ int) ([]AttrValue, int, error) {
	// Not Implement
	var a []AttrValue
	return a, -1, nil
}

func ParseLattice(token *[]string, start int) ([]AttrValue, int, error) {
	c, err := ParseInts((*token)[start : start+4]...)
	if err != nil {
		return nil, 0, err
	}
	la := AttrLattice{
		DivisionS: c[0],
		DivisionT: c[1],
		DivisionU: c[2],
	}
	la.Points = make([]AttrLatticePoint, c[3])
	for i := 0; i < c[3]*3; i += 3 {
		p, err := ParseFloats((*token)[start+4+i : start+4+i+3]...)
		if err != nil {
			return nil, 0, err
		}
		la.Points[i/3].S = p[0]
		la.Points[i/3].T = p[1]
		la.Points[i/3].U = p[2]
	}
	a := []AttrValue{&la}
	return a, 4 + (c[3] * 3), nil
}

func ParseAttr(token *[]string, start int, size *uint, attrType AttrType) ([]AttrValue, int, error) {
	switch attrType {
	case TypeShort2, TypeLong2:
		return ParseShort2Long2(token, start, size, &attrType)
	case TypeShort3, TypeLong3:
		return ParseShort3Long3(token, start, size, &attrType)
	case TypeInt32Array:
		return ParseInt32Array(token, start)
	case TypeFloat2, TypeDouble2:
		return ParseFloat2Double2(token, start, size, &attrType)
	case TypeFloat3, TypeDouble3:
		return ParseFloat3Double3(token, start, size, &attrType)
	case TypeDoubleArray:
		return ParseDoubleArray(token, start)
	case TypeMatrix:
		return ParseMatrix(token, start)
	case TypeMatrixXform:
		return ParseMatrixXform(token, start)
	case TypePointArray:
		return ParsePointArray(token, start)
	case TypeVectorArray:
		return ParseVectorArray(token, start)
	case TypeString:
		return ParseString(token, start)
	case TypeStringArray:
		return ParseStringArray(token, start)
	case TypeSphere:
		return ParseSphere(token, start)
	case TypeCone:
		return ParseCone(token, start)
	case TypeReflectanceRGB:
		return ParseReflectanceRGB(token, start)
	case TypeSpectrumRGB:
		return ParseSpectrumRGB(token, start)
	case TypeComponentList:
		return ParseComponentList(token, start)
	case TypeAttributeAlias:
		return ParseAttributeAlias(token, start)
	case TypeNurbsCurve:
		return ParseNurbsCurve(token, start)
	case TypeNurbsSurface:
		return ParseNurbsSurface(token, start)
	case TypeNurbsTrimface:
		return ParseNurbsTrimface(token, start)
	case TypePolyFaces:
		return ParsePolyFace(token, start, size)
	case TypeDataPolyComponent:
		return ParseDataPolyComponent(token, start)
	case TypeDataReferenceEdits:
		return MakeDataReferenceEdits(token, start)
	case TypeMesh:
		return ParseMesh(token, start)
	case TypeLattice:
		return ParseLattice(token, start)
	}
	return nil, 0, nil
}

func ParseAttrType(token *[]string, start int) (AttrType, error) {
	typeString := (*token)[start][1 : len((*token)[start])-1]
	switch typeString {
	case "short2":
		return TypeShort2, nil
	case "short3":
		return TypeShort3, nil
	case "long2":
		return TypeLong2, nil
	case "long3":
		return TypeLong3, nil
	case "Int32Array":
		return TypeInt32Array, nil
	case "float2":
		return TypeFloat2, nil
	case "double2":
		return TypeDouble2, nil
	case "float3":
		return TypeFloat3, nil
	case "double3":
		return TypeDouble3, nil
	case "doubleArray":
		return TypeDoubleArray, nil
	case "matrix":
		typeString2 := (*token)[start+1]
		if typeString2 == "\"xform\"" {
			return TypeMatrixXform, nil
		} else {
			return TypeMatrix, nil
		}
	case "pointArray":
		return TypePointArray, nil
	case "vectorArray":
		return TypeVectorArray, nil
	case "string":
		return TypeString, nil
	case "stringArray":
		return TypeStringArray, nil
	case "sphere":
		return TypeSphere, nil
	case "cone":
		return TypeCone, nil
	case "reflectanceRGB":
		return TypeReflectanceRGB, nil
	case "spectrumRGB":
		return TypeSpectrumRGB, nil
	case "componentList":
		return TypeComponentList, nil
	case "attributeAlias":
		return TypeAttributeAlias, nil
	case "nurbsCurve":
		return TypeNurbsCurve, nil
	case "nurbsSurface":
		return TypeNurbsSurface, nil
	case "nurbsTrimface":
		return TypeNurbsTrimface, nil
	case "polyFaces":
		return TypePolyFaces, nil
	case "dataPolyComponent":
		return TypeDataPolyComponent, nil
	case "dataReferenceEdits":
		return TypeDataReferenceEdits, nil
	case "mesh":
		return TypeMesh, nil
	case "lattice":
		return TypeLattice, nil
	}
	return TypeInvalid, errors.New("Invalid type " + typeString)
}

func ParseAddAttr(c *Cmd) *AddAttrCmd {
	aa := &AddAttrCmd{Cmd: c}
	// TODO: Do finish!
	return aa
}
