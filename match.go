package ahocorasick

type Match struct {
	pos   int64
	match []byte
}

func NewMatch(pos int64, match []byte) *Match {
	return &Match{pos, match}
}

func (m *Match) Pos() int64    { return m.pos }
func (m *Match) End() int64    { return m.pos + int64(len(m.match)) }
func (m *Match) Match() []byte { return m.match }
