package jit

// This file provides assembly for a handful of x86-64 instructions.

import "unsafe"

var (
	call_rax      = []byte{0xff, 0xd0}                   // callq *%rax
	mov_rax_xmm0  = []byte{0x66, 0x48, 0x0f, 0x6e, 0xc0} // mov %rax,%xmm0
	mov_rax_xmm1  = []byte{0x66, 0x48, 0x0f, 0x6e, 0xc8} // mov %rax,%xmm1
	mov_rsp_rbp   = []byte{0x48, 0x89, 0xe5}             // mov %rsp,%rbp
	mov_xmm0_rax  = []byte{0x66, 0x48, 0x0f, 0x7e, 0xc0} // mov %xmm0,%rax
	mov_xmm1_rax  = []byte{0x66, 0x48, 0x0f, 0x7e, 0xc8} // mov %xmm1,%rax
	pop_rax       = []byte{0x58}                         // pop %rax
	pop_rbp       = []byte{0x5d}                         // pop %rbp
	push_rax      = []byte{0x50}                         // push %rax
	push_rbp      = []byte{0x55}                         // push %rbp
	ret           = []byte{0xc3}                         // ret
	add_xmm1_xmm0 = []byte{0xf2, 0x0f, 0x58, 0xc1}       // addsd  %xmm1,%xmm0
	sub_xmm1_xmm0 = []byte{0xf2, 0x0f, 0x5c, 0xc1}       // subsd  %xmm1,%xmm0
	mul_xmm1_xmm0 = []byte{0xf2, 0x0f, 0x59, 0xc1}       // mulsd  %xmm1,%xmm0
	div_xmm1_xmm0 = []byte{0xf2, 0x0f, 0x5e, 0xc1}       // divsd  %xmm1,%xmm0
)

// returns code for movq $x,%rax
func mov_float_rax(x float64) []byte {
	return mov_imm_rax(float64Bytes(x))
}

// returns code for movq $x,%rax
func mov_uint_rax(x uintptr) []byte {
	return mov_imm_rax(uintptrBytes(x))
}

// returns code for movq $x,%rax
func mov_imm_rax(x []byte) []byte {
	return append([]byte{0x48, 0xb8}, x...)
}

// returns code for subq $x,%rsp
func sub_rsp(x uint32) []byte {
	return append([]byte{0x48, 0x81, 0xec}, uint32Bytes(x)...)
}

// returns code for addq $x,%rsp
func add_rsp(x uint32) []byte {
	return append([]byte{0x48, 0x81, 0xc4}, uint32Bytes(x)...)
}

// returns code for movq x(%rsp),%rax
func mov_x_rbp_rax(x int32) []byte {
	return append([]byte{0x48, 0x8b, 0x85}, int32Bytes(x)...)
}

// returns code for movq %rax,x(%rbp)
func mov_rax_x_rbp(x int32) []byte {
	return append([]byte{0x48, 0x89, 0x85}, int32Bytes(x)...)
}

// returns code for movq %xmm0,x(%rbp)
func mov_xmm0_x_rbp(x int32) []byte {
	return append([]byte{0x66, 0x0f, 0xd6, 0x85}, int32Bytes(x)...)
}

// returns code for movq x(%rbp),%xmm0
func mov_x_rbp_xmm0(x int32) []byte {
	return append([]byte{0xf3, 0x0f, 0x7e, 0x85}, int32Bytes(x)...)
}

// returns code for movq %r1,%r2
func mov_xmm(r1, r2 int) []byte {
	if r1 > 7 || r2 > 7 {
		panic("mov_xmm: bad register")
	}
	regs := byte(0xc0) | byte(r2)<<3 | byte(r1)
	return []byte{0xf3, 0x0f, 0x7e, regs}
}

func uint32Bytes(x uint32) []byte {
	return (*((*[4]byte)(unsafe.Pointer(&x))))[:]
}

func int32Bytes(x int32) []byte {
	return (*((*[4]byte)(unsafe.Pointer(&x))))[:]
}

func uintptrBytes(x uintptr) []byte {
	return (*((*[8]byte)(unsafe.Pointer(&x))))[:]
}

func float64Bytes(x float64) []byte {
	return (*((*[8]byte)(unsafe.Pointer(&x))))[:]
}
