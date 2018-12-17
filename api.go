package mayaascii

import (
	"io"

	"github.com/nrtkbb/go-mayaascii/object"
)

func Unmarshal(reader io.Reader) (*object.Object, error) {
	mo := &object.Object{}
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

func (ci *ConnectInfo) GetName() string {
	return ci.Name
}

func (ci *ConnectInfo) GetAttr() string {
	return ci.Attr
}

func (ci *ConnectInfo) GetType() string {
	return ci.Type
}
