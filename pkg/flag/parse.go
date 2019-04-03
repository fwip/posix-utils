package flag

import (
	"fmt"
	"os"
)

// Parser does the hard work of assigning parses
type Parser struct {
	Input   []string
	bools   map[rune]boolFlag
	strings map[rune]stringFlag
	ints    map[rune]intFlag
}

// ErrBadBundle indicates a misuse of bundling options
type ErrBadBundle error

// Parse does the parsing, returns the positional arguments and any error that occurred
// If Input is nil, Parse() will use os.Args[1:]
func (p *Parser) Parse() (positionalArgs []string, err error) {
	if p.Input == nil {
		p.Input = os.Args[1:]
	}
	for i := 0; i < len(p.Input); i++ {
		arg := p.Input[i]

		// Non-argument ends options
		if len(arg) == 0 || arg[0] != '-' || arg == "-" {
			return p.Input[i:], nil
		}
		// '--' explicitly ends options
		if arg == "--" {
			return p.Input[i+1:], nil
		}

		// Bundling is allowed, so check each character to apply it
		for j, c := range arg[1:] {
			// Check for boolean flag.
			if flag, ok := p.bools[c]; ok {
				flag.execute()
				continue
			}

			// Argument may be either bundled or separate
			value := arg[j+2:]
			if len(value) == 0 {
				i++
				value = p.Input[i]
			}

			if flag, ok := p.strings[c]; ok {
				flag.execute(value)
				continue
			}
			if flag, ok := p.ints[c]; ok {
				err := flag.execute(value)
				if err != nil {
					return nil, err
				}
				continue
			}
			return nil, fmt.Errorf("man what's happening")
			//// non-boolean flags may only be bundled if they are the last
			//if j != len(arg)-2 {
			//	return nil, ErrBadBundle(fmt.Errorf("inappropriate bundling of %c in %s", c, arg))
			//}
			//i++
			//value := p.Input[i]

			//if flag, ok := p.strings[c]; ok {
			//	flag.execute(value)
			//} else if flag, ok := p.ints[c]; ok {
			//	err := flag.execute(value)
			//	if err != nil {
			//		return nil, err
			//	}
			//}
		}
	}
	return nil, nil
}

// BoolVar creates a new boolean variable
func (p *Parser) BoolVar(addr *bool, name rune) {
	if p.bools == nil {
		p.bools = make(map[rune]boolFlag)
	}
	p.bools[name] = boolFlag{addr}
}

// IntVar creates a new int64 variable
func (p *Parser) IntVar(addr *int64, name rune) {
	if p.ints == nil {
		p.ints = make(map[rune]intFlag)
	}
	p.ints[name] = intFlag{addr}
}

// StringVar creates a new int64 variable
func (p *Parser) StringVar(addr *string, name rune) {
	if p.strings == nil {
		p.strings = make(map[rune]stringFlag)
	}
	p.strings[name] = stringFlag{addr}
}
