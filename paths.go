package quamina

import (
	"strings"
)

const PATH_SEPARATOR = "\n"

type Node interface {
	get(name string) (Node, bool)
	getFields() map[string][]byte
	nodesCount() int
	names() []string

	getOrCreate(name string) Node
	addField(name string, path []byte)
}

type pathsIndex struct {
	nodes map[string]Node

	// fields map from it's name to it's full path (as specific in quamina.Field)
	// will be present only on the leafs
	fields map[string][]byte
}

func newPathsIndex() pathsIndex {
	return pathsIndex{
		nodes:  make(map[string]Node),
		fields: make(map[string][]byte),
	}
}

func (p pathsIndex) add(path string) {
	parts := strings.Split(path, PATH_SEPARATOR)

	var node Node
	node = p.getOrCreate(parts[0])

	for i, part := range parts[1:] {
		if i == len(parts)-2 {
			node.addField(part, []byte(path))
		} else {
			node = node.getOrCreate(part)
		}
	}
}

func (p pathsIndex) get(name string) (Node, bool) {
	n, ok := p.nodes[name]
	return n, ok
}

func (p pathsIndex) getOrCreate(name string) Node {
	if _, ok := p.nodes[name]; !ok {
		p.nodes[name] = newPathsIndex()
	}

	return p.nodes[name]
}

func (p pathsIndex) addField(name string, path []byte) {
	if _, ok := p.fields[name]; !ok {
		p.fields[name] = path
	}
}

func (p pathsIndex) getFields() map[string][]byte {
	return p.fields
}

func (p pathsIndex) nodesCount() int {
	return len(p.nodes)
}

func (p pathsIndex) names() []string {
	na := make([]string, 0)

	for n := range p.nodes {
		na = append(na, n)
	}
	return na

}
