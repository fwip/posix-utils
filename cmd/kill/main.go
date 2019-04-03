package main

// TODO: This functionality appears bugged on Darwin. Check out in non-OSX platform

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var signalMap = map[string]int{
	"0":    0,
	"HUP":  1,
	"INT":  2,
	"QUIT": 3,
	"ABRT": 6,
	"KILL": 9,
	"ALRM": 14,
	"TERM": 15,
}

func errorMsg(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(2)
}

func main() {

	listSignals := false
	signal := 15
	var pid int
	var err error
	var ok bool
	for i := 1; i < len(os.Args); i++ {
		a := os.Args[i]
		switch a {
		case "-l":
			listSignals = true
		case "-s":
			i++
			signalName := strings.ToUpper(os.Args[i])
			signal, ok = signalMap[signalName]
			if !ok {
				errorMsg("%s is not a valid signal name - check kill -l for help\n", a)
			}
		default:
			if a[0] == '-' {
				signal, err = strconv.Atoi(a[1:])
				if err != nil {
					errorMsg("%s isn't a number, bro...", a)
				}
			} else {
				pid, err = strconv.Atoi(a)
				if err != nil {
					errorMsg("expected a PID here..., got %s", a)
				}
			}
		}
	}

	if listSignals {

	}
	fmt.Println("sending", signal, "to", pid)
	process, err := os.FindProcess(pid)
	if err != nil {
		errorMsg("couldn't find process with PID of %d: %s\n", pid, err)
	}
	//err = process.Signal(syscall.Signal(signal))
	err = process.Kill()
	if err != nil {
		errorMsg("couldn't send signal %d to %d, %s", signal, pid, err)
	}

	fmt.Println("vim-go")
}
