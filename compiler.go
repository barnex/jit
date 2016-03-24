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

// buf accumulates machine code.
type buf struct {
	bytes.Buffer
	usedReg                    [8]bool
	nRegistersHit, nStackSpill int
}

func (b *buf) allocReg() int {
	//for i:=2; i<len(b.usedReg); i++{
	//	if !b.usedReg[i]{
	//			b.usedReg[i]=true
	//			b.nRegistersHit++
	//			fmt.Println("alloc register", i)
	//			return i
	//	}
	//}
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
func Compile(expr string) (c *Code, e error) {
	root, err := Parse(expr)
	if err != nil {
		return nil, err
	}

	var b buf
	b.emit(push_rbp, mov_rsp_rbp)            // function preamble
	b.emit(sub_rsp(16))                      // stack space for x, y
	b.emit(mov_xmm0_rax, mov_rax_x_rbp(-8))  // x on stack
	b.emit(mov_xmm1_rax, mov_rax_x_rbp(-16)) // y on stack
	root.compile(&b)                         // function body (jit code)
	b.emit(add_rsp(16))                      // free stack space for x,y
	b.emit(pop_rbp, ret)                     // return from function

	b.dump("b.out")
	fmt.Println(b.nRegistersHit, "register hits,", b.nStackSpill, "stack spills")

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
