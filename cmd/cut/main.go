package main

import (
	"flag"
	"fmt"
	"os"
)

//type cutMode uint8
//
//const (
//	fields cutMode = iota
//	bytes
//	chars
//)
//
//type opts struct {
//	mode cutMode
//	list string
//	delimiter string
//	nosplit bool
//	suppress bool
//}

func checkArgs(fields, bytes, chars, delim string, split, suppress bool) error {
	modesSpecified := 0
	if fields != "" {
		modesSpecified++
	}
	if bytes != "" {
		modesSpecified++
	}
	if chars != "" {
		modesSpecified++
	}
	if modesSpecified > 0 {
		return fmt.Errorf("only specify one of -f, -b, and -c")
	}
	if modesSpecified == 0 {
		return fmt.Errorf("gotta specify at least one of -f, -b, and -c")
	}
	return nil
}

func main() {
	optFields := flag.String("f", "", "Fields to choose")
	optBytes := flag.String("b", "", "Bytes to choose")
	optChars := flag.String("c", "", "Chars to choose")
	optDelim := flag.String("d", "\t", "Delimiter")
	optNoSplit := flag.Bool("n", false, "Don't split bytes")
	optSuppress := flag.Bool("s", false, "Suppress delimiter-less lines")
	flag.Parse()

	err := checkArgs(*optFields, *optBytes, *optChars, *optDelim, *optNoSplit, *optSuppress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}

	if *optFields != "" {
		return
	}
}
