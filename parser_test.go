package jit

import "testing"

func TestParser(t *testing.T) {
	tests := []string{
		"x",
		"1",
		"1+1",
		"(1+1)",
		"sin(x)",
		"-x",
		"+x",
	}

	for _, test := range tests {
		if _, err := Parse(test); err != nil {
			t.Error(err)
		}
	}
}
