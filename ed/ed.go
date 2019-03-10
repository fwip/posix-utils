package ed

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fwip/posix-utils/txt"
)

// TODO: This is not comprehensive.
var multlineCmdStart = regexp.MustCompile(`^\s*[0-9.$]*\s*,?\s*[0-9.$]*\s*[aic]\s*$`)

func debug(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

// Itor is an edItor (get it) instance. It supports one open file.
type Itor struct {
	pt          txt.PieceTable
	filename    string
	currentLine int
	modified    bool
}

// NewEditor creates a new editor that reads and writes to the supplied writer
func NewEditor(in io.Reader, out io.Writer) (ed *Itor) {
	ed = &Itor{}
	//go ed.processInput(in, out)
	return ed
}

func (ed *Itor) unsaved() bool {
	return false
}

func (ed *Itor) closeFile() error {
	ed.pt.Close()
	return nil
}

func (ed *Itor) Write() error {
	data := ed.pt.String()
	ed.closeFile()

	var err error

	tmpName := ed.filename + ".swp"
	fmt.Printf("f is '%s', tmp is '%s'\n", ed.filename, tmpName)

	stat, err := os.Stat(ed.filename)
	if err != nil {
		debug("cannae stat", ed.filename)
		return err
	}
	mode := stat.Mode()
	//defer os.Remove(tmp.Name())

	err = ioutil.WriteFile(tmpName, []byte(data), os.FileMode(mode))
	if err != nil {
		os.Remove(tmpName)
		return err
	}
	err = os.Rename(tmpName, ed.filename)
	if err != nil {
		os.Remove(tmpName)
		return err
	}

	return ed.Edit(ed.filename, false)
}

// Edit opens a new file to edit!
func (ed *Itor) Edit(filename string, force bool) error {
	debug("editing", filename)
	if !force && ed.unsaved() {
		panic("oh no you haven't saved")
	}
	f, err := os.Open(filename)
	if err != nil {
		debug("cannae open", err)
		return err
	}
	stat, err := f.Stat()
	if err != nil {
		debug("cannae stat", err)
		return err
	}

	debug("setting filename")
	ed.filename = filename
	ed.pt = txt.NewPieceTable(f, int(stat.Size()))

	return nil
}

// Print prints some lines
func (ed *Itor) Print(start, end int) string {
	lines := ed.getLines()
	return strings.Join(lines[start-1:end], "\n")
}

// String returns the whole buffer
func (ed *Itor) String() string {
	s := ed.pt.String()
	// POSIX requires a file to end with a newline
	if s[len(s)-1] != '\n' {
		s += "\n"
	}
	return s
}

func (ed *Itor) insertBeforeLine(lineNum int, text string) {
	at := ed.getLineAddr(lineNum)
	ed.pt.Insert([]byte(text+"\n"), at)
}

func (ed *Itor) Delete(start, end int) error {
	realStart := ed.getLineAddr(start)
	realEnd := ed.getLineAddr(end + 1)
	debug("delete", start, end, realStart, realEnd)
	ed.pt.Delete(realEnd-realStart, realStart)
	return nil
}

// 1-indexed
func (ed *Itor) getLineAddr(num int) int {
	length := 0
	for i, line := range ed.getLines() {
		if i == num-1 {
			return length
		}
		length += len(line) + 1
	}
	return length
	//panic("too much")
}

func (ed *Itor) getLines() []string {
	return strings.Split(ed.pt.String(), "\n")
}

// Quit quits the editor
func (ed *Itor) Quit(force bool) error {
	if !force && ed.unsaved() {
		return fmt.Errorf("Unsaved, can't quit")
	}
	return nil
}

// ProcessCommands is the main entry point for an interpreter
// Blocks until all commands have been processed
func (ed *Itor) ProcessCommands(r io.Reader, w io.Writer) {

	cmds := make(chan Command)

	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			cmd := s.Text()

			// Handle multi-line commands
			if multlineCmdStart.Match([]byte(cmd)) {
				for s.Scan() {
					cmd += "\n" + s.Text()
					if s.Text() == "." {
						break
					}
				}
			}

			debug("line", cmd)
			p := &Parser{Buffer: cmd + "\n", Out: chan<- Command(cmds)}
			p.Init()
			err := p.Parse()
			if err != nil {
				//fmt.Println("Err!", err)
				w.Write([]byte("?" + err.Error() + "\n"))
				//close(cmds)
				continue
			}
			p.Execute()
		}
		if s.Err() != nil {
			fmt.Println("Err:", s.Err())
		}
		close(cmds)
	}()

	for cmd := range cmds {
		fmt.Println("cmd", cmd)
		// Special-case quit command for now
		if cmd.typ == ctquit {
			break // TODO: checking
		}
		result := ed.processCommand(cmd)
		w.Write([]byte(result + "\n"))
	}
}

func (ed *Itor) addrLine(a address) int {
	switch a.typ {
	case lCurrent:
		return ed.currentLine
	case lNum:
		n, _ := strconv.Atoi(a.text)
		return n
	case lLast:
		n := len(ed.getLines())
		return n
	}
	panic("Dunno a bout address")
}

func (ed *Itor) processCommand(cmd Command) string {
	var err error
	switch cmd.typ {
	case ctedit:
		var filename string
		if len(cmd.params) > 0 {
			filename = cmd.params[0]
		}
		err = ed.Edit(filename, false)
		if err != nil {
			return err.Error()
		}

	case ctprint:
		return ed.Print(ed.addrLine(cmd.start), ed.addrLine(cmd.end))

	case ctwrite:
		fmt.Println("writing")
		err = ed.Write()
		if err != nil {
			return "?" + err.Error()
		}

	case ctdelete:
		err := ed.Delete(ed.addrLine(cmd.start), ed.addrLine(cmd.end))
		if err != nil {
			return "?" + err.Error()
		}

	case ctappend:
		ed.insertBeforeLine(ed.addrLine(cmd.start)+1, cmd.text)

	case ctinsert:
		ed.insertBeforeLine(ed.addrLine(cmd.start), cmd.text)

	case ctchange:
		err := ed.Delete(ed.addrLine(cmd.start), ed.addrLine(cmd.end))
		if err != nil {
			return "?" + err.Error()
		}
		ed.insertBeforeLine(ed.addrLine(cmd.start), cmd.text)

	default:
		return "? (NYI)"
	}
	return ""
}
