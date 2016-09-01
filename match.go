package ahocorasick

type Match interface {
	Pos() int64
	End() int64
}

type IntMatch struct {
	pos   int64
	match []int64
}

type ByteMatch struct {
	pos   int64
	match []byte
}

type StringMatch struct {
	pos   int64
	match string
}

func NewIntMatch(pos int64, match []int64) *IntMatch {
	return &IntMatch{
		pos:   pos,
		match: match,
	}
}

func (m *IntMatch) Pos() int64     { return m.pos }
func (m *IntMatch) End() int64     { return m.pos + int64(len(m.match)) }
func (m *IntMatch) Match() []int64 { return m.match }

func NewByteMatch(match *IntMatch) *ByteMatch {
	bytes := make([]byte, len(match.match))
	for i := range bytes {
		bytes[i] = byte(match.match[i])
	}
	return &ByteMatch{pos: match.pos, match: bytes}
}

func (m *ByteMatch) Pos() int64    { return m.pos }
func (m *ByteMatch) End() int64    { return m.pos + int64(len(m.match)) }
func (m *ByteMatch) Match() []byte { return m.match }

func NewStringMatch(match *ByteMatch) *StringMatch {
	return &StringMatch{pos: match.pos, match: string(match.match)}
}
