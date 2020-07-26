package connection

import (
	"github.com/nrtkbb/go-mayaascii/cmd"
)

type Connections struct {
	source []*cmd.ConnectAttr
}

func (ci *Connections) Append(ca *cmd.ConnectAttr) {
	ci.source = append(ci.source, ca)
}

func (ci *Connections) GetSrcNames(nodeName string) []string {
	var results []string
	for _, s := range ci.source {
		if s.DstNode == nodeName {
			results = append(results, s.SrcNode)
		}
	}
	return results
}

func (ci *Connections) GetSrcNamesAttr(nodeName, attrName string) []string {
	var results []string
	for _, s := range ci.source {
		if s.DstNode == nodeName && s.DstAttr == attrName {
			results = append(results, s.SrcNode)
		}
	}
	return results
}

func (ci *Connections) GetDstNames(nodeName string) []string {
	var results []string
	for _, s := range ci.source {
		if s.SrcNode == nodeName {
			results = append(results, s.DstNode)
		}
	}
	return results
}

func (ci *Connections) GetDstNamesAttr(nodeName, attrName string) []string {
	var results []string
	for _, s := range ci.source {
		if s.SrcNode == nodeName && s.SrcAttr == attrName {
			results = append(results, s.DstNode)
		}
	}
	return results
}
