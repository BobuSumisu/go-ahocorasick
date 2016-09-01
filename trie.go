package ahocorasick

type Trie struct {
	base         []int64
	check        []int64
	dict         []bool
	fail         []int64
	suff         []int64
	alphabetSize int64
}

func (tr *Trie) MatchInts(input []int64) []*IntMatch {
	matches := make([]*IntMatch, 0)

	s := RootState

	for i, c := range input {
		s = tr.step(s, c)

		if tr.dict[s] {
			pos := int64(i) - tr.patternLen(s) + 1
			matches = append(matches, NewIntMatch(pos, input[pos:i+1]))
		}

		for f := tr.suff[s]; f > 0; f = tr.suff[f] {
			pos := int64(i) - tr.patternLen(s) + 1
			matches = append(matches, NewIntMatch(pos, input[pos:i+1]))
		}
	}

	return matches
}

func (tr *Trie) MatchBytes(input []byte) []*ByteMatch {
	intInput := make([]int64, len(input))
	for i := range input {
		intInput[i] = int64(input[i])
	}
	matches := tr.MatchInts(intInput)
	byteMatches := make([]*ByteMatch, len(matches))
	for i := range matches {
		byteMatches[i] = NewByteMatch(matches[i])
	}

	return byteMatches
}

func (tr *Trie) MatchString(input string) []*StringMatch {
	matches := tr.MatchBytes([]byte(input))
	stringMatches := make([]*StringMatch, len(matches))
	for i := range matches {
		stringMatches[i] = NewStringMatch(matches[i])
	}
	return stringMatches
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
	if tr.check[s] == -1 {
		return 0
	}

	return tr.patternLen(tr.check[s]) + 1
}

/*
func NewTrie() *Trie {
	tr := &Trie{
		base:         make([]int, 0),
		check:        make([]int, 0),
		dict:         make([]bool, 0),
		alphabetSize: 256,
	}

	// Append root state.
	tr.base = append(tr.base, 0)
	tr.check = append(tr.check, -1)
	tr.dict = append(tr.dict, false)

	return tr
}

func (tr *Trie) AddPattern(pattern string) *Trie {
	s := 0

	for _, b := range []byte(pattern) {
		c := int(b)
		// log.Printf("[status] we are in %d and got %q", s, c)
		t := tr.base[s] + c

		// // log.Printf("[status] %d --(%q)--> %d", s, c, t)

		if t >= len(tr.check) || tr.check[t] == -1 { // if cell is empty
			// log.Printf("[transition] %d --(%q)--> %d", s, c, t)
			tr.expand(t)
			tr.check[t] = s
			s = t
		} else if tr.check[t] == s { // s already has a transition on c
			s = t
		} else { // someone else is occupying this cell

			// The occupying state.
			o := tr.check[t]

			// This state may have a transition to the current state we
			// are in. In which case we must update the s value.
			cToS := -1
			for c := 0; c < tr.alphabetSize; c++ {
				if tr.base[o]+c == s && tr.check[s] == o {
					cToS = c
				}
			}

			// Move the occupier base.
			tr.relocate(o)

			// If the occupying state has a transition to s then
			// we have to update s to the new value of s resulting
			// in the relocation of o.
			if cToS != -1 {
				s = tr.base[o] + cToS
				t = tr.base[s] + c
			}

			// log.Printf("[transition] %d --(%q)--> %d", s, c, t)
			tr.check[t] = s
			s = t
		}

	}

	tr.dict[s] = true

	return tr
}

func (tr *Trie) Build() {

	// Initialize arrays
	tr.fail = make([]int, len(tr.base))
	tr.suf = make([]int, len(tr.base))

	for i := 0; i < len(tr.base); i++ {
		tr.fail[i] = -1
		tr.suf[i] = -1
	}

	tr.fail[0] = 0 // Root fails to itself.

	for s := 0; s < len(tr.base); s++ {
		tr.computeFailLink(s)
	}

	for s := 0; s < len(tr.base); s++ {
		tr.computeSufLinks(s)
	}

}

func (tr *Trie) computeFailLink(s int) {
	p := tr.check[s] // The parent of this state.
	if p == -1 {     // No transitions to this state, ignore.
		return
	}

	tr.computeFailLink(p) // Need to compute parent's fail links first.
	c := s - tr.base[p]   // The symbol into this state.

	if p == 0 {
		tr.fail[s] = 0 // Child of root fails to root
	} else {
		// Follow fail links (starting from parent) until we find a state with a transition
		// on c.
		for f := tr.fail[p]; f > 0; f = tr.fail[f] {
			tr.computeFailLink(f)
			if tr.check[tr.base[f]+c] == f {
				tr.fail[s] = tr.base[f] + c
			}
		}

		// If we didn't find any fail link.
		if tr.fail[s] == -1 {
			// Check if root has a transition on c.
			if tr.check[tr.base[0]+c] == 0 {
				tr.fail[s] = tr.base[0] + c
			} else {
				tr.fail[s] = 0 // Else fail to root.
			}
		}
	}
}

func (tr *Trie) computeSufLinks(s int) {
	for f := tr.fail[s]; f > 0; f = tr.fail[f] {
		if tr.dict[f] {
			tr.suf[s] = f
			return
		}
	}
}
func (tr *Trie) Match(input string) {
	bi := []byte(input)

	for i := 0; i < len(bi); i++ {
		s := 0
		for j, b := range bi[i:] {
			c := int(b)
			t := tr.base[s] + c
			if tr.check[t] == s {
				s = t
				if tr.dict[s] {
					p := tr.pathString(s)
					log.Printf("matched: %q at offset %d", p, i+j-len(p)+1)
				}

			} else {
				break
			}
		}
	}
}

func (tr *Trie) RealMatch(input string) {
	s := 0

	for i, b := range []byte(input) {
		c := int(b)

		s = tr.step(s, c)

		if tr.dict[s] {
			p := tr.pathString(s)
			log.Printf("matched %q at %d", p, i+1-len(p))
		}

		for u := tr.suf[s]; u > 0; u = tr.suf[u] {
			p := tr.pathString(u)
			log.Printf("matched %q at %d", p, i+1-len(p))
		}
	}
}

func (tr *Trie) step(s, c int) int {
	// If s has a transition on c, return that.
	t := tr.base[s] + c
	if t < len(tr.check) && tr.check[t] == s {
		return t
	}

	// Else walk fail links.
	for f := tr.fail[s]; f > 0; f = tr.fail[f] {
		t := tr.base[f] + c
		if t < len(tr.check) && tr.check[t] == f {
			return t // Return fail link's child if it has transition on c.
		}
	}

	// Or simply return root
	return 0
}

func (tr *Trie) Run(input string) bool {
	s := 0

	// log.Printf("*** RUNNING ***")

	for _, b := range []byte(input) {
		c := int(b)
		t := tr.base[s] + c

		if tr.check[t] == s {
			// log.Printf("[move] %d --(%q)--> %d", s, c, t)
			s = t
		} else {
			// log.Printf("[stop] %d --(%q)-->", s, c)
			for c := 0; c < tr.alphabetSize; c++ {
				t := tr.base[s] + c
				if t < len(tr.check) && tr.check[t] == s {
					// log.Printf("[info] had move %d --(%q)--> %d", s, c, t)
				}
			}
			return false
		}
	}

	// log.Printf("[accept] we ended up at %d", s)

	// log.Print("*** DONE ***")

	return true
}

func (tr *Trie) relocate(s int) {

	// log.Print("*** RELOCATION ***")
	// First find all c's and t's where check[t] == s
	ts := make([]int, 0)
	cs := make([]int, 0)
	for c := 0; c < tr.alphabetSize; c++ {
		t := tr.base[s] + c
		if t < len(tr.check) && tr.check[t] == s {
			ts = append(ts, t)
			cs = append(cs, c)
		}
	}

	// Find a new suitable base.
	b := 0
	for {
		ok := true
		for _, c := range cs {
			t_ := b + c
			if t_ < len(tr.check) && tr.check[t_] != -1 {
				ok = false
				break
			}
		}
		if ok {
			break
		}
		b++
	}

	// log.Printf("[relocation] start relocation of base[%d] from %d --> %d", s, tr.base[s], b)

	// Move base[s] to b.
	for _, c := range cs {
		t := tr.base[s] + c
		t_ := b + c

		// log.Printf("[relocation] move t %d -> %d", t, t_)

		tr.expand(t_)

		// log.Printf("[relocation] set check[%d] = %d", t_, s)
		tr.check[t_] = s // Mark owner at new t location.
		// log.Printf("[relocation] set base[%d] = base[%d] = %d", t_, t, tr.base[t])
		tr.base[t_] = tr.base[t] // Copy base value to new location
		tr.dict[t_] = tr.dict[t] // Copy dict value.

		// Update all pointers to t to new t.
		for i, x := range tr.check {
			if x == t {
				tr.check[i] = t_
			}
		}

		// log.Printf("[relocation] set check[%d] = -1", t)
		tr.check[t] = -1   // Free cell at old t.
		tr.dict[t] = false // Free dict at old t.
	}

	// Update base value.
	// log.Printf("[relocation] set base[%d] = %d", s, b)
	tr.base[s] = b

	// log.Print("*** END RELOCATION ***")
}

func (tr *Trie) debug() {
}

// Get the "path" of the state t by walking backwards
// until we hit root.
func (tr *Trie) path(t int) []int {
	if tr.check[t] == -1 {
		return nil
	}

	s := tr.check[t]
	c := t - tr.base[s]

	return append(tr.path(s), c)
}

func (tr *Trie) pathString(s int) string {
	p := tr.path(s)
	bs := make([]byte, len(p))
	for i := range p {
		bs[i] = byte(p[i])
	}
	return string(bs)
}

// Ensure our arrays are long enough.
func (tr *Trie) expand(n int) {
	for len(tr.check) <= n {
		tr.base = append(tr.base, 0)     // New state with base = 0
		tr.check = append(tr.check, -1)  // And no owner.
		tr.dict = append(tr.dict, false) // And not in dictionary.
	}
}
*/
