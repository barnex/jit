package main

import (
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

func main() {
	code := []byte{
		0x55,             // push   %rbp
		0x48, 0x89, 0xe5, // mov    %rsp,%rbp
		0xf2, 0x0f, 0x58, 0xc1, // addsd  %xmm1,%xmm0
		0x48, 0xb8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf0, 0x3f, // movabs $0x3ff0000000000000,%rax
		0x66, 0x48, 0x0f, 0x6e, 0xc0, // movq   %rax,%xmm0
		0x5d, //  pop    %rbp
		0xc3, // retq
	}
	mem, err := makeExecutable(code)
	if err != nil {
		fatal(err)
	}
	defer unix.Munmap(mem)
	result := C.run(unsafe.Pointer(&mem[0]), 3, 40)
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
