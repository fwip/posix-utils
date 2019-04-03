package regexes

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/shenwei356/util/math"
)

var shouldParse = []string{
	"",
	".",
	"a",
	"abc",
	"[ab]",
	"a*",
	"a{1}",
	"a{1,}",
	"a{1,2}",
}

var shouldNotParse = []string{
	"a{1",
	"[ab",
}

type eT struct {
	min     int
	max     int
	element breElement
}

var parseResults = []struct {
	s           string
	elements    []eT
	anchorLeft  bool
	anchorRight bool
}{
	{s: "a", elements: []eT{
		{1, 1, brePlainText(' ')},
	}},
	{s: "a*", elements: []eT{
		{0, math.MaxInt, brePlainText(' ')},
	}},
	{s: "[ab]{2,}", elements: []eT{
		{2, math.MaxInt, &breBracket{}},
	}},
	{s: "(ab){2,8}", elements: []eT{
		{2, 8, &breParen{}},
	}},
	{s: "[ab]{2,}.", elements: []eT{
		{2, math.MaxInt, &breBracket{}},
		{1, 1, breWildcard{}},
	}},
	{s: "[ab]{2,}.", elements: []eT{
		{2, math.MaxInt, &breBracket{}},
		{1, 1, breWildcard{}},
	}},
	{"^ab$", []eT{
		{1, 1, brePlainText(' ')},
		{1, 1, brePlainText(' ')},
	}, true, true},
}

var shouldMatch = []struct {
	re        string
	matches   []string
	nomatches []string
}{
	{"",
		[]string{"", "ab"},
		[]string{},
	},
	{"a",
		[]string{"a", "ab", "ba"},
		[]string{"x", ""},
	},
	{"ab",
		[]string{"ab", "cab", "abc", "baba"},
		[]string{"", "a", "b", "ba"},
	},
	{"a.b",
		[]string{"axb", "aab", "abb"},
		[]string{"ab", "ax", "axxb"},
	},
	{"a{1,2}.b",
		[]string{"abb", "aab", "aaab", "babb"},
		[]string{"ab", "ax", "axxb"},
	},
	{"a*",
		[]string{"a", "", "ba", "aabaa"},
		[]string{},
	},
	{"a{2,2}",
		[]string{"aa", "aaa"},
		[]string{"", "a", "aba"},
	},
	{".{1,2}",
		[]string{"ab", "abc"},
		[]string{""},
	},
	{"[ab]",
		[]string{"a", "b", "xby"},
		[]string{"]"},
	},
	{"[xyz]",
		[]string{"x", "y", "z"},
		[]string{"abc", ""},
	},
	{"(xyz)",
		[]string{"xyz"},
		[]string{"x", "y", "z", "xyjz"},
	},
	{"^a$",
		[]string{"a"},
		[]string{},
	},
	{"a$",
		[]string{"ahowa"},
		[]string{},
	},
	{"^",
		[]string{"", "ab"},
		[]string{},
	},
	{"$",
		[]string{"", "ab"},
		[]string{},
	},
}

func xEqualInt(expected, actual int, t *testing.T) {
	if expected != actual {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}

func shouldPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func TestParseBre(t *testing.T) {

	for _, test := range shouldParse {
		t.Run("should parse: "+test, func(t *testing.T) {
			ParseBre(test)
		})
	}
	for _, test := range shouldNotParse {
		t.Run("should not parse: "+test, func(t *testing.T) {
			shouldPanic(t, func() {
				ParseBre(test)
			})
		})
	}

	for _, test := range parseResults {
		t.Run("check parse: "+test.s, func(t *testing.T) {
			bre, err := ParseBre(test.s)
			if err != nil {
				t.Error(err)
			}
			if len(test.elements) != len(bre.elements) {
				t.Errorf("expected %d elements, got %d instead", len(test.elements), len(bre.elements))
				return
			}
			if bre.anchorLeft != test.anchorLeft {
				t.Errorf("anchorLeft should be %t, was %t", test.anchorLeft, bre.anchorLeft)
			}
			if bre.anchorRight != test.anchorRight {
				t.Errorf("anchorRight should be %t, was %t", test.anchorRight, bre.anchorRight)

			}
			for i, x := range test.elements {
				a := bre.elements[i]
				xEqualInt(x.min, a.count.min, t)
				xEqualInt(x.max, a.count.max, t)
				xtype := reflect.TypeOf(x.element)
				atype := reflect.TypeOf(a.element)
				if xtype != atype {
					t.Errorf("element %d; expected type %s, got %s", i, xtype, atype)
				}
			}
		})
	}
}

func TestMatchBre(t *testing.T) {
	for _, test := range shouldMatch {
		bre, err := ParseBre(test.re)
		for _, match := range test.matches {
			t.Run(fmt.Sprintf("/%s/ =~ '%s'", test.re, match), func(t *testing.T) {
				if err != nil {
					t.Error(err)
				}
				if !bre.Matches(match) {
					t.Errorf("does not match")
				}
			})
		}
		for _, nomatch := range test.nomatches {
			t.Run(fmt.Sprintf("/%s/ !~ '%s'", test.re, nomatch), func(t *testing.T) {
				if err != nil {
					t.Error(err)
				}
				if bre.Matches(nomatch) {
					t.Errorf("should not match")
				}
			})
		}
	}
}
