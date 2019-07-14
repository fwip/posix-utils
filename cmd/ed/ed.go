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
	input  io.Reader
	output io.Writer
	file   fileLike
}
