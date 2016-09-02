package ahocorasick

const (
	AlphabetSize int64 = 256 // The size of the alphabet is fixed to the size of a byte.
	RootState    int64 = 0   // The root state of the trie is always 0.
	EmptyCell    int64 = -1  // Represents an unused cell.
	DefaultBase  int64 = 0   // The default base for new states.
)

// A TrieBuilder must be used to properly build Tries.
type TrieBuilder struct {
	base  []int64
	check []int64
	dict  []bool
	fail  []int64
	suff  []int64
}

// Create and initialize a new TrieBuilder.
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

// Add a new pattern to be built into the resulting Trie.
func (tb *TrieBuilder) AddPattern(pattern []byte) *TrieBuilder {
	s := RootState

	for _, c := range pattern {
		t := tb.base[s] + int64(c+1)

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
				oc = EmptyCell // State o does not have a transition to s.
			}

			tb.relocate(o)

			// Update s and t if o had transitions to s.
			if oc != EmptyCell {
				s = tb.base[o] + oc
				t = tb.base[s] + int64(c+1)
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

// A helper method to make adding multiple patterns a little more comfortable.
func (tb *TrieBuilder) AddPatterns(patterns [][]byte) *TrieBuilder {
	for _, pattern := range patterns {
		tb.AddPattern(pattern)
	}
	return tb
}

// A helper method to make adding a string pattern more comfortable.
func (tb *TrieBuilder) AddString(pattern string) *TrieBuilder {
	return tb.AddPattern([]byte(pattern))
}

// A helper method to make adding multiple string patterns a little more comfortable.
func (tb *TrieBuilder) AddStrings(patterns []string) *TrieBuilder {
	for _, pattern := range patterns {
		tb.AddString(pattern)
	}
	return tb
}

// Build the trie.
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
	if tb.fail[s] != EmptyCell {
		return // Avoid computing more than one time.
	}

	p := tb.check[s]    // The parent of s.
	if p == EmptyCell { // No transitions to s, ignore.
		return
	} else if p == s {
		return // If s is it's own parent.
	}

	tb.computeFailLink(p)

	c := s - tb.base[p] // The transition symbol to this state.

	if p == RootState {
		// If parent is root, fail to root
		tb.fail[s] = RootState
	} else {

		// Follow fail links (starting from parent) until we find a state f with
		// a transition on this states symbol (c).
		for f := tb.fail[p]; f > 0; f = tb.fail[f] {
			// Set s' fail to f's child if it has a transition.
			t := tb.base[f] + c
			if t < int64(len(tb.check)) && tb.check[t] == f {
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

// Get all c's for which state s has a transition (that is, where check[base[s]+c] == s).
func (tb *TrieBuilder) transitions(s int64) []byte {
	cs := make([]byte, 0)

	for c := int64(0); c < AlphabetSize; c++ {
		t := tb.base[s] + (c + 1)
		if t < int64(len(tb.check)) && tb.check[t] == s {
			cs = append(cs, byte(c+1))
		}
	}
	return cs
}

// Check wether b is a suitable base for s given it's transitions on cs.
func (tb *TrieBuilder) suitableBase(b, s int64, cs []byte) bool {
	for _, c := range cs {
		t := b + int64(c)

		// All offsets above len(check) is of course empty.
		if t >= int64(len(tb.check)) {
			return true
		}

		if tb.check[t] != EmptyCell {
			return false
		}
	}
	return true
}

// Find a suitable (new) base for s.
func (tb *TrieBuilder) findBase(s int64, cs []byte) int64 {
	for b := DefaultBase; ; b++ {
		if tb.suitableBase(b, s, cs) {
			return b
		}
	}
	return EmptyCell
}

func (tb *TrieBuilder) relocate(s int64) {
	// First find all symbols for which s has a transition.
	cs := tb.transitions(s)

	// Find a new suitable base for s.
	b := tb.findBase(s, cs)

	// Move the base of s to b. First we must update the transitions.
	for _, c := range cs {
		// Old t and new t'.
		t := tb.base[s] + int64(c)
		t_ := b + int64(c)

		tb.expandArrays(t_) // Ensure arrays are big enough for t'.

		tb.check[t_] = s         // Mark s as owner of t'.
		tb.base[t_] = tb.base[t] // Copy base value.
		tb.dict[t_] = tb.dict[t] // As well as the dictionary value.

		// We must also update all states which had transitions from t to t'.
		for c := int64(0); c < AlphabetSize; c++ {
			u := tb.base[t] + (c + 1)

			if u >= int64(len(tb.check)) {
				break
			}

			if tb.check[u] == t {
				tb.check[u] = t_
			}
		}

		// Unset old tb.check and dictionary values for t.
		tb.check[t] = EmptyCell
		tb.dict[t] = false
	}

	// Finally we can move the base for s.
	tb.base[s] = b
}
