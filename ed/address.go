package ed

import (
	"strings"
)

type addressType byte

const (
	lNull addressType = iota
	lCurrent
	lLast
	lNum
	lOffset //?
	lRegex
	lRegexReverse
	lMark
)

type address struct {
	typ    addressType
	text   string
	offset int
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

var aCur = address{lCurrent, ".", 0}
var aFirst = address{lNum, "1", 0}
var aLast = address{lLast, "$", 0}
var aCurPlusOne = address{lCurrent, ".", 1}
var defaultAddresses = [...][2]address{
	{aCur, aCurPlusOne}, // ctnull
	{aCur, aCur},        // ctappend
	{aCur, aCur},        // ctchange
	{aCur, aCur},        // ctdelete
	{aCur, aCur},        // ctedit
	{aCur, aCur},        // cteditForce
	{aCur, aCur},        // ctfilename
	{aCur, aLast},       // ctglobal
	{aCur, aLast},       // ctinteractive
	{aCur, aCur},        // cthelp
	{aCur, aCur},        // cthelpMode
	{aCur, aCur},        // ctinsert
	{aCur, aCurPlusOne}, // ctjoin
	{aCur, aCur},        // ctmark
	{aCur, aCur},        // ctlist
	{aCur, aCur},        // ctmove
	{aCur, aCur},        // ctnumber
	{aCur, aCur},        // ctprint
	{aCur, aCur},        // ctprompt
	{aCur, aCur},        // ctquit
	{aCur, aCur},        // ctquitForce
	{aLast, aLast},      // ctread
	{aCur, aCur},        // ctsubstitute
	{aCur, aCur},        // ctcopy
	{aCur, aCur},        // ctundo
	{aCur, aLast},       // ctglobalInverse
	{aCur, aLast},       // ctinteractiveInverse
	{aCur, aLast},       // ctwrite
	{aLast, aLast},      // ctlineNumber
	{aCur, aCur},        // ctshell
}

func setDefaultAddresses(cmd Command) Command {

	if cmd.start.typ == lNull {
		cmd.start = defaultAddresses[cmd.typ][0]
		if cmd.end.typ == lNull {
			cmd.end = defaultAddresses[cmd.typ][1]
		}
	}
	return cmd
}
