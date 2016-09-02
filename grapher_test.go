package ahocorasick

import "os/exec"

func ExampleTrieGrapher_Graph() {
	trie := NewTrieBuilder().
		AddStrings([]string{"his", "hers", "he", "she"}).
		Build()

	grapher := NewTrieGrapher(trie).DrawFailLinks(true)
	grapher.Graph("trie.dot")

	exec.Command("dot", "-Tpng", "-o", "trie.png", "trie.dot").Run()

}
