package connection

import (
	"fmt"
	"log"

	"github.com/nrtkbb/go-mayaascii/cmd"
)

type Connections struct {
	SrcNodes conNodes
	DstNodes conNodes
}

func (c *Connections) Append(con *cmd.ConnectAttr) {
	c.SrcNodes.Append(con.SrcNode, con)
	c.DstNodes.Append(con.DstNode, con)
}

func (c *Connections) GetNames(
	node, attr string,
	up bool,
	filters *[]string) *[]string {
	sets := map[string]string{}
	if up {
		cons := c.DstNodes.Get(node)
		for _, con := range cons {
			if con.DstAttr != attr {
				continue
			}
			if filters == nil {
				sets[con.SrcNode] = con.SrcAttr
			} else {
				for _, filter := range *filters {
					if filter == con.SrcNode {
						sets[con.SrcNode] = con.SrcAttr
					}
				}
			}
		}
	} else {
		cons := c.SrcNodes.Get(node)
		for _, con := range cons {
			if con.SrcAttr != attr {
				continue
			}
			if filters == nil {
				sets[con.DstNode] = con.DstAttr
			} else {
				for _, filter := range *filters {
					if filter == con.SrcNode {
						sets[con.DstAttr] = con.DstAttr
					}
				}
			}
		}
	}
	var results []string
	for k, v := range sets {
		results = append(results, fmt.Sprintf("%s.%s", k, v))
	}
	return &results
}

func (c *Connections) GetNodes(
	node, attr string,
	up bool,
	inMap *map[string]*cmd.CreateNode) []*cmd.CreateNode {
	sets := map[string]*cmd.CreateNode{}
	if up {
		cons := c.DstNodes.Get(node)
		for _, con := range cons {
			if con.DstAttr != attr {
				continue
			}
			node, ok := (*inMap)[con.SrcNode]
			if ok {
				sets[node.NodeName] = node
			}
		}
	} else {
		cons := c.SrcNodes.Get(node)
		for _, con := range cons {
			log.Println(node, attr, con.SrcAttr)
			if con.SrcAttr != attr {
				continue
			}
			node, ok := (*inMap)[con.DstNode]
			if ok {
				sets[node.NodeName] = node
			}
		}
	}
	var results []*cmd.CreateNode
	for _, v := range sets {
		results = append(results, v)
	}
	return results
}

func (c *Connections) SearchNodes(
	key string,
	up bool,
	inMap *map[string]*cmd.CreateNode) []*cmd.CreateNode {
	histories := c.search(key, up)
	sets := map[string]*cmd.CreateNode{}
	if up {
		for _, history := range histories {
			value, ok := (*inMap)[history.DstNode]
			if ok {
				sets[value.NodeName] = value
			}
		}
	} else {
		for _, history := range histories {
			value, ok := (*inMap)[history.SrcNode]
			if ok {
				sets[value.NodeName] = value
			}
		}
	}
	results := make([]*cmd.CreateNode, len(sets))
	i := 0
	for _, v := range sets {
		results[i] = v
		i++
	}
	return results
}

func (c *Connections) search(key string, up bool) []*cmd.ConnectAttr {
	var histories []*cmd.ConnectAttr
	if up {
		values := c.DstNodes.Get(key)
		histories = append(histories, values...)
		for _, value := range values {
			histories = append(histories, c.search(value.SrcNode, up)...)
		}
	} else {
		values := c.SrcNodes.Get(key)
		histories = append(histories, values...)
		for _, value := range values {
			histories = append(histories, c.search(value.DstNode, up)...)
		}
	}
	return histories
}

type conNodes struct {
	m map[string][]*cmd.ConnectAttr
}

func (c *conNodes) Length() int {
	return len(c.m)
}

func (c *conNodes) String() string {
	var buf []byte
	buf = append(buf, "[\n"...)
	for k, v := range c.m {
		buf = append(buf, k...)
		buf = append(buf, ": "...)
		buf = append(buf, fmt.Sprintf("%v", v)...)
		buf = append(buf, ",\n"...)
	}
	buf = append(buf, "]\n"...)
	return string(buf)
}

func (c *conNodes) Get(key string) []*cmd.ConnectAttr {
	if c.m == nil {
		c.m = map[string][]*cmd.ConnectAttr{}
	}
	v, ok := c.m[key]
	if ok {
		return v
	} else {
		return []*cmd.ConnectAttr{}
	}
}

func (c *conNodes) Append(key string, con *cmd.ConnectAttr) {
	v := c.Get(key)
	c.m[key] = append(v, con)
}
