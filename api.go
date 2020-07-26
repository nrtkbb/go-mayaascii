package mayaascii

import (
	"io"
)

func Unmarshal(reader io.Reader) (*Object, error) {
	mo := &Object{}
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
