package mayaascii

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/nrtkbb/bufscan"
	"log"
	"os"
	"strconv"
	"strings"
)

type File struct {
	Path string
	Cmds []*Cmd
}

func (f *File) Parse() error {
	fp, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer fp.Close()

	reader := bufio.NewReader(fp)
	f.Cmds = []*Cmd{}
	cmdBuilder := &CmdBuilder{}
	err = bufscan.BufScan(reader, func(line string) error {
		cmdBuilder.Append(line)
		if cmdBuilder.IsCmdEOF() {
			cmd := cmdBuilder.Parse()
			f.Cmds = append(f.Cmds, cmd)
			cmdBuilder.Clear()
		}
		return nil
	})
	return nil
}

func (f *File) SaveSceneAs(outputPath string) error {
	if _, err := os.Stat(outputPath); err == nil {
		return errors.New("file already existed.")
	}
	fp, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer fp.Close()
	writer := bufio.NewWriter(fp)
	defer writer.Flush()

	for _, cmd := range f.Cmds {
		fmt.Fprintln(writer, cmd.Raw)
	}
	return nil
}

type CmdBuilder struct {
	cmdLine []string
}

func (c *CmdBuilder) Append(line string) {
	c.cmdLine = append(c.cmdLine, line)
}

func (c *CmdBuilder) IsCmdEOF() bool {
	if len(c.cmdLine) == 0 {
		return false
	}
	lastLine := c.cmdLine[len(c.cmdLine)-1]
	if len(lastLine) == 0 {
		return false
	}
	return lastLine[len(lastLine)-1] == byte(';')
}

func (c *CmdBuilder) Clear() {
	c.cmdLine = []string{}
}

const (
	whiteSpace        = ' '
	tabSpace          = '\t'
	enter             = '\n'
	slash             = '/'
	backSlash         = '\\'
	dQuatation        = '"'
	sQuatation        = '\''
	hyphen            = '-'
	plus              = '+'
	openRoundBracket  = '('
	closeRoundBracket = ')'
	kanma             = ','
	semiCoron         = ';'
	openBrancket      = '{'
	closeBrancket     = '}'
)

type Cmd struct {
	CmdName string   `json:"cmd"`
	Raw     string   `json:"raw"`
	Token   []string `json:"token"`
}

func (c *CmdBuilder) Parse() *Cmd {
	cmd := Cmd{Raw: strings.Join(c.cmdLine, "\n")}
	c.Clear()
	var buf []rune
	var subBuf []rune
	var subToken []string
	for _, c := range cmd.Raw {
		if 0 < len(buf) && buf[len(buf)-1] == backSlash {
			buf = append(buf, c)
			subBuf = append(subBuf, c)
			continue
		}
		if (c == openBrancket && 0 == len(buf)) ||
			(0 < len(buf) && buf[0] == openBrancket) {
			// {"attribute","alias"
			// 		"attribute","alias"}
			buf = append(buf, c)
			if len(buf) == 1 {
				cmd.Token = append(cmd.Token, string(c))
				continue
			}
			if c == dQuatation ||
				(len(subBuf) > 0 && subBuf[0] == dQuatation) {
				subBuf = append(subBuf, c)
				if len(subBuf) > 1 && c == dQuatation {
					cmd.Token = append(cmd.Token, strings.Trim(string(subBuf), "\""))
					subBuf = subBuf[:0]
				}
			}
			if c == whiteSpace ||
				c == tabSpace ||
				c == kanma ||
				c == enter {
				continue
			}
			if c == closeBrancket {
				cmd.Token = append(cmd.Token, string(c))
				buf = buf[:0]
				continue
			}
		}
		if (c == openRoundBracket && 0 == len(buf)) ||
			(0 < len(buf) && buf[0] == openRoundBracket) {
			// ( "long text" +
			//   	"long text")
			buf = append(buf, c)
			if len(buf) == 1 {
				continue
			}
			if c == dQuatation ||
				(len(subBuf) > 0 && subBuf[0] == dQuatation) {
				subBuf = append(subBuf, c)
				if len(subBuf) > 1 && c == dQuatation {
					subToken = append(subToken, strings.Trim(string(subBuf), "\""))
					subBuf = subBuf[:0]
				}
				continue
			}
			if c == whiteSpace ||
				c == tabSpace ||
				c == enter ||
				c == plus {
				continue
			}
			if c == closeRoundBracket {
				cmd.Token = append(cmd.Token, strings.Join(subToken, ""))
				buf = buf[:0]
				continue
			}
		}
		if (c == dQuatation && 0 == len(buf)) ||
			(0 < len(buf) && buf[0] == dQuatation) {
			// "text\"s"
			buf = append(buf, c)
			if len(buf) > 1 && c == dQuatation {
				cmd.Token = append(cmd.Token, string(buf))
				buf = buf[:0]
				continue
			}
			continue
		}
		if c == whiteSpace || c == tabSpace || c == enter {
			if len(buf) != 0 {
				cmd.Token = append(cmd.Token, string(buf))
				buf = buf[:0]
			}
			continue
		}
		if c == semiCoron {
			if len(buf) != 0 {
				cmd.Token = append(cmd.Token, string(buf))
			}
			break
		}
		buf = append(buf, c)
	}
	if 0 == len(cmd.Token) {
		cmd.CmdName = "None"
	} else {
		cmd.CmdName = cmd.Token[0]
	}
	return &cmd
}

type Workspace struct {
	*Cmd
	FileRule string `json:"file_rule" tag:"-fileRule"`
	Place    string `json:"place"`
}

func MakeWorkspace(cmd *Cmd) *Workspace {
	w := Workspace{Cmd: cmd}
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

type Requires struct {
	*Cmd
	PluginName string   `json:"plugin_name"`
	Version    string   `json:"version"`
	NodeTypes  []string `json:"node_types" tag:"-nodeType"`
	DataTypes  []string `json:"data_types" tag:"-dataType"`
}

func MakeRequires(cmd *Cmd) *Requires {
	// max Token = [requires, -nodeType, "typeName1", -dataType, "typeName2", "pluginName", "version"]
	// min Token = [requires, "pluginName", "version"]
	r := Requires{Cmd: cmd}
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

type ConnectAttr struct {
	*Cmd
	SrcNode       string  `json:"src_node"`
	SrcAttr       string  `json:"src_attr"`
	DstNode       string  `json:"dst_node"`
	DstAttr       string  `json:"dst_attr"`
	Force         bool    `json:"force" tag:"-f"`
	Lock          *bool   `json:"lock,omitempty" tag:"-l"`
	NextAvailable bool    `json:"next_available" tag:"-na"`
	ReferenceDest *string `json:"reference_dest,omitempty" tag:"-rd"`
}

func MakeConnectAttr(cmd *Cmd) (*ConnectAttr, error) {
	ca := &ConnectAttr{Cmd: cmd}
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

type CreateNode struct {
	*Cmd
	NodeType   string  `json:"node_type"`
	NodeName   string  `json:"node_name" short:"-n"`
	Parent     *string `json:"parent" short:"-p"`
	Shared     bool    `json:"shared" short:"-s"`
	SkipSelect bool    `json:"skip_select" short:"-ss"`
}

func MakeCreateNode(cmd *Cmd) *CreateNode {
	n := &CreateNode{Cmd: cmd}
	n.NodeType = cmd.Token[1]
	for i := 2; i < len(n.Token); i++ {
		switch n.Token[i] {
		case "-n":
			i++
			n.NodeName = n.Token[i][1 : len(cmd.Token[i])-1]
		case "-p":
			i++
			p := n.Token[i][1 : len(cmd.Token[i])-1]
			n.Parent = &p
		case "-s":
			n.Shared = true
		case "-ss":
			n.SkipSelect = true
		}
	}
	return n
}

type Rename struct {
	*Cmd
	From        *string `json:"from,omitempty"`
	To          *string `json:"to"`
	UUID        bool    `json:"uuid" short:"-uid"`
	IgnoreShape bool    `json:"ignore_shape" short:"-is"`
}

func MakeRename(cmd *Cmd) *Rename {
	r := &Rename{Cmd: cmd}
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

type SetAttr struct {
	*Cmd
	AttrName     string   `json:"attr_name"`
	AlteredValue bool     `json:"altered_value" short:"-av"`
	Caching      *bool    `json:"caching,omitempty" short:"-ca"`
	CapacityHint *uint    `json:"capacity_hint,omitempty" short:"-ch"`
	ChannelBox   *bool    `json:"channel_box,omitempty" short:"-cb"`
	Clamp        bool     `json:"clamp" short:"-c"`
	Keyable      *bool    `json:"keyable,omitempty" short:"-k"`
	Lock         *bool    `json:"lock,omitempty" short:"-l"`
	Size         *uint    `json:"size,omitempty" short:"-s"`
	AttrType     AttrType `json:"attr_type" short:"-typ"`
	Attr         Attr     `json:"attr"`
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

var PrimitiveTypes = map[AttrType]struct{}{
	TypeInvalid: struct{}{},
	TypeBool:    struct{}{},
	TypeInt:     struct{}{},
	TypeDouble:  struct{}{},
}

func MakeSetAttr(cmd *Cmd, beforeSetAttr *SetAttr) (*SetAttr, error) {
	attrNameIdx, attrName := getAttrNameFromSetAttr(&cmd.Token)
	sa := &SetAttr{Cmd: cmd}
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
				sa.Attr = appendSetAttr(sa.Attr, a)
				if count == -1 {
					break
				}
				i += count
				break
			}
			b, err := isOnYesOrOffNo(v)
			if err == nil {
				sa.AttrType = TypeBool
				sa.Attr = &b
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
				sa.Attr = appendSetAttr(sa.Attr, &intArray)
				sa.AttrType = TypeInt
				return sa, nil
			}

			floatArray, err := ParseFloats(sa.Token[i:]...)
			if err != nil {
				return nil, err
			}
			if sa.AttrType == TypeInt {
				intArray, ok := sa.Attr.(*[]int)
				if !ok {
					return nil, errors.New(
						fmt.Sprintf("invalid pattern %v and %v", sa.Attr, floatArray))
				}
				floatArrayAttr := make([]float64, len(*intArray))
				for i, v := range *intArray {
					floatArrayAttr[i] = float64(v)
				}
				sa.Attr = &floatArrayAttr
			}
			sa.AttrType = TypeDouble
			sa.Attr = appendSetAttr(sa.Attr, &floatArray)
			return sa, nil
		}
	}
	return sa, nil
}

func appendSetAttr(beforeAttr Attr, newAttr Attr) Attr {
	if beforeAttr == nil {
		return newAttr
	}
	switch newAttr.(type) {
	case *[]int:
		beforeInt, _ := beforeAttr.(*[]int)
		newInt, _ := newAttr.(*[]int)
		*beforeInt = append(*beforeInt, *newInt...)
		return beforeInt
	case *[]float64:
		beforeFloat, _ := beforeAttr.(*[]float64)
		newFloat, _ := newAttr.(*[]float64)
		*beforeFloat = append(*beforeFloat, (*newFloat)...)
		return beforeFloat
	case *[]AttrShort2:
		beforeShort2, _ := beforeAttr.(*[]AttrShort2)
		newShort2, _ := newAttr.(*[]AttrShort2)
		*beforeShort2 = append(*beforeShort2, *newShort2...)
		return beforeShort2
	case *[]AttrShort3:
		beforeShort3, _ := beforeAttr.(*[]AttrShort3)
		newShort3, _ := newAttr.(*[]AttrShort3)
		*beforeShort3 = append(*beforeShort3, *newShort3...)
		return beforeShort3
	case *AttrInt32Array:
		beforeInt32Array, _ := beforeAttr.(*AttrInt32Array)
		newInt32Array, _ := newAttr.(*AttrInt32Array)
		*beforeInt32Array = append(*beforeInt32Array, *newInt32Array...)
		return beforeInt32Array
	case *[]AttrLong2:
		beforeLong2, _ := beforeAttr.(*[]AttrLong2)
		newLong2, _ := newAttr.(*[]AttrLong2)
		*beforeLong2 = append(*beforeLong2, *newLong2...)
	case *[]AttrLong3:
		beforeLong3, _ := beforeAttr.(*[]AttrLong3)
		newLong3, _ := newAttr.(*[]AttrLong3)
		*beforeLong3 = append(*beforeLong3, *newLong3...)
	case *[]AttrFloat2:
		beforeFloat2, _ := beforeAttr.(*[]AttrFloat2)
		newFloat2, _ := newAttr.(*[]AttrFloat2)
		*beforeFloat2 = append(*beforeFloat2, *newFloat2...)
		return beforeFloat2
	case *[]AttrFloat3:
		beforeFloat3, _ := beforeAttr.(*[]AttrFloat3)
		newFloat3, _ := newAttr.(*[]AttrFloat3)
		*beforeFloat3 = append(*beforeFloat3, *newFloat3...)
		return beforeFloat3
	case *[]AttrDouble2:
		beforeDouble2, _ := beforeAttr.(*[]AttrDouble2)
		newDouble2, _ := newAttr.(*[]AttrDouble2)
		*beforeDouble2 = append(*beforeDouble2, *newDouble2...)
		return beforeDouble2
	case *[]AttrDouble3:
		beforeDouble3, _ := beforeAttr.(*[]AttrDouble3)
		newDouble3, _ := newAttr.(*[]AttrDouble3)
		*beforeDouble3 = append(*beforeDouble3, *newDouble3...)
		return beforeDouble3
	case *AttrDoubleArray:
		beforeDoubleArray, _ := beforeAttr.(*AttrDoubleArray)
		newDoubleArray, _ := newAttr.(*AttrDoubleArray)
		*beforeDoubleArray = append(*beforeDoubleArray, *newDoubleArray...)
		return beforeDoubleArray
	case *[]AttrMatrix:
		beforeMatrix, _ := beforeAttr.(*[]AttrMatrix)
		newMatrix, _ := newAttr.(*AttrMatrix)
		*beforeMatrix = append(*beforeMatrix, *newMatrix)
		return beforeMatrix
	case *[]AttrMatrixXform:
		beforeMatrixXform, _ := beforeAttr.(*[]AttrMatrixXform)
		newMatrixXform, _ := newAttr.(*AttrMatrixXform)
		*beforeMatrixXform = append(*beforeMatrixXform, *newMatrixXform)
		return beforeMatrixXform
	case *AttrPointArray:
		beforePointArray, _ := beforeAttr.(*AttrPointArray)
		newPointArray, _ := newAttr.(*AttrPointArray)
		*beforePointArray = append(*beforePointArray, *newPointArray...)
		return beforePointArray
	case *AttrVectorArray:
		beforeVectorArray, _ := beforeAttr.(*AttrVectorArray)
		newVectorArray, _ := newAttr.(*AttrVectorArray)
		*beforeVectorArray = append(*beforeVectorArray, *newVectorArray...)
		return beforeVectorArray
	case *[]AttrString:
		beforeString, _ := beforeAttr.(*[]AttrString)
		newString, _ := newAttr.(*AttrString)
		*beforeString = append(*beforeString, *newString)
		return beforeString
	case *AttrStringArray:
		beforeStringArray, _ := beforeAttr.(*AttrStringArray)
		newStringArray, _ := newAttr.(*AttrStringArray)
		*beforeStringArray = append(*beforeStringArray, *newStringArray...)
		return beforeStringArray
	case *[]AttrSphere:
		beforeSphere, _ := beforeAttr.(*[]AttrSphere)
		newSphere, _ := newAttr.(*AttrSphere)
		*beforeSphere = append(*beforeSphere, *newSphere)
		return beforeSphere
	case *[]AttrCone:
		beforeCone, _ := beforeAttr.(*[]AttrCone)
		newCone, _ := newAttr.(*[]AttrCone)
		*beforeCone = append(*beforeCone, *newCone...)
		return beforeCone
	case *[]AttrReflectanceRGB:
		beforeReflectanceRGB, _ := beforeAttr.(*[]AttrReflectanceRGB)
		newReflectanceRGB, _ := newAttr.(*AttrReflectanceRGB)
		*beforeReflectanceRGB = append(*beforeReflectanceRGB, *newReflectanceRGB)
		return beforeReflectanceRGB
	case *[]AttrSpectrumRGB:
		beforeSpectrumRGB, _ := beforeAttr.(*[]AttrSpectrumRGB)
		newSpectrumRGB, _ := newAttr.(*AttrSpectrumRGB)
		*beforeSpectrumRGB = append(*beforeSpectrumRGB, *newSpectrumRGB)
		return beforeSpectrumRGB
	case *AttrComponentList:
		beforeComponentList, _ := beforeAttr.(*AttrComponentList)
		newComponentList, _ := newAttr.(*AttrComponentList)
		*beforeComponentList = append(*beforeComponentList, *newComponentList...)
		return beforeComponentList
	case *[]AttrAttributeAlias:
		beforeAttributeAlias, _ := beforeAttr.(*[]AttrAttributeAlias)
		newAttributeAlias, _ := newAttr.(*[]AttrAttributeAlias)
		*beforeAttributeAlias = append(*beforeAttributeAlias, *newAttributeAlias...)
		return beforeAttributeAlias
	case *[]AttrNurbsCurve:
		beforeNurbsCurve, _ := beforeAttr.(*[]AttrNurbsCurve)
		newNurbsCurve, _ := newAttr.(*[]AttrNurbsCurve)
		*beforeNurbsCurve = append(*beforeNurbsCurve, *newNurbsCurve...)
		return beforeNurbsCurve
	case *[]AttrNurbsSurface:
		beforeNurbsSurface, _ := beforeAttr.(*[]AttrNurbsSurface)
		newNurbsSurface, _ := newAttr.(*AttrNurbsSurface)
		*beforeNurbsSurface = append(*beforeNurbsSurface, *newNurbsSurface)
		return beforeNurbsSurface
	case *[]AttrNurbsTrimface:
		beforeNurbsTrimface, _ := beforeAttr.(*[]AttrNurbsTrimface)
		newNurbsTrimface, _ := newAttr.(*AttrNurbsTrimface)
		*beforeNurbsTrimface = append(*beforeNurbsTrimface, *newNurbsTrimface)
		return beforeNurbsTrimface
	case *[]AttrPolyFaces:
		beforePolyFaces, _ := beforeAttr.(*[]AttrPolyFaces)
		newPolyFaces, _ := newAttr.(*[]AttrPolyFaces)
		*beforePolyFaces = append(*beforePolyFaces, *newPolyFaces...)
		return beforePolyFaces
	case *[]AttrDataPolyComponent:
		beforeDataPolyComponent, _ := beforeAttr.(*[]AttrDataPolyComponent)
		newDataPolyComponent, _ := newAttr.(*[]AttrDataPolyComponent)
		*beforeDataPolyComponent = append(*beforeDataPolyComponent, *newDataPolyComponent...)
		return beforeDataPolyComponent
	case *[]AttrLattice:
		beforeLattice, _ := beforeAttr.(*[]AttrLattice)
		newLattice, _ := newAttr.(*[]AttrLattice)
		*beforeLattice = append(*beforeLattice, *newLattice...)
		return beforeLattice
	}
	return newAttr
}

type DisconnectBehaviour uint

type AddAttr struct {
	AttributeType       *string              `json:"attribute_type,omitempty" short:"-at"`
	CachedInternally    *bool                `json:"cached_internally,omitempty" short:"-ci"`
	Category            *string              `json:"category,omitempty" short:"-ct"`
	DataType            *string              `json:"data_type,omitempty" short:"-dt"`
	DefaultValue        *float64             `json:"default_value,omitempty" short:"-dv"`
	DisconnectBehaviour *DisconnectBehaviour `json:"disconnect_behaviour,omitempty" short:"-dcb"`
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
			panic(err)
		}
		result = append(result, int(i))
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

func MakeShort2Long2(token *[]string, start int, size *uint, at *AttrType) (Attr, int, error) {
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
		var a Attr = &s2
		return a, end - start, nil
	} else {
		l2 := make([]AttrLong2, (end-start)/2)
		for i := 0; i < len(l2); i++ {
			l2[i][0] = v[i*2]
			l2[i][1] = v[i*2+1]
		}
		var a Attr = &l2
		return a, end - start, nil
	}
}

func MakeShort3Long3(token *[]string, start int, size *uint, at *AttrType) (Attr, int, error) {
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
		var a Attr = &s3
		return a, end - start, nil
	} else {
		l3 := make([]AttrLong3, (end-start)/3)
		for i := 0; i < len(l3); i++ {
			l3[i][0] = v[i*3]
			l3[i][1] = v[i*3+1]
			l3[i][2] = v[i*3+2]
		}
		var a Attr = &l3
		return a, end - start, nil
	}
}

func MakeInt32Array(token *[]string, start int) (Attr, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	var a Attr
	if numberOfArray != 0 {
		result, err := ParseInts((*token)[start+1 : start+1+numberOfArray]...)
		if err != nil {
			return nil, 0, err
		}
		ia := AttrInt32Array(result)
		a = &ia
	} else {
		ia := AttrInt32Array{}
		a = &ia
	}
	return a, 1 + numberOfArray, nil
}

func MakeFloat2Double2(token *[]string, start int, size *uint, at *AttrType) (Attr, int, error) {
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
		var a Attr = &f2
		return a, end - start, nil
	} else {
		d2 := make([]AttrDouble2, (end-start)/2)
		for i := 0; i < len(d2); i++ {
			d2[i][0] = v[i*2]
			d2[i][1] = v[i*2+1]
		}
		var a Attr = &d2
		return a, end - start, nil
	}
}

func MakeFloat3Double3(token *[]string, start int, size *uint, at *AttrType) (Attr, int, error) {
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
		var a Attr = &f3
		return a, end - start, nil
	} else {
		d3 := make([]AttrDouble3, (end-start)/3)
		for i := 0; i < len(d3); i++ {
			d3[i][0] = v[i*3]
			d3[i][1] = v[i*3+1]
			d3[i][2] = v[i*3+2]
		}
		var a Attr = &d3
		return a, end - start, nil
	}
}

func MakeDoubleArray(token *[]string, start int) (Attr, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	var a Attr
	if numberOfArray != 0 {
		f, err := ParseFloats((*token)[start+1 : start+1+numberOfArray]...)
		if err != nil {
			return nil, 0, err
		}
		da := AttrDoubleArray(f)
		a = &da
	} else {
		da := AttrDoubleArray{}
		a = &da
	}
	return a, 1 + numberOfArray, nil
}

func MakeMatrix(token *[]string, start int) (Attr, int, error) {
	mat4x4, err := ParseFloats((*token)[start : start+16]...)
	if err != nil {
		return nil, 0, err
	}
	var a Attr = &AttrMatrix{
		mat4x4[0], mat4x4[1], mat4x4[2], mat4x4[3],
		mat4x4[4], mat4x4[5], mat4x4[6], mat4x4[7],
		mat4x4[8], mat4x4[9], mat4x4[10], mat4x4[11],
		mat4x4[12], mat4x4[13], mat4x4[14], mat4x4[15],
	}
	return a, 16, nil
}

func MakeMatrixXform(token *[]string, start int) (Attr, int, error) {
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
		Scale:                    AttrVector{floats[0], floats[1], floats[2]},
		Rotate:                   AttrVector{floats[3], floats[4], floats[5]},
		RotateOrder:              rotateOrder,
		Translate:                AttrVector{floats[7], floats[8], floats[9]},
		Shear:                    AttrShear{floats[10], floats[11], floats[12]},
		ScalePivot:               AttrVector{floats[13], floats[14], floats[15]},
		ScaleTranslate:           AttrVector{floats[16], floats[17], floats[18]},
		RotatePivot:              AttrVector{floats[19], floats[20], floats[21]},
		RotateTranslation:        AttrVector{floats[22], floats[23], floats[24]},
		RotateOrient:             AttrOrient{floats[25], floats[26], floats[27], floats[28]},
		JointOrient:              AttrOrient{floats[29], floats[30], floats[31], floats[32]},
		InverseParentScale:       AttrVector{floats[33], floats[34], floats[35]},
		CompensateForParentScale: componentForParentScale,
	}
	var a Attr = &mx
	return a, 38, nil
}

func MakePointArray(token *[]string, start int) (Attr, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	var a Attr
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
		a = &paa
	} else {
		paa := AttrPointArray{}
		a = &paa
	}
	return a, 1 + (numberOfArray * 4), nil
}

func MakeVectorArray(token *[]string, start int) (Attr, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	var a Attr
	if numberOfArray != 0 {
		f, err := ParseFloats((*token)[start+1 : start+1+numberOfArray]...)
		if err != nil {
			return nil, 0, err
		}
		va := make([]AttrVector, numberOfArray)
		for i := 0; i < numberOfArray; i += 3 {
			va[i] = AttrVector{
				X: f[i],
				Y: f[i+1],
				Z: f[i+2],
			}
		}
		vaa := AttrVectorArray(va)
		a = &vaa
	} else {
		vaa := AttrVectorArray{}
		a = &vaa
	}
	return a, 1 + (numberOfArray * 3), nil
}

func MakeString(token *[]string, start int) (Attr, int, error) {
	s := AttrString((*token)[start][1 : len((*token)[start])-1])
	var a Attr = &s
	return a, 1, nil
}

func MakeStringArray(token *[]string, start int) (Attr, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	sa := AttrStringArray((*token)[start+1 : start+1+numberOfArray])
	for i, s := range sa {
		sa[i] = s[1 : len(s)-1]
	}
	var a Attr = &sa
	return a, 1 + numberOfArray, nil
}

func MakeSphere(token *[]string, start int) (Attr, int, error) {
	s, err := strconv.ParseFloat((*token)[start], 64)
	if err != nil {
		return nil, 0, err
	}
	sp := AttrSphere(s)
	var a Attr = &sp
	return a, 1, nil
}

func MakeCone(token *[]string, start int) (Attr, int, error) {
	f, err := ParseFloats((*token)[start : start+2]...)
	if err != nil {
		return nil, 0, err
	}
	c := []AttrCone{
		{
			ConeAngle: f[0],
			ConeCap:   f[1],
		},
	}
	var a Attr = &c
	return a, 2, nil
}

func MakeReflectanceRGB(token *[]string, start int) (Attr, int, error) {
	f, err := ParseFloats((*token)[start : start+3]...)
	if err != nil {
		return nil, 0, err
	}
	var a Attr = &AttrReflectanceRGB{
		RedReflect:   f[0],
		GreenReflect: f[1],
		BlueReflect:  f[2],
	}
	return a, 3, nil
}

func MakeSpectrumRGB(token *[]string, start int) (Attr, int, error) {
	f, err := ParseFloats((*token)[start : start+3]...)
	if err != nil {
		return nil, 0, err
	}
	var a Attr = &AttrSpectrumRGB{
		RedSpectrum:   f[0],
		GreenSpectrum: f[1],
		BlueSpectrum:  f[2],
	}
	return a, 3, nil
}

func MakeComponentList(token *[]string, start int) (Attr, int, error) {
	numberOfArray, err := strconv.Atoi((*token)[start])
	if err != nil {
		return nil, 0, err
	}
	var cl AttrComponentList
	for _, c := range (*token)[start+1 : start+1+numberOfArray] {
		cl = append(cl, strings.Trim(c, "\""))
	}
	var a Attr = &cl
	return a, 1 + numberOfArray, nil
}

func MakeAttributeAlias(token *[]string, start int) (Attr, int, error) {
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
	var a Attr = &aaa
	return a, 2 + (len(aaa) * 2), nil
}

func MakeNurbsCurve(token *[]string, start int) (Attr, int, error) {
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
	var a Attr = &[]AttrNurbsCurve{{
		Degree:     degree,
		Spans:      spans,
		Form:       form,
		IsRational: isRational,
		Dimension:  dimension,
		KnotValues: kv,
		CvValues:   cvValues,
	}}
	count := 7 + knotCount + (cvCount * divideCv)
	return a, count, nil
}

func MakeNurbsSurface(token *[]string, start int) (Attr, int, error) {
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
	var a Attr = &[]AttrNurbsSurface{
		{
			UDegree:     uDegree,
			VDegree:     vDegree,
			UForm:       uForm,
			VForm:       vForm,
			IsRational:  isRational,
			UKnotValues: uKnotValues,
			VKnotValues: vKnotValues,
			IsTrim:      isTrim,
			CvValues:    cvValue,
		},
	}
	count := (cvStart + (cvCount * divideCv)) - start
	return a, count, nil
}

func MakeNurbsTrimface(token *[]string, start int) (Attr, int, error) {
	// TODO: Waiting for Autodesk
	var a Attr = &AttrNurbsTrimface{}
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

func MakePolyFace(token *[]string, start int, size *uint) (Attr, int, error) {
	switchNumber := start
	var pfs []AttrPolyFaces
	var fCount uint
	for _, v := range (*token)[start:] {
		if v == "f" {
			fCount ++
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
			fe, err := MakeCountInt(token, switchNumber+1)
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
	var a Attr = &pfs
	return a, switchNumber - start, nil
}

func MakeDataPolyComponent(token *[]string, start int) (Attr, int, error) {
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
	adpc := []AttrDataPolyComponent{dpc,}
	var a Attr = &adpc
	return a, 3 + (count * 2), nil
}

func MakeMesh(_ *[]string, _ int) (Attr, int, error) {
	// Not Implement
	var a Attr
	return a, -1, nil
}

func MakeLattice(token *[]string, start int) (Attr, int, error) {
	c, err := ParseInts((*token)[start : start+4]...)
	if err != nil {
		return nil, 0, err
	}
	la := []AttrLattice{
		{
			DivisionS: c[0],
			DivisionT: c[1],
			DivisionU: c[2],
		},
	}
	la[0].Points = make([]AttrLaticePoint, c[3])
	for i := 0; i < c[3]*3; i += 3 {
		p, err := ParseFloats((*token)[start+4+i : start+4+i+3]...)
		if err != nil {
			return nil, 0, err
		}
		la[0].Points[i/3].S = p[0]
		la[0].Points[i/3].T = p[1]
		la[0].Points[i/3].U = p[2]
	}
	var a Attr = &la
	return a, 4 + (c[3] * 3), nil
}

func MakeAttr(token *[]string, start int, size *uint, attrType AttrType) (Attr, int, error) {
	switch attrType {
	case TypeShort2, TypeLong2:
		return MakeShort2Long2(token, start, size, &attrType)
	case TypeShort3, TypeLong3:
		return MakeShort3Long3(token, start, size, &attrType)
	case TypeInt32Array:
		return MakeInt32Array(token, start)
	case TypeFloat2, TypeDouble2:
		return MakeFloat2Double2(token, start, size, &attrType)
	case TypeFloat3, TypeDouble3:
		return MakeFloat3Double3(token, start, size, &attrType)
	case TypeDoubleArray:
		return MakeDoubleArray(token, start)
	case TypeMatrix:
		return MakeMatrix(token, start)
	case TypeMatrixXform:
		return MakeMatrixXform(token, start)
	case TypePointArray:
		return MakePointArray(token, start)
	case TypeVectorArray:
		return MakeVectorArray(token, start)
	case TypeString:
		return MakeString(token, start)
	case TypeStringArray:
		return MakeStringArray(token, start)
	case TypeSphere:
		return MakeSphere(token, start)
	case TypeCone:
		return MakeCone(token, start)
	case TypeReflectanceRGB:
		return MakeReflectanceRGB(token, start)
	case TypeSpectrumRGB:
		return MakeSpectrumRGB(token, start)
	case TypeComponentList:
		return MakeComponentList(token, start)
	case TypeAttributeAlias:
		return MakeAttributeAlias(token, start)
	case TypeNurbsCurve:
		return MakeNurbsCurve(token, start)
	case TypeNurbsSurface:
		return MakeNurbsSurface(token, start)
	case TypeNurbsTrimface:
		return MakeNurbsTrimface(token, start)
	case TypePolyFaces:
		return MakePolyFace(token, start, size)
	case TypeDataPolyComponent:
		return MakeDataPolyComponent(token, start)
	case TypeMesh:
		return MakeMesh(token, start)
	case TypeLattice:
		return MakeLattice(token, start)
	}
	return nil, 0, nil
}

func MakeAttrType(token *[]string, start int) (AttrType, error) {
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
	case "mesh":
		return TypeMesh, nil
	case "lattice":
		return TypeLattice, nil
	}
	return TypeInvalid, errors.New("Invalid type " + typeString)
}

type Attr interface{}

type AttrShort2 [2]int

type AttrShort3 [3]int

type AttrLong2 [2]int

type AttrLong3 [3]int

type AttrInt32Array []int

type AttrFloat2 [2]float64

type AttrFloat3 [3]float64

type AttrDouble2 [2]float64

type AttrDouble3 [3]float64

type AttrDoubleArray []float64

type AttrShear struct {
	XY float64 `json:"xy"`
	XZ float64 `json:"xz"`
	YZ float64 `json:"yz"`
}

type AttrOrient struct {
	W float64 `json:"w"`
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type AttrRotateOrder int

const (
	// 0=XYZ, 1=YZX, 2=ZXY, 3=XZY, 4=YXZ, 5=ZYX
	RotateOrderXYZ AttrRotateOrder = iota
	RotateOrderYZX
	RotateOrderZXY
	RotateOrderXZY
	RotateOrderYXZ
	RotateOrderZYX
)

var rotateOrderList = [6]AttrRotateOrder{
	RotateOrderXYZ,
	RotateOrderYZX,
	RotateOrderZXY,
	RotateOrderXZY,
	RotateOrderYXZ,
	RotateOrderZYX,
}

func ConvertAttrRotateOrder(i int) (AttrRotateOrder, error) {
	if len(rotateOrderList) <= i {
		return AttrRotateOrder(i),
			errors.New(
				fmt.Sprintf("this is not AttrRotateOrder number. \"%d\"", i))
	}
	return rotateOrderList[i], nil
}

type AttrMatrix [16]float64 // mat4x4

type AttrMatrixXform struct {
	Scale                    AttrVector      `json:"scale"`
	Rotate                   AttrVector      `json:"rotate"`
	RotateOrder              AttrRotateOrder `json:"rotate_order"`
	Translate                AttrVector      `json:"translate"`
	Shear                    AttrShear       `json:"shear"`
	ScalePivot               AttrVector      `json:"scale_pivot"`
	ScaleTranslate           AttrVector      `json:"scale_translate"`
	RotatePivot              AttrVector      `json:"rotate_pivot"`
	RotateTranslation        AttrVector      `json:"rotate_translation"`
	RotateOrient             AttrOrient      `json:"rotate_orient"`
	JointOrient              AttrOrient      `json:"joint_orient"`
	InverseParentScale       AttrVector      `json:"inverse_parent_scale"`
	CompensateForParentScale bool            `json:"compensate_for_parent_scale"`
}

type AttrPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	W float64 `json:"w"`
}

type AttrPointArray []AttrPoint

type AttrVector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type AttrVectorArray []AttrVector

type AttrString string

type AttrStringArray []string

type AttrSphere float64

type AttrCone struct {
	ConeAngle float64 `json:"cone_angle"`
	ConeCap   float64 `json:"cone_cap"`
}

type AttrReflectanceRGB struct {
	RedReflect   float64 `json:"red_reflect"`
	GreenReflect float64 `json:"green_reflect"`
	BlueReflect  float64 `json:"blue_reflect"`
}

type AttrSpectrumRGB struct {
	RedSpectrum   float64 `json:"red_spectrum"`
	GreenSpectrum float64 `json:"green_spectrum"`
	BlueSpectrum  float64 `json:"blue_spectrum"`
}

type AttrComponentList []string

type AttrAttributeAlias struct {
	NewAlias    string `json:"new_alias"`
	CurrentName string `json:"current_name"`
}

type AttrFormType int

const (
	AttrFormOpen AttrFormType = iota
	AttrFormClosed
	AttrFormPeriodic
)

var attrFormTypeList = [3]AttrFormType{
	AttrFormOpen,
	AttrFormClosed,
	AttrFormPeriodic,
}

func ConvertAttrFormType(i int) (AttrFormType, error) {
	if len(attrFormTypeList) <= i {
		return AttrFormType(i),
			errors.New(
				fmt.Sprintf("this number is not FormType number. \"%d\"", i))
	}
	return attrFormTypeList[i], nil
}

type AttrCvValue struct {
	X float64  `json:"x"`
	Y float64  `json:"y"`
	Z *float64 `json:"z,omitempty"`
	W *float64 `json:"w,omitempty"`
}

type AttrNurbsCurve struct {
	Degree     int           `json:"degree"`
	Spans      int           `json:"spans"`
	Form       AttrFormType  `json:"form"`
	IsRational bool          `json:"is_rational"`
	Dimension  int           `json:"dimension"`
	KnotValues []float64     `json:"knot_values"`
	CvValues   []AttrCvValue `json:"cv_values"`
}

type AttrNurbsSurface struct {
	UDegree     int           `json:"u_degree"`
	VDegree     int           `json:"v_degree"`
	UForm       AttrFormType  `json:"u_form"`
	VForm       AttrFormType  `json:"v_form"`
	IsRational  bool          `json:"is_rational"`
	UKnotValues []float64     `json:"u_knot_values"`
	VKnotValues []float64     `json:"v_knot_values"`
	IsTrim      *bool         `json:"is_trim,omitempty"`
	CvValues    []AttrCvValue `json:"cv_values"`
}

type AttrNurbsTrimface struct {
}

type AttrFaceUV struct {
	UVSet  int   `json:"uv_set"`
	FaceUV []int `json:"face_uv"`
}

type AttrMultiColor struct {
	ColorIndex int   `json:"color_index"`
	ColorIDs   []int `json:"color_ids"`
}

type AttrPolyFaces struct {
	FaceEdge   []int            `json:"face_edge"`
	HoleEdge   []int            `json:"hole_edge"`
	FaceUV     []AttrFaceUV     `json:"face_uv"`
	FaceColor  []int            `json:"face_color"`
	MultiColor []AttrMultiColor `json:"multi_color"`
}

type AttrDPCType int

const (
	DPCedge AttrDPCType = iota
	DPCface
	DPCvertex
	DPCuv
)

type AttrDataPolyComponent struct {
	PolyComponentType AttrDPCType     `json:"poly_component_type"`
	IndexValue        map[int]float64 `json:"index_value"`
}

type AttrLaticePoint struct {
	S float64 `json:"s"`
	T float64 `json:"t"`
	U float64 `json:"u"`
}

type AttrLattice struct {
	DivisionS int               `json:"division_s"`
	DivisionT int               `json:"division_t"`
	DivisionU int               `json:"division_u"`
	Points    []AttrLaticePoint `json:"points"`
}
