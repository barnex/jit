package jit

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

// C wrapper for generated machine code,
// allows us to use C calling conventions.

//#cgo LDFLAGS: -lm
//
//#include <math.h>
//
//void* func_sqrt = sqrt;
//void* func_sin  = sin;
//void* func_cos  = cos;
//
//double eval(void *code, double x, double y) {
//	double (*func)(double, double) = code;
//	return func(x, y);
//}
//
//void eval_2d(void *code, double *dst, double xmin, double xmax, int nx, double ymin, double ymax, int ny){
//	int i, ix, iy;
//	double x, y;
//	double (*func)(double, double) = code;
//	for(iy=0; iy<ny; iy++){
//		y = ymin + ((ymax-ymin)*(iy+0.5))/ny;
//		for(ix=0; ix<nx; ix++){
//			x = xmin + ((xmax-xmin)*(ix+0.5))/nx;
//			dst[iy*nx+ix] = func(x, y);
//		}
//	}
//}
import "C"

var funcs = map[string]uintptr{
	"sqrt": uintptr(C.func_sqrt),
	"sin":  uintptr(C.func_sin),
	"cos":  uintptr(C.func_cos),
}

// makeExecutable copies machine code to executable memory.
func makeExecutable(code []byte) ([]byte, error) {
	length := len(code)
	prot := unix.PROT_WRITE
	flags := unix.MAP_ANON | unix.MAP_PRIVATE
	const fd = -1
	const offset = 0
	mem, err := unix.Mmap(fd, offset, length, prot, flags)
	if err != nil {
		return nil, err
	}
	copy(mem, code)

	err = unix.Mprotect(mem, unix.PROT_READ|unix.PROT_EXEC)
	if err != nil {
		return nil, err
		unix.Munmap(mem)
	}
	return mem, nil
}

// call calls the machine code, which must hold a function of two float64s,
// and returns the result.
func eval(code []byte, x, y float64) float64 {
	return float64(C.eval(unsafe.Pointer(&code[0]), C.double(x), C.double(y)))
}

func eval2D(code []byte, dst []float64, xmin, xmax float64, nx int, ymin, ymax float64, ny int) {
	if len(dst) != nx*ny {
		panic(fmt.Sprintf("eval2D: nx=%v, ny=%v does not match len(dst)=%v", nx, ny, len(dst)))
	}
	C.eval_2d(unsafe.Pointer(&code[0]), (*C.double)(&dst[0]),
		C.double(xmin), C.double(xmax), C.int(nx),
		C.double(ymin), C.double(ymax), C.int(ny))
}

// Code stores JIT compiled machine code and allows to evaluate it.
type Code struct {
	instr []byte
}

// Eval executes the code, passing values for the variables x and y,
// and returns the result.
func (c *Code) Eval(x, y float64) float64 {
		if len(c.instr)==0{
		panic("eval called on nil code")	
		}
	return eval(c.instr, x, y)
}

func (c *Code) Eval2D(dst []float64, xmin, xmax float64, nx int, ymin, ymax float64, ny int) {
	eval2D(c.instr, dst, xmin, xmax, nx, ymin, ymax, ny)
}

// Free unmaps the code, after which Eval cannot be called anymore.
func (c *Code) Free() {
	unix.Munmap(c.instr)
	c.instr = nil
}
