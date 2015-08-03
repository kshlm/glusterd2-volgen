package main

import (
	"bytes"
	"fmt"

	"github.com/kshlm/glusterd2-volgen/volgen"
)

func main() {
	graph := generateDummyGraph()

	var b bytes.Buffer

	graph.DumpGraph(&b)

	fmt.Print(b.String())
}

func generateDummyGraph() volgen.Xlator {
	dht := volgen.Xlator{
		Name:    "distribute0",
		Type:    "cluster/dht",
		Options: map[string]string{"transport-type": "tcp", "ping-timeout": "42"},
	}

	afr0 := volgen.Xlator{
		Name:    "replicate0",
		Type:    "cluster/afr",
		Options: map[string]string{"transport-type": "tcp", "ping-timeout": "42"},
	}

	afr1 := volgen.Xlator{
		Name:    "replicate1",
		Type:    "cluster/afr",
		Options: map[string]string{"transport-type": "tcp", "ping-timeout": "42"},
	}

	afr2 := volgen.Xlator{
		Name:    "replicate2",
		Type:    "cluster/afr",
		Options: map[string]string{"transport-type": "tcp", "ping-timeout": "42"},
	}

	client0 := volgen.Xlator{
		Name:    "client0",
		Type:    "cluster/afr",
		Options: map[string]string{"transport-type": "tcp", "ping-timeout": "42"},
	}

	client1 := volgen.Xlator{
		Name:    "client1",
		Type:    "cluster/afr",
		Options: map[string]string{"transport-type": "tcp", "ping-timeout": "42"},
	}

	client2 := volgen.Xlator{
		Name:    "client2",
		Type:    "cluster/afr",
		Options: map[string]string{"transport-type": "tcp", "ping-timeout": "42"},
	}

	dht.Children = append(dht.Children, afr0, afr1, afr2)
	afr0.Children = append(afr0.Children, client0)
	afr1.Children = append(afr1.Children, client1)
	afr2.Children = append(afr2.Children, client2)

	return dht
}
