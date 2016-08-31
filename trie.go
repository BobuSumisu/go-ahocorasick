package ahocorasick

import "log"

type Trie struct {
	base         []int
	check        []int
	dict         []bool
	alphabetSize int
}

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

func (tr *Trie) Match(input string) {
	bi := []byte(input)

	for i := 0; i < len(bi); i++ {
		s := 0

		for j, b := range bi[i:] {
			c := int(b)
			t := tr.base[s] + c

			if tr.dict[s] {
				p := tr.path(s)
				log.Printf("matched: %q at offset %d", p, i+j-len(p))
			}

			if tr.check[t] == s {
				s = t
			} else {
				break
			}
		}
	}

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

		// Update all pointers to t to new t.
		for i, x := range tr.check {
			if x == t {
				tr.check[i] = t_
			}
		}

		// log.Printf("[relocation] set check[%d] = -1", t)
		tr.check[t] = -1 // Free cell at old t.
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

// Ensure our arrays are long enough.
func (tr *Trie) expand(n int) {
	for len(tr.check) <= n {
		tr.base = append(tr.base, 0)     // New state with base = 0
		tr.check = append(tr.check, -1)  // And no owner.
		tr.dict = append(tr.dict, false) // And not in dictionary.
	}
}
