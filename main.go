package main

// This is just a test application for the `volgen` package

import (
	"log"
	"os"

	"github.com/kshlm/glusterd2-volgen/volgen"
)

func main() {
	if len(os.Args) < 3 {
		log.Panic("Not enough args: ", os.Args[0], " <xlators-dir> <target>")
	}

	if e := volgen.LoadXlators(os.Args[1]); e != nil {
		log.Panic("failed to load xlators: ", e)
	}

	t, e := volgen.LoadTarget(path)
	if e != nil {
		log.Panic("failed to load target: ", e)
	}

	g, e := t.BuildGraph()
	if e != nil {
		log.Panic("building graph failed: ", e)
	}

	_ = g.Write(os.Stdout)

	return
}
