package ahocorasick

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestWiki(t *testing.T) {
	patterns := []string{
		"a", "ab", "bab", "bc", "bca", "c", "caa",
	}
	input := "Aho-Corasick"
	expected := []*Match{
		&Match{match: []byte("a"), pos: 7},
		&Match{match: []byte("c"), pos: 10},
	}

	matches := NewTrieBuilder().AddStrings(patterns).Build().MatchString(input)

	if len(expected) != len(matches) {
		t.Errorf("expected %d matches, got %d", len(expected), len(matches))
	} else {
		for i := range expected {
			if !bytes.Equal(expected[i].Match(), matches[i].Match()) ||
				expected[i].Pos() != matches[i].Pos() {
				t.Errorf("expected %v, got %v", expected[i], matches[i])
			}
		}
	}
}

func TestPrefix(t *testing.T) {
	trie := NewTrieBuilder().
		AddStrings([]string{"Aho-Corasick", "Aho-Cora", "Aho", "A"}).
		Build()
	matches := trie.MatchString("Aho-Corasick")

	if len(matches) != 4 {
		t.Errorf("expected %d matches, got %d", 4, len(matches))
	}
}

func TestSuffix(t *testing.T) {
	trie := NewTrieBuilder().
		AddStrings([]string{"Aho-Corasick", "Corasick", "rasick", "k"}).
		Build()
	matches := trie.MatchString("Aho-Corasick")

	if len(matches) != 4 {
		t.Errorf("expected %d matches, got %d", 4, len(matches))
	}
}

func TestInfix(t *testing.T) {
	trie := NewTrieBuilder().
		AddStrings([]string{"Aho-Corasick", "ho-Corasi", "-Cora", "-"}).
		Build()
	matches := trie.MatchString("Aho-Corasick")

	if len(matches) != 4 {
		t.Errorf("expected %d matches, got %d", 4, len(matches))
	}
}

func TestOverlap(t *testing.T) {
	trie := NewTrieBuilder().
		AddStrings([]string{"Aho-Co", "ho-Cora", "o-Coras", "-Corasick"}).
		Build()
	matches := trie.MatchString("Aho-Corasick")

	if len(matches) != 4 {
		t.Errorf("expected %d matches, got %d", 4, len(matches))
	}
}

func TestAdjacent(t *testing.T) {
	trie := NewTrieBuilder().
		AddStrings([]string{"Ah", "o-Co", "ras", "ick"}).
		Build()
	matches := trie.MatchString("Aho-Corasick")

	if len(matches) != 4 {
		t.Errorf("expected %d matches, got %d", 4, len(matches))
	}
}

func TestSingleSymbol(t *testing.T) {
	trie := NewTrieBuilder().
		AddStrings([]string{"o"}).
		Build()
	matches := trie.MatchString("Aho-Corasick")

	if len(matches) != 2 {
		t.Errorf("expected %d matches, got %d", 2, len(matches))
	}
}

func TestNoMatch(t *testing.T) {
	trie := NewTrieBuilder().
		AddStrings([]string{"Gazorpazorpfield", "Knuth", "b"}).
		Build()
	matches := trie.MatchString("Aho-Corasick")

	if len(matches) != 0 {
		t.Errorf("expected %d matches, got %d", 0, len(matches))
	}
}

func TestUtf8(t *testing.T) {
	trie := NewTrieBuilder().
		AddStrings([]string{"Øyvind", "lærer", "å", "♡"}).
		Build()
	matches := trie.MatchString("Øyvind lærer seg å programmere ♡")

	if len(matches) != 4 {
		t.Errorf("expected %d matches, got %d", 0, len(matches))
	}
}

func TestZeroes(t *testing.T) {
	trie := NewTrieBuilder().
		AddPattern([]byte{0x00, 0x00}).
		Build()

	matches := trie.MatchString("\x00\x00Aho\x00\x00Cora\x00\x00sick\x00\x00\x00\x00")

	if len(matches) != 6 {
		t.Errorf("expected %d matches, got %d", 6, len(matches))
	}
}

func TestMatchFirst(t *testing.T) {
	trie := NewTrieBuilder().AddString("foo").Build()

	match := trie.MatchStringFirst("foo foo foo foo foo foo foo foo foo foo foo foo")

	if match.MatchString() != "a" || match.Pos() != 0 {
		fmt.Errorf("expected match %q at %d, got match %q at %d",
			"a", 0, match.MatchString(), match.Pos())
	}
}

func BenchmarkBuildNSF(b *testing.B) {
	patterns, err := ReadStrings("./test_data/NSF-ordlisten.cleaned.txt")
	if err != nil {
		b.Error(err)
	}

	for n := 0; n < b.N; n++ {
		NewTrieBuilder().AddPatterns(patterns[:200]).Build()
	}
}

func BenchmarkMatchIbsen(b *testing.B) {
	patterns, err := ReadStrings("./test_data/NSF-ordlisten.cleaned.txt")
	if err != nil {
		b.Error(err)
	}

	input, err := ioutil.ReadFile("./test_data/Ibsen.txt")
	if err != nil {
		b.Error(err)
	}

	trie := NewTrieBuilder().AddPatterns(patterns[:200]).Build()

	for n := 0; n < b.N; n++ {
		trie.Match(input[:1000])
	}
}

func ExampleReadme() {
	trie := NewTrieBuilder().
		AddStrings([]string{"hers", "his", "he", "she"}).
		Build()

	matches := trie.MatchString("I have never tasted a hershey bar.")
	fmt.Printf("We got %d matches.\n", len(matches))
	for _, match := range matches {
		fmt.Printf("Matched %q at offset %d.\n", match.Match(), match.Pos())
	}

	NewTrieGrapher(trie).DrawFailLinks(true).Graph("example.dot")

	exec.Command("dot", "-Tpng", "-o", "example.png", "example.dot").Run()

	// Output:
	// We got 4 matches.
	// Matched "he" at offset 22.
	// Matched "hers" at offset 22.
	// Matched "she" at offset 25.
	// Matched "he" at offset 26.
}

func ExampleTrie_MatchString() {
	trie := NewTrieBuilder().
		AddStrings([]string{"hers", "his", "he", "she"}).
		Build()
	matches := trie.MatchString("she is here")
	fmt.Println(len(matches))
	// Output: 3
}
