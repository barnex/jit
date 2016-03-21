package jit

import "unsafe"

// amd64 machine code
var (
	push_rbp      = []byte{0x55}
	push_rax      = []byte{0x50}
	pop_rbp       = []byte{0x5d}
	pop_rax       = []byte{0x58}
	ret           = []byte{0xc3}
	mov_rsp_rbp   = []byte{0x48, 0x89, 0xe5}
	mov_xmm0_rax  = []byte{0x66, 0x48, 0x0f, 0x7e, 0xc0}
	mov_xmm1_rax  = []byte{0x66, 0x48, 0x0f, 0x7e, 0xc8}
	mov_xmm2_rax  = []byte{0x66, 0x48, 0x0f, 0x7e, 0xd0}
	mov_xmm3_rax  = []byte{0x66, 0x48, 0x0f, 0x7e, 0xd8}
	mov_rax_xmm0  = []byte{0x66, 0x48, 0x0f, 0x6e, 0xc0}
	mov_rax_xmm1  = []byte{0x66, 0x48, 0x0f, 0x6e, 0xc8}
	mov_rax_xmm2  = []byte{0x66, 0x48, 0x0f, 0x6e, 0xd0}
	mov_rax_xmm3  = []byte{0x66, 0x48, 0x0f, 0x6e, 0xd8}
	add_xmm3_xmm2 = []byte{0xf2, 0x0f, 0x58, 0xd3}
	mul_xmm3_xmm2 = []byte{0xf2, 0x0f, 0x59, 0xd3}
	sub_xmm3_xmm2 = []byte{0xf2, 0x0f, 0x5c, 0xd3}
	div_xmm3_xmm2 = []byte{0xf2, 0x0f, 0x5e, 0xd3}
	call_rax      = []byte{0xff, 0xd0}
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
func mov_x_rsp_rax(x int32) []byte {
	return append([]byte{0x48, 0x8b, 0x84, 0x24}, int32Bytes(x)...)
}

// returns code for movq %rax,x(%rsp)
func mov_rax_x_rsp(x int32) []byte {
	return append([]byte{0x48, 0x89, 0x84, 0x24}, int32Bytes(x)...)
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
