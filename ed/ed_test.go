package ed

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/fwip/posix-utils/txt"
)

type e2etest struct {
	cmds       string
	initialBuf string
	endBuf     string
}

var e2etests = []e2etest{
	{"2a\nhello\n.",
		"a\nb\nc",
		"a\nb\nhello\nc",
	},
	{"2,2d\n",
		"a\nb\nc",
		"a\nc",
	},
	{"0a\nz\n.",
		"a\nb\nc",
		"z\na\nb\nc",
	},
	{"3i\nhello\n.",
		"a\nb\nc",
		"a\nb\nhello\nc",
	},
}

func TestEndToEnd(t *testing.T) {

	for _, e := range e2etests {

		t.Run("e2e:"+e.cmds, func(t *testing.T) {
			input := strings.NewReader(e.cmds)
			buf := strings.NewReader(e.initialBuf)
			devnull := ioutil.Discard

			ed := NewEditor(input, devnull)
			ed.pt = txt.NewPieceTable(buf, len(e.initialBuf))

			ed.ProcessCommands(input, os.Stderr)

			actual := ed.String()
			if e.endBuf+"\n" != actual {
				t.Errorf("expected\n%q\ngot\n%q\n", e.endBuf+"\n", actual)
			}
		})
	}
}
