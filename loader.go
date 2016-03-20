package jit

import (
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
//double run(void *code, double x, double y) {
//  double (*func)(double, double) = code;
//  return func(x, y);
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
func call(code []byte, x, y float64) float64 {
	return float64(C.run(unsafe.Pointer(&code[0]), C.double(x), C.double(y)))
}
