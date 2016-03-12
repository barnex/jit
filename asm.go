package main

import (
	"bytes"
	"unsafe"
)

type Buf struct {
	bytes.Buffer
	instr []byte
}

var (
	mov_xmm0_rax = []byte{0x66, 0x48, 0x0f, 0x7e, 0xc0}
	mov_xmm1_rax = []byte{0x66, 0x48, 0x0f, 0x7e, 0xc8}
	mov_rax_xmm0 = []byte{0x66, 0x48, 0x0f, 0x6e, 0xc0}
	push_rax     = []byte{0x50}
	pop_rax      = []byte{0x58}

	openFunc  = []byte{0x55, 0x48, 0x89, 0xe5} // push  %rbp; mov  %rsp,%rbp
	closeFunc = []byte{0x5d, 0xc3}             // pop   %rbp; retq
	//add10     = []byte{0xf2, 0x0f, 0x58, 0xc1}             // addsd %xmm1,%xmm0
	//sub10     = []byte{0xf2, 0x0f, 0x5c, 0xc1}             // subsd %xmm1,%xmm0
	//mul10     = []byte{0xf2, 0x0f, 0x59, 0xc1}             // subsd %xmm1,%xmm0
	//div10     = []byte{0xf2, 0x0f, 0x5e, 0xc1}             // divsd %xmm1,%xmm0
	//push0     = []byte{0x66, 0x48, 0x0f, 0x7e, 0xc0, 0x50} // movq  %xmm0,%rax; push %rax
	//pop0      = []byte{0x58, 0x66, 0x48, 0x0f, 0x6e, 0xc0} // pop   %rax; movq %rax,%xmm0
	//pop1      = []byte{0x58, 0x66, 0x48, 0x0f, 0x6e, 0xc8} // pop   %rax; movq %rax,%xmm1
)

func (b *Buf) emit(ops ...[]byte) {
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

// load immediate x into rax
func mov_imm_rax(x float64) []byte {
	movabs := []byte{0x48, 0xb8} // movabs ...
	imm := *((*[8]byte)(unsafe.Pointer(&x)))
	return cat(movabs, imm[:])
}
