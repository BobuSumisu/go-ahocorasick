package ahocorasick

type Trie struct {
	base  []int64
	check []int64
	dict  []bool
	fail  []int64
	suff  []int64
}

func (tr *Trie) Match(input []byte) []*Match {
	matches := make([]*Match, 0)

	s := RootState

	for i, c := range input {
		s = tr.step(s, c)

		if tr.dict[s] {
			pos := int64(i) - tr.patternLen(s) + 1
			matches = append(matches, NewMatch(pos, input[pos:i+1]))
		}

		for f := tr.suff[s]; f > 0; f = tr.suff[f] {
			pos := int64(i) - tr.patternLen(f) + 1
			matches = append(matches, NewMatch(pos, input[pos:i+1]))
		}
	}

	return matches
}

func (tr *Trie) MatchString(input string) []*Match {
	return tr.Match([]byte(input))
}

func (tr *Trie) step(s int64, c byte) int64 {
	t := tr.base[s] + int64(c)
	if t < int64(len(tr.check)) && tr.check[t] == s {
		return t
	}

	for f := tr.fail[s]; f > 0; f = tr.fail[f] {
		t := tr.base[f] + int64(c)
		if t < int64(len(tr.check)) && tr.check[t] == f {
			return t
		}
	}

	t = tr.base[RootState] + int64(c)
	if t < int64(len(tr.check)) && tr.check[t] == RootState {
		return t
	}

	return RootState
}

func (tr *Trie) pathString(t int64) string {
	return string(tr.pathBytes(t))
}

func (tr *Trie) pathBytes(t int64) []byte {
	pathInts := tr.path(t)
	pathBytes := make([]byte, len(pathInts))
	for i := range pathInts {
		pathBytes[i] = byte(pathInts[i])
	}
	return pathBytes
}

func (tr *Trie) path(t int64) []int64 {
	if tr.check[t] == -1 {
		return nil
	}
	s := tr.check[t]
	c := t - tr.base[s]
	return append(tr.path(s), c)
}

func (tr *Trie) patternLen(s int64) int64 {
	if tr.check[s] == 0 {
		return 1
	}

	return tr.patternLen(tr.check[s]) + 1
}
