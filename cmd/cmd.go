package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/nrtkbb/bufscan"
)

type Type string

const (
	FILE        Type = "file"
	WORKSPACE   Type = "workspace"
	REQUIRES    Type = "requires"
	CONNECTATTR Type = "connectAttr"
	CREATENODE  Type = "createNode"
	RENAME      Type = "rename"
	SETATTR     Type = "setAttr"
	ADDATTR     Type = "addAttr"
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
		return errors.New("file already existed")
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
	lineNo  uint
}

func (c *CmdBuilder) Append(line string) {
	c.cmdLine = append(c.cmdLine, line)
	c.lineNo++
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
	Type   Type     `json:"cmd"`
	Raw    string   `json:"raw"`
	Token  []string `json:"token"`
	LineNo uint
}

func (c *CmdBuilder) Parse() *Cmd {
	cmd := Cmd{Raw: strings.Join(c.cmdLine, "\n"), LineNo: c.lineNo}
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
		cmd.Type = Type("None")
	} else {
		cmd.Type = Type(cmd.Token[0])
	}
	return &cmd
}

type Workspace struct {
	*Cmd
	FileRule string `json:"file_rule" tag:"-fileRule"`
	Place    string `json:"place"`
}

type Requires struct {
	*Cmd
	PluginName string   `json:"plugin_name"`
	Version    string   `json:"version"`
	NodeTypes  []string `json:"node_types" tag:"-nodeType"`
	DataTypes  []string `json:"data_types" tag:"-dataType"`
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

type CreateNode struct {
	*Cmd
	NodeType   string  `json:"node_type"`
	NodeName   string  `json:"node_name" short:"-n"`
	Parent     *string `json:"parent" short:"-p"`
	Shared     bool    `json:"shared" short:"-s"`
	SkipSelect bool    `json:"skip_select" short:"-ss"`
}

type Rename struct {
	*Cmd
	From        *string `json:"from,omitempty"`
	To          *string `json:"to"`
	UUID        bool    `json:"uuid" short:"-uid"`
	IgnoreShape bool    `json:"ignore_shape" short:"-is"`
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
	Attr         []Attr   `json:"attr"`
}

type DisconnectBehaviour uint

type AddAttr struct {
	*Cmd
	AttributeType       *string              `json:"attribute_type,omitempty" short:"-at"`
	CachedInternally    *bool                `json:"cached_internally,omitempty" short:"-ci"`
	Category            *string              `json:"category,omitempty" short:"-ct"`
	DataType            *string              `json:"data_type,omitempty" short:"-dt"`
	DefaultValue        *float64             `json:"default_value,omitempty" short:"-dv"`
	DisconnectBehaviour *DisconnectBehaviour `json:"disconnect_behaviour,omitempty" short:"-dcb"`
}

type Attr interface {
	String() string
}

type AttrBool bool

func ToAttrBool(attrs []Attr) ([]*AttrBool, error) {
	ret := make([]*AttrBool, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrBool)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (ab *AttrBool) String() string {
	return fmt.Sprint(*ab)
}

func (ab *AttrBool) Bool() bool {
	return bool(*ab)
}

type AttrInt int

func ToAttrInt(attrs []Attr) ([]*AttrInt, error) {
	ret := make([]*AttrInt, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrInt)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (ai *AttrInt) String() string {
	return strconv.Itoa(int(*ai))
}

func (ai *AttrInt) Int() int {
	return int(*ai)
}

type AttrFloat float64

func ToAttrFloat(attrs []Attr) ([]*AttrFloat, error) {
	ret := make([]*AttrFloat, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrFloat)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (af *AttrFloat) String() string {
	return strconv.FormatFloat(float64(*af), 'f', -1, 64)
}

func (af *AttrFloat) Float() float64 {
	return float64(*af)
}

type AttrShort2 [2]int

func ToAttrShort2(attrs []Attr) ([]*AttrShort2, error) {
	ret := make([]*AttrShort2, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrShort2)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (as2 *AttrShort2) String() string {
	var out bytes.Buffer

	out.WriteString(strconv.Itoa(as2[0]))
	out.WriteString(" ")
	out.WriteString(strconv.Itoa(as2[1]))

	return out.String()
}

type AttrShort3 [3]int

func ToAttrShort3(attrs []Attr) ([]*AttrShort3, error) {
	ret := make([]*AttrShort3, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrShort3)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (as3 *AttrShort3) String() string {
	var out bytes.Buffer

	out.WriteString(strconv.Itoa(as3[0]))
	out.WriteString(" ")
	out.WriteString(strconv.Itoa(as3[1]))
	out.WriteString(" ")
	out.WriteString(strconv.Itoa(as3[2]))

	return out.String()
}

type AttrLong2 [2]int

func ToAttrLong2(attrs []Attr) ([]*AttrLong2, error) {
	ret := make([]*AttrLong2, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrLong2)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (al2 *AttrLong2) String() string {
	var out bytes.Buffer

	out.WriteString(strconv.Itoa(al2[0]))
	out.WriteString(" ")
	out.WriteString(strconv.Itoa(al2[1]))

	return out.String()
}

type AttrLong3 [3]int

func ToAttrLong3(attrs []Attr) ([]*AttrLong3, error) {
	ret := make([]*AttrLong3, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrLong3)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (al3 *AttrLong3) String() string {
	var out bytes.Buffer

	out.WriteString(strconv.Itoa(al3[0]))
	out.WriteString(" ")
	out.WriteString(strconv.Itoa(al3[1]))
	out.WriteString(" ")
	out.WriteString(strconv.Itoa(al3[2]))

	return out.String()
}

type AttrInt32Array []int

func ToAttrInt32Array(attrs []Attr) ([]*AttrInt32Array, error) {
	ret := make([]*AttrInt32Array, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrInt32Array)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (ai32a *AttrInt32Array) String() string {
	var s []string
	for _, i := range *ai32a {
		s = append(s, strconv.Itoa(i))
	}
	return strings.Join(s, " ")
}

type AttrFloat2 [2]float64

func ToAttrFloat2(attrs []Attr) ([]*AttrFloat2, error) {
	ret := make([]*AttrFloat2, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrFloat2)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (af2 *AttrFloat2) String() string {
	return fmt.Sprintf("%f %f", af2[0], af2[1])
}

type AttrFloat3 [3]float64

func ToAttrFloat3(attrs []Attr) ([]*AttrFloat3, error) {
	ret := make([]*AttrFloat3, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrFloat3)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (af3 *AttrFloat3) String() string {
	return fmt.Sprintf("%f %f %f", af3[0], af3[1], af3[2])
}

type AttrDouble2 [2]float64

func ToAttrDouble2(attrs []Attr) ([]*AttrDouble2, error) {
	ret := make([]*AttrDouble2, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrDouble2)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (ad2 *AttrDouble2) String() string {
	return fmt.Sprintf("%f %f", ad2[0], ad2[1])
}

type AttrDouble3 [3]float64

func ToAttrDouble3(attrs []Attr) ([]*AttrDouble3, error) {
	ret := make([]*AttrDouble3, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrDouble3)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (ad3 *AttrDouble3) String() string {
	return fmt.Sprintf("%f %f %f", ad3[0], ad3[1], ad3[2])
}

type AttrDoubleArray []float64

func ToAttrDoubleArray(attrs []Attr) ([]*AttrDoubleArray, error) {
	ret := make([]*AttrDoubleArray, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrDoubleArray)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (ada *AttrDoubleArray) String() string {
	var s []string
	for _, d := range *ada {
		s = append(s, strconv.FormatFloat(d, 'f', -1, 64))
	}
	return strings.Join(s, " ")
}

type AttrShear struct {
	XY float64 `json:"xy"`
	XZ float64 `json:"xz"`
	YZ float64 `json:"yz"`
}

func (as *AttrShear) String() string {
	return fmt.Sprintf("xy: %f, xz: %f, yz: %f",
		as.XY, as.XZ, as.YZ)
}

type AttrOrient struct {
	W float64 `json:"w"`
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (ao *AttrOrient) String() string {
	return fmt.Sprintf("x: %f, y: %f, z: %f, w: %f",
		ao.X, ao.Y, ao.Z, ao.W)
}

type AttrRotateOrder int

func ToAttrRotateOrder(attrs []Attr) ([]*AttrRotateOrder, error) {
	ret := make([]*AttrRotateOrder, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrRotateOrder)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (aro *AttrRotateOrder) String() string {
	switch *aro {
	case RotateOrderXYZ:
		return "XYZ"
	case RotateOrderYZX:
		return "YZX"
	case RotateOrderZXY:
		return "ZXY"
	case RotateOrderXZY:
		return "XZY"
	case RotateOrderYXZ:
		return "YXZ"
	case RotateOrderZYX:
		return "ZYX"
	default:
		return "XYZ"
	}
}

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

func ToAttrMatrix(attrs []Attr) ([]*AttrMatrix, error) {
	ret := make([]*AttrMatrix, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrMatrix)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (am *AttrMatrix) String() string {
	var s []string
	for _, m := range *am {
		s = append(s, strconv.FormatFloat(m, 'f', -1, 64))
	}
	return strings.Join(s, " ")
}

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

func ToAttrMatrixXform(attrs []Attr) ([]*AttrMatrixXform, error) {
	ret := make([]*AttrMatrixXform, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrMatrixXform)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (amx *AttrMatrixXform) String() string {
	var out bytes.Buffer

	out.WriteString("scale: ")
	out.WriteString(amx.Scale.String())
	out.WriteString("\nrotate: ")
	out.WriteString(amx.Rotate.String())
	out.WriteString("\nrotateOrder: ")
	out.WriteString(amx.RotateOrder.String())
	out.WriteString("\ntranslate: ")
	out.WriteString(amx.Translate.String())
	out.WriteString("\nshear: ")
	out.WriteString(amx.Shear.String())
	out.WriteString("\nscalePivot: ")
	out.WriteString(amx.ScalePivot.String())
	out.WriteString("\nscaleTranslate: ")
	out.WriteString(amx.ScaleTranslate.String())
	out.WriteString("\nrotatePivot: ")
	out.WriteString(amx.RotatePivot.String())
	out.WriteString("\nrotateTranslate: ")
	out.WriteString(amx.RotateTranslation.String())
	out.WriteString("\nrotateOrient: ")
	out.WriteString(amx.RotateOrient.String())
	out.WriteString("\njointOrient: ")
	out.WriteString(amx.JointOrient.String())
	out.WriteString("\ninverseParentScale: ")
	out.WriteString(amx.InverseParentScale.String())
	out.WriteString("\ncompensateForParentScale: ")
	out.WriteString(fmt.Sprint(amx.CompensateForParentScale))

	return out.String()
}

type AttrPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	W float64 `json:"w"`
}

func ToAttrPoint(attrs []Attr) ([]*AttrPoint, error) {
	ret := make([]*AttrPoint, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrPoint)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (ap *AttrPoint) String() string {
	return fmt.Sprintf("x: %f, y: %f, z: %f, w: %f",
		ap.X, ap.Y, ap.Z, ap.W)
}

type AttrPointArray []AttrPoint

func ToAttrPointArray(attrs []Attr) ([]*AttrPointArray, error) {
	ret := make([]*AttrPointArray, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrPointArray)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (apa *AttrPointArray) String() string {
	var s []string
	for _, ap := range *apa {
		s = append(s, ap.String())
	}

	return strings.Join(s, ", ")
}

type AttrVector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func ToAttrVector(attrs []Attr) ([]*AttrVector, error) {
	ret := make([]*AttrVector, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrVector)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (av *AttrVector) String() string {
	return fmt.Sprintf("x: %f, y: %f, z: %f",
		av.X, av.Y, av.Z)
}

type AttrVectorArray []AttrVector

func ToAttrVectorArray(attrs []Attr) ([]*AttrVectorArray, error) {
	ret := make([]*AttrVectorArray, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrVectorArray)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (ava *AttrVectorArray) String() string {
	var s []string
	for _, av := range *ava {
		s = append(s, av.String())
	}
	return strings.Join(s, ", ")
}

type AttrString string

func ToAttrString(attrs []Attr) ([]*AttrString, error) {
	ret := make([]*AttrString, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrString)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (as *AttrString) String() string {
	return string(*as)
}

type AttrStringArray []string

func ToAttrStringArray(attrs []Attr) ([]*AttrStringArray, error) {
	ret := make([]*AttrStringArray, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrStringArray)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (asa *AttrStringArray) String() string {
	var s []string
	for _, as := range *asa {
		s = append(s, as)
	}
	return strings.Join(s, ", ")
}

type AttrSphere float64

func ToAttrSphere(attrs []Attr) ([]*AttrSphere, error) {
	ret := make([]*AttrSphere, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrSphere)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (as *AttrSphere) String() string {
	return fmt.Sprint(*as)
}

type AttrCone struct {
	ConeAngle float64 `json:"cone_angle"`
	ConeCap   float64 `json:"cone_cap"`
}

func ToAttrCone(attrs []Attr) ([]*AttrCone, error) {
	ret := make([]*AttrCone, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrCone)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (ac *AttrCone) String() string {
	var out bytes.Buffer

	out.WriteString("coneAngle: ")
	out.WriteString(fmt.Sprint(ac.ConeAngle))
	out.WriteString("\nconeCap: ")
	out.WriteString(fmt.Sprint(ac.ConeCap))

	return out.String()
}

type AttrReflectanceRGB struct {
	RedReflect   float64 `json:"red_reflect"`
	GreenReflect float64 `json:"green_reflect"`
	BlueReflect  float64 `json:"blue_reflect"`
}

func ToAttrReflectanceRGB(attrs []Attr) ([]*AttrReflectanceRGB, error) {
	ret := make([]*AttrReflectanceRGB, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrReflectanceRGB)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (ar *AttrReflectanceRGB) String() string {
	var out bytes.Buffer

	out.WriteString("redReflect: ")
	out.WriteString(fmt.Sprint(ar.RedReflect))
	out.WriteString("\ngreenReflect: ")
	out.WriteString(fmt.Sprint(ar.GreenReflect))
	out.WriteString("\nblueReflect: ")
	out.WriteString(fmt.Sprint(ar.BlueReflect))

	return out.String()
}

type AttrSpectrumRGB struct {
	RedSpectrum   float64 `json:"red_spectrum"`
	GreenSpectrum float64 `json:"green_spectrum"`
	BlueSpectrum  float64 `json:"blue_spectrum"`
}

func ToAttrSpectrumRGB(attrs []Attr) ([]*AttrSpectrumRGB, error) {
	ret := make([]*AttrSpectrumRGB, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrSpectrumRGB)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (as *AttrSpectrumRGB) String() string {
	var out bytes.Buffer

	out.WriteString("redSpectrum: ")
	out.WriteString(fmt.Sprint(as.RedSpectrum))
	out.WriteString("\ngreenSpectrum: ")
	out.WriteString(fmt.Sprint(as.GreenSpectrum))
	out.WriteString("\nblueSpectrum: ")
	out.WriteString(fmt.Sprint(as.BlueSpectrum))

	return out.String()
}

type AttrComponentList []string

func ToAttrComponentList(attrs []Attr) ([]*AttrComponentList, error) {
	ret := make([]*AttrComponentList, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrComponentList)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (acl *AttrComponentList) String() string {
	var s []string
	for _, ac := range *acl {
		s = append(s, ac)
	}
	return strings.Join(s, ", ")
}

type AttrAttributeAlias struct {
	NewAlias    string `json:"new_alias"`
	CurrentName string `json:"current_name"`
}

func ToAttrAttributeAlias(attrs []Attr) ([]*AttrAttributeAlias, error) {
	ret := make([]*AttrAttributeAlias, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrAttributeAlias)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (aaa *AttrAttributeAlias) String() string {
	var out bytes.Buffer

	out.WriteString("newAlias: ")
	out.WriteString(aaa.NewAlias)
	out.WriteString("\ncurrentName: ")
	out.WriteString(aaa.CurrentName)

	return out.String()
}

type AttrFormType int

func ToAttrFormType(attrs []Attr) ([]*AttrFormType, error) {
	ret := make([]*AttrFormType, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrFormType)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (aft *AttrFormType) String() string {
	switch *aft {
	case AttrFormOpen:
		return "Open"
	case AttrFormClosed:
		return "Close"
	case AttrFormPeriodic:
		return "Periodic"
	default:
		return "Open"
	}
}

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

func ToAttrCvValue(attrs []Attr) ([]*AttrCvValue, error) {
	ret := make([]*AttrCvValue, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrCvValue)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (acv *AttrCvValue) String() string {
	return fmt.Sprintf("x: %f, y: %f, z: %+v, w: %+v",
		acv.X, acv.Y, acv.Z, acv.W)
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

func ToAttrNurbsCurve(attrs []Attr) ([]*AttrNurbsCurve, error) {
	ret := make([]*AttrNurbsCurve, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrNurbsCurve)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (anc *AttrNurbsCurve) String() string {
	var out bytes.Buffer

	out.WriteString("degree: ")
	out.WriteString(strconv.Itoa(anc.Degree))
	out.WriteString("\nspans: ")
	out.WriteString(strconv.Itoa(anc.Spans))
	out.WriteString("\nform: ")
	out.WriteString(anc.Form.String())
	out.WriteString("\nisRational: ")
	out.WriteString(fmt.Sprint(anc.IsRational))

	return out.String()
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

func ToAttrNurbsSurface(attrs []Attr) ([]*AttrNurbsSurface, error) {
	ret := make([]*AttrNurbsSurface, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrNurbsSurface)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (ans *AttrNurbsSurface) String() string {
	var out bytes.Buffer

	out.WriteString("uDegree: ")
	out.WriteString(strconv.Itoa(ans.UDegree))
	out.WriteString("\nvDegree: ")
	out.WriteString(strconv.Itoa(ans.VDegree))
	out.WriteString("\nuForm: ")
	out.WriteString(ans.UForm.String())
	out.WriteString("\nvForm: ")
	out.WriteString(ans.VForm.String())
	out.WriteString("\nisRational: ")
	out.WriteString(fmt.Sprint(ans.IsRational))

	out.WriteString("\nuKnotValues: {")
	var uKnot []string
	for _, u := range ans.UKnotValues {
		uKnot = append(uKnot, strconv.FormatFloat(u, 'f', -1, 64))
	}
	out.WriteString(strings.Join(uKnot, ", "))

	out.WriteString("}\nvKnotValues: {")
	var vKnot []string
	for _, v := range ans.VKnotValues {
		vKnot = append(vKnot, strconv.FormatFloat(v, 'f', -1, 64))
	}
	out.WriteString(strings.Join(vKnot, ", "))

	out.WriteString("}\nisTrim: ")
	if ans.IsTrim == nil {
		out.WriteString("nil")
	} else {
		out.WriteString(fmt.Sprint(*ans.IsTrim))
	}

	out.WriteString("\ncvValues: ")
	var cvs []string
	for _, cv := range ans.CvValues {
		cvs = append(cvs, cv.String())
	}
	out.WriteString(strings.Join(cvs, ", "))

	return out.String()
}

type AttrNurbsTrimface struct{}

func ToAttrNurbsTrimface(attrs []Attr) ([]*AttrNurbsTrimface, error) {
	ret := make([]*AttrNurbsTrimface, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrNurbsTrimface)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (ant *AttrNurbsTrimface) String() string {
	return ""
}

type AttrFaceUV struct {
	UVSet  int   `json:"uv_set"`
	FaceUV []int `json:"face_uv"`
}

func ToAttrFaceUV(attrs []Attr) ([]*AttrFaceUV, error) {
	ret := make([]*AttrFaceUV, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrFaceUV)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (af *AttrFaceUV) String() string {
	var s []string
	for _, uv := range af.FaceUV {
		s = append(s, strconv.Itoa(uv))
	}
	return fmt.Sprintf("uvSet: %d, faceUV: %s",
		af.UVSet, strings.Join(s, ", "))
}

type AttrMultiColor struct {
	ColorIndex int   `json:"color_index"`
	ColorIDs   []int `json:"color_ids"`
}

func ToAttrMultiColor(attrs []Attr) ([]*AttrMultiColor, error) {
	ret := make([]*AttrMultiColor, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrMultiColor)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (amc *AttrMultiColor) String() string {
	var s []string
	for _, id := range amc.ColorIDs {
		s = append(s, strconv.Itoa(id))
	}
	return fmt.Sprintf("colorIndex: %d, colorIDs: %s",
		amc.ColorIndex, strings.Join(s, ", "))
}

type AttrPolyFaces struct {
	FaceEdge   []int            `json:"face_edge"`
	HoleEdge   []int            `json:"hole_edge"`
	FaceUV     []AttrFaceUV     `json:"face_uv"`
	FaceColor  []int            `json:"face_color"`
	MultiColor []AttrMultiColor `json:"multi_color"`
}

func ToAttrPolyFaces(attrs []Attr) ([]*AttrPolyFaces, error) {
	ret := make([]*AttrPolyFaces, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrPolyFaces)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (apf *AttrPolyFaces) String() string {
	var out bytes.Buffer

	out.WriteString("faceEdge: ")
	out.WriteString(intArrayToString(apf.FaceEdge))

	out.WriteString(", holeEdge: ")
	out.WriteString(intArrayToString(apf.HoleEdge))

	out.WriteString(", faceUV: ")
	var faceUV []string
	for _, fuv := range apf.FaceUV {
		faceUV = append(faceUV, fuv.String())
	}
	out.WriteString(strings.Join(faceUV, ", "))

	out.WriteString(", faceColor: ")
	out.WriteString(intArrayToString(apf.FaceColor))

	out.WriteString(", multiColor: ")
	var multiColor []string
	for _, mc := range apf.MultiColor {
		multiColor = append(multiColor, mc.String())
	}
	out.WriteString(strings.Join(multiColor, ", "))

	return out.String()
}

type AttrDPCType int

func ToAttrDPCType(attrs []Attr) ([]*AttrDPCType, error) {
	ret := make([]*AttrDPCType, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrDPCType)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (dpc *AttrDPCType) String() string {
	switch *dpc {
	case DPCedge:
		return "edge"
	case DPCface:
		return "face"
	case DPCvertex:
		return "vertex"
	case DPCuv:
		return "uv"
	default:
		return "edge"
	}
}

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

func ToAttrDataPolyComponent(attrs []Attr) ([]*AttrDataPolyComponent, error) {
	ret := make([]*AttrDataPolyComponent, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrDataPolyComponent)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (adpc *AttrDataPolyComponent) String() string {
	var indecies []int
	for k := range adpc.IndexValue {
		indecies = append(indecies, k)
	}
	sort.Ints(indecies)
	var s []string
	for _, i := range indecies {
		s = append(s, fmt.Sprintf("%d %f", i, adpc.IndexValue[i]))
	}
	return fmt.Sprintf("polyComponentType: %s, indexValue: %s",
		adpc.PolyComponentType.String(), strings.Join(s, ", "))
}

type AttrLatticePoint struct {
	S float64 `json:"s"`
	T float64 `json:"t"`
	U float64 `json:"u"`
}

func ToAttrLatticePoint(attrs []Attr) ([]*AttrLatticePoint, error) {
	ret := make([]*AttrLatticePoint, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrLatticePoint)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (alp *AttrLatticePoint) String() string {
	return fmt.Sprintf("s: %f, t: %f, u: %f",
		alp.S, alp.T, alp.U)
}

type AttrLattice struct {
	DivisionS int                `json:"division_s"`
	DivisionT int                `json:"division_t"`
	DivisionU int                `json:"division_u"`
	Points    []AttrLatticePoint `json:"points"`
}

func ToAttrLattice(attrs []Attr) ([]*AttrLattice, error) {
	ret := make([]*AttrLattice, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrLattice)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (al *AttrLattice) String() string {
	var out bytes.Buffer

	out.WriteString("divisionS: ")
	out.WriteString(strconv.Itoa(al.DivisionS))
	out.WriteString(", divisionT: ")
	out.WriteString(strconv.Itoa(al.DivisionT))
	out.WriteString(", divisionU: ")
	out.WriteString(strconv.Itoa(al.DivisionU))

	return out.String()
}

func intArrayToString(intArray []int) string {
	var s []string
	for _, i := range intArray {
		s = append(s, strconv.Itoa(i))
	}
	return strings.Join(s, ", ")
}
