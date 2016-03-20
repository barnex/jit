package jit

import (
	"math"
	"testing"
)

func sqrt(x float64) float64 { return math.Sqrt(x) }

func TestJIT(t *testing.T) {
	for _, x := range []float64{-1e9, -123.4, -1, 0, 1, 123.4, 1e9}{
	for _, y := range []float64{-1e9, -123.4, -1, 0, 1, 123.4, 1e9}{
	tests := []struct {
		expr string
		want float64
	}{
		{"x",  x},
		{"y",  y},
		{"1",  1},
		{"1.0",  1},
		{"1+2",  1 + 2},
		{"1-2",  1 - 2},
		{"2*3",  2 * 3},
		{"5/2",  5. / 2.},
		{"2*(x+y)*(x-y)/2",  2 * (x + y) * (x - y) / 2},
		{"1+1+1+1+1+1+1+1+1+1+1+1",  1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1},
		{"sqrt(x)", sqrt(x)},
		{"sqrt(9)",  sqrt(9)},
		{"sqrt(x+y)", sqrt(x + y)},
	}

	for _, test := range tests {
		code, err := Compile(test.expr)
		if err != nil {
			t.Error(err)
			continue
		}
		have := code.Eval(x, y)
		code.Free()
		if !equal(have,test.want) {
			t.Errorf("%v with x=%v,y=%v: have %v, want: %v", test.expr, x, y, have, test.want)
		}
	}
	}}
}

// equal returns whether x == y, treating NaN's as equal
func equal(x, y float64)bool{
	if math.IsNaN(x) && math.IsNaN(y){
		return true
	}
	return x == y
}

func BenchmarkJIT(b *testing.B) {
	code, err := Compile("(x+y)*2 + (1+x) / y")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		code.Eval(2, 3)
	}
}

func nativeGo(x, y float64) float64 {
	return (x+y)*2 + (1+x)/y
}

func BenchmarkNativeGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nativeGo(2, 3)
	}
}
