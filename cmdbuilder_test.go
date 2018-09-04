package mayaascii

import (
	"testing"
	"log"
)

func TestCmdBuilder_Append(t *testing.T) {
	c := CmdBuilder{}
	c.Append("append")
	if len(c.cmdLine) != 1 {
		t.Fatal("Append failed")
	}
	if c.cmdLine[0] != "append" {
		t.Fatal("Append failed")
	}
}

func TestCmdBuilder_IsCmdEOF(t *testing.T) {
	c := CmdBuilder{}
	if c.IsCmdEOF() {
		t.Fatal("IsCmdEOF failed")
	}
	c.Append("")
	if c.IsCmdEOF() {
		t.Fatal("IsCmdEOF failed")
	}
	c.Append("test;")
	if !c.IsCmdEOF() {
		t.Fatal("IsCmdEOF failed")
	}
}

func TestCmdBuilder_Clear(t *testing.T) {
	c := CmdBuilder{}
	if len(c.cmdLine) != 0 {
		t.Fatal("CmdBuilder{} failed")
	}
	c.Append("test")
	if len(c.cmdLine) != 1 {
		t.Fatal("Append failed")
	}
	c.Clear()
	if len(c.cmdLine) != 0 {
		t.Fatal("Clear failed")
	}
}

func testParse(t *testing.T, raw string, cmdName string, token []string) {
	c := &CmdBuilder{}
	c.Append(raw)
	if !c.IsCmdEOF() {
		t.Fatal("IsCmdEOF failed")
	}
	cmd := c.Parse()
	if cmd.Raw != raw {
		t.Fatal("Parsed Raw was failed")
	}
	if cmd.CmdName != cmdName {
		log.Printf("got CmdName %s\nwant CmdName %s", cmd.CmdName, cmdName)
		t.Fatal("Parsed CmdName was failed")
	}
	if len(cmd.Token) != len(token) {
		log.Printf("got Token len %d\nwant token len %d",
			len(cmd.Token),
			len(token))
		log.Printf("got cmd.Token is %v", cmd.Token)
		log.Printf("want token is %v", token)
		t.Fatal("Parsed Token was failed")
	}
	for i := 0; i < len(cmd.Token); i++ {
		if cmd.Token[i] != token[i] {
			log.Printf("got %v\nwant %v", cmd.Token[i], token[i])
			t.Fatal("Parsed Token was failed")
		}
	}
}

func TestCmdBuilder_Parse(t *testing.T) {
	testParse(
		t,
		`createNode transform -s -n "persp";`,
		"createNode",
		[]string{
			"createNode",
			"transform",
			"-s",
			"-n",
			`"persp"`,
		},
	)
	testParse(
		t,
		`rename -uid "B0F0F886-4CD8-DA88-AC80-C1B83173300D";`,
		"rename",
		[]string{
			"rename",
			"-uid",
			`"B0F0F886-4CD8-DA88-AC80-C1B83173300D"`,
		},
	)
	testParse(
		t,
		`setAttr ".v" no;`,
		"setAttr",
		[]string{
			"setAttr",
			`".v"`,
			"no",
		},
	)
	testParse(
		t,
		`setAttr ".t" -type "double3" 2278.3519359468883 -3416.658654922082 6789.7722233495588 ;`,
		"setAttr",
		[]string{
			"setAttr",
			`".t"`,
			"-type",
			`"double3"`,
			"2278.3519359468883",
			"-3416.658654922082",
			"6789.7722233495588",
		},
	)
	testParse(
		t,
		`	setAttr ".rpt" -type "double3" -1.7624290720136436e-013 -1.0748060174015983e-013
	9.6419181914696075e-015 ;`,

		"setAttr",
		[]string{
			"setAttr",
			`".rpt"`,
			"-type",
			`"double3"`,
			"-1.7624290720136436e-013",
			"-1.0748060174015983e-013",
			"9.6419181914696075e-015",
		},
	)
	testParse(
		t,
		`setAttr ".hc" -type "string" "viewSet -t %camera";`,
		"setAttr",
		[]string{
			"setAttr",
			`".hc"`,
			"-type",
			`"string"`,
			`"viewSet -t %camera"`,
		},
	)
	testParse(
		t,
		`setAttr ".dcc" -type "string" "Ambient+Diffuse";`,
		"setAttr",
		[]string{
			"setAttr",
			`".dcc"`,
			"-type",
			`"string"`,
			`"Ambient+Diffuse"`,
		},
	)
	testParse(
		t,
		`setAttr ".ixp" -type "string" "float $detonateFrame = .I[0];\n\n.O[0] = 0;\nif( frame >= $detonateFrame && frame < $detonateFrame + 1) {\n\t.O[0] = 1000;\n}";`,
		"setAttr",
		[]string{
			"setAttr",
			`".ixp"`,
			"-type",
			`"string"`,
			`"float $detonateFrame = .I[0];\n\n.O[0] = 0;\nif( frame >= $detonateFrame && frame < $detonateFrame + 1) {\n\t.O[0] = 1000;\n}"`,
		},
	)
	testParse(
		t,
		`setAttr ".aal" -type "attributeAlias" {"detonationFrame","borderConnections[0]","incandescence"
,"borderConnections[1]","color","borderConnections[2]","nucleusSolver","publishedNodeInfo[0]"
} ;`,
		"setAttr",
		[]string{
			"setAttr",
			`".aal"`,
			`-type`,
			`"attributeAlias"`,
			"{",
			"detonationFrame",
			"borderConnections[0]",
			"incandescence",
			"borderConnections[1]",
			"color",
			"borderConnections[2]",
			"nucleusSolver",
			"publishedNodeInfo[0]",
			"}",
		},
	)
	testParse(
		t,
		`setAttr ".b" -type "string" (
		"long text "
		+ "long text");`,
		"setAttr",
		[]string{
			"setAttr",
			`".b"`,
			"-type",
			`"string"`,
			"long text long text",
		},
	)
}
