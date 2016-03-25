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

func TestWalk(t *testing.T) {
	tests := []struct {
		expr string
		want int
	}{
		{"x", 1},
		{"x+y", 3},
		{"(x+y)", 3},
		{"sin(x-y)", 4},
		{"sin(cos(x)*cos(y))", 6},
		{"sin(cos(x)/cos(y))", 6},
		{"+x", 1},
		{"-x", 3}, // 0 - x
	}

	for _, test := range tests {
		root, err := Parse(test.expr)
		if err != nil {
			t.Error(err)
			continue
		}
		n := 0
		walk(root, func(expr) { n++ })
		if n != test.want {
			t.Errorf("walk %q: have %v nodes, want %v", test.expr, n, test.want)
		}
	}
}

func TestRecordCalls(t *testing.T) {
	tests := []struct {
		expr string
		want bool
	}{
		{"x", false},
		{"x+y", false},
		{"(x+y)", false},
		{"sin(x-y)", true},
		{"sin(cos(x)*cos(y))", true},
		{"sin(cos(x)/cos(y))", true},
		{"1+sin(cos(x)/cos(y))", true},
		{"(1+sin(cos(x)/cos(y)))+1", true},
		{"1+sin(x)", true},
		{"1+(2*(sin(x)))", true},
		{"1+(2*sin(x))", true},
		{"1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1", false},
		{"+x", false},
		{"-x", false}, // 0 - x
	}

	for _, test := range tests {
		root, err := Parse(test.expr)
		if err != nil {
			t.Error(err)
			continue
		}
		m := make(map[expr]bool)
		recordCalls(root, m)
		if m[root] != test.want {
			t.Errorf("has calls %q: have %v, want %v", test.expr, m[root], test.want)
		}
	}
}
