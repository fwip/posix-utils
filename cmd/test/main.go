package main

import (
	"fmt"
	"os"
)

func main() {
	ok, err := parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(2)
	}
	fmt.Fprintf(os.Stderr, "%t\n", ok)
	if !ok {
		os.Exit(1)
	}
}

// Parse returns true if the expression is true, false if not true.
// err is non-nil if it couldn't parse
func parse(args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}
	if len(args) == 1 {
		return args[0] != "", nil
	}

	arg := args[0]
	if arg == "!" {
		ok, err := parse(args[1:])
		return !ok, err
	}
	if arg == "" {
		return false, nil
	}
	if len(arg) == 2 && arg[0] == '-' {
		if f, ok := flagMap[arg[1]]; ok {
			return f(args[1])
		}
	}

	if len(args) == 3 {
		if f, ok := infixMap[args[1]]; ok {
			return f(args[0], args[2])
		}
	}

	return false, fmt.Errorf("idk what to do about this")
}
