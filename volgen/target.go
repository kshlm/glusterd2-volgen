package volgen

import (
	"errors"
	"os"
	"path/filepath"
)

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

func LoadTarget(path string) (*Target, err) {
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
		*Node: n,
	}

	// Load target xlators
	e := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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

	return t, nil
}

// BuildGraph will resolve dependencies and generate a graph from the
// xlators/nodes listed in t.Xlators
func (t *Target) BuildGraph(volume string) (*Node, error) {
	// TODO: Everything

	return nil, nil
}
