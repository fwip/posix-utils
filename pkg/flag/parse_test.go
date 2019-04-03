package flag

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

const (
	tBool = iota
	tString
	tInt
)

type parsetest struct {
	name           string
	args           string
	flags          []flagtest
	positionalArgs string
}

type flagtest struct {
	name rune
	typ  int
	val  string
}

var parsetests = []parsetest{
	{"empty", "", nil, ""},
	{"positional", "hello world", nil, "hello world"},
	{"boolean", "-v", []flagtest{{name: 'v'}}, ""},
	{"bundled", "-vx", []flagtest{{name: 'v'}, {name: 'x'}}, ""},
	{"-- ends", "-v -- -p hello", []flagtest{{name: 'v'}, {name: 'p', typ: tString, val: ""}}, "-p hello"},
	{"string", "-p hello world", []flagtest{{'p', tString, "hello"}}, "world"},
	{"int", "-n 32 world", []flagtest{{'n', tInt, "32"}}, "world"},
	{"bundled_complex", "-vxp hi", []flagtest{{name: 'v'}, {name: 'x'}, {'p', tString, "hi"}}, ""},
	//{"unset", "-vx", []flagtest{{name: 'v'}, {name: 'x'}, {name: 'y'}}, ""},
}

func TestParse(t *testing.T) {
	for _, pt := range parsetests {
		t.Run("test parse:"+pt.name, func(t *testing.T) {
			args := strings.Fields(pt.args)
			p := Parser{Input: args}

			mbool := make(map[rune]*bool)
			mint := make(map[rune]*int64)
			mstring := make(map[rune]*string)
			for _, f := range pt.flags {
				switch f.typ {
				case tBool:
					var x bool
					mbool[f.name] = &x
					p.BoolVar(&x, f.name)
				case tInt:
					var x int64
					mint[f.name] = &x
					p.IntVar(&x, f.name)
				case tString:
					var x string
					mstring[f.name] = &x
					p.StringVar(&x, f.name)
				}
			}

			pos, err := p.Parse()

			for _, f := range pt.flags {
				switch f.typ {
				case tBool:
					if !*mbool[f.name] {
						t.Error(fmt.Errorf("-%c should be true", f.name))
					}
				case tString:
					if *mstring[f.name] != f.val {
						t.Error(fmt.Errorf("-%c expected %q, got %q", f.name, f.val, *mstring[f.name]))
					}
				case tInt:
					x, _ := strconv.ParseInt(f.val, 0, 64)
					if *mint[f.name] != x {
						t.Error(fmt.Errorf("-%c expected %d, got %d", f.name, x, *mint[f.name]))
					}
				}
			}

			if err != nil {
				t.Error(err)
			}
			if strings.Join(pos, " ") != pt.positionalArgs {
				t.Error(fmt.Errorf("positional args : expected %q got %q", pt.positionalArgs, pos))
			}

		})
	}
}

func expectError(err error, typ reflect.Type, t *testing.T) {
	if err == nil {
		t.Errorf("no error thrown")
	}
	if reflect.TypeOf(err) != typ {
		t.Errorf("error of wrong type, expected %q got %q", typ, reflect.TypeOf(err))
	}
}

//func TestInappropriateBundling(t *testing.T) {
//	p := Parser{Input: strings.Fields("-ajx hello")}
//	var a bool
//	var j string
//	var x string
//
//	p.BoolVar(&a, 'a')
//	p.StringVar(&j, 'j')
//	p.StringVar(&x, 'x')
//
//	_, err := p.Parse()
//	if err == nil {
//		t.Errorf("no error thrown")
//	}
//	if _, ok := err.(ErrBadBundle); !ok {
//		t.Errorf("error of wrong type, expected ErrBadBundle got %q", reflect.TypeOf(err))
//	}
//}

func TestNonInt(t *testing.T) {
	p := Parser{Input: strings.Fields("-a hello")}
	var a int64
	p.IntVar(&a, 'a')
	_, err := p.Parse()
	if err == nil {
		t.Errorf("no error thrown")
	}
	if _, ok := err.(ErrNonInt); !ok {
		t.Errorf("error of wrong type, expected ErrNonInt, got %q", reflect.TypeOf(err))
	}

}
