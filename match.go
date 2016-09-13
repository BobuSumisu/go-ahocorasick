package ahocorasick

import (
	"bytes"
	"fmt"
)

// Represents a matched pattern.
type Match struct {
	pos   int64
	match []byte
}

func newMatch(pos int64, match []byte) *Match {
	return &Match{pos, match}
}

func newMatchString(pos int64, match string) *Match {
	return &Match{pos, []byte(match)}
}

func (m *Match) String() string {
	return fmt.Sprintf("{%v %q}", m.pos, m.match)
}

// Get the position (offset) of the matched pattern.
func (m *Match) Pos() int64 { return m.pos }

// Get the end position of the matched pattern.
func (m *Match) End() int64 { return m.pos + int64(len(m.match)) }

// Get the matched byte pattern.
func (m *Match) Match() []byte { return m.match }

// Just to make working with strings a little more comfortable.
func (m *Match) MatchString() string { return string(m.match) }

// Check if two matches are equal.
func MatchEqual(m1, m2 *Match) bool {
	return bytes.Equal(m1.match, m2.match) && m1.pos == m2.pos
}
