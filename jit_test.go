package jit

import (
	"math"
	"testing"
)

func sqrt(x float64) float64 { return math.Sqrt(x) }
func sin(x float64) float64  { return math.Sin(x) }
func cos(x float64) float64  { return math.Cos(x) }

func TestJIT(t *testing.T) {
	for _, x := range []float64{3} { //, -1e9, -123.4, -1, 0, 1, 123.4, 1e9} {
		for _, y := range []float64{5} { //, -1e9, -123.4, -1, 0, 1, 123.4, 1e9} {
			tests := []struct {
				expr string
				want float64
			}{
				{"x", x},
				{"y", y},
				{"-x", -x},
				{"x+y", x + y},
				{"2+x+y+1", 2 + x + y + 1},
				{"1", 1},
				{"1.0", 1},
				{"1+2", 1 + 2},
				{"1-2", 1 - 2},
				{"2*3", 2 * 3},
				{"5/2", 5. / 2.},
				{"2*(x+y)*(x-y)/2", 2 * (x + y) * (x - y) / 2},
				{"1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1", 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1},
				{"sqrt(x)", sqrt(x)},
				{"sqrt(9)", sqrt(9)},
				{"sqrt(x+y)", sqrt(x + y)},
				{"sin(2/x)+cos(sqrt(x+y+1))", sin(2/x) + cos(sqrt(x+y+1))},
				{"cos(9)", cos(9)},
				{"sin(x+y)", sin(x + y)},
				{"sqrt(sqrt(sqrt(x)))", sqrt(sqrt(sqrt(x)))},
				{"1+2+(3+2*4+((((5+6*2)+7)+sqrt(8))+9)+10*sin(2-x+y/3))+11", 1 + 2 + (3 + 2*4 + ((((5 + 6*2) + 7) + sqrt(8)) + 9) + 10*sin(2-x+y/3)) + 11},
				{"1+x+(y+2*4+((((5+y*2)+7)+sqrt(x))+y)+10*sin(2-x+y/x))+y", 1 + x + (y + 2*4 + ((((5 + y*2) + 7) + sqrt(x)) + y) + 10*sin(2-x+y/x)) + y},
				{"1+2+(3+2*4+((((5+6*2)+7)+(8))+9)+10*sin(2-x+y/3))+11", 1 + 2 + (3 + 2*4 + ((((5 + 6*2) + 7) + (8)) + 9) + 10*sin(2-x+y/3)) + 11},
				{"1+2+(3+2*4+((((5+6*2)+7)+(8))+9)+10*(2-x+y/3))+11", 1 + 2 + (3 + 2*4 + ((((5 + 6*2) + 7) + (8)) + 9) + 10*(2-x+y/3)) + 11},
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
	return math.Abs((x-y)/(x+y)) < 1e-14
}

func TestEval2D(t *testing.T) {
	nx := 100
	ny := 50
	dst := make([]float64, nx*ny)
	matrix := make([][]float64, ny)
	for iy := range matrix {
		matrix[iy] = dst[iy*nx : (iy+1)*nx]
	}
	//code, err := Compile("x+y")
	code, err := Compile("x+y")
	if err != nil {
		t.Fatal(err)
	}
	xmin := -0.5
	xmax := 0.5
	ymin := -0.2
	ymax := 0.2
	// half cell offsets: must eval at center.
	dx := 0.5 * (xmax - xmin) / float64(nx)
	dy := 0.5 * (ymax - ymin) / float64(ny)
	code.Eval2D(dst, xmin, xmax, nx, ymin, ymax, ny)

	tests := []struct {
		ix, iy int
		want   float64
	}{
		{0, 0, xmin + dx + ymin + dy},
		{nx - 1, 0, xmax - dx + ymin + dy},
		{0, ny - 1, xmin + dx + ymax - dy},
		{nx - 1, ny - 1, xmax - dx + ymax - dy},
	}

	for _, test := range tests {
		have := matrix[test.iy][test.ix]
		if !equal(have, test.want) {
			t.Errorf("eval2D dst[%v][%v]: want %v, have %v", test.iy, test.ix, test.want, have)
		}
	}
}

func BenchmarkJIT(b *testing.B) {
	code, err := Compile("(x+y)*2 + (1+x) / y")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	n := b.N / 1000 // loader does 1000 loops
	for i := 0; i < n; i++ {
		code.Eval(2, 3)
	}
}

func BenchmarkJITBig(b *testing.B) {
	code, err := Compile("1+x+(3+y*4+((((x+y*2)+x)+sqrt(8))+y)+10*sin(2-x+y/3))+11")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	n := b.N / 1000 // loader does 1000 loops
	for i := 0; i < n; i++ {
		code.Eval(2, 3)
	}
}

func BenchmarkNativeGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nativeGo(2, 3)
	}
}
func nativeGo(x, y float64) float64 {
	return (x+y)*2 + (1+x)/y
}

func BenchmarkNativeGoBig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nativeGoBig(2, 3)
	}
}
func nativeGoBig(x, y float64) float64 {
	return 1 + x + (3 + y*4 + ((((x + y*2) + x) + sqrt(8)) + y) + 10*sin(2-x+y/3)) + 11
}
