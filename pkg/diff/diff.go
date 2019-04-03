package diff

import (
	"fmt"
	"strings"
)

const (
	equal = iota
	add
	minus
)

var symbols = map[int]string{
	equal: " ",
	add:   "+",
	minus: "-",
}

type comparison struct {
	kind   int
	values []string
}

func (c comparison) String() string {
	output := make([]string, 0, len(c.values))
	symbol := symbols[c.kind]
	for _, v := range c.values {
		output = append(output, fmt.Sprintf("%s %s", symbol, v))
	}
	return strings.Join(output, "\n")
}

func Diff(old, new []string) []comparison {
	if len(old) == 0 {
		if len(old) == len(new) {
			return []comparison{}
		}
		return []comparison{{kind: add, values: new}}
	}
	if len(new) == 0 {
		return []comparison{{kind: minus, values: old}}
	}

	oldIndices := make(map[string][]int, 0)
	for i, k := range old {
		oldIndices[k] = append(oldIndices[k], i)
	}

	overlap := make(map[int]int, 0)
	subStartOld := 0
	subStartNew := 0
	subLen := 0
	for inew, val := range new {
		newOverlap := make(map[int]int, 0)
		for _, iold := range oldIndices[val] {
			if iold == 0 {
				newOverlap[iold] = 1
			} else {
				newOverlap[iold] = overlap[iold-1] + 1
			}
			//newOverlap[iold] = (iold) // ????
			if newOverlap[iold] > subLen {
				subLen = newOverlap[iold]
				subStartOld = iold - subLen + 1
				subStartNew = inew - subLen + 1
			}
		}
		overlap = newOverlap
	}

	if subLen == 0 {
		return []comparison{
			{
				kind:   minus,
				values: old,
			},
			{kind: add,
				values: new,
			},
		}

	}
	comparisons := Diff(old[:subStartOld], new[:subStartNew])
	comparisons = append(comparisons, comparison{
		kind:   equal,
		values: new[subStartNew : subStartNew+subLen],
	})
	comparisons = append(comparisons, Diff(old[subStartOld+subLen:], new[subStartNew+subLen:])...)
	return comparisons
}
