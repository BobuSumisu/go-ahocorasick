package ahocorasick

import "testing"

func TestReadStrings(t *testing.T) {
	patterns, err := ReadStrings("./test_data/NSF-ordlisten.cleaned.txt")
	if err != nil {
		t.Error(err)
	}

	if len(patterns) != 622115 {
		t.Errorf("expected %d patterns, got %d", 622115, len(patterns))
	}

	if string(patterns[7]) != "abandonerende" {
		t.Errorf("expected %q, got %q", "abandonerende", patterns[7])
	}
}
