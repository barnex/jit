package jit

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// MakeExecutable copies machine code to executable memory.
// The returned memory must be freed with unix.Munmap.
func MakeExecutable(code []byte) ([]byte, error) {
	exec, err := unix.Mmap(-1, 0, len(code), unix.PROT_WRITE, unix.MAP_ANON|unix.MAP_PRIVATE)
	if err != nil {
		return nil, err
	}
	copy(exec, code)
	err = unix.Mprotect(exec, unix.PROT_READ|unix.PROT_EXEC)
	if err != nil {
		unix.Munmap(exec)
		return nil, err
	}
	return exec, nil
}

func MakeFunc(instr []byte) func(float64, float64) float64 {
	fStruct := &funcData{funcPtr: unsafe.Pointer(&instr[0])}
	return *(*func(float64, float64) float64)(unsafe.Pointer(&fStruct))
}

// funcData is a data structure that a Go function value may refer to.
// This is tied to the Go implementation details. See
// 	https://docs.google.com/document/d/1bMwCey-gmqZVTpRax-ESeVuZGmjwbocYs1iHplK-cjo/pub
type funcData struct {
	funcPtr unsafe.Pointer // C-style function pointer
	_       [8]byte        // Closure data. Unused but being addressed.
}
