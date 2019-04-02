package main

import (
	"fmt"
	"os"
	"strconv"
)

type prefixTest func(in string) (bool, error)

func modeIs(fm os.FileMode) prefixTest {
	return func(in string) (bool, error) {
		fi, err := os.Stat(in)
		if err != nil {
			return false, nil
		}
		return (fi.Mode() & fm) == fm, nil
	}
}

func info(f func(os.FileInfo) bool) prefixTest {
	return func(in string) (bool, error) {
		fi, err := os.Stat(in)
		if err != nil {
			return false, nil
		}
		return f(fi), nil
	}
}

func noerr(test func(in string) bool) prefixTest {
	return func(in string) (bool, error) {
		return test(in), nil
	}
}

func nyi(in string) (bool, error) {
	return false, fmt.Errorf("NYI")
}
func isLink(in string) (bool, error) {
	i, err := os.Lstat(in)
	if err != nil {
		return false, nil
	}
	return i.Mode()&os.ModeSymlink == os.ModeSymlink, nil
}

var flagMap = map[byte]prefixTest{
	'a': modeIs(os.ModeDevice),
	'c': modeIs(os.ModeCharDevice),
	'd': info(func(i os.FileInfo) bool { return i.Mode().IsDir() }),
	'e': info(func(_ os.FileInfo) bool { return true }),
	'f': info(func(i os.FileInfo) bool { return i.Mode().IsRegular() }),
	'g': modeIs(os.ModeSetgid),
	'h': isLink,
	'L': isLink,
	'n': noerr(func(in string) bool { return in != "" }),
	'p': modeIs(os.ModeNamedPipe),
	'r': nyi,
	'S': modeIs(os.ModeSocket),
	's': info(func(i os.FileInfo) bool { return i.Size() > 0 }),
	't': nyi, // ????????
	'u': modeIs(os.ModeSetuid),
	'w': nyi,
	'x': nyi,
	'z': noerr(func(in string) bool { return in == "" }),
}

type infixTest func(a, b string) (bool, error)

func numcmp(cmp func(a, b int) bool) infixTest {
	return func(a, b string) (bool, error) {
		an, err := strconv.Atoi(a)
		if err != nil {
			return false, err
		}
		bn, err := strconv.Atoi(b)
		if err != nil {
			return false, err
		}
		return cmp(an, bn), nil
	}
}

var infixMap = map[string]infixTest{
	"=":   func(a, b string) (bool, error) { return a == b, nil },
	"!=":  func(a, b string) (bool, error) { return a != b, nil },
	"-eq": numcmp(func(a, b int) bool { return a == b }),
	"-ne": numcmp(func(a, b int) bool { return a != b }),
	"-gt": numcmp(func(a, b int) bool { return a > b }),
	"-ge": numcmp(func(a, b int) bool { return a >= b }),
	"-lt": numcmp(func(a, b int) bool { return a < b }),
	"-le": numcmp(func(a, b int) bool { return a <= b }),
}
