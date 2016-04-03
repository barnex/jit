package jit

import (
	"bytes"
	"fmt"
	"os"
)

var (
	useRegisters    = true
	useCallDepth    = true
	useConstFolding = true
)

// Compile compiles an arithmetic expression, which may contain the variables x and y. E.g.:
// 	(x+1) * (y-2)
// If no longer needed, the returned code must be explicitly freed with Free().
func Compile(ex string) (c *Code, e error) {
	root, err := Parse(ex)
	if err != nil {
		return nil, err
	}

	if useConstFolding {
		root = FoldConst(root)
	}

	b := buf{hasCall: make(map[expr]bool), callDepth: make(map[expr]int)}
	recordCalls(root, b.hasCall)
	if useCallDepth {
		recordDepth(root, b.callDepth)
	}

	b.emit(push_rbp, mov_rsp_rbp)            // function preamble
	b.emit(sub_rsp(16))                      // stack space for x, y
	b.emit(mov_xmm0_rax, mov_rax_x_rbp(-8))  // x on stack
	b.emit(mov_xmm1_rax, mov_rax_x_rbp(-16)) // y on stack
	b.compileExpr(root)                      // function body (jit code)
	b.emit(add_rsp(16))                      // free stack space for x,y
	b.emit(pop_rbp, ret)                     // return from function

	b.dump("b.out")
	fmt.Println(ex, ":", b.nRegistersHit, "reg hits,", b.maxReg, "highest register used, ", b.nStackSpill, "stack spills")

	instr, err := makeExecutable(b.Bytes())
	if err != nil {
		return nil, err
	}
	return &Code{instr}, nil
}

// buf accumulates machine code.
type buf struct {
	bytes.Buffer
	usedReg                            [8]bool
	nRegistersHit, nStackSpill, maxReg int
	hasCall                            map[expr]bool
	callDepth                          map[expr]int
}

// emit writes machine code to the buffer.
func (b *buf) emit(ops ...[]byte) {
	for _, op := range ops {
		b.Write(op)
	}
}

func (b *buf) stash(destroyRegs bool) int {
	// stash result
	reg := -1
	if !destroyRegs {
		reg = b.allocReg()
	} else {
		b.nStackSpill++
	}
	if reg == -1 {
		b.emit(mov_xmm0_rax, push_rax)
	} else {
		b.emit(mov_xmm(0, reg))
	}
	return reg
}

func (b *buf) unstash(reg, dest int) {
	switch {
	case reg == -1 && dest == 1:
		b.emit(pop_rax, mov_rax_xmm1)
	case reg == -1 && dest == 0:
		b.emit(pop_rax, mov_rax_xmm0)
	case reg != -1:
		b.emit(mov_xmm(reg, dest))
	default:
		panic("bug")
	}
	b.freeReg(reg)
}

func (b *buf) allocReg() int {
	if !useRegisters {
		b.nStackSpill++
		return -1
	}
	for i := 2; i < len(b.usedReg); i++ {
		if !b.usedReg[i] {
			b.usedReg[i] = true
			b.nRegistersHit++
			if i > b.maxReg {
				b.maxReg = i
			}
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

func (b *buf) dump(fname string) {
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(b.Bytes())
}

func (b *buf) compileExpr(e expr) {
	switch e := e.(type) {
	default:
		panic(fmt.Sprintf("compileExpr %T", e))
	case binexpr:
		b.compileBinexpr(e)
	case callexpr:
		b.compileCallexpr(e)
	case constant:
		b.compileConstant(e)
	case variable:
		b.compileVariable(e)
	}
}

func (b *buf) compileVariable(e variable) {
	switch e.name {
	default:
		panic("undefined variable:" + e.name)
	case "x":
		b.emit(mov_x_rbp_rax(-8), mov_rax_xmm0)
	case "y":
		b.emit(mov_x_rbp_rax(-16), mov_rax_xmm0)
	}
}

func (b *buf) compileConstant(e constant) {
	b.emit(mov_float_rax(e.value), mov_rax_xmm0)
}

func (b *buf) compileBinexpr(e binexpr) {
	// Determine which side of the binary expression to evaluate first:
	//  * prefer deeper branch first, so we use least registers
	//  * however, avoid function calls in the second branch,
	// 	  as those destroy the registers.
	var first, second expr
	if b.callDepth[e.x] > b.callDepth[e.y] && !b.hasCall[e.y] {
		first, second = e.x, e.y
	} else {
		first, second = e.y, e.x
	}

	b.compileExpr(first)
	stash := b.stash(b.hasCall[second])
	b.compileExpr(second)

	// Move the results back:
	// y -> xmm0
	// x -> xmm1
	if first == e.y {
		b.unstash(stash, 1)
	} else {
		b.emit(mov_xmm(0, 1))
		b.unstash(stash, 0)
	}

	switch e.op {
	case "+":
		b.emit(add_xmm1_xmm0)
	case "-":
		b.emit(sub_xmm1_xmm0)
	case "*":
		b.emit(mul_xmm1_xmm0)
	case "/":
		b.emit(div_xmm1_xmm0)
	default:
		panic(e.op)
	}
}

func (b *buf) compileCallexpr(e callexpr) {
	//if len(e.args) != 1 {
	//	panic(fmt.Sprintf("%v arguments not supported", len(e.args)))
	//}
	fptr := funcs[e.fun]
	if fptr == 0 {
		panic(fmt.Sprintf("undefined:", e.fun))
	}

	b.compileExpr(e.arg)
	b.emit(mov_uint_rax(fptr), call_rax)
}
