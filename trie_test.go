package ahocorasick

import (
	"log"
	"testing"
)

func TestFuckAll(t *testing.T) {
	trie := NewTrieBuilder(256).AddStringPatterns([]string{
		"foo",
		"bar",
		"baz",
		"football",
		"bazar",
		"øyvind",
		"foot",
		"ball",
	}).Build()
	trie.MatchString("hei du")
	// tr.Build()
	// tr.Match("footballz and bazars are bar foobar")
}

func TestWiki(t *testing.T) {
	trie := NewTrieBuilder(256).AddStringPatterns([]string{
		"a", "ab", "bab", "bc", "bca", "c", "caa",
	}).Build()
	trie.MatchString("alle barna gikk i skogen caa bca c")
}

func TestPrefix(t *testing.T) {
	trie := NewTrieBuilder(256).
		AddStringPatterns([]string{"Aho-Corasick", "Aho-Cora", "Aho", "A"}).
		Build()

	matches := trie.MatchString("Aho-Corasick")

	if len(matches) != 4 {
		t.Errorf("expected %d matches, got %d", 4, len(matches))
	}

	for _, m := range matches {
		log.Print(m)
	}
}
