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

	for _, str := range testStrings {
		if !tr.Run(str) {
			t.Error("didn't accept")
		}
	}

	tr.Match("footballz")

}
