package ed

import (
	"strings"
)

type addressType byte

const (
	lCurrent addressType = iota
	lLast
	lNum
	lOffset //?
	lRegex
	lRegexReverse
	lMark
)

type address struct {
	typ  addressType
	text string
	//? offset int
}

func parseAddresses(cmd string) (start, end address, remainder string) {
	cmd = strings.TrimSpace(cmd)
	runes := []rune(cmd)

	if len(runes) == 0 {

	}

	//if len(runes) == 0 {
	//	// what do i return
	//	return
	//}
	//r := runes[0]
	//switch r {
	//case '.':
	//case '$':
	//case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
	//case '\'':
	//case '/':
	//case '?':
	//case '+':
	//case '-':
	//default:
	//	// what do i return
	//	return
	//}
	return
}

func parseAddress(runes []rune) address {
	if len(runes) == 0 {
		panic("idk what to do w/ 0-len addresses")
	}

	return address{}
}
