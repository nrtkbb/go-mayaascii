package mayaascii

import (
	"github.com/nrtkbb/go-mayaascii/cmd"
	"github.com/nrtkbb/go-mayaascii/connection"
	"io"
)

func Unmarshal(reader io.Reader) (*Object, error) {
	mo := &Object{
		Files: []*File{},
		Requires: []*Require{},
		Nodes: map[string]*Node{},

		cmds: []*cmd.Cmd{},
		connections: connection.NewConnections(),
	}
	err := mo.Unmarshal(reader)
	if err != nil {
		return nil, err
	}

	return mo, nil
}

func UnmarshalFocus(reader io.Reader, focusCommands []cmd.Type) (*Object, error) {
	mo := &Object{
		Files: []*File{},
		Requires: []*Require{},
		Nodes: map[string]*Node{},

		cmds: []*cmd.Cmd{},
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
