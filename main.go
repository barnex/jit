package main

import (
	"bytes"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
)

//double run(void *code, double x, double y) {
//  double (*func)(double, double) = code;
//  return func(x, y);
//}
import "C"
import "unsafe"

var (
	openFunc = []byte{
		0x55,             // push   %rbp
		0x48, 0x89, 0xe5, // mov    %rsp,%rbp
	}
	closeFunc = []byte{
		0x5d, // pop    %rbp
		0xc3, // retq
	}
	add10 = []byte{
		0xf2, 0x0f, 0x58, 0xc1, // addsd  %xmm1,%xmm0
	}
	sub10 = []byte{
		0xf2, 0x0f, 0x5c, 0xc1, // subsd  %xmm1,%xmm0
	}
	mul10 = []byte{
		0xf2, 0x0f, 0x59, 0xc1, // subsd  %xmm1,%xmm0
	}
	div10 = []byte{
		0xf2, 0x0f, 0x5e, 0xc1, // divsd  %xmm1,%xmm0
	}
	push0 = []byte{
		0x66, 0x48, 0x0f, 0x7e, 0xc0, // movq   %xmm0,%rax
		0x50, // push   %rax
	}
	pop0 = []byte{
		0x58,                         // pop    %rax
		0x66, 0x48, 0x0f, 0x6e, 0xc0, // movq   %rax,%xmm0
	}
	pop1 = []byte{
		0x58,                         // pop    %rax
		0x66, 0x48, 0x0f, 0x6e, 0xc8, // movq   %rax,%xmm1
	}
)

type codeBuffer struct{ bytes.Buffer }

func (b *codeBuffer) emit(ops ...[]byte) {
	for _, op := range ops {
		b.Write(op)
	}
}

func cat(ops ...[]byte) []byte {
	var cat []byte
	for _, op := range ops {
		cat = append(cat, op...)
	}
	return cat
}

// load immediate x into xmm0
func imm0(x float64) []byte {
	movabs := []byte{0x48, 0xb8} // movabs ...
	imm := *((*[8]byte)(unsafe.Pointer(&x)))
	abs0 := []byte{0x66, 0x48, 0x0f, 0x6e, 0xc0} // movq   %rax,%xmm0
	return cat(movabs, imm[:], abs0)
}

func main() {
	var code codeBuffer

	code.emit(openFunc, mul10, closeFunc)

	mem, err := makeExecutable(code.Bytes())
	if err != nil {
		fatal(err)
	}
	defer unix.Munmap(mem)
	result := C.run(unsafe.Pointer(&mem[0]), 3, 30)
	fmt.Println(result)
}

func makeExecutable(code []byte) ([]byte, error) {
	length := len(code)
	prot := unix.PROT_WRITE | unix.PROT_EXEC
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
	}
	return mem, nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
