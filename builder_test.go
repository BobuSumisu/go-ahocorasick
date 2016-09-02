package ahocorasick

import (
	"fmt"
	"log"
	"testing"
)

func ExampleTrieBuilder_Build() {
	builder := NewTrieBuilder()

	builder.AddPattern([]byte{0x44, 0x22, 0x31, 0x52, 0x32, 0x00, 0x01, 0x01})
	builder.AddStrings([]string{"hello", "world"})

	trie := builder.Build()

	fmt.Println(len(trie.MatchString("hello!")))
	// Output: 1
}

func TestTrieBuilder_Build_zeroes(t *testing.T) {
	trie := NewTrieBuilder().
		AddPattern([]byte{0, 1, 2, 1}).
		AddPattern([]byte{0, 0}).
		AddPattern([]byte{0, 0, 0, 0}).
		Build()

	matches := trie.Match([]byte{0, 0, 0, 1, 2, 1, 0, 0, 1, 2, 0, 0, 0, 0})

	if len(matches) != 8 {
		t.Errorf("expected %d matches, got %d", 8, len(matches))

		for _, match := range matches {
			log.Printf("Matched %q at %d", match.Match(), match.Pos())
		}
	}

}
