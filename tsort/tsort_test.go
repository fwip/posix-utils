package tsort

import (
	"fmt"
	"strings"
	"testing"
)

var tests = [][]string{
	[]string{"a b c", "b e"},
	[]string{"j", "a b c", "a x", "b x", "a y b", "z"},
}

func indexOf(haystack []string, needle string) (int, error) {

	for i, h := range haystack {
		if needle == h {
			return i, nil
		}
	}
	return 0, fmt.Errorf("needle not found")
}

func validateOrdering(ordering []string, partials []string) error {

	for _, partial := range partials {
		var position = -1
		for _, f := range strings.Fields(partial) {
			index, err := indexOf(ordering, f)
			if err != nil {
				return fmt.Errorf("%s not in ordering", f)
			}
			if index <= position {
				return fmt.Errorf("%s too soon", f)
			}
			position = index
		}
	}
	return nil
}

func containsDuplicates(items []string) bool {
	seen := make(map[string]struct{})
	for _, item := range items {
		if _, in := seen[item]; in {
			return true
		}
		seen[item] = struct{}{}
	}
	return false
}

func TestTsort(t *testing.T) {
	for _, test := range tests {
		sorter := Sorter{}
		for _, ordering := range test {
			sorter.Add(strings.Fields(ordering))
		}
		order, err := sorter.Order()
		if err != nil {
			t.Error(err)
		}
		if containsDuplicates(order) {
			t.Errorf("output contains duplicates")
		}
		if err = validateOrdering(order, test); err != nil {
			t.Errorf("Invalid ordering: %s, %s", order, err)
		}
	}
}
