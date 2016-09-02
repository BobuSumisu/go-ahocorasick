package ahocorasick

import "fmt"

func ExampleTrieBuilder_Build() {
	builder := NewTrieBuilder()

	builder.AddPattern([]byte{0x44, 0x22, 0x31, 0x52, 0x32, 0x00, 0x01, 0x01})
	builder.AddStrings([]string{"hello", "world"})

	trie := builder.Build()

	fmt.Println(len(trie.MatchString("hello!")))
	// Output: 1
}
