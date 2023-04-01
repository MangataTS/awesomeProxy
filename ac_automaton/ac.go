package ac_automaton

import "unicode/utf8"

type ACAutomaton struct {
	Trie
	fail []int
}

func NewACAutomaton() *ACAutomaton {
	ac := new(ACAutomaton)
	ac.Trie = *NewTrie()
	ac.fail = make([]int, 0)
	return ac
}

func (ac *ACAutomaton) Build() {
	if ac.count == 0 {
		return
	}

	ac.fail = make([]int, ac.count)
	queue := make([]int, 0)
	for _, u := range ac.trans[0] {
		queue = append(queue, u)
	}

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		for c, v := range ac.trans[u] {
			failU := ac.fail[u]
			_, ok := ac.trans[failU][c]
			for failU > 0 && !ok {
				failU = ac.fail[failU]
				_, ok = ac.trans[failU][c]
			}

			if ok {
				ac.fail[v] = ac.trans[failU][c]
			}
			queue = append(queue, v)
		}
	}
}

func (ac *ACAutomaton) FindMatches(s string) map[string]int {
	res := make(map[string]int)
	endMap := make(map[int]int)

	var u int = 0
	var ok bool = false
	for i, c := range s {
		_, ok = ac.trans[u][c]
		for u > 0 && !ok {
			u = ac.fail[u]
			_, ok = ac.trans[u][c]
		}

		if ok {
			u = ac.trans[u][c]
			cLen := utf8.RuneLen(c)
			for v := u; v > 0; v = ac.fail[v] {
				if pattern, exist := ac.patterns[v]; exist {
					start := i - len(pattern) + cLen
					end := endMap[v]
					if start >= end {
						res[pattern]++
						endMap[v] = i + cLen
					}
				}
			}
		}
	}
	return res
}
