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

	tr.Build()

	tr.Match("footballz and bazars are bar foobar")

}
