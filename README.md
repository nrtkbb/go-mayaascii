# go-mayaascii

[![Build Status](https://travis-ci.org/nrtkbb/go-mayaascii.svg?branch=master)](https://travis-ci.org/nrtkbb/go-mayaascii)


# For example.

```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	ma "github.com/nrtkbb/go-mayaascii"
)

func main() {
	// Read file Maya Ascii.
	fp, err := os.Open("basic.ma")
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	reader := bufio.NewReader(fp)

	// Unmarshal Maya Ascii File to 'mo' (Maya Ascii Object).
	mo, err := ma.Unmarshal(reader)
	if err != nil {
		log.Fatal(err)
	}

	// requires command parsed data.
	for _, r := range mo.Requires {
		fmt.Printf("%s version is %s. %d nodeTypes, %d dataTypes, %d Plugin's nodes.\n",
			r.Name, r.Version, len(r.NodeTypes), len(r.DataTypes), len(r.Nodes))
	}

	// Print all nodes.
	for _, n := range mo.Nodes {
		fmt.Printf("%d : %s\n", n.LineNo, n.Name)
	}

	// Specify node name.
	persp, err := mo.GetNode("persp")
	if err != nil {
		log.Fatal(err)
	}

	// Get attribute (Must be short name) and cast to string.
	ow, err := persp.Attr(".ow").asString()  // or .asInt() or .asFloat() etc..
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s.ow is %s", persp.Name, ow)

	// Print Node's all attrs.
	for _, a := range persp.Attrs {
		fmt.Printf("%s%s is %d type is %s\n", persp.Name, a.Name, len(a.Values), a.Type)
	}

	// Print Node's all children.
	for _, c := range persp.Children {
		fmt.Printf("%s child is %s\n", persp.Name, c.Name)
	}

	// Print Node's parent node.
	if persp.Parent != nil {
		fmt.Printf("%s parent is %s\n", persp.Name, persp.Parent.Name)
	}

	// Get nodes by nodeType.
	transforms, err := mo.GetNodes("transform") // Specify node type.
	if err != nil {
		log.Fatal(err)
	}

	// Print transform nodes.
	for _, t := range transforms {
		fmt.Printf("%d : %s\n", t.LineNo, t.Name)
	}

	// Get specified source connection nodes.
	srcNodes, err := persp.Src(&ma.ConnectInfo{
		Name: "topShape",  // (Optional)
		Attr: ".ow",       // (Optional) Must be specified with a short name
		Type: "camera",    // (Optional)
	})

	// Print src nodes.
	for _, t := range srcNodes {
		fmt.Printf("%d : %s\n", t.LineNo, t.Name)
	}

	// Get all destination connection nodes.
	dstNodes, err := persp.Dst(nil) // nil will return all connection nodes.

	// Print dst nodes.
	for _, t := range dstNodes {
		fmt.Printf("%d : %s\n", t.LineNo, t.Name)
	}
}
```


## TODO

- [x] Get Requires
- [x] Get Parent Node
- [x] Get Children Node
- [x] Get Node
- [x] Get Nodes
- [x] Get Attr
- [x] Get Attrs
- [ ] Get Src Connection
- [ ] Get Src Connections
- [ ] Get History Connections
- [ ] Get Dst Connection
- [ ] Get Dst Connections
- [ ] Get Future Connections
- [ ] Get Default Node
- [ ] Get Default Node Attr
- [ ] Get FileInfo
- [ ] Get currentUnit
- [ ] Get LockNode
- [ ] Get Relationship
- [ ] Remove Require
- [ ] Add Require
- [ ] Remove Node
- [ ] Add node
- [ ] Remove AddAttr
- [ ] Add AddAttr
- [ ] Remove SetAttr
- [ ] Add SetAttr
- [ ] Remove Connection
- [ ] Add Connection
- [ ] Save As

done 7 / 30
