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
		{"1+2", 666, 666, 3},
		{"1-2", 666, 666, -1},
		{"2*3", 666, 666, 6},
		{"5/2", 666, 666, 2.5},
		{"2*(x+y)*(x-y)/2", 2, 3, -5},
		{"1+1+1+1+1+1+1+1+1+1+1+1", 666, 666, 12},
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
