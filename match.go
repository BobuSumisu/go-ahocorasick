package ahocorasick

// Represents a matched pattern.
type Match struct {
	pos   int64
	match []byte
}

func newMatch(pos int64, match []byte) *Match {
	return &Match{pos, match}
}

// Get the position (offset) of the matched pattern.
func (m *Match) Pos() int64 { return m.pos }

// Get the end position of the matched pattern.
func (m *Match) End() int64 { return m.pos + int64(len(m.match)) }

// Get the matched byte pattern.
func (m *Match) Match() []byte { return m.match }
