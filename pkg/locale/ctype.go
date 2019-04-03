package locale

type Ctype struct {
	// TODO: Should these be 'character' structs rather than strings?
	upper   []string
	lower   []string
	alpha   []string
	digit   []string
	alnum   []string
	space   []string
	cntrl   []string
	punct   []string
	graph   []string
	print   []string
	xdigit  []string
	blank   []string
	other   map[string][]string // This maps locale-defined class names to characters
	toupper map[string]string
	tolower map[string]string
}
