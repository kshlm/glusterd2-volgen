package volgen

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/apex/log"
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
	Xlators map[string]*Node
}

var (
	ErrPathNotTarget = errors.New("provided path is not a target")
	ErrMultipleRoots = errors.New("multiple roots for target")
)

const (
	targetNodeFile = "info"
	NoneTarget     = "NONE"
)

func LoadTarget(path string) (*Target, error) {
	log.WithField("path", path).Debug("attempting to load target")

	// Ensure the path has targetExt as extension and is a directory
	if filepath.Ext(path) != targetExt {
		return nil, ErrPathNotTarget
	}

	d, e := os.Stat(path)
	if e != nil {
		return nil, e
	}
	if !d.IsDir() {
		return nil, ErrPathNotTarget
	}

	// Load node file
	nf := filepath.Join(path, targetNodeFile)
	n, e := NodeFromFile(nf)
	if e != nil {
		return nil, e
	}

	t := &Target{
		Node:    n,
		Xlators: make(map[string]*Node),
	}

	// Load target xlators
	log.WithField("target", t.ID).Debug("loading target xlators and targets")
	e = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == xlatorExt {
			id := filepath.Base(path)
			x, e := FindXlator(id)
			if e != nil {
				log.WithError(e).WithField("path", path).Error("couldn't load target xlator")
				return e
			}
			t.Xlators[id] = x
			log.WithFields(log.Fields{
				"xlator": x.ID,
				"target": t.ID,
			}).Debug("added xlator to target")

		} else if filepath.Ext(path) == targetExt {
			//TODO: Find and add target to xlator list

		} else {
			log.WithField("path", path).Error("path is not valid xlator or target")
		}

		return nil
	})
	if e != nil {
		log.WithError(e).Error("failed to load target")
		return nil, e
	}

	// TODO: Need to make sure stuff in Target.Requires are in Target.Xlators
	// TODO: Need to make sure stuff in Target.Conflicts are not in Target.Xlators

	return t, nil
}

// BuildGraph will resolve dependencies and generate a graph from the
// xlators/nodes listed in t.Xlators
func (t *Target) BuildGraph(volume string) (*volGraph, error) {
	graph := new(volGraph)
	graph.members = make(map[string]*gNode)

	for nid, n := range t.Xlators {
		// If the current node cannot come after any other node, it is the root
		// Panic if root already exists
		if n.After != nil && n.After[0] == NoneTarget {
			if graph.root == nil {
				if e := graph.setRoot(nid); e != nil {
					log.WithError(e).WithField("node", nid).Error("failed to set node as root of graph")
				}
			} else {
				log.WithFields(log.Fields{
					"oldroot": graph.root.ID,
					"newroot": n.ID,
				}).Fatal("multiple roots found")
			}
		} else {
			log.WithFields(log.Fields{
				"node":  nid,
				"after": n.After,
			}).Debug("AFTER dependecies for node")
			// Add node as child of all AFTER dependencies
			for _, t := range n.After {
				log.WithFields(log.Fields{
					"node":       nid,
					"other":      t,
					"dependency": "AFTER",
				}).Debug("setting dependency for node")
				if e := graph.addChild(nid, t); e != nil {
					log.WithFields(log.Fields{
						"node":       nid,
						"other":      t,
						"dependency": "AFTER",
					}).WithError(e).Error("setting dependency for node failed")
					return nil, e
				}
			}
		}
		// Add all BEFORE dependencies as children of n
		log.WithFields(log.Fields{
			"xlator": n.ID,
			"before": n.Before,
		}).Debug("BEFORE dependecies for xlator")
		for _, t := range n.Before {
			if t == NoneTarget {
				continue
			}
			log.WithFields(log.Fields{
				"node":       nid,
				"other":      t,
				"dependency": "BEFORE",
			}).Debug("setting dependency for node")
			if e := graph.addChild(t, nid); e != nil {
				log.WithFields(log.Fields{
					"node":       n.ID,
					"other":      t,
					"dependency": "BEFORE",
				}).WithError(e).Error("setting dependency for node failed")
				return nil, e
			}
		}
		// TODO: Handle requires,conflicts,parent,child
	}
	// Set root as parent for any node that doesn't have a parent
	for _, m := range graph.members {
		if m != graph.root && m.Parent == nil {
			m.Parent = graph.root
		}
	}
	// TODO: Topologically Sort the temporary graph to linearize it
	// This will provide the graph order
	//graph.TopoSort()

	// TODO: Resolve any target nodes if present

	// TODO: With the ordered graph, generate the final graph by
	// enabling/disabling nodes based on the volume information

	// TODO: Handle nodes which have multiple children and nodes which can appear
	// multiple times in the graph (ie. cluster xlators and client xlators)

	return graph, nil
}
