package ed

import (
	"testing"
)

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
		//fmt.Printf("cmd: %#v\n", c)
		cmds = append(cmds, c)
	}
	return
}

func TestSimpleCommandParser(t *testing.T) {
	// These commands take no addresses and accept no parameters
	var psimpleCmds = []string{
		"h",
		"H",
		"P",
		"q",
		"Q",
		"u",
	}
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
	var pParamCmds = []string{"e",
		"e file",
		"E",
		"E file",
		"f",
		"f file",
	}
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
	var pRangeCmds = []string{
		"d",
		".,$d",
		"3,d",
		",10d",
		"2;.+2d",
		";/hi/d",
		"/abc/;?123?d",
	}
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
	var pTextCmds = []string{
		"c\nhi\n.",
		"10,12c\nhello\nsome\ntext\n.",
		"c\n. some text .\n.",
		"3a\nhello\n.",
		"3i\nhello\n.",
	}
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
	var pMiscCmds = []string{
		"2kf",
		"w",
		"1,4w",
		"2w newfile",
		"2a\nhi\n.",
		",p",
	}
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
