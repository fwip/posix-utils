package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type settings struct {
	numLines  int
	filenames []string
}

func parseSettings(args []string) (settings, error) {
	s := settings{
		numLines: 10,
		//filenames: []string{"-"},
	}

	if len(args) == 0 {
		return s, nil
		//return settings{}, fmt.Errorf("No arguments")
	}

	for i := 0; i < len(args); i++ {
		if args[i] == "-n" {
			i++
			n, err := strconv.Atoi(args[i])
			if err != nil {
				return s, err
			}
			s.numLines = n
		} else {
			s.filenames = append(s.filenames, args[i])
		}
	}

	if len(s.filenames) == 0 {
		s.filenames = []string{"-"}
	}

	return s, nil
}

func process(s settings) {

	for i, filename := range s.filenames {
		if len(s.filenames) > 1 {
			if i > 0 {
				fmt.Print("\n")
			}
			fmt.Printf("==> %s <==\n", filename)
		}
		var f io.ReadCloser
		var err error
		if filename == "-" {
			f = os.Stdin
		} else {
			f, err = os.Open(filename)
			if err != nil {
				panic(err)
			}
		}

		scan := bufio.NewScanner(f)
		j := 0
		for scan.Scan() && j < s.numLines {
			fmt.Println(scan.Text())
			j++
		}
		if scan.Err() != nil {
			panic(scan.Err())
		}
	}

}

func main() {
	s, err := parseSettings(os.Args[1:])
	if err != nil {
		panic(err)
	}
	process(s)
}
