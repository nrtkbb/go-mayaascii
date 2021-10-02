package mayaascii

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/nrtkbb/bufscan"
)

// Object ...
type Object struct {
	Files         []*File
	FileInfos     []*FileInfo
	Requires      []*Require
	Nodes         map[string]*Node
	Selects       []*Select
	LineComments  []*LineComment
	BlockComments []*BlockComment

	cmds        []*Cmd
	connections Connections
}

func (o *Object) Unmarshal(reader io.Reader) error {
	br := bufio.NewReader(reader)

	cmdBuilder := &CmdBuilder{}
	lineCommentBuilder := &CmdBuilder{}
	err := bufscan.BufScan(br, func(line string) error {
		// Normally, it should be parsed with a proper
		// tokenizer and lexer, but here it is simply
		// treated as one token per line.
		if LineCommentType.HasPrefix(line) {
			lineCommentBuilder.Append(line)
			c := lineCommentBuilder.Parse()
			o.cmds = append(o.cmds, c)
			lineCommentBuilder.Clear()
			cmdBuilder.lineNo++
			return nil
		}
		cmdBuilder.Append(line)
		lineCommentBuilder.lineNo++
		if cmdBuilder.IsCmdEOF() {
			c := cmdBuilder.Parse()
			o.cmds = append(o.cmds, c)
			cmdBuilder.Clear()
		}
		return nil
	})
	if err != nil {
		return err
	}

	p := New(o.cmds)
	p.o = o
	p.ParseCmds()
	if !p.CheckErrors() {
		return nil
	}

	return nil
}

func (o *Object) UnmarshalFocus(reader io.Reader, focusCommands CommandTypes) error {
	if 0 == len(focusCommands) {
		return errors.New("focusCommands must one type")
	}

	br := bufio.NewReader(reader)

	cmdBuilder := &CmdBuilder{}
	lineCommentBuilder := &CmdBuilder{}
	isFocus := false
	err := bufscan.BufScan(br, func(line string) error {
		if LineCommentType.HasPrefix(line) {
			lcc := "//"
			if focusCommands.InHasPrefix(&lcc) {
				lineCommentBuilder.Append(line)
				c := lineCommentBuilder.Parse()
				o.cmds = append(o.cmds, c)
				lineCommentBuilder.Clear()
			}
			cmdBuilder.lineNo++
			return nil
		}
		if cmdBuilder.IsClear() {
			isFocus = focusCommands.InHasPrefix(&line)
		}
		if !isFocus {
			cmdBuilder.lineNo++
			return nil
		}
		cmdBuilder.Append(line)
		lineCommentBuilder.lineNo++
		if cmdBuilder.IsCmdEOF() {
			c := cmdBuilder.Parse()
			o.cmds = append(o.cmds, c)
			cmdBuilder.Clear()
		}
		return nil
	})
	if err != nil {
		return err
	}

	p := New(o.cmds)
	p.o = o
	p.ParseCmds()
	if !p.CheckErrors() {
		return nil
	}

	return nil
}

func (o *Object) GetNode(n string) (*Node, error) {
	node, ok := o.Nodes[n]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s node was not found", n))
	}
	return node, nil
}

func (o *Object) GetNodes(nodeType string) ([]*Node, error) {
	var results []*Node
	for _, node := range o.Nodes {
		if node.Type == nodeType {
			results = append(results, node)
		}
	}
	if len(results) == 0 {
		return nil, errors.New(fmt.Sprintf("%s type was not found", nodeType))
	}
	return results, nil
}

type LineComment struct {
	lineCommentCmd *LineCommentCmd
}

func (lc LineComment) GetComment() string {
	return lc.lineCommentCmd.Comment
}

type BlockComment struct {
	blockCommentCmd *BlockCommentCmd
}

func (bc BlockComment) GetComment() string {
	return bc.blockCommentCmd.Comment
}

type File struct {
	Parent   *File
	Children []*File

	fileCmd *FileCmd
}

func (f File) GetLineNo() uint {
	return f.fileCmd.LineNo
}

func (f File) GetPath() string {
	return f.fileCmd.Path
}

func (f File) GetNamespace() string {
	return f.fileCmd.Namespace
}

type FileInfo struct {
	fileInfoCmd *FileInfoCmd
}

func (fi *FileInfo) GetName() string {
	return fi.fileInfoCmd.Name
}

func (fi FileInfo) GetValue() string {
	return fi.fileInfoCmd.Value
}

type Require struct {
	Nodes []*Node
	Data  []*Node // TODO: 何もセットしてない

	requireCmd *RequiresCmd
}

func (r Require) GetPluginName() string {
	return r.requireCmd.PluginName
}

func (r Require) GetVersion() string {
	return r.requireCmd.Version
}

func (r Require) GetNodeTypes() []string {
	return r.requireCmd.NodeTypes
}

func (r Require) GetDataTypes() []string {
	return r.requireCmd.DataTypes
}

type Node struct {
	object     *Object
	LineNo     uint
	Type       string
	Name       string
	Attrs      []*Attr
	Shared     bool
	SkipSelect bool
	Parent     *Node
	Children   []*Node

	isDeleted     bool
	createNodeCmd *CreateNodeCmd
	renameCmd     *RenameCmd
}

func (n *Node) GetType() string {
	return n.createNodeCmd.NodeType
}

func (n *Node) GetName() string {
	return n.createNodeCmd.NodeName
}

func (n *Node) GetAttr(name string) *Attr {
	for _, a := range n.Attrs {
		if a.GetName() == name {
			return a
		}
	}
	return nil // not found.
}

func (n *Node) GetUUID() (string, error) {
	if n.renameCmd.UUID && n.renameCmd.To != nil {
		return *n.renameCmd.To, nil
	}
	return "", errors.New(fmt.Sprintf("%s has not UUID", n.GetName()))
}

func (n *Node) Remove() error {
	if n.isDeleted {
		return errors.New(fmt.Sprintf("%s was already deleted", n.Name))
	}
	for _, c := range n.Children {
		err := c.Remove()
		if err != nil {
			return err
		}
	}
	for _, a := range n.Attrs {
		err := a.Remove()
		if err != nil {
			return err
		}
	}
	n.isDeleted = true
	return nil
}

type Select struct {
	Name string

	Select   *SelectCmd
	SetAttrs []*SetAttrCmd
	AddAttrs []*AddAttrCmd
}

type ConnectionArgs struct {
	Source      bool
	Destination bool
	Type        string
	AttrName    string
}

func (n *Node) ListConnections(ca *ConnectionArgs) []*Node {
	if ca == nil {
		// default settings.
		ca = &ConnectionArgs{
			Source:      true,
			Destination: true,
			Type:        "",
			AttrName:    "",
		}
	}
	var names []string
	if ca.Source {
		if ca.AttrName != "" {
			names = n.object.connections.GetSrcNamesAttr(n.Name, ca.AttrName)
		} else {
			names = n.object.connections.GetSrcNames(n.Name)
		}
	}
	if ca.Destination {
		if ca.AttrName != "" {
			names = n.object.connections.GetDstNamesAttr(n.Name, ca.AttrName)
		} else {
			names = n.object.connections.GetDstNames(n.Name)
		}
	}
	var nodes []*Node
	for _, name := range names {
		srcNode, ok := n.object.Nodes[name]
		if !ok {
			// name is maybe default node.
			defaultNode := &Node{
				object:        n.object,
				LineNo:        0,
				Type:          "Unknown",
				Name:          name,
				Attrs:         nil,
				Shared:        true,
				SkipSelect:    false,
				Parent:        nil,
				Children:      nil,
				isDeleted:     false,
				createNodeCmd: nil,
				renameCmd:     nil,
			}
			nodes = append(nodes, defaultNode)
		} else {
			nodes = append(nodes, srcNode)
		}
	}
	if ca.Type == "" {
		return nodes
	}
	var typeFiltered []*Node
	for _, node := range nodes {
		if node.Type == ca.Type {
			typeFiltered = append(typeFiltered, node)
		}
	}
	return typeFiltered
}

type Attr struct {
	Node *Node

	attrCmd   AttrCmd
	err       error
	isDeleted bool
}

func (a *Attr) GetName() string {
	return a.attrCmd.GetName()
}

func (a *Attr) IsChannelBox() bool {
	return a.attrCmd.IsChannelBox()
}

func (a *Attr) IsKeyable() bool {
	return a.attrCmd.IsKeyable()
}

func (a *Attr) GetAttrType() AttrType {
	return a.attrCmd.GetAttrType()
}

func (a *Attr) GetAttrValue() []AttrValue {
	return a.attrCmd.GetAttrValue()
}

func (a *Attr) Remove() error {
	if a.isDeleted {
		return errors.New(fmt.Sprintf("%s.%s was already deleted",
			a.Node.Name, a.GetName()))
	}
	a.isDeleted = true
	return nil
}

func (a *Attr) String() (string, error) {
	return a.GetName(), a.err
}

type Parser struct {
	o    *Object
	errs []string
	cmds []*Cmd
	cur  int

	CurCmd  *Cmd
	PeekCmd *Cmd
}

func New(cmds []*Cmd) *Parser {
	p := &Parser{
		cmds: cmds,
		cur:  -1,
	}
	p.NextCmd()
	p.NextCmd()
	return p
}

func (p *Parser) ParseCmds() {
	for p.CurCmd != nil {
		var err error
		switch p.CurCmd.Type {
		case LineCommentType:
			err = p.parseLineComments()
		case BlockCommentType:
			err = p.parseBlockComments()
		case FileType:
			err = p.parseFiles()
		case FileInfoType:
			err = p.parseFileInfos()
		case RequiresType:
			err = p.parseRequires()
		case CreateNodeType:
			err = p.parseCreateNode()
		case ConnectAttrType:
			err = p.parseConnectAttr()
		case SelectType:
			err = p.parseSelect()
		default:
			err = errors.New(fmt.Sprintf(
				"%s type parser is not found. %v", p.CurCmd.Type, *p.CurCmd))
		}
		if err != nil {
			p.errs = append(p.errs, err.Error())
		}
		p.NextCmd()
	}
}

func (p *Parser) CheckErrors() bool {
	// if 0 < len(p.errs) {
	//	for _, e := range p.errs {
	//		log.Println(e)
	//	}
	//	return false
	// }
	return true
}

func (p *Parser) parseLineComments() error {
	lc := ParseLineComment(p.CurCmd)
	lineComment := &LineComment{
		lineCommentCmd: lc,
	}
	p.o.LineComments = append(p.o.LineComments, lineComment)
	return nil
}

func (p *Parser) parseBlockComments() error {
	bc := ParseBlockComment(p.CurCmd)
	blockComment := &BlockComment{
		blockCommentCmd: bc,
	}
	p.o.BlockComments = append(p.o.BlockComments, blockComment)
	return nil
}

func (p *Parser) parseFiles() error {
	f := ParseFile(p.CurCmd)
	file := &File{
		Parent:   nil,
		Children: nil,
		fileCmd:  f,
	}
	p.o.Files = append(p.o.Files, file)

	if len(p.o.Files) == 1 {
		return nil
	}

	if f.ReferenceDepthInfo <= 1 {
		return nil
	}

	for i := len(p.o.Files) - 2; i >= 0; i-- {
		prevFile := p.o.Files[i]
		if prevFile.fileCmd.ReferenceDepthInfo < f.ReferenceDepthInfo {
			if prevFile.Children == nil {
				prevFile.Children = []*File{}
			}
			prevFile.Children = append(prevFile.Children, file)
			file.Parent = prevFile
			break
		}
	}
	return nil
}

func (p *Parser) parseFileInfos() error {
	fi := ParseFileInfo(p.CurCmd)
	fileInfo := &FileInfo{
		fileInfoCmd: fi,
	}
	p.o.FileInfos = append(p.o.FileInfos, fileInfo)
	return nil
}

func (p *Parser) parseRequires() error {
	rq := ParseRequires(p.CurCmd)
	requires := &Require{
		Nodes:      []*Node{},
		Data:       []*Node{},
		requireCmd: rq,
	}
	p.o.Requires = append(p.o.Requires, requires)
	return nil
}

func (p *Parser) parseCreateNode() error {
	cn := ParseCreateNode(p.CurCmd)
	node := &Node{
		object:        p.o,
		Type:          cn.NodeType,
		Name:          cn.NodeName,
		Shared:        cn.Shared,
		SkipSelect:    cn.SkipSelect,
		createNodeCmd: cn,
		LineNo:        cn.LineNo,
	}
	if _, ok := p.o.Nodes[node.Name]; ok {
		return errors.New(fmt.Sprintf("Already found node ... %s", node.Name))
	}
	p.o.Nodes[node.Name] = node

	if cn.Parent != nil {
		// reverse loop.
		parentNode, ok := p.o.Nodes[*cn.Parent]
		if !ok {
			return errors.New(fmt.Sprintf("Not found parent %s. node is %s",
				*cn.Parent, node.Name))
		}
		node.Parent = parentNode
		parentNode.Children = append(parentNode.Children, node)
	}

	if p.PeekCmdIs(RenameType) {
		p.NextCmd()
		node.renameCmd = ParseRename(p.CurCmd)
	}

	for p.PeekCmdIs(AddAttrType) {
		p.NextCmd()
		ad := ParseAddAttr(p.CurCmd)
		a := &Attr{
			Node:    node,
			attrCmd: ad,
		}
		node.Attrs = append(node.Attrs, a)
	}

	var setAttrCmds []*SetAttrCmd
	for p.PeekCmdIs(SetAttrType) {
		p.NextCmd()
		var sa *SetAttrCmd
		var err error
		if len(setAttrCmds) == 0 {
			sa, err = ParseSetAttr(p.CurCmd, nil)
		} else {
			sa, err = ParseSetAttr(p.CurCmd, setAttrCmds[len(setAttrCmds)-1])
		}
		if err != nil {
			return err
		}
		setAttrCmds = append(setAttrCmds, sa)
		a := &Attr{
			Node:    node,
			attrCmd: sa,
		}
		node.Attrs = append(node.Attrs, a)
	}

	isPluginsNode := false
	if len(p.o.Requires) != 0 {
		for _, r := range p.o.Requires {
			for _, nt := range r.GetNodeTypes() {
				if node.Type == nt {
					r.Nodes = append(r.Nodes, node)
					isPluginsNode = true
					break
				}
			}
			if isPluginsNode {
				break
			}
		}
	}

	return nil
}

func (p *Parser) parseConnectAttr() error {
	ca, err := ParseConnectAttr(p.CurCmd)
	if err != nil {
		return err
	}
	p.o.connections.Append(ca)
	return nil
}

func (p *Parser) parseSelect() error {
	s := ParseSelect(p.CurCmd)
	if len(s.Names) > 1 {
		return errors.New(fmt.Sprintf("un-support bulk select. [%s], %v", strings.Join(s.Names, ", "), *s))
	} else if len(s.Names) == 0 {
		return errors.New(fmt.Sprintf("un-support zero select. %v", *s))
	}
	sel := &Select{
		Name: s.Names[0],

		Select:   s,
		SetAttrs: []*SetAttrCmd{},
		AddAttrs: []*AddAttrCmd{},
	}
	p.o.Selects = append(p.o.Selects, sel)

	for p.PeekCmdIs(AddAttrType) {
		p.NextCmd()
		ad := ParseAddAttr(p.CurCmd)
		sel.AddAttrs = append(sel.AddAttrs, ad)
	}

	for p.PeekCmdIs(SetAttrType) {
		p.NextCmd()
		var at *SetAttrCmd
		var err error
		if len(sel.SetAttrs) == 0 {
			at, err = ParseSetAttr(p.CurCmd, nil)
		} else {
			at, err = ParseSetAttr(p.CurCmd, sel.SetAttrs[len(sel.SetAttrs)-1])
		}
		if err != nil {
			return err
		}
		sel.SetAttrs = append(sel.SetAttrs, at)
	}

	return nil
}

func (p *Parser) NextCmd() {
	p.CurCmd = p.PeekCmd
	p.cur++
	if p.cur < len(p.cmds) {
		p.PeekCmd = p.cmds[p.cur]
	} else {
		p.PeekCmd = nil
	}
}

func (p *Parser) CurCmdIs(t Type) bool {
	return p.CurCmd.Type == t
}

func (p *Parser) PeekCmdIs(t Type) bool {
	if p.PeekCmd == nil {
		return false
	}
	return p.PeekCmd.Type == t
}
