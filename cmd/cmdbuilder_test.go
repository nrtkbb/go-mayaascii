package cmd

import (
	"log"
	"testing"
)

func TestCmdBuilder_Append(t *testing.T) {
	cb := CmdBuilder{}
	cb.Append("append")
	if len(cb.cmdLine) != 1 {
		t.Fatal("Append failed")
	}
	if cb.cmdLine[0] != "append" {
		t.Fatal("Append failed")
	}
}

func TestCmdBuilder_IsCmdEOF(t *testing.T) {
	cb := CmdBuilder{}
	if cb.IsCmdEOF() {
		t.Fatal("IsCmdEOF failed")
	}
	cb.Append("")
	if cb.IsCmdEOF() {
		t.Fatal("IsCmdEOF failed")
	}
	cb.Append("test;")
	if !cb.IsCmdEOF() {
		t.Fatal("IsCmdEOF failed")
	}
}

func TestCmdBuilder_Clear(t *testing.T) {
	cb := CmdBuilder{}
	if len(cb.cmdLine) != 0 {
		t.Fatal("CmdBuilder{} failed")
	}
	cb.Append("test")
	if len(cb.cmdLine) != 1 {
		t.Fatal("Append failed")
	}
	cb.Clear()
	if len(cb.cmdLine) != 0 {
		t.Fatal("Clear failed")
	}
}

func testParse(t *testing.T, raw string, cmdName Type, token []string) {
	cb := &CmdBuilder{}
	cb.Append(raw)
	if !cb.IsCmdEOF() {
		t.Fatal("IsCmdEOF failed")
	}
	c := cb.Parse()
	if c.Raw != raw {
		t.Fatal("Parsed Raw was failed")
	}
	if c.Type != cmdName {
		log.Printf("got Type %s\nwant Type %s", c.Type, cmdName)
		t.Fatal("Parsed Type was failed")
	}
	if len(c.Token) != len(token) {
		log.Printf("got Token len %d\nwant token len %d",
			len(c.Token),
			len(token))
		log.Printf("got c.Token is %v", c.Token)
		log.Printf("want token is %v", token)
		t.Fatal("Parsed Token was failed")
	}
	for i := 0; i < len(c.Token); i++ {
		if c.Token[i] != token[i] {
			log.Printf("got %v\nwant %v", c.Token[i], token[i])
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
