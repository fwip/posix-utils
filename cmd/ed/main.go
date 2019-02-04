package main

import (
	"fmt"
	"os"

	"github.com/fwip/posix-utils/ed"
)

func processCommand(cmd string) string {

	return ""
}

func main() {
	fmt.Println("its on")
	e := &ed.Itor{}
	e.ProcessCommands(os.Stdin, os.Stdout)
	os.Stdout.Close()
}
