package mayaascii

import (
	"github.com/nrtkbb/go-mayaascii/cmd"
	"github.com/nrtkbb/go-mayaascii/connection"
	"io"
	"strings"
)

const (
	FileCommand        = "file "
	FileInfoCommand    = "fileInfo "
	WorkspaceCommand   = "workspace "
	RequiresCommand    = "requires "
	ConnectAttrCommand = "connectAttr "
	CreateNodeCommand  = "createNode "
	RenameCommand      = "rename "
	SetAttrCommand     = "setAttr "
	AddAttrCommand     = "addAttr "
	SelectCommand      = "select "
)

type CommandTypes []string

func (ct *CommandTypes) InHasPrefix(line *string) bool {
	for i := 0; i < len(*ct); i ++ {
		if strings.HasPrefix(*line, (*ct)[i]) {
			return true
		}
	}
	return false
}

func Unmarshal(reader io.Reader) (*Object, error) {
	mo := &Object{
		Files:     []*File{},
		FileInfos: []*FileInfo{},
		Requires:  []*Require{},
		Nodes:     map[string]*Node{},

		cmds:        []*cmd.Cmd{},
		connections: connection.NewConnections(),
	}
	err := mo.Unmarshal(reader)
	if err != nil {
		return nil, err
	}

	return mo, nil
}

func UnmarshalFocus(reader io.Reader, focusCommands CommandTypes) (*Object, error) {
	mo := &Object{
		Files:     []*File{},
		FileInfos: []*FileInfo{},
		Requires:  []*Require{},
		Nodes:     map[string]*Node{},

		cmds:        []*cmd.Cmd{},
		connections: connection.NewConnections(),
	}
	err := mo.UnmarshalFocus(reader, focusCommands)
	if err != nil {
		return nil, err
	}

	return mo, nil
}

type ConnectInfo struct {
	Name string
	Attr string
	Type string
}
