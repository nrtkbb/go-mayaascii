package main

import (
	"bufio"
	ma "github.com/nrtkbb/go-mayaascii"
	"log"
	"os"
	"testing"
)

func TestSynopsis(t *testing.T) {
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
		FmtPrintf("%s version is %s. %d nodeTypes, %d dataTypes, %d Plugin's nodes.\n",
			r.GetPluginName(),
			r.GetVersion(),
			len(r.GetNodeTypes()),
			len(r.GetDataTypes()),
			len(r.Nodes))
	}

	// Print all nodes.
	for _, n := range mo.Nodes {
		FmtPrintf("%d : %s\n", n.LineNo, n.Name)
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
	FmtPrintf("%s.t is %s", persp.Name, ow)

	// Print Node's all attrs.
	for _, a := range persp.Attrs {
		FmtPrintf("%s%s is %d type is %s\n",
			persp.Name,
			a.GetName(),
			len(a.GetAttrValue()),
			a.GetAttrType())
	}

	// Print Node's all children.
	for _, c := range persp.Children {
		FmtPrintf("%s child is %s\n", persp.Name, c.Name)
	}

	// Print Node's parent node.
	if persp.Parent != nil {
		FmtPrintf("%s parent is %s\n", persp.Name, persp.Parent.Name)
	}

	// Get nodes by nodeType.
	transforms, err := mo.GetNodes("transform") // Specify node type.
	if err != nil {
		log.Fatal(err)
	}

	// Print transform nodes.
	for _, t := range transforms {
		FmtPrintf("%d : %s\n", t.LineNo, t.Name)
	}

	// Get specified source connection nodes.
	srcNodes := persp.ListConnections(&ma.ConnectionArgs{
		Source: true,
		Type: "camera",
		AttrName: ".ow", // Specific persp's attr name.
	})

	// Print src nodes.
	for _, t := range srcNodes {
		FmtPrintf("%d : %s\n", t.LineNo, t.Name)
	}

	// Get all destination connection nodes.
	dstNodes := persp.ListConnections(&ma.ConnectionArgs{
		Destination: true,
	})

	// Print dst nodes.
	for _, t := range dstNodes {
		FmtPrintf("%d : %s\n", t.LineNo, t.Name)
	}
}

func FmtPrintf(format string, arguments ...interface{}) {
	//fmt.Printf(format, arguments...)
}