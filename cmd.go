package mayaascii

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

type Type string

const (
	NoneType         Type = ""
	LineCommentType  Type = "//"
	BlockCommentType Type = "/*"
	FileType         Type = "file"
	FileInfoType     Type = "fileInfo"
	WorkspaceType    Type = "workspace"
	RequiresType     Type = "requires"
	ConnectAttrType  Type = "connectAttr"
	CreateNodeType   Type = "createNode"
	RenameType       Type = "rename"
	SetAttrType      Type = "setAttr"
	AddAttrType      Type = "addAttr"
	SelectType       Type = "select"
)

func (t Type) String() string {
	return string(t)
}

func (t Type) HasPrefix(line string) bool {
	line = strings.TrimLeft(line, " \t\n")
	return strings.HasPrefix(line, t.String())
}

func (t Type) HasPrefixWithSpace(line string) bool {
	line = strings.TrimLeft(line, " \t\n")
	return strings.HasPrefix(line, t.String()+" ")
}

type CmdBuilder struct {
	cmdLine        []string
	lineNo         uint
	isBlockComment bool
}

func (c *CmdBuilder) Append(line string) {
	if !c.isBlockComment && BlockCommentType.HasPrefix(line) {
		c.isBlockComment = true
	}
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
	if c.isBlockComment {
		return strings.HasSuffix(
			strings.TrimRight(lastLine, " \t\r\n"), "*/")
	} else {
		return lastLine[len(lastLine)-1] == byte(';')
	}
}

func (c *CmdBuilder) Clear() {
	c.cmdLine = []string{}
}

func (c *CmdBuilder) IsClear() bool {
	return len(c.cmdLine) == 0
}

const (
	whiteSpace        = ' '
	tabSpace          = '\t'
	enter             = '\n'
	slash             = '/'
	asterisk          = '*'
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
	cmd := Cmd{
		Raw:    strings.Join(c.cmdLine, "\n"),
		LineNo: c.lineNo,
		Type:   NoneType,
	}
	// fmt.Println("Raw", cmd.Raw)
	c.Clear()
	var buf []rune
	var subBuf []rune
	var subToken []string
	for _, c := range cmd.Raw {
		// comment
		if (0 == len(buf) && c == slash) ||
			(1 == len(buf) && buf[0] == slash && (c == slash || c == asterisk)) {
			buf = append(buf, c)
			if 2 == len(buf) {
				cmd.Token = append(cmd.Token, string(buf))
				buf = buf[:0]
				if c == slash {
					cmd.Type = LineCommentType
				} else if c == asterisk {
					cmd.Type = BlockCommentType
				}
			}
			continue
		}
		if cmd.Type == BlockCommentType {
			if 2 <= len(buf) && buf[len(buf)-1] == asterisk && c == slash {
				cmd.Token = append(cmd.Token, string(buf[:len(buf)-1]))
				break
			}
			buf = append(buf, c)
			continue
		}
		if cmd.Type == LineCommentType {
			buf = append(buf, c)
			continue
		}
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
				// buf = [:0] はしない
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
					if len(subToken) == 0 {
						subToken = append(subToken, "\"")
					}
					subToken = append(subToken, string(subBuf)[1:len(string(subBuf))-1])
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
				subToken = append(subToken, "\"")
				cmd.Token = append(cmd.Token, strings.Join(subToken, ""))
				buf = buf[:0]
				subToken = subToken[:0]
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
	if 0 != len(buf) && cmd.Type == LineCommentType {
		cmd.Token = append(cmd.Token, string(buf))
	}
	if 0 == len(cmd.Token) {
		cmd.Type = NoneType
	} else {
		cmd.Type = Type(cmd.Token[0])
	}
	return &cmd
}

type LineCommentCmd struct {
	*Cmd
	Comment string `json:"comment"`
}

type BlockCommentCmd struct {
	*Cmd
	Comment string `json:"comment"`
}

type FileCmd struct {
	*Cmd
	Path               string `json:"path"`
	ReferenceDepthInfo int    `json:"reference_depth_info" type:"-referenceDepthInfo"`
	Namespace          string `json:"namespace" type:"-namespace"`
	ReferenceNode      string `json:"reference_node" type:"-referenceNode"`
	Options            string `json:"options" type:"-options"`
	Type               string `json:"type" type:"-type"`
	Reference          bool   `json:"reference" type:"-reference"`
	DeferReference     bool   `json:"defer_reference" type:"-deferReference"`
}

func (f *FileCmd) String() string {
	var buf bytes.Buffer
	buf.WriteString("file ")
	if f.Reference {
		buf.WriteString("-r ")
	}
	if f.ReferenceDepthInfo != 0 {
		buf.WriteString(fmt.Sprintf("-rdi %d ", f.ReferenceDepthInfo))
	}
	buf.WriteString("-ns \"")
	buf.WriteString(f.Namespace)
	buf.WriteString("\"")
	if f.DeferReference {
		buf.WriteString(" -dr 1")
	}
	buf.WriteString(" -rfn \"")
	buf.WriteString(f.ReferenceNode)
	buf.WriteString("\" -op \"")
	buf.WriteString(f.Options)
	buf.WriteString("\" -typ \"")
	buf.WriteString(f.Type)
	buf.WriteString("\" ")
	if len(buf.String())+len(f.Path) > 160 {
		buf.WriteString("\n\t\t")
	}
	buf.WriteString("\"")
	buf.WriteString(f.Path)
	buf.WriteString("\";")
	return buf.String()
}

type FileInfoCmd struct {
	*Cmd
	Name  string `json:"name" tag:"-fileInfo"`
	Value string `json:"value" tag:"-value"`
}

func (fi *FileInfoCmd) StringWrite(writer io.StringWriter) (int, error) {
	n, err := writer.WriteString("fileInfo \"")
	if err != nil {
		return 0, err
	}
	na, err := writer.WriteString(fi.Name)
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString("\" \"")
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString(fi.Value)
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString("\";\n")
	if err != nil {
		return 0, err
	}
	n += na
	return n, err
}

func (fi *FileInfoCmd) String() string {
	return strings.Join([]string{"fileInfo \"",
		fi.Name,
		"\" \"",
		fi.Value,
		"\";"},
		"")
}

type WorkspaceCmd struct {
	*Cmd
	FileRule string `json:"file_rule" tag:"-fileRule"`
	Place    string `json:"place"`
}

type RequiresCmd struct {
	*Cmd
	PluginName string   `json:"plugin_name"`
	Version    string   `json:"version"`
	NodeTypes  []string `json:"node_types" tag:"-nodeType"`
	DataTypes  []string `json:"data_types" tag:"-dataType"`
}

type ConnectAttrCmd struct {
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

type CreateNodeCmd struct {
	*Cmd
	NodeType   string  `json:"node_type"`
	NodeName   string  `json:"node_name" short:"-n"`
	Parent     *string `json:"parent" short:"-p"`
	Shared     bool    `json:"shared" short:"-s"`
	SkipSelect bool    `json:"skip_select" short:"-ss"`
}

type RenameCmd struct {
	*Cmd
	From        *string `json:"from,omitempty"`
	To          *string `json:"to"`
	UUID        bool    `json:"uuid" short:"-uid"`
	IgnoreShape bool    `json:"ignore_shape" short:"-is"`
}

type SelectCmd struct {
	*Cmd
	Names              []string `json:"names"`
	Add                bool     `json:"add" short:"-add"`
	AddFirst           bool     `json:"add_first" short:"-af"`
	All                bool     `json:"all" short:"-all"`
	AllDagObjects      bool     `json:"all_dag_objects" short:"-ado"`
	AllDependencyNodes bool     `json:"all_dependency_nodes" short:"-adn"`
	Clear              bool     `json:"clear" short:"-cl"`
	ContainerCentric   bool     `json:"container_centric" short:"-cc"`
	Deselect           bool     `json:"deselect" short:"-d"`
	Hierarchy          bool     `json:"hierarchy" short:"-hi"`
	NoExpand           bool     `json:"no_expand" short:"-ne"`
	Replace            bool     `json:"replace" short:"-r"`
	Symmetry           bool     `json:"symmetry" short:"-sym"`
	SymmetrySide       bool     `json:"symmetry_side" short:"-sys"`
	Toggle             bool     `json:"toggle" short:"-tgl"`
	Visible            bool     `json:"visible" short:"-vis"`
}

type SetAttrCmd struct {
	*Cmd
	AttrName     string      `json:"attr_name"`
	AlteredValue bool        `json:"altered_value" short:"-av"`
	Caching      *bool       `json:"caching,omitempty" short:"-ca"`
	CapacityHint *uint       `json:"capacity_hint,omitempty" short:"-ch"`
	ChannelBox   *bool       `json:"channel_box,omitempty" short:"-cb"`
	Clamp        bool        `json:"clamp" short:"-c"`
	Keyable      *bool       `json:"keyable,omitempty" short:"-k"`
	Lock         *bool       `json:"lock,omitempty" short:"-l"`
	Size         *uint       `json:"size,omitempty" short:"-s"`
	AttrType     AttrType    `json:"attr_type" short:"-typ"`
	Attr         []AttrValue `json:"attr"`
}

func (sa *SetAttrCmd) StringWrite(writer io.StringWriter) (int, error) {
	n, err := writer.WriteString("\tsetAttr ")
	if err != nil {
		return 0, err
	}
	if sa.AlteredValue {
		na, err := writer.WriteString("-av ")
		if err != nil {
			return 0, err
		}
		n += na
	}
	if sa.Clamp {
		na, err := writer.WriteString("-c ")
		if err != nil {
			return 0, err
		}
		n += na
	}
	if sa.Caching != nil {
		na, err := writer.WriteString("-ca ")
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writeOnOffBw(writer, sa.Caching)
		if err != nil {
			return 0, err
		}
		n += na
	}
	if sa.ChannelBox != nil {
		na, err := writer.WriteString("-cb ")
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writeOnOffBw(writer, sa.ChannelBox)
		if err != nil {
			return 0, err
		}
		n += na
	}
	if sa.Keyable != nil {
		na, err := writer.WriteString("-k ")
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writeOnOffBw(writer, sa.Keyable)
		if err != nil {
			return 0, err
		}
		n += na
	}
	if sa.Lock != nil {
		na, err := writer.WriteString("-l ")
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writeOnOffBw(writer, sa.Lock)
		if err != nil {
			return 0, err
		}
		n += na
	}
	if sa.Size != nil {
		na, err := writer.WriteString("-s ")
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writer.WriteString(strconv.FormatUint(uint64(*sa.Size), 10))
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writer.WriteString(" ")
		if err != nil {
			return 0, err
		}
		n += na
	}
	if sa.CapacityHint != nil {
		na, err := writer.WriteString("-ch ")
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writer.WriteString(strconv.FormatUint(uint64(*sa.CapacityHint), 10))
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writer.WriteString(" ")
		if err != nil {
			return 0, err
		}
		n += na
	}
	na, err := writer.WriteString("\"")
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString(sa.AttrName)
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString("\" ")
	if err != nil {
		return 0, err
	}
	n += na
	if sa.AttrType == TypeBool || sa.AttrType == TypeInt || sa.AttrType == TypeDouble {
		// nothing
	} else {
		na, err = writer.WriteString("-type \"")
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writer.WriteString(sa.AttrType.Name())
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writer.WriteString("\" ")
		if err != nil {
			return 0, err
		}
		n += na
	}
	for i, a := range sa.Attr {
		if i != 0 {
			na, err = writer.WriteString("\n")
			if err != nil {
				return 0, err
			}
			n += na
		}
		na, err = a.StringWrite(writer)
		if err != nil {
			return 0, err
		}
		n += na
	}
	na, err = writer.WriteString(";\n")
	if err != nil {
		return 0, err
	}
	n += na
	return n, nil
}

func (sa *SetAttrCmd) String() string {
	var buf bytes.Buffer
	buf.WriteString("\tsetAttr ")
	if sa.AlteredValue {
		buf.WriteString("-av ")
	}
	if sa.Clamp {
		buf.WriteString("-c ")
	}
	if sa.Caching != nil {
		buf.WriteString("-ca ")
		writeOnOff(&buf, sa.Caching)
	}
	if sa.ChannelBox != nil {
		buf.WriteString("-cb ")
		writeOnOff(&buf, sa.ChannelBox)
	}
	if sa.Keyable != nil {
		buf.WriteString("-k ")
		writeOnOff(&buf, sa.Keyable)
	}
	if sa.Lock != nil {
		buf.WriteString("-l ")
		writeOnOff(&buf, sa.Lock)
	}
	if sa.Size != nil {
		buf.WriteString("-s ")
		buf.WriteString(strconv.FormatUint(uint64(*sa.Size), 10))
		buf.WriteString(" ")
	}
	if sa.CapacityHint != nil {
		buf.WriteString("-ch ")
		buf.WriteString(strconv.FormatUint(uint64(*sa.CapacityHint), 10))
		buf.WriteString(" ")
	}
	buf.WriteString("\"")
	buf.WriteString(sa.AttrName)
	buf.WriteString("\" ")
	if sa.AttrType == TypeBool || sa.AttrType == TypeInt || sa.AttrType == TypeDouble {
		// nothing
	} else {
		buf.WriteString("-type \"")
		buf.WriteString(sa.AttrType.Name())
		buf.WriteString("\" ")
	}
	for i, a := range sa.Attr {
		if i != 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(a.String())
	}
	buf.WriteString(";\n")
	return buf.String()
}

type DisconnectBehaviour uint

type AddAttrCmd struct {
	*Cmd
	AttributeType       *string              `json:"attribute_type,omitempty" short:"-at"`
	CachedInternally    *bool                `json:"cached_internally,omitempty" short:"-ci"`
	Category            *string              `json:"category,omitempty" short:"-ct"`
	DataType            *string              `json:"data_type,omitempty" short:"-dt"`
	DefaultValue        *float64             `json:"default_value,omitempty" short:"-dv"`
	DisconnectBehaviour *DisconnectBehaviour `json:"disconnect_behaviour,omitempty" short:"-dcb"`
}

type AttrValue interface {
	String() string
	StringWrite(writer io.StringWriter) (int, error)
}

type AttrBool bool

func ToAttrBool(attrs []AttrValue) ([]*AttrBool, error) {
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

func (ab *AttrBool) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(fmt.Sprint(*ab))
}

func (ab *AttrBool) Bool() bool {
	return bool(*ab)
}

type AttrInt int

func ToAttrInt(attrs []AttrValue) ([]*AttrInt, error) {
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

func (ai *AttrInt) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(strconv.Itoa(int(*ai)))
}

func (ai *AttrInt) Int() int {
	return int(*ai)
}

type AttrFloat float64

func ToAttrFloat(attrs []AttrValue) ([]*AttrFloat, error) {
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

func (af *AttrFloat) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(strconv.FormatFloat(float64(*af), 'f', -1, 64))
}

func (af *AttrFloat) Float() float64 {
	return float64(*af)
}

type AttrShort2 [2]int

func ToAttrShort2(attrs []AttrValue) ([]*AttrShort2, error) {
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

func (as2 *AttrShort2) StringWrite(writer io.StringWriter) (int, error) {
	n, err := writer.WriteString(strconv.Itoa(as2[0]))
	if err != nil {
		return 0, err
	}
	na, err := writer.WriteString(" ")
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString(strconv.Itoa(as2[1]))
	if err != nil {
		return 0, err
	}
	n += na
	return n, nil
}

type AttrShort3 [3]int

func ToAttrShort3(attrs []AttrValue) ([]*AttrShort3, error) {
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

func (as3 *AttrShort3) StringWrite(writer io.StringWriter) (int, error) {
	n, err := writer.WriteString(strconv.Itoa(as3[0]))
	if err != nil {
		return 0, err
	}
	na, err := writer.WriteString(" ")
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString(strconv.Itoa(as3[1]))
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString(" ")
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString(strconv.Itoa(as3[2]))
	if err != nil {
		return 0, err
	}
	n += na
	return n, nil
}

type AttrLong2 [2]int

func ToAttrLong2(attrs []AttrValue) ([]*AttrLong2, error) {
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

func (al2 *AttrLong2) StringWrite(writer io.StringWriter) (int, error) {
	n, err := writer.WriteString(strconv.Itoa(al2[0]))
	if err != nil {
		return 0, err
	}
	na, err := writer.WriteString(" ")
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString(strconv.Itoa(al2[1]))
	if err != nil {
		return 0, err
	}
	return n, nil
}

type AttrLong3 [3]int

func ToAttrLong3(attrs []AttrValue) ([]*AttrLong3, error) {
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

func (al3 *AttrLong3) StringWrite(writer io.StringWriter) (int, error) {
	n, err := writer.WriteString(strconv.Itoa(al3[0]))
	if err != nil {
		return 0, err
	}
	na, err := writer.WriteString(" ")
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString(strconv.Itoa(al3[1]))
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString(" ")
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString(strconv.Itoa(al3[2]))
	if err != nil {
		return 0, err
	}
	n += na
	return n, nil
}

type AttrInt32Array []int

func ToAttrInt32Array(attrs []AttrValue) ([]*AttrInt32Array, error) {
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

func (ai32a *AttrInt32Array) StringWrite(writer io.StringWriter) (int, error) {
	n := 0
	for i, ai := range *ai32a {
		if i != 0 {
			na, err := writer.WriteString(" ")
			if err != nil {
				return 0, err
			}
			n += na
		}
		na, err := writer.WriteString(strconv.Itoa(ai))
		if err != nil {
			return 0, err
		}
		n += na
	}
	return n, nil
}

type AttrFloat2 [2]float64

func ToAttrFloat2(attrs []AttrValue) ([]*AttrFloat2, error) {
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

func (af2 *AttrFloat2) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(af2.String())
}

type AttrFloat3 [3]float64

func ToAttrFloat3(attrs []AttrValue) ([]*AttrFloat3, error) {
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

func (af3 *AttrFloat3) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(af3.String())
}

type AttrDouble2 [2]float64

func ToAttrDouble2(attrs []AttrValue) ([]*AttrDouble2, error) {
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

func (ad2 *AttrDouble2) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(ad2.String())
}

type AttrDouble3 [3]float64

func ToAttrDouble3(attrs []AttrValue) ([]*AttrDouble3, error) {
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

func (ad3 *AttrDouble3) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(ad3.String())
}

type AttrDoubleArray []float64

func ToAttrDoubleArray(attrs []AttrValue) ([]*AttrDoubleArray, error) {
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

func (ada *AttrDoubleArray) StringWrite(writer io.StringWriter) (int, error) {
	n := 0
	for i, d := range *ada {
		if i != 0 {
			na, err := writer.WriteString(" ")
			if err != nil {
				return 0, err
			}
			n += na
		}
		na, err := writer.WriteString(strconv.FormatFloat(d, 'f', -1, 64))
		if err != nil {
			return 0, err
		}
		n += na
	}
	return n, nil
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

func (as *AttrShear) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(as.String())
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

func (ao *AttrOrient) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(ao.String())
}

type AttrRotateOrder int

func ToAttrRotateOrder(attrs []AttrValue) ([]*AttrRotateOrder, error) {
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

func (aro *AttrRotateOrder) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(aro.String())
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

func ToAttrMatrix(attrs []AttrValue) ([]*AttrMatrix, error) {
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

func (am *AttrMatrix) StringWrite(writer io.StringWriter) (int, error) {
	n := 0
	for i, m := range *am {
		if i != 0 {
			na, err := writer.WriteString(" ")
			if err != nil {
				return 0, err
			}
			n += na
		}
		na, err := writer.WriteString(strconv.FormatFloat(m, 'f', -1, 64))
		if err != nil {
			return 0, err
		}
		n += na
	}
	return n, nil
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

func ToAttrMatrixXform(attrs []AttrValue) ([]*AttrMatrixXform, error) {
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

func (amx *AttrMatrixXform) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(amx.String())
}

type AttrPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	W float64 `json:"w"`
}

func ToAttrPoint(attrs []AttrValue) ([]*AttrPoint, error) {
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

func (ap *AttrPoint) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(ap.String())
}

type AttrPointArray []AttrPoint

func ToAttrPointArray(attrs []AttrValue) ([]*AttrPointArray, error) {
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

func (apa *AttrPointArray) StringWrite(writer io.StringWriter) (int, error) {
	n := 0
	for i, ap := range *apa {
		if i != 0 {
			na, err := writer.WriteString(", ")
			if err != nil {
				return 0, err
			}
			n += na
		}
		na, err := writer.WriteString(ap.String())
		if err != nil {
			return 0, err
		}
		n += na
	}
	return n, nil
}

type AttrVector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func ToAttrVector(attrs []AttrValue) ([]*AttrVector, error) {
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

func (av *AttrVector) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(av.String())
}

type AttrVectorArray []AttrVector

func ToAttrVectorArray(attrs []AttrValue) ([]*AttrVectorArray, error) {
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

func (ava *AttrVectorArray) StringWrite(writer io.StringWriter) (int, error) {
	n := 0
	for i, av := range *ava {
		if i != 0 {
			na, err := writer.WriteString(", ")
			if err != nil {
				return 0, err
			}
			n += na
		}
		na, err := writer.WriteString(av.String())
		if err != nil {
			return 0, err
		}
		n += na
	}
	return n, nil
}

type AttrString string

func ToAttrString(attrs []AttrValue) ([]*AttrString, error) {
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

func (as *AttrString) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(as.String())
}

type AttrStringArray []string

func ToAttrStringArray(attrs []AttrValue) ([]*AttrStringArray, error) {
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

func (asa *AttrStringArray) StringWrite(writer io.StringWriter) (int, error) {
	n := 0
	for i, as := range *asa {
		if i != 0 {
			na, err := writer.WriteString(", ")
			if err != nil {
				return 0, err
			}
			n += na
		}
		na, err := writer.WriteString(as)
		if err != nil {
			return 0, err
		}
		n += na
	}
	return n, nil
}

type AttrSphere float64

func ToAttrSphere(attrs []AttrValue) ([]*AttrSphere, error) {
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

func (as *AttrSphere) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(as.String())
}

type AttrCone struct {
	ConeAngle float64 `json:"cone_angle"`
	ConeCap   float64 `json:"cone_cap"`
}

func ToAttrCone(attrs []AttrValue) ([]*AttrCone, error) {
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

func (ac *AttrCone) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(ac.String())
}

type AttrReflectanceRGB struct {
	RedReflect   float64 `json:"red_reflect"`
	GreenReflect float64 `json:"green_reflect"`
	BlueReflect  float64 `json:"blue_reflect"`
}

func ToAttrReflectanceRGB(attrs []AttrValue) ([]*AttrReflectanceRGB, error) {
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

func (ar *AttrReflectanceRGB) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(ar.String())
}

type AttrSpectrumRGB struct {
	RedSpectrum   float64 `json:"red_spectrum"`
	GreenSpectrum float64 `json:"green_spectrum"`
	BlueSpectrum  float64 `json:"blue_spectrum"`
}

func ToAttrSpectrumRGB(attrs []AttrValue) ([]*AttrSpectrumRGB, error) {
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

func (as *AttrSpectrumRGB) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(as.String())
}

type AttrComponentList []string

func ToAttrComponentList(attrs []AttrValue) ([]*AttrComponentList, error) {
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

func (acl *AttrComponentList) StringWrite(writer io.StringWriter) (int, error) {
	n := 0
	for i, ac := range *acl {
		if i != 0 {
			na, err := writer.WriteString(", ")
			if err != nil {
				return 0, err
			}
			n += na
		}
		na, err := writer.WriteString(ac)
		if err != nil {
			return 0, err
		}
		n += na
	}
	return n, nil
}

type AttrAttributeAlias struct {
	NewAlias    string `json:"new_alias"`
	CurrentName string `json:"current_name"`
}

func ToAttrAttributeAlias(attrs []AttrValue) ([]*AttrAttributeAlias, error) {
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

func (aaa *AttrAttributeAlias) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(aaa.String())
}

type AttrFormType int

func ToAttrFormType(attrs []AttrValue) ([]*AttrFormType, error) {
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

func (aft *AttrFormType) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(aft.String())
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

func ToAttrCvValue(attrs []AttrValue) ([]*AttrCvValue, error) {
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

func (acv *AttrCvValue) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(acv.String())
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

func ToAttrNurbsCurve(attrs []AttrValue) ([]*AttrNurbsCurve, error) {
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

func (anc *AttrNurbsCurve) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(anc.String())
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

func ToAttrNurbsSurface(attrs []AttrValue) ([]*AttrNurbsSurface, error) {
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

func (ans *AttrNurbsSurface) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(ans.String())
}

type AttrNurbsTrimface struct{}

func ToAttrNurbsTrimface(attrs []AttrValue) ([]*AttrNurbsTrimface, error) {
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

func (ant *AttrNurbsTrimface) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(ant.String())
}

type AttrFaceUV struct {
	UVSet  int   `json:"uv_set"`
	FaceUV []int `json:"face_uv"`
}

func ToAttrFaceUV(attrs []AttrValue) ([]*AttrFaceUV, error) {
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

func (af *AttrFaceUV) StringWrite(writer io.StringWriter) (int, error) {
	n, err := writer.WriteString("uvSet: ")
	if err != nil {
		return 0, err
	}
	na, err := writer.WriteString(strconv.Itoa(af.UVSet))
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString(", faceUV: ")
	if err != nil {
		return 0, err
	}
	n += na
	for i, uv := range af.FaceUV {
		if i != 0 {
			na, err = writer.WriteString(", ")
			if err != nil {
				return 0, err
			}
			n += na
		}
		na, err = writer.WriteString(strconv.Itoa(uv))
		if err != nil {
			return 0, err
		}
		n += na
	}
	return n, nil
}

type AttrMultiColor struct {
	ColorIndex int   `json:"color_index"`
	ColorIDs   []int `json:"color_ids"`
}

func ToAttrMultiColor(attrs []AttrValue) ([]*AttrMultiColor, error) {
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

func (amc *AttrMultiColor) StringWrite(writer io.StringWriter) (int, error) {
	n, err := writer.WriteString("colorIndex: ")
	if err != nil {
		return 0, err
	}
	na, err := writer.WriteString(strconv.Itoa(amc.ColorIndex))
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString(", colorIDs: ")
	if err != nil {
		return 0, err
	}
	n += na
	for i, id := range amc.ColorIDs {
		if i != 0 {
			na, err = writer.WriteString(", ")
			if err != nil {
				return 0, err
			}
			n += na
		}
		na, err = writer.WriteString(strconv.Itoa(id))
		if err != nil {
			return 0, err
		}
		n += na
	}
	return n, nil

}

type AttrPolyFaces struct {
	FaceEdge   []int            `json:"face_edge"`
	HoleEdge   []int            `json:"hole_edge"`
	FaceUV     []AttrFaceUV     `json:"face_uv"`
	FaceColor  []int            `json:"face_color"`
	MultiColor []AttrMultiColor `json:"multi_color"`
}

func ToAttrPolyFaces(attrs []AttrValue) ([]*AttrPolyFaces, error) {
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

func (apf *AttrPolyFaces) StringWrite(writer io.StringWriter) (int, error) {
	n, err := writer.WriteString("faceEdge: ")
	if err != nil {
		return 0, err
	}
	na, err := intArrayToStringWrite(apf.FaceEdge, writer)
	if err != nil {
		return 0, err
	}
	n += na

	na, err = writer.WriteString(", holeEdge: ")
	if err != nil {
		return 0, err
	}
	n += na
	na, err = intArrayToStringWrite(apf.HoleEdge, writer)
	if err != nil {
		return 0, err
	}
	n += na

	na, err = writer.WriteString(", faceUV: ")
	if err != nil {
		return 0, err
	}
	n += na
	for i, fuv := range apf.FaceUV {
		if i != 0 {
			na, err = writer.WriteString(", ")
			if err != nil {
				return 0, err
			}
			n += na
		}
		na, err = fuv.StringWrite(writer)
		if err != nil {
			return 0, err
		}
		n += na
	}

	na, err = writer.WriteString(", faceColor: ")
	if err != nil {
		return 0, err
	}
	n += na
	na, err = intArrayToStringWrite(apf.FaceColor, writer)
	if err != nil {
		return 0, err
	}
	n += na

	na, err = writer.WriteString(", multiColor: ")
	if err != nil {
		return 0, err
	}
	n += na
	for i, mc := range apf.MultiColor {
		if i != 0 {
			na, err = writer.WriteString(", ")
			if err != nil {
				return 0, err
			}
			n += na
		}
		na, err = mc.StringWrite(writer)
		if err != nil {
			return 0, err
		}
		n += na
	}

	return n, nil
}

type AttrDPCType int

func ToAttrDPCType(attrs []AttrValue) ([]*AttrDPCType, error) {
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

func (dpc *AttrDPCType) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(dpc.String())
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

func ToAttrDataPolyComponent(attrs []AttrValue) ([]*AttrDataPolyComponent, error) {
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
	var indices []int
	for k := range adpc.IndexValue {
		indices = append(indices, k)
	}
	sort.Ints(indices)
	var s []string
	for _, i := range indices {
		s = append(s, fmt.Sprintf("%d %f", i, adpc.IndexValue[i]))
	}
	return fmt.Sprintf("polyComponentType: %s, indexValue: %s",
		adpc.PolyComponentType.String(), strings.Join(s, ", "))
}

func (adpc *AttrDataPolyComponent) StringWrite(writer io.StringWriter) (int, error) {
	n, err := writer.WriteString("polyComponentType: ")
	if err != nil {
		return 0, err
	}
	na, err := writer.WriteString(adpc.PolyComponentType.String())
	if err != nil {
		return 0, err
	}
	n += na
	na, err = writer.WriteString(", indexValue: ")
	if err != nil {
		return 0, err
	}
	n += na
	var indices []int
	for k := range adpc.IndexValue {
		indices = append(indices, k)
	}
	sort.Ints(indices)
	for i, index := range indices {
		if i != 0 {
			na, err = writer.WriteString(", ")
			if err != nil {
				return 0, err
			}
			n += na
		}
		na, err = writer.WriteString(fmt.Sprintf("%d %f",
			index, adpc.IndexValue[index]))
		if err != nil {
			return 0, err
		}
		n += na
	}
	return n, nil
}

type ReferenceEditsCmdType string

const (
	RETypePArent         ReferenceEditsCmdType = "0"
	RETypeAddAttr        ReferenceEditsCmdType = "1"
	RETypeSetAttr        ReferenceEditsCmdType = "2"
	RETypeDisconnectAttr ReferenceEditsCmdType = "3"
	RETypeDeleteAttr     ReferenceEditsCmdType = "4"
	RETypeConnectAttr    ReferenceEditsCmdType = "5"
	RETypeRelationship   ReferenceEditsCmdType = "7"
	RETypeLock           ReferenceEditsCmdType = "8"
	RETypeUnlock         ReferenceEditsCmdType = "9"
)

func isNeedNewline(args ...string) string {
	count := 0
	for _, arg := range args {
		count += len(arg)
	}
	if count > 80 {
		return "\n\t\t"
	}
	return " "
}

type RECmdParent struct {
	NodeA     string `json:"node_a"`
	NodeB     string `json:"node_b"`
	Arguments string `json:"arguments"`
}

func (c *RECmdParent) String() string {
	needNewLine := isNeedNewline(c.NodeA, c.NodeB)
	return fmt.Sprintf("\t\t0 \"%s\"%s\"%s\" \"%s\"",
		c.NodeA, needNewLine, c.NodeB, c.Arguments)
}

type RECmdAddAttr struct {
	Node      string `json:"node"`
	LongAttr  string `json:"long_attr"`
	ShortAttr string `json:"short_attr"`
	Arguments string `json:"arguments"`
}

func (c *RECmdAddAttr) String() string {
	needNewline := isNeedNewline(c.Node, c.LongAttr, c.ShortAttr)
	return fmt.Sprintf("\t\t1 %s \"%s\" \"%s\"%s\"%s\"",
		c.Node, c.LongAttr, c.ShortAttr, needNewline, c.Arguments)
}

type RECmdSetAttr struct {
	Node      string `json:"node"`
	Attr      string `json:"attr"`
	Arguments string `json:"arguments"`
}

func (c *RECmdSetAttr) String() string {
	needNewLine := isNeedNewline(c.Node, c.Attr)
	return fmt.Sprintf("\t\t2 \"%s\" \"%s\"%s\"%s\"",
		c.Node, c.Attr, needNewLine, c.Arguments)
}

type RECmdDisconnectAttr struct {
	SourcePlug string `json:"source_plug"`
	DistPlug   string `json:"dist_plug"`
	Arguments  string `json:"arguments"`
}

func (c *RECmdDisconnectAttr) String() string {
	needNewline := isNeedNewline(c.SourcePlug, c.DistPlug)
	return fmt.Sprintf("\t\t3 \"%s\" \"%s\"%s\"%s\"",
		c.SourcePlug, c.DistPlug, needNewline, c.Arguments)
}

type RECmdDeleteAttr struct {
	Node      string `json:"node"`
	Attr      string `json:"attr"`
	Arguments string `json:"arguments"`
}

func (c *RECmdDeleteAttr) String() string {
	needNewLine := isNeedNewline(c.Node, c.Attr)
	return fmt.Sprintf("\t\t4 \"%s\" \"%s\"%s\"%s\"",
		c.Node, c.Attr, needNewLine, c.Arguments)
}

type RECmdConnectAttr struct {
	MagicNumber   int     `json:"magic_number"` // 0 or 3 or 4
	ReferenceNode string  `json:"reference_node"`
	SourcePlug    string  `json:"source_plug"`
	DistPlug      string  `json:"dist_plug"`
	SourcePHL     *string `json:"source_phl"` // only MagicNumber is 0
	DistPHL       *string `json:"dist_phl"`   // only MagicNumber is 0
	Arguments     string  `json:"arguments"`
}

func (c *RECmdConnectAttr) String() string {
	needNewLine := isNeedNewline(c.SourcePlug)
	phl := ""
	if c.SourcePHL != nil && c.DistPHL != nil {
		phlNewLine := ""
		if needNewLine == " " {
			phlNewLine = isNeedNewline(c.SourcePlug, c.DistPlug)
		}
		phl = fmt.Sprintf("%s\"%s\" \"%s\" ",
			phlNewLine, *c.SourcePHL, *c.DistPHL)
	}
	return fmt.Sprintf("\t\t5 %d \"%s\" \"%s\"%s\"%s\" %s\"%s\"",
		c.MagicNumber, c.ReferenceNode, c.SourcePlug, needNewLine,
		c.DistPlug, phl, c.Arguments)
}

type RECmdLock struct {
	Node string `json:"node"`
	Attr string `json:"attr"`
}

func (c *RECmdLock) String() string {
	needNewLine := isNeedNewline(c.Node, c.Attr)
	return fmt.Sprintf("\t\t8 \"%s\"%s\"%s\"",
		c.Node, needNewLine, c.Attr)
}

type RECmdRelationship struct {
	Type       string   `json:"type"`
	NodeName   string   `json:"node_name"`
	Commands   []string `json:"commands"`
	CommandNum int      `json:"command_num"`
}

func (c *RECmdRelationship) String() string {
	needNewLine := isNeedNewline(c.Type, c.NodeName)
	var buf bytes.Buffer
	for i, c := range c.Commands {
		if i == 0 {
			buf.WriteString("\"")
			buf.WriteString(c)
			buf.WriteString("\"")
		} else {
			buf.WriteString(" \"")
			buf.WriteString(c)
			buf.WriteString("\"")
		}
	}
	if buf.Len() > 0 {
		buf.WriteString(" ")
	}
	return fmt.Sprintf("\t\t7 \"%s\" \"%s\" %d%s%s0",
		c.Type, c.NodeName, len(c.Commands), needNewLine, buf.String())
}

type RECmdUnlock struct {
	Node string `json:"node"`
	Attr string `json:"attr"`
}

func (c *RECmdUnlock) String() string {
	needNewLine := isNeedNewline(c.Node, c.Attr)
	return fmt.Sprintf("\t\t9 \"%s\"%s\"%s\"",
		c.Node, needNewLine, c.Attr)
}

type ReferenceEdit struct {
	ReferenceNode   string                 `json:"reference_node"`
	CommandNum      int                    `json:"command_num"`
	Parents         []*RECmdParent         `json:"parents"`
	AddAttrs        []*RECmdAddAttr        `json:"add_attrs"`
	SetAttrs        []*RECmdSetAttr        `json:"set_attrs"`
	DisconnectAttrs []*RECmdDisconnectAttr `json:"disconnect_attrs"`
	DeleteAttrs     []*RECmdDeleteAttr     `json:"delete_attrs"`
	ConnectAttrs    []*RECmdConnectAttr    `json:"connect_attrs"`
	Relationships   []*RECmdRelationship   `json:"relationships"`
	Locks           []*RECmdLock           `json:"locks"`
	Unlocks         []*RECmdUnlock         `json:"unlocks"`
}

func (re *ReferenceEdit) String(buf *strings.Builder) {
	buf.WriteString(fmt.Sprintf("\t\t\"%s\" %d\n", re.ReferenceNode, re.CommandNum))
	for _, parent := range re.Parents {
		buf.WriteString(parent.String())
		buf.WriteRune('\n')
	}
	for _, addAttr := range re.AddAttrs {
		buf.WriteString(addAttr.String())
		buf.WriteRune('\n')
	}
	for _, setAttr := range re.SetAttrs {
		buf.WriteString(setAttr.String())
		buf.WriteRune('\n')
	}
	for _, disconnectAttr := range re.DisconnectAttrs {
		buf.WriteString(disconnectAttr.String())
		buf.WriteRune('\n')
	}
	for _, deleteAttr := range re.DeleteAttrs {
		buf.WriteString(deleteAttr.String())
		buf.WriteRune('\n')
	}
	for _, connectAttr := range re.ConnectAttrs {
		buf.WriteString(connectAttr.String())
		buf.WriteRune('\n')
	}
	for _, relationship := range re.Relationships {
		buf.WriteString(relationship.String())
		buf.WriteRune('\n')
	}
	for _, lock := range re.Locks {
		buf.WriteString(lock.String())
		buf.WriteRune('\n')
	}
	for _, unlock := range re.Unlocks {
		buf.WriteString(unlock.String())
		buf.WriteRune('\n')
	}
}

func (re *ReferenceEdit) StringWrite(writer io.StringWriter) (int, error) {
	n, err := writer.WriteString(fmt.Sprintf("\t\t\"%s\" %d\n", re.ReferenceNode, re.CommandNum))
	if err != nil {
		return 0, err
	}
	for _, parent := range re.Parents {
		na, err := writer.WriteString(parent.String())
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writer.WriteString("\n")
		if err != nil {
			return 0, err
		}
		n += na
	}
	for _, addAttr := range re.AddAttrs {
		na, err := writer.WriteString(addAttr.String())
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writer.WriteString("\n")
		if err != nil {
			return 0, err
		}
		n += na
	}
	for _, setAttr := range re.SetAttrs {
		na, err := writer.WriteString(setAttr.String())
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writer.WriteString("\n")
		if err != nil {
			return 0, err
		}
		n += na
	}
	for _, disconnectAttr := range re.DisconnectAttrs {
		na, err := writer.WriteString(disconnectAttr.String())
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writer.WriteString("\n")
		if err != nil {
			return 0, err
		}
		n += na
	}
	for _, deleteAttr := range re.DeleteAttrs {
		na, err := writer.WriteString(deleteAttr.String())
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writer.WriteString("\n")
		if err != nil {
			return 0, err
		}
		n += na
	}
	for _, connectAttr := range re.ConnectAttrs {
		na, err := writer.WriteString(connectAttr.String())
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writer.WriteString("\n")
		if err != nil {
			return 0, err
		}
		n += na
	}
	for _, relationship := range re.Relationships {
		na, err := writer.WriteString(relationship.String())
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writer.WriteString("\n")
		if err != nil {
			return 0, err
		}
		n += na
	}
	for _, lock := range re.Locks {
		na, err := writer.WriteString(lock.String())
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writer.WriteString("\n")
		if err != nil {
			return 0, err
		}
		n += na
	}
	for _, unlock := range re.Unlocks {
		na, err := writer.WriteString(unlock.String())
		if err != nil {
			return 0, err
		}
		n += na
		na, err = writer.WriteString("\n")
		if err != nil {
			return 0, err
		}
		n += na
	}
	return n, nil
}

type AttrDataReferenceEdits struct {
	TopReferenceNode string           `json:"top_reference_node"`
	ReferenceEdits   []*ReferenceEdit `json:"reference_edits"`
}

func ToAttrDataReferenceEdits(attrs []AttrValue) ([]*AttrDataReferenceEdits, error) {
	ret := make([]*AttrDataReferenceEdits, len(attrs))
	for i, a := range attrs {
		aa, ok := a.(*AttrDataReferenceEdits)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot cast %T", a))
		}
		ret[i] = aa
	}
	return ret, nil
}

func (re *AttrDataReferenceEdits) String() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("\n\t\t\"%s\"\n", re.TopReferenceNode))
	for _, edits := range re.ReferenceEdits {
		edits.String(&buf)
	}
	return buf.String()
}

func (re *AttrDataReferenceEdits) StringWrite(writer io.StringWriter) (int, error) {
	n, err := writer.WriteString(fmt.Sprintf("\n\t\t\"%s\"\n", re.TopReferenceNode))
	if err != nil {
		return 0, err
	}
	for _, edits := range re.ReferenceEdits {
		na, err := edits.StringWrite(writer)
		if err != nil {
			return 0, err
		}
		n += na
	}
	return n, nil
}

type AttrLatticePoint struct {
	S float64 `json:"s"`
	T float64 `json:"t"`
	U float64 `json:"u"`
}

func ToAttrLatticePoint(attrs []AttrValue) ([]*AttrLatticePoint, error) {
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

func (alp *AttrLatticePoint) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(alp.String())
}

type AttrLattice struct {
	DivisionS int                `json:"division_s"`
	DivisionT int                `json:"division_t"`
	DivisionU int                `json:"division_u"`
	Points    []AttrLatticePoint `json:"points"`
}

func ToAttrLattice(attrs []AttrValue) ([]*AttrLattice, error) {
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

func (al *AttrLattice) StringWrite(writer io.StringWriter) (int, error) {
	return writer.WriteString(al.String())
}

func intArrayToString(intArray []int) string {
	var s []string
	for _, i := range intArray {
		s = append(s, strconv.Itoa(i))
	}
	return strings.Join(s, ", ")
}

func intArrayToStringWrite(intArray []int, writer io.StringWriter) (int, error) {
	n := 0
	for i, ia := range intArray {
		if i != 0 {
			na, err := writer.WriteString(", ")
			if err != nil {
				return 0, err
			}
			n += na
		}
		na, err := writer.WriteString(strconv.Itoa(ia))
		if err != nil {
			return 0, err
		}
		n += na
	}
	return n, nil
}

func writeOnOff(buf *bytes.Buffer, value *bool) {
	if *value {
		buf.WriteString("on ")
	} else {
		buf.WriteString("off ")
	}
}

func writeOnOffBw(writer io.StringWriter, value *bool) (n int, err error) {
	if *value {
		n, err = writer.WriteString("on ")
	} else {
		n, err = writer.WriteString("off ")
	}
	return
}
