package ed

import (
	"strings"
	"testing"

	"github.com/fwip/posix-utils/pkg/txt"
)

type e2etest struct {
	cmds       string
	initialBuf string
	endBuf     string
	output     string
}

var abc = "a\nb\nc"
var e2etests = []e2etest{
	{"2a\nhello\n.",
		abc,
		"a\nb\nhello\nc",
		"",
	},
	{"2,2d\n",
		abc,
		"a\nc",
		"",
	},
	{"0a\nz\n.",
		abc,
		"z\na\nb\nc",
		"",
	},
	{"3i\nhello\n.",
		abc,
		"a\nb\nhello\nc",
		"",
	},
	{"2,3c\nx\n.",
		abc,
		"a\nx",
		"",
	},
	{"a\na\nb\nc\n.",
		"",
		abc,
		"",
	},
	{".,.d",
		abc,
		"a\nb",
		"",
	},
	{"n",
		abc,
		"",
		"3\tc\n",
	},
	{",p",
		abc,
		"",
		abc,
	},
	{";p",
		abc,
		"",
		"c",
	},
	{"2;p",
		abc,
		"",
		"b",
	},
	{",2p",
		abc,
		"",
		"a\nb",
	},
	{"1\n \n\n",
		abc,
		"",
		abc,
	},
}

func TestEndToEnd(t *testing.T) {

	for _, e := range e2etests {

		t.Run("e2e:"+e.cmds, func(t *testing.T) {
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
