package regexes

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/shenwei356/util/math"
)

var breCountRe = regexp.MustCompile(`^(\d+)(?:(,)(\d+)?)?$`)

func mustBeInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

// Bre is a Basic Regular Expression
type Bre struct {
	text        string
	anchorLeft  bool
	anchorRight bool

	elements []breDuplElement
}

type breDuplElement struct {
	element breElement
	count   breCount
}

func (bde *breDuplElement) String() string {
	return bde.element.Text() + bde.count.String()
}

type breCount struct {
	min     int
	max     int
	thrifty bool // thrifty is the opposite of greedy
}

func (bc breCount) String() string {
	return fmt.Sprintf("{%d,%d}", bc.min, bc.max)
}

type breElement interface {
	Text() string
	Matches(rune) bool
}

type brePlainText rune

func (bpt brePlainText) Text() string {
	return string([]rune{rune(bpt)})
}
func (bpt brePlainText) Matches(r rune) bool {
	return rune(bpt) == r
}

type breWildcard struct{}

func (bw breWildcard) Text() string        { return "." }
func (bw breWildcard) Matches(r rune) bool { return true }

type breBracket struct {
	chars []rune
}

func (bb *breBracket) Text() string {
	return "[" + string(bb.chars) + "]"
}
func (bb *breBracket) Matches(r rune) bool {
	for _, c := range bb.chars {
		if r == c {
			return true
		}
	}
	return false
}

type breParen struct {
	chars []rune
}

func (bp *breParen) Text() string {
	return "(" + string(bp.chars) + ")"
}
func (bp *breParen) Matches(r rune) bool {
	for _, c := range bp.chars {
		if r == c {
			return true
		}
	}
	return false
}

func parseCount(runes []rune) (count breCount, countLen int) {
	count.min = 1
	count.max = 1
	if len(runes) == 0 {
		return count, 0
	}
	switch runes[0] {
	case '{':
		n, err := readUntil(runes[1:], '}')
		if err != nil {
			panic(err)
		}

		matches := breCountRe.FindStringSubmatch(string(runes[1 : 1+n]))
		if matches == nil {
			panic(fmt.Errorf("don't understand count '%s'", string(runes[1:1+n])))
		}
		// Hack
		for i := len(matches) - 1; i >= 0; i-- {
			if matches[i] == "" {
				matches = matches[:i]
			}
		}
		count.min = mustBeInt(matches[1])
		switch len(matches) {
		case 2:
			count.max = count.min
		case 3:
			count.max = math.MaxInt
		case 4:
			count.max = mustBeInt(matches[3])
		}

		countLen = 2 + n
		// // Check for thriftiness
		// if countLen < len(runes) && runes[countLen] == '?' {
		// 	count.thrifty = true
		// 	countLen++
		// }

	//case '?':
	//	count.min = 0
	//	count.max = 1
	//	count.thrifty = true
	//	countLen = 1
	case '*':
		count.min = 0
		count.max = math.MaxInt
		countLen = 1
		//case '+':
		//	count.min = 1
		//	count.max = math.MaxInt
		//	countLen = 1
		//}
	}
	return count, countLen
}

// ParseBre will return a new Bre object
func ParseBre(s string) (*Bre, error) {
	b := Bre{}
	runes := []rune(s)
	if len(runes) == 0 {
		return &b, nil
	}

	if runes[0] == '^' {
		b.anchorLeft = true
		runes = runes[1:]
	}
	if len(runes) == 0 {
		return &b, nil
	}
	if runes[len(runes)-1] == '$' {
		b.anchorRight = true
		runes = runes[:len(runes)-1]
	}

	// Parse element
	var element breElement
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch r {
		case '[':
			n, err := readUntil(runes[i+1:], ']')
			if err != nil {
				panic(err)
			}
			element = &breBracket{chars: runes[i+1 : i+n+1]}
			i += n + 1
		//case '(':
		//	n, err := readUntil(runes[i+1:], ')')
		//	if err != nil {
		//		panic(err)
		//	}
		//	element = &breParen{chars: runes[i+1 : i+n]}
		//	i += n + 1
		case '(', ')': //TODO: Implement this properly
			continue
		case '.':
			element = breWildcard{}
		default:
			element = brePlainText(r)
		}

		// Check for count specification
		count, countLen := parseCount(runes[i+1:])
		i += countLen
		dupl := breDuplElement{count: count, element: element}
		b.elements = append(b.elements, dupl)
	}

	return &b, nil
}

func readUntil(runes []rune, stop rune) (length int, err error) {
	if runes[0] == stop {
		return 0, nil
	}
	for j := 0; j < len(runes); j++ {
		r := runes[j]
		if r == '\\' {
			j++
			continue
		}
		if runes[j] == stop {
			return j, nil
		}
	}
	return 0, fmt.Errorf("stop char %c not found", stop)
}
