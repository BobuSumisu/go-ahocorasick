package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BobuSumisu/go-ahocorasick"
	"github.com/pkg/profile"
)

func main() {
	var err error

	defer profile.Start(profile.ProfilePath(".")).Stop()

	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <patterns-file> <input-file> [<num-patterns-to-use>]", os.Args[0])
		os.Exit(1)
	}

	numPatterns := 10000
	if len(os.Args) > 3 {
		numPatterns, err = strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
	}

	// Read patterns.
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	patterns := make([][]byte, 0)

	for s.Scan() {
		patterns = append(patterns, []byte(strings.TrimSpace(s.Text())))
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// Read input data.
	input, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Running trie on %d bytes of data using %d patterns", len(input), numPatterns)

	start := time.Now()

	// Build trie.
	trie := ahocorasick.NewTrieBuilder().AddPatterns(patterns[:numPatterns]).Build()

	for n := 0; n < 1000; n++ {
		trie.Match(input)
	}

	log.Printf("Done in %v.", time.Since(start))
}
