package ahocorasick

import (
	"bufio"
	"os"
	"strings"
)

// Read string patterns from a text file, one pattern on each line.
func ReadStrings(path string) ([][]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	patterns := make([][]byte, 0)

	for s.Scan() {
		pattern := strings.TrimSpace(s.Text())
		patterns = append(patterns, []byte(pattern))
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return patterns, nil
}

// Read patterns in hex format, one pattern on each line.
func ReadHex(path string) ([][]byte, error) {
	// TODO: implement
	return nil, nil
}
