package jit

import "testing"

func TestSSA(t *testing.T) {
	for expr, want := range tests {
		code, err := SSADump(expr)
		if err != nil {
			t.Fatal(err)
		}
		for _, x := range []float64{3, -1e9, -123.4, -1, 0, 1, 123.4, 1e9} {
			for _, y := range []float64{5, -1e9, -123.4, -1, 0, 1, 123.4, 1e9} {
				have := code.Eval(x, y)
				if !equal(have, want(x, y)) {
					t.Errorf("%v with x=%v,y=%v: have %v, want: %v", expr, x, y, have, want(x, y))
				}
			}
		}

		code.Free()
	}
}
