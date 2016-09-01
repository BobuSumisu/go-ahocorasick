package ahocorasick

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestWiki(t *testing.T) {
	trie := NewTrieBuilder().AddStrings([]string{
		"a", "ab", "bab", "bc", "bca", "c", "caa",
	}).Build()
	matches := trie.MatchString("abracadabra") // 5 a's, 2 ab's, 1 c's = 8

	if len(matches) != 8 {
		t.Errorf("expected %d matches, got %d", 8, len(matches))
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

func readPatterns(path string) ([][]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	patterns := make([][]byte, 0)

	for s.Scan() {
		patterns = append(patterns, []byte(strings.TrimSpace(s.Text())))
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return patterns, nil
}

func BenchmarkBuildNSF(b *testing.B) {
	patterns, err := readPatterns("./test_data/NSF-ordlisten.cleaned.txt")
	if err != nil {
		b.Error(err)
	}

	for n := 0; n < b.N; n++ {
		NewTrieBuilder().AddPatterns(patterns[:200]).Build()
	}
}

func BenchmarkMatchIbsen(b *testing.B) {
	patterns, err := readPatterns("./test_data/NSF-ordlisten.cleaned.txt")
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
