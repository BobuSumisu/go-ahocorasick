package ahocorasick

const (
	AlphabetSize int64 = 256
	RootState    int64 = 0
	EmptyCell    int64 = -1
	DefaultBase  int64 = 0
)

type TrieBuilder struct {
	base  []int64
	check []int64
	dict  []bool
	fail  []int64
	suff  []int64
}

func NewTrieBuilder() *TrieBuilder {
	tb := &TrieBuilder{
		base:  make([]int64, 0),
		check: make([]int64, 0),
		dict:  make([]bool, 0),
		fail:  make([]int64, 0),
		suff:  make([]int64, 0),
	}

	// Add the root state.
	tb.addState()

	return tb
}

func (tb *TrieBuilder) AddPattern(pattern []byte) *TrieBuilder {
	s := RootState

	for _, c := range pattern {
		t := tb.base[s] + int64(c)

		if t >= int64(len(tb.check)) || tb.check[t] == EmptyCell {
			// Cell is empty: expand arrays and set transition.
			tb.expandArrays(t)
			tb.check[t] = s
		} else if tb.check[t] == s {
			// Cell is in use by s, simply move on.
		} else {
			// Someone is occupying the cell. Move the occupier.
			o := tb.check[t]

			// Relocating o changes its states. So if o has a transition to s,
			// we must update s after relocating o. First check if o actually has
			// a transition to s.
			oc := s - tb.base[o]
			if tb.check[tb.base[o]+oc] != o {
				oc = -1 // State o does not have a transition to s.
			}

			tb.relocate(o)

			// Update s and t if o had transitions to s.
			if oc != -1 {
				s = tb.base[o] + oc
				t = tb.base[s] + int64(c)
			}

			// Set transition.
			tb.check[t] = s
		}

		// Move to next state.
		s = t
	}

	// Mark s as in dictionary.
	tb.dict[s] = true

	return tb
}

func (tb *TrieBuilder) AddPatterns(patterns [][]byte) *TrieBuilder {
	for _, pattern := range patterns {
		tb.AddPattern(pattern)
	}
	return tb
}

func (tb *TrieBuilder) AddString(pattern string) *TrieBuilder {
	return tb.AddPattern([]byte(pattern))
}

func (tb *TrieBuilder) AddStrings(patterns []string) *TrieBuilder {
	for _, pattern := range patterns {
		tb.AddString(pattern)
	}
	return tb
}

func (tb *TrieBuilder) Build() *Trie {

	// Initialize link arrays.
	tb.fail = make([]int64, len(tb.base))
	tb.suff = make([]int64, len(tb.base))

	for i := 0; i < len(tb.base); i++ {
		tb.fail[i] = EmptyCell
		tb.suff[i] = EmptyCell
	}

	// Root fails to itself.
	tb.fail[RootState] = RootState

	for s := int64(0); s < int64(len(tb.base)); s++ {
		tb.computeFailLink(s)
	}

	for s := int64(0); s < int64(len(tb.base)); s++ {
		tb.computeSuffLink(s)
	}

	// Should I copy these slices over or?
	return &Trie{
		base:  tb.base,
		check: tb.check,
		dict:  tb.dict,
		fail:  tb.fail,
		suff:  tb.suff,
	}
}

func (tb *TrieBuilder) computeFailLink(s int64) {
	p := tb.check[s]    // The parent of s.
	if p == EmptyCell { // No transitions to s, ignore.
		return
	}

	c := s - tb.base[p] // The transition symbol to this state.

	if p == RootState {
		// If parent is root, fail to root
		tb.fail[s] = RootState
	} else {

		tb.computeFailLink(p)

		// Follow fail links (starting from parent) until we find a state f with
		// a transition on this states symbol (c).
		for f := tb.fail[p]; f > 0; f = tb.fail[f] {

			// Set s' fail to f's child if it has a transition.
			t := tb.base[f] + c
			if tb.check[t] == f {
				tb.fail[s] = t
				break
			}

			// Compute f's fail link before the next iteration.
			tb.computeFailLink(f)
		}

		// If for some reason we didn't find any fail link.
		if tb.fail[s] == EmptyCell {
			// Check if root has transition on this s' symbol.
			t := tb.base[RootState] + c
			if t < int64(len(tb.check)) && tb.check[t] == RootState {
				tb.fail[s] = t
			} else {
				// Else fail to root.
				tb.fail[s] = RootState
			}
		}
	}
}

func (tb *TrieBuilder) computeSuffLink(s int64) {
	// Follow fail links until we (possibly) find a state in the dictionary.
	for f := tb.fail[s]; f > 0; f = tb.fail[f] {
		if tb.dict[f] {
			tb.suff[s] = f
			return
		}
	}
}

func (tb *TrieBuilder) addState() {
	tb.base = append(tb.base, DefaultBase)
	tb.check = append(tb.check, EmptyCell)
	tb.dict = append(tb.dict, false)
}

func (tb *TrieBuilder) expandArrays(n int64) {
	for int64(len(tb.base)) <= n {
		tb.addState()
	}
}

func (tb *TrieBuilder) relocate(s int64) {
	// First find all symbols for which s has a transition.
	cs := make([]int64, 0)
	for c := int64(0); c < AlphabetSize; c++ {
		t := tb.base[s] + c
		if t < int64(len(tb.check)) && tb.check[t] == s {
			cs = append(cs, c)
		}
	}

	// Find a new suitable base for s.
	var b int64 = 0
	for {
		foundIt := true

		// Check if the offset b + c is available for every c.
		for _, c := range cs {
			t_ := b + c
			if t_ < int64(len(tb.check)) && tb.check[t_] != EmptyCell {
				foundIt = false
				break
			}
		}

		if foundIt {
			// Current base b is OK.
			break
		}

		// Test next b.
		b++
	}

	// Move the base of s to b. First we must update the transitions.
	for _, c := range cs {
		// Old t and new t'.
		t := tb.base[s] + c
		t_ := b + c

		tb.expandArrays(t_) // Ensure arrays are big enough for t'.

		tb.check[t_] = s         // Mark s as owner of t'.
		tb.base[t_] = tb.base[t] // Copy base value.
		tb.dict[t_] = tb.dict[t] // As well as the dictionary value.

		// We must also update all states which had transitions from t to t'.
		// TODO: This is probabably not the most efficient way of doing this.
		for i := range tb.check {
			if tb.check[i] == t {
				tb.check[i] = t_
			}
		}

		// Unset old tb.check and dictionary values for t.
		tb.check[t] = EmptyCell
		tb.dict[t] = false
	}

	// Finally we can move the base for s.
	tb.base[s] = b
}
