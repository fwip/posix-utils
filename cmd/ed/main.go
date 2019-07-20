package main

import (
	"io"
	"os"
	"strings"

	"github.com/fwip/posix-utils/pkg/ed"
)

func process(in io.Reader, out io.Writer, edit *string) {
	e := &ed.Itor{}
	if edit != nil {
		// If an argument is supplid on the command-line, open it in the editor
		editCmd := strings.NewReader("e " + *edit + "\n")
		in = io.MultiReader(editCmd, in)
	}
	e.ProcessCommands(in, out)
	if w, ok := out.(io.WriteCloser); ok {
		w.Close()
	}
}

func main() {
	var editFile *string
	if len(os.Args) > 1 {
		editFile = &os.Args[1]
	}
	process(os.Stdin, os.Stdout, editFile)
}
