package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type settings struct {
	ignoreTrailingWhitespace bool
	provideContext           bool
	contextSize              int
	edFormat                 bool
	fFormat                  bool
	recursive                bool
	unifiedContext           bool
	file1                    string
	file2                    string
}

func parseSettings(args []string) (settings, error) {
	s := settings{contextSize: 3}
	var err error

	for i := 0; i < len(args); i++ {
		a := args[i]

		switch a {
		case "-b":
			s.ignoreTrailingWhitespace = true
		case "-c":
			s.provideContext = true
		case "-C":
			if i+1 >= len(args) {
				return s, fmt.Errorf("ah heck I though that a number was gonna follow %s", a)
			}
			s.provideContext = true
			s.contextSize, err = strconv.Atoi(args[i+1])
			i++
		case "-e":
			s.edFormat = true
		case "-f":
			s.fFormat = true
		case "-r":
			s.recursive = true
		case "-u":
			s.unifiedContext = true
		case "-U":
			if i+1 >= len(args) {
				return s, fmt.Errorf("ah heck I though that a number was gonna follow %s", a)
			}
			s.unifiedContext = true
			s.contextSize, err = strconv.Atoi(args[i+1])
			if err != nil {
				return s, err
			}
			i++
		default:
			if s.file1 == "" {
				s.file1 = a
			} else if s.file2 == "" {
				s.file2 = a
			} else {
				return s, fmt.Errorf("aww jeez you can only compare two files, %s is the third", a)
			}
		}

	}

	if s.file1 == "" {
		return s, fmt.Errorf("oh no you forgot to give any files")
	}
	if s.file2 == "" {
		return s, fmt.Errorf("oh no you forgot to give the second file")
	}

	return s, nil
}

func readLines(fn string) ([]string, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	settings, err := parseSettings(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stdout, "Oh heck an error: %s\n", err)
		os.Exit(1)
	}

	l1, _ := readLines(settings.file1)
	l2, _ := readLines(settings.file2)
	changes := diff(l1, l2)

	fmt.Println(output(settings, changes))
}
