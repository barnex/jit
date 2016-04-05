package jit

// C wrapper for generated machine code,
// allows us to use C calling conventions
// and not worry about providing sufficient stack space.

import (
	"fmt"
	"unsafe"
)

//#cgo LDFLAGS: -lm
//#include "shim.h"
import "C"

var funcs = map[string]uintptr{
	"acos":  uintptr(C.func_acos),
	"asin":  uintptr(C.func_asin),
	"atan":  uintptr(C.func_atan),
	"cos":   uintptr(C.func_cos),
	"cosh":  uintptr(C.func_cosh),
	"sin":   uintptr(C.func_sin),
	"sinh":  uintptr(C.func_sinh),
	"tan":   uintptr(C.func_tan),
	"tanh":  uintptr(C.func_tanh),
	"exp":   uintptr(C.func_exp),
	"log":   uintptr(C.func_log),
	"log10": uintptr(C.func_log10),
	"sqrt":  uintptr(C.func_sqrt),
	"fabs":  uintptr(C.func_fabs),
}

// call calls the machine code, which must hold a function of two float64s,
// and returns the result.
func eval(code []byte, x, y float64) float64 {
	return float64(C.eval(unsafe.Pointer(&code[0]), C.double(x), C.double(y)))
}

// callCFunc calls a C function with one double argument.
// Used for constant folding, like sqrt(2).
func callCFunc(f uintptr, x float64) float64 {
	return float64(C.call_func(unsafe.Pointer(f), C.double(x)))
}

// eval2D evaluates the code nx * ny times
// while varying x between xmin, xmax and y between ymin, ymax.
// The result is stored in dst.
func eval2D(code []byte, dst []float64, xmin, xmax float64, nx int, ymin, ymax float64, ny int) {
	if len(dst) != nx*ny {
		panic(fmt.Sprintf("eval2D: nx=%v, ny=%v does not match len(dst)=%v", nx, ny, len(dst)))
	}
	C.eval_2d(unsafe.Pointer(&code[0]), (*C.double)(&dst[0]),
		C.double(xmin), C.double(xmax), C.int(nx),
		C.double(ymin), C.double(ymax), C.int(ny))
}
