package ahocorasick

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestTrie(t *testing.T) {
	cases := []struct {
		name     string
		patterns []string
		input    string
		expected []*Match
	}{
		{
			"Wikipedia",
			[]string{"a", "ab", "bab", "bc", "bca", "c", "caa"},
			"Aho-Corasick",
			[]*Match{
				newMatchString(7, "a"),
				newMatchString(10, "c"),
			},
		},
		{
			"Prefix",
			[]string{"Aho-Corasick", "Aho-Cora", "Aho", "A"},
			"Aho-Corasick",
			[]*Match{
				newMatchString(0, "A"),
				newMatchString(0, "Aho"),
				newMatchString(0, "Aho-Cora"),
				newMatchString(0, "Aho-Corasick"),
			},
		},
		{
			"Suffix",
			[]string{"Aho-Corasick", "Corasick", "sick", "k"},
			"Aho-Corasick",
			[]*Match{
				newMatchString(0, "Aho-Corasick"),
				newMatchString(4, "Corasick"),
				newMatchString(8, "sick"),
				newMatchString(11, "k"),
			},
		},
		{
			"Infix",
			[]string{"Aho-Corasick", "ho-Corasi", "o-Co", "-"},
			"Aho-Corasick",
			[]*Match{
				newMatchString(3, "-"),
				newMatchString(2, "o-Co"),
				newMatchString(1, "ho-Corasi"),
				newMatchString(0, "Aho-Corasick"),
			},
		},
		{
			"Overlap",
			[]string{"Aho-Co", "ho-Cora", "o-Coras", "-Corasick"},
			"Aho-Corasick",
			[]*Match{
				newMatchString(0, "Aho-Co"),
				newMatchString(1, "ho-Cora"),
				newMatchString(2, "o-Coras"),
				newMatchString(3, "-Corasick"),
			},
		},
		{
			"Adjacent",
			[]string{"Ah", "o-Co", "ras", "ick"},
			"Aho-Corasick",
			[]*Match{
				newMatchString(0, "Ah"),
				newMatchString(2, "o-Co"),
				newMatchString(6, "ras"),
				newMatchString(9, "ick"),
			},
		},
		{
			"SingleSymbol",
			[]string{"o"},
			"Aho-Corasick",
			[]*Match{
				newMatchString(2, "o"),
				newMatchString(5, "o"),
			},
		},
		{
			"NoMatch",
			[]string{"Gazorpazopfield", "Knuth", "O"},
			"Aho-Corasick",
			[]*Match{},
		},
		{
			"Zeroes",
			[]string{"\x00\x00"},
			"\x00\x00Aho\x00\x00-\x00\x00Corasick\x00\x00",
			[]*Match{
				newMatchString(0, "\x00\x00"),
				newMatchString(5, "\x00\x00"),
				newMatchString(8, "\x00\x00"),
				newMatchString(18, "\x00\x00"),
			},
		},
	}

	for _, c := range cases {
		tr := NewTrieBuilder().AddStrings(c.patterns).Build()
		matches := tr.MatchString(c.input)

		if len(matches) != len(c.expected) {
			t.Errorf("%s: expected %d matches, got %d", c.name, len(c.expected), len(matches))
			continue
		}

		for i := range matches {
			if !MatchEqual(matches[i], c.expected[i]) {
				t.Errorf("%s: expected %v, got %v", c.name, matches[i], c.expected[i])
			}
		}
	}
}

func TestMatchFirst(t *testing.T) {
	trie := NewTrieBuilder().AddString("o").Build()
	match := trie.MatchStringFirst("Aho-Corasick")
	expected := newMatchString(2, "o")

	if !MatchEqual(match, expected) {
		fmt.Errorf("expected %v, got %v", expected, match)
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
