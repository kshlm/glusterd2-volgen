package volgen

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/apex/log"
)

type gNode struct {
	*Node
	Parent   *gNode
	Children map[string]*gNode
}

type volGraph struct {
	root    *gNode
	members map[string]*gNode
}

var (
	ErrNodeNotFound        = errors.New("node not found")
	ErrNodeMultipleParents = errors.New("node has more than one parent")
)

const (
	volfileTemplate = `
{{define "volume"}}
volume {{.Node.Name}}
	type {{.Node.ID}}
	{{range $opt := .Node.Options -}}
	option {{$opt.Key}} {{$opt.Default}}
	{{- end}}
	subvolumes {{range $child := .Children}}{{$child.Node.ID}} {{end}}
{{end}}
`
	dotfileTemplate = `
{{define "volume"}}
{{- $node := . -}}
{{- range $child := .Children}}
"{{$node.ID}}" -> "{{$child.Node.ID}}"{{end -}}
{{end}}
`
	dotHeader  = "digraph {"
	dotTrailer = "\n}"
)

var (
	volTmpl = template.Must(template.New("volume").Parse(volfileTemplate))
	dotTmpl = template.Must(template.New("volume").Parse(dotfileTemplate))
)

// Write will write the graph to the given writer
func (n *gNode) Write(w io.Writer) error {
	return n.writeVol(w, make(map[string]bool))
}

func (n *gNode) writeVol(w io.Writer, processed map[string]bool) error {
	for _, c := range n.Children {
		c.writeVol(w, processed)
	}
	if !processed[n.ID] {
		if e := volTmpl.Execute(w, n); e != nil {
			return e
		}
		processed[n.ID] = true
	}
	return nil
}

// WriteDot writes a dot graph of volume to the writer
func (n *gNode) WriteDot(w io.Writer) error {
	w.Write([]byte(dotHeader))
	n.writeDot(w, make(map[string]bool))
	w.Write([]byte(dotTrailer))

	return nil
}

func (n *gNode) writeDot(w io.Writer, processed map[string]bool) error {
	for _, c := range n.Children {
		c.writeDot(w, processed)
	}
	if !processed[n.ID] {
		if e := dotTmpl.Execute(w, n); e != nil {
			return e
		}
		processed[n.ID] = true
	}
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
	log.WithField("node", id).Debug("finding node in graph")
	// Find a match if present in g.members
	n, ok := g.members[id]
	if ok {
		log.WithField("node", n.ID).Debug("found node")
		return n
	}
	log.WithField("node", id).Debug("node not found in existing members")

	switch ext := filepath.Ext(id); ext {
	case xlatorExt:
		log.WithField("node", id).Debug("finding node in global xlator list")
		xl, err := FindXlator(id)
		if err != nil {
			log.WithField("node", id).WithError(err).Error("could not find node")
			return nil
		}
		n := &gNode{xl, nil, make(map[string]*gNode)}
		g.members[id] = n
		if g.root != nil {
			n.Parent = g.root
			g.root.Children[n.ID] = n
		}
		log.WithField("node", n.ID).Debug("found node")
		return n

	case targetExt:
		log.WithField("node", id).Debug("TODO: return target")
		return nil

	default:
		return nil

	}

	return nil
}

// setRoot sets the given node as the root of the graph
func (g *volGraph) setRoot(id string) error {
	log.WithField("xlatorid", id).Debug("attempting to set graph root")

	r := g.node(id)
	if r == nil {
		return ErrNodeNotFound
	}

	// If a root node exists already make it the child of the new root
	if g.root != nil {
		log.WithFields(log.Fields{
			"existingroot": g.root.ID,
			"newroot":      r.ID,
		}).Debug("existing root found")
		r.Children[g.root.ID] = g.root
	}
	g.root = r
	log.WithField("xlator", g.root.ID).Debug("root of graph set")
	log.WithField("children", g.root.Children).Debug("root's children")

	return nil
}

// addChild adds node as a child of the given parent
func (g *volGraph) addChild(nid, pid string) error {
	log.WithFields(log.Fields{
		"parent": pid,
		"child":  nid,
	}).Debug("attempting to add child to parent")

	n := g.node(nid)
	if n == nil {
		return ErrNodeNotFound
	}

	p := g.node(pid)
	if p == nil {
		return ErrNodeNotFound
	}

	//if n.Parent != nil {
	//return ErrNodeMultipleParents
	//}
	n.Parent = p
	p.Children[nid] = n

	log.WithFields(log.Fields{
		"parent": p.ID,
		"child":  n.ID,
	}).Debug("child added to parent")
	log.WithFields(log.Fields{
		"parent":   p.ID,
		"children": p.Children,
	}).Debug("children of parent")
	return nil
}

func (g *volGraph) Write(w io.Writer) error {
	w.Write([]byte(dotHeader))
	for _, n := range g.members {
		dotTmpl.Execute(w, n)
	}
	w.Write([]byte(dotTrailer))
	return nil
}
