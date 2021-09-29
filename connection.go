package mayaascii

type Connections struct {
	source []*ConnectAttrCmd
}

func NewConnections() Connections {
	return Connections{
		source: []*ConnectAttrCmd{},
	}
}

func (ci *Connections) Append(ca *ConnectAttrCmd) {
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
