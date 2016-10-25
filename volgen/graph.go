package volgen

import (
	"io"
	"os"
)

type gNode struct {
	*Node
	parent   *gNode
	children []*gNode
}

// Write will write the graph to the given writer
func (g *gNode) Write(w io.Writer) error {
	for _, n := range g.children {
		n.Write(w)
	}
	// TODO: Write actually

	return nil
}

// WriteToFile writes the graph to the given path, creating the volfile.
// NOTE: Any existing file at the path is truncated.
func (g *gNode) WriteToFile(path string) error {
	f, e := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if e != nil {
		return e
	}
	defer f.Close()
	defer f.Sync()

	return g.Write(f)
}
