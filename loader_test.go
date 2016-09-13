package ahocorasick

import "testing"

func TestSaveTrie(t *testing.T) {
	patterns, err := ReadStrings("./test_data/NSF-ordlisten.cleaned.uniq.txt")
	if err != nil {
		t.Error(err)
	}

	trie := NewTrieBuilder().AddPatterns(patterns[:100]).Build()

	if err = SaveTrie(trie, "test.trie"); err != nil {
		t.Error(err)
	}
}

func TestLoadTrie(t *testing.T) {
	// Assumes TestSaveTrie has been run.

	trie, err := LoadTrie("test.trie")
	if err != nil {
		t.Error(err)
	}

	if trie == nil {
		t.Error("trie was nil")
	}

	if trie.NumPatterns() != 100 {
		t.Errorf("expected %d patterns, got %d", 100, trie.NumPatterns())
	}
}
