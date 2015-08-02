package volgen

import (
	"fmt"
	"io"
)

func (x Xlator) DumpGraph(w io.Writer) {

	for _, xl := range x.Children {
		xl.DumpGraph(w)
	}

	str := fmt.Sprintf("volume %s\n", x.Name)
	str += fmt.Sprintf("   type %s", x.Type)
	flag := true
	for _, ch_xl := range x.Children {
		if flag {
			str += fmt.Sprintf("\n   subvolumes")
			flag = false
		}
		str += fmt.Sprintf(" %s", ch_xl.Name)
	}
	str += fmt.Sprintf("\nend-volume\n\n")

	_, _ = w.Write([]byte(str)) //Write into given interface
}
