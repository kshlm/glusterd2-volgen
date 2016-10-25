# Prototype Volgen for GD2

This repository implements a prototype flexible `volgen` pacakge for [GD2][1].

The design of this package is described in the [GD2 wiki][2].
This package is a WIP and may be ahead of the wiki in terms of design.
I'll try to keep the wiki updated regularly.

# Things to be done

The main things still to be done are (in decreasing order of priority),

- Implement the dependency resolution to be actually able to build graphs.
- Decide how to handle nodes with multiple children
- Actually implement the graph writer, to be able to create volfiles.

[1]: https://github.com/gluster/glusterd2
[2]: https://github.com/gluster/glusterd2/wiki/Flexible-Volgen
