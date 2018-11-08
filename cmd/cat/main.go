package main

import (
	//"bufio"
	"fmt"
	"io"
	"os"
)

// TODO: Support environment variables

func write(r io.Reader, w io.Writer) (err error) {
	_, err = io.Copy(w, r)
	return err
}

func writeFile(filename string) error {
	var r io.ReadCloser
	var err error
	// Special filename '-' means stdin
	if filename == "-" {
		r = os.Stdin
	} else {
		r, err = os.Open(filename)
		if err != nil {
			return err
		}
		defer r.Close()
	}

	return write(r, os.Stdout)
}

func main() {
	filenames := make([]string, 0)
	for _, s := range os.Args[1:] {
		// "-u" is an POSIX option that forces unbuffered output. Output is already unbuffered.
		if s != "-u" {
			filenames = append(filenames, s)
		}
	}

	// If called with no arguments, read from stdin
	if len(filenames) == 0 {
		filenames = append(filenames, "-")
	}

	for _, fn := range filenames {
		err := writeFile(fn)
		if err != nil {
			os.Stderr.Write([]byte(fmt.Sprintf("Error encountered while attempting to read from %s: %s\n", fn, err)))
			break
		}
	}
}
