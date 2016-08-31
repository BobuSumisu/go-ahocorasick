package ahocorasick

import "testing"

func TestFuckAll(t *testing.T) {
	tr := NewTrie()

	testStrings := []string{
		"foo",
		"bar",
		"baz",
		"football",
		"bazar",
		"øyvind",
		"foot",
		"ball",
	}

	for _, str := range testStrings {
		tr.AddPattern(str)
	}

	// tr.Build()
	// tr.Match("footballz and bazars are bar foobar")
}

func TestWiki(t *testing.T) {
	tr := NewTrie()

	patterns := []string{"a", "ab", "bab", "bc", "bca", "c", "caa"}
	for _, p := range patterns {
		tr.AddPattern(p)
	}
	tr.Build()

	tr.RealMatch("alle barna gikk i skogen caa bca c")
}

func TestPrefix(t *testing.T) {
	tr := NewTrie()
	for _, s := range []string{"Aho-Corasick", "Aho-Cora", "Aho", "A"} {
		tr.AddPattern(s)
	}
	tr.Build()
	tr.RealMatch("Aho-Corasick")
}
