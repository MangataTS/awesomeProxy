package ac_automaton

type Trie struct {
	trans    []map[rune]int
	count    int
	patterns map[int]string
}

func NewTrie() *Trie {
	t := new(Trie)
	t.count = 1
	t.patterns = make(map[int]string)
	t.trans = make([]map[rune]int, 1)
	t.trans[0] = map[rune]int{}
	return t
}

func (t *Trie) Insert(s string) {
	var u int = 0
	for _, c := range s {
		if _, ok := t.trans[u][c]; !ok {
			t.trans[u][c] = t.count
			t.count++
			t.trans = append(t.trans, map[rune]int{})
		}
		u = t.trans[u][c]
	}
	t.patterns[u] = s
}

func (t *Trie) Exists(s string) bool {
	var u int = 0
	for _, c := range s {
		if _, ok := t.trans[u][c]; !ok {
			return false
		}
		u = t.trans[u][c]
	}
	return true
}
