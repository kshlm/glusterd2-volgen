package volgen

import (
	"errors"
	"io"
	"os"
)

type gNode struct {
	*Node
	parent   *gNode
	children []*gNode
}

type volGraph struct {
	root    *gNode
	members []*gNode
}

var (
	ErrNodeNotFound        = errors.New("node not found")
	ErrNodeMultipleParents = errors.New("node has more than one parent")
)

// Write will write the graph to the given writer
func (n *gNode) Write(w io.Writer) error {
	for _, c := range n.children {
		c.Write(w)
	}
	// TODO: Write actually

	return nil
}

// WriteToFile writes the graph to the given path, creating the volfile.
// NOTE: Any existing file at the path is truncated.
func (n *gNode) WriteToFile(path string) error {
	f, e := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if e != nil {
		return e
	}
	defer f.Close()
	defer f.Sync()

	return n.Write(f)
}

// node returns the matching *gNode from the graph.
// If a matching gNode isn't present, it loads one by searching in the xlatorMap
// If a match is not found nil is returned
func (g *volGraph) node(id string) *gNode {
	// Find a match if present in g.members
	for _, m := range g.members {
		if m.ID == id {
			return m
		}
	}
	// Find xlator with given ID in the xlator map
	for _, xl := range xlatorMap {
		if xl.ID == id {
			n := &gNode{xl, nil, nil}
			if g.root != nil {
				n.parent = g.root
				g.root.children = append(g.root.children, n)
			}
			return n
		}
	}

	return nil
}

// setRoot sets the given node as the root of the graph
func (g *volGraph) setRoot(id string) error {
	r := g.node(id)
	if r == nil {
		return ErrNodeNotFound
	}
	// If a root node exists already make it the child of the new root
	if g.root != nil {
		g.root.parent = r
		r.children = append(r.children, g.root)
	}
	g.root = r
	g.members = append(g.members, r)

	return nil
}

// addChild adds node as a child of the given parent
func (g *volGraph) addChild(nid, pid string) error {
	n := g.node(nid)
	if n == nil {
		return ErrNodeNotFound
	}

	p := g.node(pid)
	if n == nil {
		return ErrNodeNotFound
	}

	if n.parent != nil {
		return ErrNodeMultipleParents
	}
	n.parent = p
	p.children = append(p.children, n)

	return nil
}
