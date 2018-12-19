package object

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/nrtkbb/bufscan"
	"github.com/nrtkbb/go-mayaascii/cmd"
	"github.com/nrtkbb/go-mayaascii/connection"
	"github.com/nrtkbb/go-mayaascii/parser"
)

// Object ...
type Object struct {
	Requires []*Require
	Nodes    []*Node

	cmds []*cmd.Cmd
	Cons connection.Connections
}

func (o *Object) Unmarshal(reader io.Reader) error {
	br := bufio.NewReader(reader)

	var cmds []*cmd.Cmd
	cmdBuilder := &cmd.CmdBuilder{}
	err := bufscan.BufScan(br, func(line string) error {
		cmdBuilder.Append(line)
		if cmdBuilder.IsCmdEOF() {
			c := cmdBuilder.Parse()
			cmds = append(cmds, c)
			cmdBuilder.Clear()
		}
		return nil
	})
	if err != nil {
		return err
	}

	p := New(cmds)
	p.ParseCmds(o)
	if !p.CheckErrors() {
		return nil
	}

	return nil
}

func (o *Object) GetNode(n string) (*Node, error) {
	for _, node := range o.Nodes {
		if node.Name == n {
			return node, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("%s node was not found", n))
}

func (o *Object) GetNodes(tn string) ([]*Node, error) {
	var results []*Node
	for _, node := range o.Nodes {
		if node.Type == tn {
			results = append(results, node)
		}
	}
	if len(results) == 0 {
		return nil, errors.New(fmt.Sprintf("%s type was not found", tn))
	}
	return results, nil
}

type Require struct {
	LineNo    uint
	Name      string
	Version   string
	NodeTypes []string
	DataTypes []string
	Nodes     []*Node
	Datas     []*Node // TODO: 何もセットしてない
}

type Node struct {
	LineNo     uint
	Type       string
	Name       string
	Attrs      []*Attr
	Shared     bool
	SkipSelect bool
	Parent     *Node
	Children   []*Node

	isDeleted bool
	CN        *cmd.CreateNode
	RN        *cmd.Rename
	AT        []*cmd.SetAttr
	AD        []*cmd.AddAttr
}

type ConnectInfo interface {
	GetName() string
	GetAttr() string
	GetType() string
}

func (n *Node) Src(ci ConnectInfo) ([]*Node, error) {
	return nil, nil
}

func (n *Node) Dst(ci ConnectInfo) ([]*Node, error) {
	return nil, nil
}

func (n *Node) Attr(name string) *Attr {
	for _, a := range n.Attrs {
		if a.Name == name {
			return a
		}
	}

	// not found.
	a := &Attr{err: errors.New(
		fmt.Sprintf("%s attr is not found", name))}
	return a
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

type Attr struct {
	LineNo uint
	Name   string
	Node   *Node
	Values []cmd.Attr
	Type   cmd.AttrType

	SA        *cmd.SetAttr
	err       error
	isDeleted bool
}

func (a *Attr) Remove() error {
	if a.isDeleted {
		return errors.New(fmt.Sprintf("%s.%s was already deleted",
			a.Node.Name, a.Name))
	}
	a.isDeleted = true
	return nil
}

func (a *Attr) String() (string, error) {
	return a.Name, a.err
}

type Parser struct {
	o    *Object
	errs []string
	cmds []*cmd.Cmd
	cur  int

	CurCmd  *cmd.Cmd
	PeekCmd *cmd.Cmd
}

func New(cmds []*cmd.Cmd) *Parser {
	p := &Parser{
		cmds: cmds,
		cur:  0,
	}
	p.NextCmd()
	p.NextCmd()
	return p
}

func (p *Parser) ParseCmds(o *Object) {
	p.o = o
	var err error
	for p.CurCmd != nil {
		switch p.CurCmd.Type {
		case cmd.REQUIRES:
			err = p.parseRequires()
		case cmd.CREATENODE:
			err = p.parseCreateNode()
		case cmd.CONNECTATTR:
			err = p.parseConnectAttr()
		default:
			err = errors.New(fmt.Sprintf(
				"%s type parser is not found", p.CurCmd.Type))
		}
		if err != nil {
			p.errs = append(p.errs, err.Error())
			//return err
		}
		p.NextCmd()
	}
}

func (p *Parser) CheckErrors() bool {
	if 0 < len(p.errs) {
		for _, e := range p.errs {
			log.Println(e)
		}
		return false
	}
	return true
}

func (p *Parser) parseRequires() error {
	rq := parser.MakeRequires(p.CurCmd)
	requires := &Require{
		Name:      rq.PluginName,
		Version:   rq.Version,
		NodeTypes: rq.NodeTypes,
		DataTypes: rq.DataTypes,
		LineNo:    rq.LineNo,
	}
	p.o.Requires = append(p.o.Requires, requires)
	return nil
}

func (p *Parser) parseCreateNode() error {
	cn := parser.MakeCreateNode(p.CurCmd)
	node := &Node{
		Type:       cn.NodeType,
		Name:       cn.NodeName,
		Shared:     cn.Shared,
		SkipSelect: cn.SkipSelect,
		CN:         cn,
		LineNo:     cn.LineNo,
	}
	p.o.Nodes = append(p.o.Nodes, node)

	if cn.Parent != nil {
		// reverse loop.
		for i := len(p.o.Nodes) - 1; i >= 0; i-- {
			n := p.o.Nodes[i]
			if n.Name == *cn.Parent {
				node.Parent = n
				n.Children = append(n.Children, node)
				break
			}
		}
	}

	if p.PeekCmdIs(cmd.RENAME) {
		p.NextCmd()
		node.RN = parser.MakeRename(p.CurCmd)
	}

	for p.PeekCmdIs(cmd.ADDATTR) {
		p.NextCmd()
		ad := parser.MakeAddAttr(p.CurCmd)
		node.AD = append(node.AD, ad)
	}

	for p.PeekCmdIs(cmd.SETATTR) {
		p.NextCmd()
		var at *cmd.SetAttr
		var err error
		if len(node.AT) == 0 {
			at, err = parser.MakeSetAttr(p.CurCmd, nil)
		} else {
			at, err = parser.MakeSetAttr(p.CurCmd, node.AT[len(node.AT)-1])
		}
		if err != nil {
			return err
		}
		node.AT = append(node.AT, at)
		a := &Attr{
			Name:   at.AttrName,
			Node:   node,
			Values: at.Attr,
			Type:   at.AttrType,
			SA:     at,
			LineNo: at.LineNo,
		}
		node.Attrs = append(node.Attrs, a)
	}

	isPluginsNode := false
	if len(p.o.Requires) != 0 {
		for _, r := range p.o.Requires {
			for _, nt := range r.NodeTypes {
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
	ca, err := parser.MakeConnectAttr(p.CurCmd)
	if err != nil {
		return err
	}
	p.o.Cons.Append(ca)
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

func (p *Parser) CurCmdIs(t cmd.Type) bool {
	return p.CurCmd.Type == t
}

func (p *Parser) PeekCmdIs(t cmd.Type) bool {
	return p.PeekCmd.Type == t
}
