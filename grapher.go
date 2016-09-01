package ahocorasick

import (
	"fmt"
	"io"
	"os"
	"unicode"
)

type TrieGrapher struct {
	trie          *Trie
	w             io.Writer
	drawFailLinks bool
}

func NewTrieGrapher(trie *Trie) *TrieGrapher {
	return &TrieGrapher{
		trie: trie,
	}
}

func (tg *TrieGrapher) DrawFailLinks(b bool) *TrieGrapher {
	tg.drawFailLinks = b
	return tg
}

func (tg *TrieGrapher) Graph(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	tg.w = f

	fmt.Fprintln(f, "digraph T {")
	fmt.Fprintln(f, "\tnodesep=0.2; ranksep=0.4; splines=false; outputorder=edgesfirst;")
	fmt.Fprintln(f, "\tnode [shape=circle, style=filled, fillcolor=white, fixedsize=true];")
	fmt.Fprintln(f, "\tedge [arrowsize=0.5];")

	// Will recursivelly call graphState on every state (which is in use).
	tg.graphState(RootState, EmptyCell)

	fmt.Fprintln(f, "}")

	return nil
}

func (tg *TrieGrapher) graphState(s, c int64) {

	if tg.trie.dict[s] {
		fmt.Fprintf(tg.w, "\t%d [label=%q, shape=doublecircle];\n", s, label(c))
	} else {
		fmt.Fprintf(tg.w, "\t%d [label=%q];\n", s, label(c))
	}

	for c := int64(0); c < tg.trie.alphabetSize; c++ {
		t := tg.trie.base[s] + c
		if t < int64(len(tg.trie.check)) && tg.trie.check[t] == s {
			tg.graphState(t, c)
			fmt.Fprintf(tg.w, "\t%d -> %d;\n", s, t)
		}
	}

	if f := tg.trie.fail[s]; tg.drawFailLinks && f != EmptyCell && f != RootState {
		fmt.Fprintf(tg.w, "\t%d -> %d [color=red, constraint=false];\n", s, f)
	}

	if f := tg.trie.suff[s]; f != EmptyCell {
		fmt.Fprintf(tg.w, "\t%d -> %d [color=darkgreen, constraint=false];\n", s, f)
	}
}

func label(c int64) string {
	if c == -1 {
		return ""
	}

	if unicode.IsPrint(rune(byte(c))) {
		return fmt.Sprintf("%c", byte(c))
	}

	return fmt.Sprintf("%x", c)
}