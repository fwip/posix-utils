package regexes

// Matches returns true if the input string matches the regex
func (bre Bre) Matches(input string) bool {
	runes := []rune(input)

	if bre.anchorLeft || len(input) == 0 {
		match, length := elementsMatch(bre.elements, runes)
		if !match {
			return false
		}
		if bre.anchorRight {
			return length == len(runes)
		}
		return true
	}

	for i := range runes {
		match, length := elementsMatch(bre.elements, runes[i:])
		if !match {
			continue
		}
		if bre.anchorRight && length != len(runes[i:]) {
			continue
		}
		return true
	}

	return false
}

// recursive
func elementsMatch(elements []breDuplElement, runes []rune) (matches bool, matchlen int) {

	// No elements matches any string
	if len(elements) == 0 {
		return true, len(runes)
	}

	e := elements[0]
	matched := 0

	// Find how many runes match the current element
	for i := 0; i < e.count.max && i < len(runes); i++ {
		if e.element.Matches(runes[i]) {
			matched++
		} else if matched < e.count.min {
			// Return early if we don't have at least the minimum matches
			return false, 0
		}
	}

	// Check remainder of matches, backtracking when necessary
	for i := matched; i >= e.count.min; i-- {
		if match, length := elementsMatch(elements[1:], runes[i:]); match {
			return match, length + i
		}
	}

	return false, 0
}
