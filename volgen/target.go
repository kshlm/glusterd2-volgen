package volgen

import (
	"errors"
	"os"
	"path/filepath"
)

// Target describes a single graph/volfile for eg. the brick graph/volfile
// A target is a directory with name with a postfix of `targetExt`.
// This directory will contain sym-links to xlator files of xlators or target
// directories to be contained in the graph. The directory MUST contain a file
// named `targetNodeFile` with information about the graph.
// The targetNodeFile can explicitly set 'Requires' and 'Conflicts'
// dependencies to ensure that specifi xlators are or are not loaded into this
// graph.
type Target struct {
	*Node
	Xlators []string
}

var (
	ERR_PATH_NOT_TARGET = errors.New("provided path is not a target")
)

const (
	targetNodeFile = "info"
)

func LoadTarget(path string) (*Target, error) {
	// Ensure the path has targetExt as extension and is a directory
	if filepath.Ext(path) != targetExt {
		return nil, ERR_PATH_NOT_TARGET
	}

	d, e := os.Stat(path)
	if e != nil {
		return nil, e
	}
	if !d.IsDir() {
		return nil, ERR_PATH_NOT_TARGET
	}

	// Load node file
	nf := filepath.Join(path, targetNodeFile)
	n, e := NodeFromFile(nf)
	if e != nil {
		return nil, e
	}

	t := &Target{
		Node: n,
	}

	// Load target xlators
	e = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != xlatorExt {
			return nil
		}
		// Adding the filenames as the Xlator ids
		// TODO: Add actual Node from xlatorMap
		t.Xlators = append(t.Xlators, filepath.Base(path))

		return nil
	})
	if e != nil {
		return nil, e
	}

	// TODO: Need to make sure stuff in Target.Requires are in Target.Xlators
	// TODO: Need to make sure stuff in Target.Conflicts are not in Target.Xlators

	return t, nil
}

// BuildGraph will resolve dependencies and generate a graph from the
// xlators/nodes listed in t.Xlators
func (t *Target) BuildGraph(volume string) (*gNode, error) {
	graph := new(volGraph)

	// Set first xlator in list as root just to get started
	graph.setRoot(t.Xlators[0])

	for _, nid := range t.Xlators[1:] {
		n := graph.node(nid)

		// Add node as child of all after dependencies
		for _, t := range n.After {
			graph.addChild(nid, t)
		}
		// Add all before dependencies as children of n
		for _, t := range n.Before {
			graph.addChild(t, nid)
		}
		// TODO: Handle requires,conflicts,parent,child
	}
	// TODO: Resolve any target nodes if present

	// TODO: Topologically Sort the temporary graph to linearize it
	// This will provide the graph order

	// TODO: With the ordered graph, generate the final graph by
	// enabling/disabling nodes based on the volume information

	// TODO: Handle nodes which have multiple children and nodes which can appear
	// multiple times in the graph (ie. cluster xlators and client xlators)

	return nil, nil
}
