package ahocorasick

import (
	"fmt"
	"os/exec"
)

func ExampleTrieGrapher_Graph() {
	trie := NewTrieBuilder().
		AddStrings([]string{"his", "hers", "he", "she"}).
		Build()

	grapher := NewTrieGrapher(trie).DrawFailLinks(true)
	grapher.Graph("trie.dot")

	if _, err := exec.LookPath("dot"); err == nil {
		if err := exec.Command("dot", "-Tpng", "-o", "trie.png", "trie.dot").Run(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("OK")
		}
		// Output: OK
	} else {
		// Output:
	}
}
