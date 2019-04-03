package main

import (
	"io"
	"os"

	"github.com/fwip/posix-utils/ed"
)

func process(in io.Reader, out io.Writer) {
	e := &ed.Itor{}
	e.ProcessCommands(in, out)
	if w, ok := out.(io.WriteCloser); ok {
		w.Close()
	}
}

func main() {
	process(os.Stdin, os.Stdout)
}
