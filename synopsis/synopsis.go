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
	defer func() {
		err := fp.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	reader := bufio.NewReader(fp)

	// Unmarshal Maya Ascii FileCmd to 'mo' (Maya Ascii Object).
	mo, err := ma.Unmarshal(reader)
	if err != nil {
		log.Fatal(err)
	}

	// requires command parsed data.
	for _, r := range mo.Requires {
		fmt.Printf("%s version is %s. %d nodeTypes, %d dataTypes, %d Plugin's nodes.\n",
			r.GetPluginName(),
			r.GetVersion(),
			len(r.GetNodeTypes()),
			len(r.GetDataTypes()),
			len(r.Nodes))
	}

	// Print all nodes.
	for _, n := range mo.Nodes {
		fmt.Printf("%s\n", n.GetName())
	}

	// Specify node name.
	persp, err := mo.GetNode("perspShape")
	if err != nil {
		log.Fatal(err)
	}

	// Get attribute (Must be short name) and cast to string.
	ow, err := persp.GetAttr(".imn").String() // or .Int() or .Float() etc..
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s.t is %s", persp.GetName(), ow)

	// Print Node's all attrs.
	for _, a := range persp.Attrs {
		fmt.Printf("%s%s is %d type is %s\n",
			persp.GetName(),
			a.GetName(),
			len(a.GetAttrValue()),
			a.GetAttrType())
	}

	// Print Node's all children.
	for _, c := range persp.Children {
		fmt.Printf("%s child is %s\n", persp.GetName(), c.GetName())
	}

	// Print Node's parent node.
	if persp.Parent != nil {
		fmt.Printf("%s parent is %s\n", persp.GetName(), persp.Parent.GetName())
	}

	// Get nodes by nodeType.
	transforms, err := mo.GetNodes("transform") // Specify node type.
	if err != nil {
		log.Fatal(err)
	}

	// Print transform nodes.
	for _, t := range transforms {
		fmt.Printf("%s\n", t.GetName())
	}

	// Get specified source connection nodes.
	srcNodes := persp.ListConnections(&ma.ConnectionArgs{
		Source: true,
		Type: "camera",
		AttrName: ".ow", // Specific persp's attr name.
	})

	// Print src nodes.
	for _, t := range srcNodes {
		fmt.Printf("%s\n", t.GetName())
	}

	// Get all destination connection nodes.
	dstNodes := persp.ListConnections(&ma.ConnectionArgs{
		Destination: true,
	})

	// Print dst nodes.
	for _, t := range dstNodes {
		fmt.Printf("%s\n", t.GetName())
	}
}
