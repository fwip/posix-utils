package ed

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fwip/posix-utils/pkg/txt"
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
	ed = &Itor{pt: txt.NewPieceTable(&bytes.Reader{}, 0)}
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
	ed.currentLine = len(ed.getLines())

	return nil
}

// Print prints some lines
func (ed *Itor) Print(start, end int) string {
	lines := ed.getLines()
	return strings.Join(lines[start-1:end], "\n")
}

func (ed *Itor) number(start, end int) string {
	lines := ed.getLines()
	out := ""
	for i := start; i <= end; i++ {
		out += fmt.Sprintf("%d\t%s\n", i, lines[i-1])
	}
	return out
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

// Delete will delete from the starting line to the end line
// This is inclusive, so Delete(3, 3) will delete the fourth line of the buffer
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
	lines := ed.getLines()
	if num > len(lines) {
		for _, l := range lines {
			length += len(l) + 1
		}
		return length
	}
	for i := 0; i < num-1; i++ {
		length += len(lines[i]) + 1
	}
	return length
}

func (ed *Itor) getLines() (lines []string) {
	// Retrieve the piece table's output and tokenize into lines
	s := bufio.NewScanner(strings.NewReader(ed.pt.String()))
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if s.Err() != nil {
		panic(s.Err())
	}

	// If there's an empty line at the end, remove it
	// This happens when the file ends in a newline character
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
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

	ed.currentLine = len(ed.getLines())
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
	var line int
	switch a.typ {
	case lCurrent:
		line = ed.currentLine
	case lNum:
		n, _ := strconv.Atoi(a.text)
		line = n
	case lLast:
		n := len(ed.getLines())
		line = n
	default:
		panic("Dunno a bout address")
	}

	return line + a.offset
}

func linesIn(s string) int {
	n := 1
	for _, r := range s {
		if r == '\n' {
			n++
		}
	}
	return n
}

func (ed *Itor) processCommand(cmd Command) string {
	var err error

	cmd = setDefaultAddresses(cmd)
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

	case ctnumber:
		return ed.number(ed.addrLine(cmd.start), ed.addrLine(cmd.end))
	case ctnull:
		lineNum := ed.addrLine(cmd.start)
		ed.currentLine = lineNum
		return ed.Print(lineNum, lineNum)

	case ctwrite:
		err = ed.Write()
		if err != nil {
			return "?" + err.Error()
		}

	case ctdelete:
		err := ed.Delete(ed.addrLine(cmd.start), ed.addrLine(cmd.end))
		if err != nil {
			return "?" + err.Error()
		}
		ed.currentLine = ed.addrLine(cmd.start)

	case ctappend:
		start := ed.addrLine(cmd.start) + 1
		ed.insertBeforeLine(start, cmd.text)
		ed.currentLine = start + linesIn(cmd.text)

	case ctinsert:
		ed.insertBeforeLine(ed.addrLine(cmd.start), cmd.text)

	case ctchange:
		err := ed.Delete(ed.addrLine(cmd.start), ed.addrLine(cmd.end))
		if err != nil {
			return "?" + err.Error()
		}
		ed.insertBeforeLine(ed.addrLine(cmd.start), cmd.text)
		ed.currentLine = ed.addrLine(cmd.start) + linesIn(cmd.text)

	case ctlineNumber:
		return strconv.Itoa(ed.addrLine(cmd.start))

	default:
		return "? (NYI)"
	}
	return ""
}
