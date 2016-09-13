package ahocorasick

import (
	"encoding/binary"
	"fmt"
	"os"
)

const MagicNumber int32 = 0x45495254

// Save a Trie to file.
func SaveTrie(tr *Trie, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	binary.Write(f, binary.LittleEndian, MagicNumber)

	// Write each of the arrays to the file (preceded by its length).
	for _, arr := range [][]int64{tr.base, tr.check, tr.dict, tr.fail, tr.suff} {
		if err = binary.Write(f, binary.LittleEndian, int64(len(arr))); err != nil {
			return err
		}

		if err = binary.Write(f, binary.LittleEndian, arr); err != nil {
			return err
		}
	}

	return nil
}

// Read a Trie from file.
func LoadTrie(path string) (*Trie, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read and check magic number.
	var magic int32
	if err = binary.Read(f, binary.LittleEndian, &magic); err != nil {
		return nil, err
	}
	if magic != MagicNumber {
		return nil, fmt.Errorf("Not a valid trie file (magic mismatch: 0x%08x).", magic)
	}

	tr := new(Trie)

	// Read arrays.

	for _, arr := range []*[]int64{&tr.base, &tr.check, &tr.dict, &tr.fail, &tr.suff} {
		var n int64
		if err = binary.Read(f, binary.LittleEndian, &n); err != nil {
			return nil, err
		}

		*arr = make([]int64, n)
		if err = binary.Read(f, binary.LittleEndian, arr); err != nil {
			return nil, err
		}
	}

	return tr, nil
}
