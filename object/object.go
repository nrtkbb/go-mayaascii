package object

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/nrtkbb/bufscan"
	"github.com/nrtkbb/go-mayaascii/cmd"
	mac "github.com/nrtkbb/go-mayaascii/connection"
	"github.com/nrtkbb/go-mayaascii/parser"
)

// Object ...
type Object struct {
	Requires []*Require
	Nodes    []*Node

	cmds []*cmd.Cmd
	Cons mac.Connections
}

func (o *Object) Unmarshal(reader io.Reader) error {
	br := bufio.NewReader(reader)

	var cmds []*cmd.Cmd
	cmdBuilder := &cmd.CmdBuilder{}
	err := bufscan.BufScan(br, func(line string) error {
		cmdBuilder.Append(line)
		if cmdBuilder.IsCmdEOF() {
			cmd := cmdBuilder.Parse()
			cmds = append(cmds, cmd)
			cmdBuilder.Clear()
		}
		return nil
	})
	if err != nil {
		return err
	}

	p := parser.New(cmds)
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
