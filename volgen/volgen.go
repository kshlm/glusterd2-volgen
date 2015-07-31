package volgen

import (
	"io"
)

func (x Xlator) DumpGraph(w io.Writer) {
	_, _ = w.Write([]byte("Dump function not implemented yet"))
}
