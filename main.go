package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
)

//int run(void *code) {
//  int (*func)() = code;
//  return func();
//}
import "C"
import "unsafe"

func main() {
	// mov eax, $5
	// ret
	code := []byte{0xb8, 0x05, 0x00, 0x00, 0x00, 0xc3}
	mem, err := makeExecutable(code)
	if err != nil {
		fatal(err)
	}
	defer unix.Munmap(mem)
	fmt.Println(C.run(unsafe.Pointer(&mem[0])))
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
