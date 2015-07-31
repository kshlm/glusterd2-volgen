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
		Name: "distribute0",
		Type: "cluster/dht",
	}

	afr0 := volgen.Xlator{
		Name: "replicate0",
		Type: "cluster/afr",
	}

	afr1 := volgen.Xlator{
		Name: "replicate1",
		Type: "cluster/afr",
	}

	afr2 := volgen.Xlator{
		Name: "replicate2",
		Type: "cluster/afr",
	}

	client0 := volgen.Xlator{
		Name: "client0",
		Type: "cluster/afr",
	}

	client1 := volgen.Xlator{
		Name: "client1",
		Type: "cluster/afr",
	}

	client2 := volgen.Xlator{
		Name: "client2",
		Type: "cluster/afr",
	}

	dht.Children = append(dht.Children, afr0, afr1, afr2)
	afr0.Children = append(afr0.Children, client0)
	afr1.Children = append(afr1.Children, client1)
	afr2.Children = append(afr2.Children, client2)

	return dht
}
