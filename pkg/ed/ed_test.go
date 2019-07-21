package ed

import (
	"strings"
	"testing"

	"github.com/fwip/posix-utils/pkg/txt"
)

type e2etest struct {
	cmds        string
	initialBuf  string
	endBuf      string
	output      string
	description string
}

var abc = "a\nb\nc"
var e2etests = []e2etest{
	{"2a\nhello\n.",
		abc,
		"a\nb\nhello\nc",
		"",
		"append lines to middle of file",
	},
	{"2,2d\n",
		abc,
		"a\nc",
		"",
		"delete line from middle of file",
	},
	{"0a\nz\n.",
		abc,
		"z\na\nb\nc",
		"",
		"append line at beginning of file",
	},
	{"3i\nhello\n.",
		abc,
		"a\nb\nhello\nc",
		"",
		"insert line in middle of file",
	},
	{"2,3c\nx\n.",
		abc,
		"a\nx",
		"",
		"change lines in middle of file",
	},
	{"a\na\nb\nc\n.",
		"",
		abc,
		"",
		"append lines to an empty file",
	},
	{".,.d",
		abc,
		"a\nb",
		"",
		"delete current (last) line of file",
	},
	{"n",
		abc,
		"",
		"3\tc\n",
		"number last line of file",
	},
	{",p",
		abc,
		"",
		abc,
		"print whole file",
	},
	{";p",
		abc,
		"",
		"c",
		"print last line of file",
	},
	{"2;p",
		abc,
		"",
		"b",
		"print second line of file",
	},
	{",2p",
		abc,
		"",
		"a\nb",
		"print first two lines of file",
	},
	{"1\n \n\n",
		abc,
		"",
		abc,
		"?????",
	},
	{"/b/",
		abc,
		abc,
		"b",
		"regex works",
	},
}

func TestEndToEnd(t *testing.T) {

	for _, e := range e2etests {

		t.Run("e2e:"+e.description, func(t *testing.T) {
			input := strings.NewReader(e.cmds)
			buf := strings.NewReader(e.initialBuf)

			output := &strings.Builder{}

			ed := NewEditor(input, output)
			ed.pt = txt.NewPieceTable(buf, len(e.initialBuf))

			ed.ProcessCommands(input, output)

			actual := ed.String()
			if e.endBuf != "" && e.endBuf+"\n" != actual {
				t.Errorf("expected buffer\n%q\ngot\n%q\n", e.endBuf+"\n", actual)
			}
			if e.output != "" && e.output+"\n" != output.String() {
				t.Errorf("expected output\n%q\ngot\n%q\n", e.output+"\n", output.String())

			}
		})
	}
}
