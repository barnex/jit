package jit

import (
	"math"
	"testing"
)

func sqrt(x float64) float64 { return math.Sqrt(x) }
func sin(x float64) float64  { return math.Sin(x) }
func cos(x float64) float64  { return math.Cos(x) }

func TestJIT(t *testing.T) {
	for _, x := range []float64{-1e9, -123.4, -1, 0, 1, 123.4, 1e9} {
		for _, y := range []float64{-1e9, -123.4, -1, 0, 1, 123.4, 1e9} {
			tests := []struct {
				expr string
				want float64
			}{
				{"x", x},
				{"y", y},
				{"x+y", x + y},
				{"2+x+y+1", 2 + x + y + 1},
				{"1", 1},
				{"1.0", 1},
				{"1+2", 1 + 2},
				{"1-2", 1 - 2},
				{"2*3", 2 * 3},
				{"5/2", 5. / 2.},
				{"2*(x+y)*(x-y)/2", 2 * (x + y) * (x - y) / 2},
				{"1+1+1+1+1+1+1+1+1+1+1+1", 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1},
				{"sqrt(x)", sqrt(x)},
				{"sqrt(9)", sqrt(9)},
				{"sqrt(x+y)", sqrt(x + y)},
				{"sin(2/x)+cos(sqrt(x+y+1))", sin(2/x) + cos(sqrt(x+y+1))},
				{"cos(9)", cos(9)},
				{"sin(x+y)", sin(x + y)},
				{"sqrt(sqrt(sqrt(x)))", sqrt(sqrt(sqrt(x)))},
			}

			for _, test := range tests {
				code, err := Compile(test.expr)
				if err != nil {
					t.Error(err)
					continue
				}
				have := code.Eval(x, y)
				code.Free()
				if !equal(have, test.want) {
					t.Errorf("%v with x=%v,y=%v: have %v, want: %v", test.expr, x, y, have, test.want)
				}
			}
		}
	}
}

// equal returns whether x and y are approximately equal
func equal(x, y float64) bool {
	if math.IsNaN(x) && math.IsNaN(y) {
		return true
	}
	if x == y {
		return true
	}
	return math.Abs((x-y)/(x+y)) < 1e-15
}

func BenchmarkJIT(b *testing.B) {
	code, err := Compile("(x+y)*2 + (1+x) / y")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	n := b.N/1000// loader does 1000 loops
	for i := 0; i < n; i++ {
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
