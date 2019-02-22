package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/fwip/posix-utils/tsort"
)

func die(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "ahh jeez, "+msg+"\n", args)
	os.Exit(1)
}

func run(input io.Reader, output io.Writer) {
	var sorter tsort.Sorter
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		w1 := scanner.Text()

		if !scanner.Scan() {
			die("unmatched word! %s", w1)
		}
		w2 := scanner.Text()
		if w1 == w2 {
			sorter.Add([]string{w1})
		} else {
			sorter.Add([]string{w1, w2})
		}
	}
	if scanner.Err() != nil {
		die("scan error! %s", scanner.Err())
	}

	out, err := sorter.Order()
	if err != nil {
		die("sort error! %s", err)
	}
	for _, item := range out {
		fmt.Fprintln(output, item)
	}

	if closer, ok := output.(io.Closer); ok {
		closer.Close()
	}

}

func main() {
	var input io.Reader = os.Stdin
	if len(os.Args) > 1 {
		fname := os.Args[1]
		f, err := os.Open(fname)
		if err != nil {
			die("can't open %s", fname)
		}
		defer f.Close()
		input = f
	}

	run(input, os.Stdout)
}
