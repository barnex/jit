package jit

import (
	"runtime"
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
		runtime.GC()
		code, err := Compile(test.expr)
		runtime.GC()
		if err != nil {
			t.Error(err)
			continue
		}
		have := code.Eval(test.x, test.y)
		runtime.GC()
		code.Free()
		runtime.GC()
		if have != test.want {
			t.Errorf("%v with x=%v,y=%v: have %v, want: %v", test.expr, test.x, test.y, have, test.want)
		}
	}

	//note: calling runtime.GC() in an attempt to detect some types of memory corruption
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
