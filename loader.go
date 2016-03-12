package main

import "golang.org/x/sys/unix"

//double run(void *code, double x, double y) {
//  double (*func)(double, double) = code;
//  return func(x, y);
//}
import "C"
import "unsafe"

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

func (b *Buf) call(x, y float64) float64 {
	return float64(C.run(unsafe.Pointer(&b.instr[0]), C.double(x), C.double(y)))
}

func (b *Buf) Free() {
	unix.Munmap(b.instr)
	b.instr = nil
	b.Buffer.Reset()
}
