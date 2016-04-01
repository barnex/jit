package jit

import "testing"

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
