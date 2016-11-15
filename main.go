package main

// This is just a test application for the `volgen` package

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"

	"github.com/kshlm/glusterd2-volgen/volgen"
)

func main() {
	log.SetHandler(text.New(os.Stderr))
	log.SetLevel(log.FatalLevel)

	if len(os.Args) < 3 {
		log.Fatal("Not enough args: " + os.Args[0] + " <xlators-dir> <target>")
	}

	if e := volgen.LoadXlators(os.Args[1]); e != nil {
		log.WithError(e).Fatal("failed to load xlators")
	}

	t, e := volgen.LoadTarget(os.Args[2])
	if e != nil {
		log.WithError(e).Fatal("failed to load target")
	}

	g, e := t.BuildGraph("test")
	if e != nil {
		log.WithError(e).Fatal("building graph failed")
	}

	//g.WriteDot(os.Stdout)
	g.Write(os.Stdout)

	return
}
