package jit

import "testing"

const (
	nx = 1000
	ny = 1000
)

func BenchmarkCompileSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Compile("(x+y)*2 + (1+x) / y")
	}
}

func BenchmarkCompileBig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Compile("1+x+(3+y*4+((((x+y*2)+x)+sqrt(8))+y)+10*sin(2-x+y/3))+11")
	}
}

func BenchmarkSmallJIT(b *testing.B) {
	code, err := Compile("(x+y)*2 + (1+x) / y")
	if err != nil {
		b.Fatal(err)
	}

	dst := make([]float64, nx*ny)
	matrix := make([][]float64, ny)
	for iy := range matrix {
		matrix[iy] = dst[iy*nx : (iy+1)*nx]
	}
	n := b.N / (nx * ny)
	b.ResetTimer()
	for i := 0; i < n; i++ {
		code.Eval2D(dst, -1, 1, nx, -1, 1, ny)
	}
}

func BenchmarkSmallGo(b *testing.B) {
	dst := make([]float64, nx*ny)
	matrix := make([][]float64, ny)
	for iy := range matrix {
		matrix[iy] = dst[iy*nx : (iy+1)*nx]
	}
	n := b.N / (nx * ny)
	b.ResetTimer()
	for i := 0; i < n; i++ {
		eval2dGo(nativeGo, matrix, -1, 1, nx, -1, 1, ny)
	}
}

func BenchmarkBigJIT(b *testing.B) {
	code, err := Compile("1+x+(3+y*4+((((x+y*2)+x)+sqrt(8))+y)+10*sin(2-x+y/3))+11")
	if err != nil {
		b.Fatal(err)
	}

	dst := make([]float64, nx*ny)
	matrix := make([][]float64, ny)
	for iy := range matrix {
		matrix[iy] = dst[iy*nx : (iy+1)*nx]
	}
	n := b.N / (nx * ny)
	b.ResetTimer()
	for i := 0; i < n; i++ {
		code.Eval2D(dst, -1, 1, nx, -1, 1, ny)
	}
}

func BenchmarkBigGo(b *testing.B) {
	dst := make([]float64, nx*ny)
	matrix := make([][]float64, ny)
	for iy := range matrix {
		matrix[iy] = dst[iy*nx : (iy+1)*nx]
	}
	n := b.N / (nx * ny)
	b.ResetTimer()
	for i := 0; i < n; i++ {
		eval2dGo(nativeGoBig, matrix, -1, 1, nx, -1, 1, ny)
	}
}

func nativeGo(x, y float64) float64 {
	return (x+y)*2 + (1+x)/y
}

func nativeGoBig(x, y float64) float64 {
	return 1 + x + (3 + y*4 + ((((x + y*2) + x) + sqrt(8)) + y) + 10*sin(2-x+y/3)) + 11
}

func eval2dGo(f func(float64, float64) float64, dst [][]float64, xmin, xmax float64, nx int, ymin, ymax float64, ny int) {
	var ix, iy int
	var x, y float64
	for iy = 0; iy < ny; iy++ {
		y = ymin + ((ymax-ymin)*(float64(iy)+0.5))/float64(ny)
		for ix = 0; ix < nx; ix++ {
			x = xmin + ((xmax-xmin)*(float64(ix)+0.5))/float64(nx)
			dst[iy][ix] = f(x, y)
		}
	}
}
