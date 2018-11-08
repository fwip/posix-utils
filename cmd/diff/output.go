package main

import (
	"fmt"
	"strings"
)

func fmtRange(start, length int) string {
	start++
	if length > 1 {
		return fmt.Sprintf("%d,%d", start, start+length-1)
	}
	return fmt.Sprintf("%d", start)
}

func output(s settings, changes []comparison) string {
	out := make([]string, 0)
	lineOld := 0
	lineNew := 0
	for i := 0; i < len(changes); i++ {
		c := changes[i]
		lineCount := len(c.values)
		switch c.kind {
		case equal:
			// Don't print anything, just keep track of where we are
			lineOld += lineCount
			lineNew += lineCount

		case add:
			out = append(out, fmt.Sprintf("%da%s", lineOld, fmtRange(lineNew, lineCount)))
			for _, l := range c.values {
				out = append(out, fmt.Sprintf("> %s", l))
			}
			lineNew += lineCount

		case minus:
			// Merge remove/add into a single change instruction
			if i+1 < len(changes) && changes[i+1].kind == add {
				c2 := changes[i+1]
				lineCount2 := len(c2.values)
				out = append(out, fmt.Sprintf("%sc%s",
					fmtRange(lineOld, lineCount),
					fmtRange(lineNew, lineCount2)))
				for _, l := range c.values {
					out = append(out, fmt.Sprintf("< %s", l))
				}
				out = append(out, "---")
				for _, l := range c2.values {
					out = append(out, fmt.Sprintf("> %s", l))
				}
				i++
				lineOld += lineCount
				lineNew += lineCount2

			} else {
				out = append(out, fmt.Sprintf("%sd%d", fmtRange(lineOld, lineCount), lineNew))
				for _, l := range c.values {
					out = append(out, fmt.Sprintf("< %s", l))
				}
				lineOld += lineCount
			}
		}
	}

	return strings.Join(out, "\n")
}
