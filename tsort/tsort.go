package tsort

import "fmt"

type Sorter struct {
	orderings map[string][]string
}

func (s *Sorter) Add(items []string) {
	if s.orderings == nil {
		s.orderings = make(map[string][]string)
	}
	if len(items) == 0 {
		return
	}
	if len(items) == 1 {
		s.orderings[items[0]] = nil
		return
	}
	for i := 0; i < len(items)-1; i++ {
		s.orderings[items[i]] = append(s.orderings[items[i]], items[i+1])
	}
}

func (s *Sorter) heads() []string {
	var out []string
	candidates := make(map[string]struct{})
	for k := range s.orderings {
		candidates[k] = struct{}{}
	}
	for _, values := range s.orderings {
		for _, v := range values {
			if _, ok := candidates[v]; ok {
				delete(candidates, v)
			}
		}
	}
	for k := range candidates {
		out = append(out, k)
	}

	return out
}

func (s *Sorter) isBlocked(n string) bool {
	for _, items := range s.orderings {
		for _, i := range items {
			if i == n {
				return true
			}
		}
	}
	return false
}

func (s *Sorter) Order() ([]string, error) {
	var out []string
	var candidates = s.heads()
	for len(candidates) > 0 {
		var head string
		head, candidates = candidates[0], candidates[1:]
		out = append(out, head)
		unblocked := s.orderings[head]
		delete(s.orderings, head)
		for _, u := range unblocked {
			if !s.isBlocked(u) {
				candidates = append(candidates, u)
			}
		}
	}
	if len(s.orderings) > 0 {
		return nil, fmt.Errorf("Hmm. Still items left. Ordering impossible?")
	}

	return out, nil
}
