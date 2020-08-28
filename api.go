package mayaascii

import (
	"github.com/nrtkbb/go-mayaascii/cmd"
	"github.com/nrtkbb/go-mayaascii/connection"
	"io"
)

func Unmarshal(reader io.Reader) (*Object, error) {
	mo := &Object{
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

type ConnectInfo struct {
	Name string
	Attr string
	Type string
}
