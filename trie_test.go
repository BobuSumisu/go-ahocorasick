package ahocorasick

import (
	"log"
	"testing"
)

func TestWiki(t *testing.T) {
	t.Skip()
	trie := NewTrieBuilder(256).AddStringPatterns([]string{
		"a", "ab", "bab", "bc", "bca", "c", "caa",
	}).Build()
	trie.MatchString("abba")
}

func TestPrefix(t *testing.T) {
	t.Skip()
	trie := NewTrieBuilder(256).
		AddStringPatterns([]string{"Aho-Corasick", "Aho-Cora", "Aho", "A"}).
		Build()
	matches := trie.MatchString("Aho-Corasick")

	if len(matches) != 4 {
		t.Errorf("expected %d matches, got %d", 4, len(matches))
	}
}

func TestSuffix(t *testing.T) {
	trie := NewTrieBuilder(256).
		AddStringPatterns([]string{"Aho-Corasick", "Corasick", "rasick"}).
		Build()

	NewTrieGrapher(trie).DrawFailLinks(true).Graph("trie.dot")

	matches := trie.MatchString("Aho-Corasick")

	if len(matches) != 4 {
		t.Errorf("expected %d matches, got %d", 4, len(matches))
	}

	for _, m := range matches {
		log.Print(m)
	}
}
