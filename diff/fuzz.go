// +build gofuzz

package diff

import (
	"strings"
)

func assembleChanges(changes []comparison) (old, new []string) {
	for _, c := range changes {
		for _, l := range c.values {
			switch c.kind {
			case add:
				new = append(new, l)
			case minus:
				old = append(old, l)
			case equal:
				old = append(old, l)
				new = append(new, l)
			default:
				panic("Unexpected change type")
			}
		}
	}
	return
}
func areIdentical(expected, actual []string) bool {
	if len(expected) != len(actual) {
		return false
	}
	for i := range expected {
		if expected[i] != actual[i] {
			return false
		}
	}
	return true
}

// Fuzz tests for basic correctness
func Fuzz(data []byte) int {
	dataStr := string(data)
	inputs := strings.Split(dataStr, "\n==========\n")
	if len(inputs) != 2 {
		return -1
	}
	old := strings.Split(inputs[0], "\n")
	new := strings.Split(inputs[1], "\n")
	changes := Diff(old, new)

	actualOld, actualNew := assembleChanges(changes)

	if !areIdentical(old, actualOld) {
		panic("Did not reconstruct 'old' input")
	}
	if !areIdentical(new, actualNew) {
		panic("Did not reconstruct 'new' input")
	}
	return 1
}
