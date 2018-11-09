package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type settings struct {
	writeAll bool
	silent   bool
	file1    string
	file2    string
}

func parseSettings(args []string) (settings, error) {
	s := settings{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "-l":
			s.writeAll = true
		case "-s":
			s.silent = true
		default:
			if s.file1 == "" {
				s.file1 = a
			} else if s.file2 == "" {
				s.file2 = a
			} else {
				return s, fmt.Errorf("too many files: %s", a)
			}
		}
	}
	return s, nil
}

func open(filename string) (io.ReadCloser, error) {
	if filename == "-" {
		return os.Stdin, nil
	}
	return os.Open(filename)
}

func cmp(s settings) error {
	f1, err := open(s.file1)
	if err != nil {
		return fmt.Errorf("could not open %s", s.file1)
	}
	defer f1.Close()
	f2, err := open(s.file2)
	if err != nil {
		return fmt.Errorf("could not open %s", s.file2)
	}
	defer f2.Close()

	buf1 := bufio.NewReader(f1)
	buf2 := bufio.NewReader(f2)
	l := 1
	for i := 1; err == nil; i++ {
		b1, err1 := buf1.ReadByte()
		b2, err2 := buf2.ReadByte()

		if err1 == nil && err2 == nil {
			if b1 != b2 {
				if !s.silent {
					if s.writeAll {
						fmt.Printf("%d %o %o\n", i, b1, b2)
					} else {
						fmt.Printf("%s %s differ: char %d, line %d\n", s.file1, s.file2, i, l)
					}
				}
				if !s.writeAll {
					os.Exit(1)
				}
			}
			if b1 == byte('\n') {
				l++
			}
			continue
		}
		if err1 == io.EOF && err2 == io.EOF {
			break
		}
		if err1 != nil {
			if err1 == io.EOF {
				fmt.Fprintf(os.Stderr, "cmp: EOF on %s%s\n", s.file1, "")
				os.Exit(1)
			}
			return fmt.Errorf("Error reading %s: %s", s.file1, err1)
		}
		if err2 != nil {
			if err2 == io.EOF {
				fmt.Fprintf(os.Stderr, "cmp: EOF on %s%s\n", s.file2, "")
				os.Exit(1)
			}
			return fmt.Errorf("Error reading %s: %s", s.file2, err2)
		}

	}

	return nil
}

func main() {

	s, err := parseSettings(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(2)
	}
	err = cmp(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(2)
	}
}
