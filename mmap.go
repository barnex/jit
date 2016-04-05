package jit

import "golang.org/x/sys/unix"

// functionality for making the generated machine code executable, and executing it.

// MakeExecutable copies machine code to executable memory.
func MakeExecutable(code []byte) ([]byte, error) {
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

// Code stores JIT compiled machine code and allows to evaluate it.
type Code struct {
	instr []byte
}

// Eval executes the code, passing values for the variables x and y,
// and returns the result.
func (c *Code) Eval(x, y float64) float64 {
	if len(c.instr) == 0 {
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
