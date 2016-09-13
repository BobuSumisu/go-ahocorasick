// Package ahocorasick implements the Aho-Corasick string searching algorithm in Go.
//
// The algorithm is implemented using a double array trie for increased access speed and reduced memory consumption.
//
// The algorithm uses an alphabet size of 256, so can only be used to match byte patterns.
package ahocorasick

// Trie implementing the Aho-Corasick algorithm. Uses two arrays (base and check) for transitions
// (as described by Aho).
//
// A transition in the trie from state s to state t on symbol c is described by:
//
//     base[s] + c = t
//     check[t] = s
//
// Note that the symbol c is actually stored as c + 1 in this implementation for easier handling of
// transition on 0.
type Trie struct {
	base  []int64 // base[s] holds the state s' pointer into check.
	check []int64 // check holds the "owner" of states.
	dict  []int64 // Holds the pattern length of s (if it is in the dictionary).
	fail  []int64 // Holds the fail link for s.
	suff  []int64 // Holds the dictionary suffix link for s.
}

// Run the Trie against the provided input and returns potentially matches.
func (tr *Trie) Match(input []byte) []*Match {
	matches := make([]*Match, 0)

	s := RootState

	for i, c := range input {
		s = tr.step(s, encodeByte(c))

		if tr.dict[s] != 0 {
			pos := int64(i+1) - tr.dict[s]
			matches = append(matches, newMatch(pos, input[pos:i+1]))
		}

		for f := tr.suff[s]; f != EmptyCell; f = tr.suff[f] {
			pos := int64(i+1) - tr.dict[f]
			matches = append(matches, newMatch(pos, input[pos:i+1]))
		}
	}

	return matches
}

// Same as Match, but returns immediately after the first matched pattern.
func (tr *Trie) MatchFirst(input []byte) *Match {
	s := RootState

	for i, c := range input {
		s = tr.step(s, encodeByte(c))

		if tr.dict[s] != 0 {
			pos := int64(i+1) - tr.dict[s]
			return newMatch(pos, input[pos:i+1])
		}

		if f := tr.suff[s]; f != EmptyCell {
			pos := int64(i+1) - tr.dict[f]
			return newMatch(pos, input[pos:i+1])
		}
	}

	return nil
}

// Helper method to make matching strings a little more comfortable.
func (tr *Trie) MatchString(input string) []*Match {
	return tr.Match([]byte(input))
}

// Helper method to make matching a string a little more comfortable.
func (tr *Trie) MatchStringFirst(input string) *Match {
	return tr.MatchFirst([]byte(input))
}

func (tr *Trie) NumPatterns() int64 {
	var c int64 = 0
	for _, n := range tr.dict {
		if n != 0 {
			c += 1
		}
	}
	return c
}

func (tr *Trie) step(s, c int64) int64 {
	t := tr.base[s] + c
	if t < int64(len(tr.check)) && tr.check[t] == s {
		return t
	}

	for f := tr.fail[s]; f > 0; f = tr.fail[f] {
		t := tr.base[f] + c
		if t < int64(len(tr.check)) && tr.check[t] == f {
			return t
		}
	}

	t = tr.base[RootState] + c
	if t < int64(len(tr.check)) && tr.check[t] == RootState {
		return t
	}

	return RootState
}
