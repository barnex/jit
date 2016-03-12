package main

import (
	"testing"
)

func TestJIT(t *testing.T) {
	tests := []struct {
		expr string
		x, y float64
		want float64
	}{
		{"x", 42, 666, 42},
		{"y", 42, 666, 666},
		{"1", 42, 666, 1},
		{"1.0", 42, 666, 1},
	}

	for _, test := range tests {
		var b Buf
		err := b.Compile(test.expr)
		if err != nil {
			t.Error(err)
			continue
		}
		have := b.call(test.x, test.y)
		b.Free()
		if have != test.want {
			t.Errorf("%v with x=%v,y=%v: have %v, want: %v", test.expr, test.x, test.y, have, test.want)
		}
	}
}
