package ed

import (
	"fmt"
	"strings"
)

type cmdType byte

const (
	ctnull cmdType = iota
	ctappend
	ctchange
	ctdelete
	ctedit
	cteditForce
	ctfilename
	ctglobal
	ctinteractive
	cthelp
	cthelpMode
	ctinsert
	ctjoin
	ctmark
	ctlist
	ctmove
	ctnumber
	ctprint
	ctprompt
	ctquit
	ctquitForce
	ctread
	ctsubstitute
	ctcopy
	ctundo
	ctglobalInverse
	ctinteractiveInverse
	ctwrite
	ctlineNumber
	ctshell
)

var cmdMap = map[rune]struct {
	typ          cmdType
	numAddresses int
	takesParam   bool
	takesRegex   bool
	takesText    bool
}{
	'a': {ctappend, 1, false, false, true},             // TODO
	'c': {ctchange, 2, false, false, true},             // TODO
	'd': {ctdelete, 2, false, false, false},            // TODO
	'e': {ctedit, 0, true, false, false},               // TODO
	'E': {cteditForce, 0, true, false, false},          // TODO
	'f': {ctfilename, 0, true, false, false},           // TODO
	'g': {ctglobal, 2, false, true, false},             // TODO
	'G': {ctinteractive, 2, false, true, false},        // TODO
	'h': {cthelp, 0, false, false, false},              // Good
	'H': {cthelpMode, 0, false, false, false},          // Good
	'i': {ctinsert, 1, false, false, true},             // TODO
	'j': {ctjoin, 2, false, false, false},              // TODO
	'k': {ctmark, 1, true, false, false},               // TODO
	'l': {ctlist, 2, false, false, false},              // TODO
	'm': {ctmove, 2, true, false, false},               // TODO
	'n': {ctnumber, 2, false, false, false},            // TODO
	'p': {ctprint, 2, false, false, false},             // TODO
	'P': {ctprompt, 0, false, false, false},            // Good
	'q': {ctquit, 0, false, false, false},              // Good
	'Q': {ctquitForce, 0, false, false, false},         // Good
	'r': {ctread, 1, true, false, false},               // TODO
	's': {ctsubstitute, 2, false, true, false},         // TODO
	't': {ctcopy, 2, true, false, false},               // TODO
	'u': {ctundo, 0, false, false, false},              // Good
	'v': {ctglobalInverse, 2, false, false, true},      // TODO
	'V': {ctinteractiveInverse, 2, false, false, true}, // TODO
	'w': {ctwrite, 2, true, false, false},              // TODO
	'=': {ctlineNumber, 1, false, false, false},        // TODO
	'!': {ctshell, 0, true, false, false},              // TODO
}

// Command is magick
type Command struct {
	typ    cmdType
	start  address
	end    address
	dest   address
	text   string
	params []string
}

//type address struct {
//	num  int64
//	text string
//}

// Execute a command for ed
func (ed *Itor) Execute(cmd string) string {
	words := strings.Fields(cmd)
	c := words[0]
	var err error
	var out string
	switch c {
	case "=":
		out = fmt.Sprintf("%d", ed.currentLine)
	case "e":
		err = ed.Edit(words[1], false)
	case "E":
		err = ed.Edit(words[1], true)
	case "w":
		err = ed.Write()
	default:
		return "?"
	}
	if err != nil {
		return "?"
	}
	return out
}

//func (ed *Itor) processInput(r io.Reader, w io.Writer) {
//	s := bufio.NewScanner(r)
//	for s.Scan() {
//		out := ed.Execute(s.Text())
//		_, err := w.Write([]byte(out))
//		if err != nil {
//			// TODO: robusitfy
//			panic(err)
//		}
//	}
//	if s.Err() != nil {
//		panic(s.Err())
//	}
//}
//
//func parseCommands(r io.Reader, out chan Command) error {
//	s := bufio.NewScanner(r)
//	for s.Scan() {
//		cmd := s.Text()
//		p := &Parser{Buffer: cmd, Out: out}
//		p.Init()
//		err := p.Parse()
//		if err != nil {
//			close(out)
//			panic("can't parse that baby")
//		}
//		p.Execute()
//
//		fmt.Println("---")
//
//		//cmd := strings.TrimSpace(s.Text())
//		//var simpleType cmdType
//		//// Simple commands
//		//switch cmd {
//		//case "e":
//		//	simpleType = ctedit
//		//case "E":
//		//	simpleType = cteditForce
//		//case "f":
//		//	simpleType = ctfilename
//		//case "h":
//		//	simpleType = cthelp
//		//case "H":
//		//	simpleType = cthelpMode
//		//}
//		//if simpleType != ctnull {
//		//	out <- command{typ: simpleType}
//		//	continue
//		//}
//
//		//return fmt.Errorf("can't understand command %s", cmd)
//	}
//	close(out)
//	return s.Err()
//}
