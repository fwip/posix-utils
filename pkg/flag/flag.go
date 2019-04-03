package flag

import "strconv"

// TODO: Is this right?
//type flag interface {
//execute(s string) error
//}

type stringFlag struct {
	dest *string
}

type boolFlag struct {
	dest *bool
}

type intFlag struct {
	dest *int64
}

// ErrNonInt indicates a non-integer value passed to an integer flag
type ErrNonInt error

func (bflag boolFlag) execute() {
	*(bflag.dest) = true
}
func (iflag intFlag) execute(s string) error {
	x, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return ErrNonInt(err)
	}
	*(iflag.dest) = x
	return nil
}

func (sflag stringFlag) execute(s string) {
	*(sflag.dest) = s
}
