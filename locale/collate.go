package locale

type Collate struct {
	elements []string
	symbols  []string

	order    []ordering
	backward bool
	position bool // ???????
}

type ordering struct {
	id      string
	weights []int
}
