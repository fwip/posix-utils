package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func run(filename string) {
	reader := os.Stdin
	if filename != "-" {
		f, err := os.Open(filename)
		if err != nil {
			die("couldn't open file %s", err)
		}
		defer f.Close()
		reader = f
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		die("couldn't read input %s", err)
	}

	fmt.Printf("%d %d %s\n", memcrc(bytes), len(bytes), filename)
}

func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		files = append(files, "-")
	}
	for _, f := range files {
		run(f)
	}
}
