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
//  for(int i=0; i<999; i++){
//  	func(x, y);
//  }
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

// Code stores JIT compiled machine code and allows to evaluate it.
type Code struct {
	instr []byte
}

// Eval executes the code, passing values for the variables x and y,
// and returns the result.
func (c *Code) Eval(x, y float64) float64 {
	return call(c.instr, x, y)
}

// Free unmaps the code, after which Eval cannot be called anymore.
func (c *Code) Free() {
	unix.Munmap(c.instr)
	c.instr = nil
}

