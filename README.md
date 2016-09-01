# Aho-Corasick

Aho-Corasick string search algorithm implemented in Go.

Uses a [double array trie](https://linux.thai.net/~thep/datrie/datrie.html) for improved
access speeds and reduced memory consumption.

## Usage


Use a `TrieBuilder` to create a `Trie`:

```go
trie := NewTrieBuilder().
    AddStrings([]string{"hers", "his", "he", "she"}).
    Build()
```

Match something:

```go
matches := trie.MatchString("I have never tasted a hershey bar.")
fmt.Printf("We got %d matches.\n", len(matches))

// => We got 4 matches.
```

Examine matches:

```go
for _, match := range matches {
    fmt.Printf("Matched %q at offset %d.\n", match.Match(), match.Pos())
}

// => Matched "he" at offset 22.
// => Matched "hers" at offset 22.
// => Matched "she" at offset 25.
// => Matched "he" at offset 26.
```

For debugging you may output the trie in DOT format:

```go
NewTrieGrapher(trie).DrawFailLinks(true).Graph("example.dot")
```

And convert to image, e.g.:

```bash
    $ dot -Tpng -o example.png example.dot
```

![example-trie](example.png)
