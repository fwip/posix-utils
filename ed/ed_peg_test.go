package ed

import (
	"fmt"
	"testing"
)

// These commands take no addresses and accept no parameters
var psimpleCmds = []string{
	"h",
	"H",
	"P",
	"q",
	"Q",
	"u",
}

var pParamCmds = []string{
	"e",
	"e file",
	"E",
	"E file",
	"f",
	"f file",
}

var pAddrCmds = []string{
	"=",
	"10=",
	"'x=",
	".=",
	"$=",
	"/123/=",
	`/12\/3/=`,
	"?123?=",
	".+1=",
	"12-=",
	"$-10=",
}

var pRangeCmds = []string{
	"d",
	".,$d",
	"3,d",
	",10d",
	"2;.+2d",
	";/hi/d",
	"/abc/;?123?d",
}

var pTextCmds = []string{
	"c\nhi\n.",
	"10,12c\nhello\nsome\ntext\n.",
	"c\n. some text .\n.",
	"3a\nhello\n.",
	"3i\nhello\n.",
}

var pMiscCmds = []string{
	"2kf",
	"w",
	"1,4w",
	"2w newfile",
}

func parse(cmd string) (cmds []Command, err error) {

	out := make(chan Command)
	go func() {
		defer close(out)
		p := &Parser{Buffer: cmd, Out: out}
		p.Init()
		err = p.Parse()
		if err != nil {
			return
		}
		p.Execute()
	}()
	for c := range out {
		fmt.Printf("cmd: %#v\n", c)
		cmds = append(cmds, c)
	}
	return
}

func TestSimpleCommandParser(t *testing.T) {
	for _, s := range psimpleCmds {
		cmds, err := parse(s + "\n")
		if err != nil {
			t.Error(err)
		}
		if len(cmds) != 1 {
			t.Errorf("Expected 1 command, got %d", len(cmds))
		}
	}
}

func TestParamCommandParser(t *testing.T) {
	for _, s := range pParamCmds {
		t.Run("paramParse:"+s, func(t *testing.T) {
			cmds, err := parse(s + "\n")
			if err != nil {
				t.Error(err)
			}
			if len(cmds) != 1 {
				t.Errorf("Expected 1 command, got %d", len(cmds))
			}

		})
	}
}

func TestAddrCommandParser(t *testing.T) {
	for _, s := range pAddrCmds {
		t.Run("addrParse:"+s, func(t *testing.T) {
			cmds, err := parse(s + "\n")
			if err != nil {
				t.Error(err)
			}
			if len(cmds) != 1 {
				t.Errorf("Expected 1 command, got %d", len(cmds))
			}

		})
	}
}

func TestRangeCommandParser(t *testing.T) {
	for _, s := range pRangeCmds {
		t.Run("rangeParse:"+s, func(t *testing.T) {
			cmds, err := parse(s + "\n")
			if err != nil {
				t.Error(err)
			}
			if len(cmds) != 1 {
				t.Errorf("Expected 1 command, got %d", len(cmds))
			}

		})
	}
}

func TestTextCommandParser(t *testing.T) {
	for _, s := range pTextCmds {
		t.Run("textParse:"+s, func(t *testing.T) {
			cmds, err := parse(s + "\n")
			if err != nil {
				t.Error(err)
			}
			if len(cmds) != 1 {
				t.Errorf("Expected 1 command, got %d", len(cmds))
			}

		})
	}
}

func TestMiscCommandParser(t *testing.T) {
	for _, s := range pMiscCmds {
		t.Run("miscParse:"+s, func(t *testing.T) {
			cmds, err := parse(s + "\n")
			if err != nil {
				t.Error(err)
			}
			if len(cmds) != 1 {
				t.Errorf("Expected 1 command, got %d", len(cmds))
			}

		})
	}
}
