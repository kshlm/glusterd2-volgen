#!/bin/sh

# Requires Graphviz to be installed

go run main.go examples/xlators examples/brick.target 2>/dev/null | dot -Tsvg > brick-graph.svg
