package main

import (
	"io"
)

// A fileLike is something you can read from and write to.
type fileLike interface {
	io.ReadSeeker
	io.WriteCloser
}

type editor struct {
	input  chan command
	output chan string
	file   fileLike
}
