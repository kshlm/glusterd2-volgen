package volgen

// The following go:generate instruction requires github.com/alvaroloes/enumer
// to be installed.
//go:generate enumer -type OptionType,NodeType -output node_enumer.go

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/naoina/toml"
)

// Node describes a node in the volume graph.
// Nodes can be of 2 types:
// - Xlators: Which describe an individual Xlator
// - Targets: Which describe a collection of xlators
// A Node can be defined using TOML in a file.
// Below is an example TOML declaration for a Xlator node.
//   Name = "Disperse"
//   ID = "ec"
//   Before = ["client.xlator"]
//   After = ["dht.xlator"]
//   Conflicts = ["afr.xlator"]
//   Options = [
//     { Key = "count", Type = 2, Default = 1 },
//     { Key = "redundancy", Type = 2, Default = 1 },
//     { Key = "disperse-count", Type = 2, Default = 3 }
//   ]
// Options can also be specified as,
//   [[Options]]
//     Key = "count"
//     Type = 2
//     Default = 1
//
//   [[Options]]
//     Key = "redundancy"
//     Type = 2
//     Default = 1
//
//   [[Options]]
//     Key = "disperse-count"
//     Type = 2
//     Default = 3

type Node struct {
	// Name is the name of the xlator or target
	Name string
	// ID is a short name that will be used when specifying dependencies.
	// The TOML files are expected to be named `<ID>.<Type>`
	ID string
	// Before gives a list of node-ids that this node should appear before in the graph
	Before []string
	// After gives a list of node-ids that this node should appear after in the graph
	After []string
	// Requires gives a list of nodes that are required for this node to work
	Requires []string
	// Conflicts gives a list of nodes that this node cannot work with
	Conflicts []string
	// Options gives a list of xlator options that belong to the xlator described by this node.
	// Applies only to an xlator type Node.
	Options []Option
	// Type specifies the node type, either xlator or target.
	// This will be determined by the file extension, and should not be defined in a node file.
	Type NodeType `toml:"-"`
}

type NodeType int

const (
	TYPE_NONE NodeType = iota
	TYPE_XLATOR
	TYPE_TARGET
	TYPE_MAX
)

type Option struct {
	Key     string
	Type    string
	Default string
}

type OptionType int

const (
	OPT_NONE OptionType = iota
	OPT_STRING
	OPT_INT
	OPT_DOUBLE
	OPT_BOOL
	OPT_MAX
)

const (
	xlatorExt = ".xlator"
	targetExt = ".target"
)

var (
	extToType map[string]NodeType = map[string]NodeType{
		xlatorExt: TYPE_XLATOR,
		targetExt: TYPE_TARGET,
	}
	ERR_UNKNOWN_TYPE error = errors.New("unknown node file type")
)

func NodeFromFile(path string) (*Node, error) {
	var e error

	n := new(Node)

	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	d := toml.NewDecoder(f)
	if e = d.Decode(n); e != nil {
		return nil, e
	}

	if n.Type != TYPE_NONE {
		log.Print("WARNING: ", "Ignoring Node type set in file")
	}

	ext := filepath.Ext(path)
	t, ok := extToType[ext]
	if !ok {
		return nil, ERR_UNKNOWN_TYPE
	}
	n.Type = t

	return n, nil
}
