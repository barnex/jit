/*
Package jit provides a toy just-in-time compiler for arithmetic expressions of variables x and y. E.g.:
	code, err := Compile("(x+1) * (y+2) / 3")
	x, y := 1.0, 2.0
	z := code.Eval(x, y)

Works on 64-bit linux only.

Inspired by the book "The Go Programming Language" by Alan A. A. Donovan and Brian W. Kernighan, section 7.9: Example: Expression Evaluator.
*/
package jit

import (
	"bytes"
	"fmt"
	"os"
)

func (e *variable) compile(b *buf) {
	switch e.name {
	default:
		panic("undefined variable:" + e.name)
	case "x":
		b.emit(mov_x_rbp_rax(-8), mov_rax_xmm0)
	case "y":
		b.emit(mov_x_rbp_rax(-16), mov_rax_xmm0)
	}
}

func (e *constant) compile(b *buf) {
	b.emit(mov_float_rax(e.value), mov_rax_xmm0)
}

func (e *binexpr) compileArgs(b *buf) {
	e.y.compile(b) // y in xmm0

	// stash result
	reg := -1
	if !b.hasCall[e.x] {
		reg = b.allocReg()
	} else {
		b.nStackSpill++
	}
	if reg == -1 {
		b.emit(mov_xmm0_rax, push_rax)
	} else {
		b.emit(mov_xmm(0, reg))
	}

	e.x.compile(b) // x in xmm0

	if reg == -1 {
		b.emit(pop_rax, mov_rax_xmm1) // y in xmm1
	} else {
		b.emit(mov_xmm(reg, 1))
	}
	b.freeReg(reg)
}

func (e *add) compile(b *buf) {
	e.compileArgs(b)
	b.emit(add_xmm1_xmm0)
}

func (e *sub) compile(b *buf) {
	e.compileArgs(b)
	b.emit(sub_xmm1_xmm0)
}

func (e *quo) compile(b *buf) {
	e.compileArgs(b)
	b.emit(div_xmm1_xmm0)
}

func (e *mul) compile(b *buf) {
	e.compileArgs(b)
	b.emit(mul_xmm1_xmm0)
}

func (e *callexpr) compile(b *buf) {
	if len(e.args) != 1 {
		panic(fmt.Sprintf("%v arguments not supported", len(e.args)))
	}
	fptr := funcs[e.fun]
	if fptr == 0 {
		panic(fmt.Sprintf("undefined:", e.fun))
	}

	e.args[0].compile(b)
	b.emit(mov_uint_rax(fptr), call_rax)
}

// buf accumulates machine code.
type buf struct {
	bytes.Buffer
	usedReg                    [8]bool
	nRegistersHit, nStackSpill int
	hasCall                    map[expr]bool
}

func (b *buf) allocReg() int {
	for i := 2; i < len(b.usedReg); i++ {
		if !b.usedReg[i] {
			b.usedReg[i] = true
			b.nRegistersHit++
			return i
		}
	}
	b.nStackSpill++
	return -1
}

func (b *buf) freeReg(reg int) {
	if reg == -1 {
		return
	}
	if !b.usedReg[reg] {
		panic(fmt.Sprint("register double free", reg))
	}
	b.usedReg[reg] = false
}

// Compile compiles an arithmetic expression, which may contain the variables x and y. E.g.:
// 	(x+1) * (y-2)
// If no longer needed, the returned code must be explicitly freed with Free().
func Compile(ex string) (c *Code, e error) {
	root, err := Parse(ex)
	if err != nil {
		return nil, err
	}

	b := buf{hasCall: make(map[expr]bool)}
	recordCalls(root, b.hasCall)

	b.emit(push_rbp, mov_rsp_rbp)            // function preamble
	b.emit(sub_rsp(16))                      // stack space for x, y
	b.emit(mov_xmm0_rax, mov_rax_x_rbp(-8))  // x on stack
	b.emit(mov_xmm1_rax, mov_rax_x_rbp(-16)) // y on stack
	root.compile(&b)                         // function body (jit code)
	b.emit(add_rsp(16))                      // free stack space for x,y
	b.emit(pop_rbp, ret)                     // return from function

	b.dump("b.out")
	fmt.Println(ex, ":", b.nRegistersHit, "reg hits,", b.nStackSpill, "stack spills")

	instr, err := makeExecutable(b.Bytes())
	if err != nil {
		return nil, err
	}
	return &Code{instr}, nil
}

// emit writes machine code to the buffer.
func (b *buf) emit(ops ...[]byte) {
	for _, op := range ops {
		b.Write(op)
	}
}

func (b *buf) dump(fname string) {
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(b.Bytes())
}
