package ac_automaton

import "testing"

func TestTrie(t *testing.T) {
	trie := NewTrie()
	ss := []string{"she", "her", "his", "he", "is"}
	for _, s := range ss {
		trie.Insert(s)
	}

	if !trie.Exists("she") {
		t.Error("1")
	}
	if trie.Exists("sher") {
		t.Error("2")
	}
}
