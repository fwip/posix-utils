package main

import "fmt"

type settings struct {
	filename  string
	numBytes  int
	numLines  int
	fromStart bool
	follow    bool
}

func parseSettings(args []string) (settings, error) {
	s := settings{
		filename: "-",
		numLines: 10,
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "-c":
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

func main() {
	fmt.Println("vim-go")
}
