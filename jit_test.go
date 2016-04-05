package jit

import (
	"math"
	"testing"
)

var tests = map[string]func(float64, float64) float64{
	"x": func(x float64, y float64) float64 {
		return x
	},
	"y": func(x float64, y float64) float64 {
		return y
	},
	"-x": func(x float64, y float64) float64 {
		return -x
	},
	"x+y": func(x float64, y float64) float64 {
		return x + y
	},
	"2+x+y+1": func(x float64, y float64) float64 {
		return 2 + x + y + 1
	},
	"1": func(x float64, y float64) float64 {
		return 1
	},
	"1.0": func(x float64, y float64) float64 {
		return 1
	},
	"1+2": func(x float64, y float64) float64 {
		return 1 + 2
	},
	"1-2": func(x float64, y float64) float64 {
		return 1 - 2
	},
	"2*3": func(x float64, y float64) float64 {
		return 2 * 3
	},
	"5/2": func(x float64, y float64) float64 {
		return 5. / 2.
	},
	"2*(x+y)*(x-y)/2": func(x float64, y float64) float64 {
		return 2 * (x + y) * (x - y) / 2
	},
	"1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1": func(x float64, y float64) float64 {
		return 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1
	},
	"sqrt(x)": func(x float64, y float64) float64 {
		return sqrt(x)
	},
	"sqrt(9)": func(x float64, y float64) float64 {
		return sqrt(9)
	},
	"sqrt(x+y)": func(x float64, y float64) float64 {
		return sqrt(x + y)
	},
	"sin(2/x)+cos(sqrt(x+y+1))": func(x float64, y float64) float64 {
		return sin(2/x) + cos(sqrt(x+y+1))
	},
	"cos(9)": func(x float64, y float64) float64 {
		return cos(9)
	},
	"sin(x+y)": func(x float64, y float64) float64 {
		return sin(x + y)
	},
	"sqrt(sqrt(sqrt(x)))": func(x float64, y float64) float64 {
		return sqrt(sqrt(sqrt(x)))
	},
	"1+2+(3+2*4+((((5+6*2)+7)+sqrt(8))+9)+10*sin(2-x+y/3))+11": func(x float64, y float64) float64 {
		return 1 + 2 + (3 + 2*4 + ((((5 + 6*2) + 7) + sqrt(8)) + 9) + 10*sin(2-x+y/3)) + 11
	},
	"1+x+(y+2*4+((((5+y*2)+7)+sqrt(x))+y)+10*sin(2-x+y/x))+y": func(x float64, y float64) float64 {
		return 1 + x + (y + 2*4 + ((((5 + y*2) + 7) + sqrt(x)) + y) + 10*sin(2-x+y/x)) + y
	},
	"1+2+(3+2*4+((((5+6*2)+7)+(8))+9)+10*sin(2-x+y/3))+11": func(x float64, y float64) float64 {
		return 1 + 2 + (3 + 2*4 + ((((5 + 6*2) + 7) + (8)) + 9) + 10*sin(2-x+y/3)) + 11
	},
	"1+2+(3+2*4+((((5+6*2)+7)+(8))+9)+10*(2-x+y/3))+11": func(x float64, y float64) float64 {
		return 1 + 2 + (3 + 2*4 + ((((5 + 6*2) + 7) + (8)) + 9) + 10*(2-x+y/3)) + 11
	},
	"(1+2+(3+2*4+((((5+6*2)+7)+(8))+9)+10*(2-x+y/3))+11)*(sin(x*y*2+1)*cos(1+2+x+y)+sin(2/x)+cos(sqrt(x+y+1)))": func(x float64, y float64) float64 {
		return (1 + 2 + (3 + 2*4 + ((((5 + 6*2) + 7) + (8)) + 9) + 10*(2-x+y/3)) + 11) * (sin(x*y*2+1)*cos(1+2+x+y) + sin(2/x) + cos(sqrt(x+y+1)))
	},
}

func TestJIT(t *testing.T) {
	defer func() {
		useConstFolding = true
		useCallDepth = true
		useRegisters = true
	}()
	for expr, want := range tests {
		for _, useConstFolding = range []bool{true, false} {
			for _, useCallDepth = range []bool{true, false} {
				for _, useRegisters = range []bool{true, false} {
					code, err := Compile(expr)
					if err != nil {
						t.Fatal(err)
					}
					for _, x := range []float64{3, -1e3, -123.4, -1, 0, 1, 123.4, 1e3} {
						for _, y := range []float64{5, -1e3, -123.4, -1, 0, 1, 123.4, 1e3} {
							have := code.Eval(x, y)
							if !equal(have, want(x, y)) {
								t.Errorf("%v with x=%v,y=%v: have %v, want: %v", expr, x, y, have, want(x, y))
							}
						}
					}

					code.Free()
				}
			}
		}
	}
}

func TestErrors(t *testing.T) {
	tests := []string{
		"",
		"notafunc(x)",
		"a.b",
		"a.",
		"1||2",
	}

	for _, test := range tests {
		_, err := Compile(test)
		if err == nil {
			t.Errorf("Compile %q: expected error, got nil", test)
		}
	}
}

func TestEval2D(t *testing.T) {
	nx := 100
	ny := 50
	dst := make([]float64, nx*ny)
	matrix := make([][]float64, ny)
	for iy := range matrix {
		matrix[iy] = dst[iy*nx : (iy+1)*nx]
	}
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

func sqrt(x float64) float64 { return math.Sqrt(x) }
func sin(x float64) float64  { return math.Sin(x) }
func cos(x float64) float64  { return math.Cos(x) }
