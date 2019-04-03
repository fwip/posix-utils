package flag

import (
	"fmt"
	"strings"
)

func Example() {
	input := "-nrk 2 -x hello from ExampleLand"
	p := Parser{Input: strings.Fields(input)}

	var n, r, z bool
	var k int64
	var x string
	p.BoolVar(&n, 'n')
	p.BoolVar(&r, 'r')
	p.BoolVar(&z, 'z')
	p.IntVar(&k, 'k')
	p.StringVar(&x, 'x')

	args, _ := p.Parse()

	fmt.Println("n:", n)
	fmt.Println("r:", r)
	fmt.Println("z:", z)
	fmt.Println("k:", k)
	fmt.Println("x:", x)
	fmt.Println("args:", args)

	// Output:
	// n: true
	// r: true
	// z: false
	// k: 2
	// x: hello
	// args: [from ExampleLand]
}
