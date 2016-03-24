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
	"go/ast"
	"go/token"
	"os"
	"reflect"
	"strconv"
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
	root.compile(&b)                          // function body (jit code)
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

// emitExpr compiles an expression and stores the machine code.
func (b *buf) emitExpr(e ast.Expr) {
	switch e := e.(type) {
	default:
		panic(err(e.Pos(), "syntax error:", typ(e)))
	case *ast.Ident:
		b.emitIdent(e)
	case *ast.BasicLit:
		b.emitBasicLit(e)
	case *ast.BinaryExpr:
		b.emitBinaryExpr(e)
	case *ast.ParenExpr:
		b.emitExpr(e.X)
	//case *ast.UnaryExpr:
	//	return w.compileUnaryExpr(e)
	case *ast.CallExpr:
		b.emitCall(e)
	}
}

func (b *buf) emitCall(e *ast.CallExpr) {
	name := e.Fun.(*ast.Ident).Name
	fptr := funcs[name]
	if fptr == 0 {
		panic(err(e.Pos(), "undefined:", name))
	}

	if len(e.Args) != 1 {
		panic(err(e.Pos(), "need one argument"))
	}
	b.emitExpr(e.Args[0])
	b.emit(mov_uint_rax(fptr), call_rax)
}

// emitIdent compiles an identifier (x or y) and stores the machine code.
func (b *buf) emitIdent(e *ast.Ident) {
	switch e.Name {
	default:
		panic(err(e.Pos(), "undefined variable:", e.Name))
	case "x":
		b.emit(mov_x_rbp_rax(-8), mov_rax_xmm0)
	case "y":
		b.emit(mov_x_rbp_rax(-16), mov_rax_xmm0)
	}
}

// emitBasicLit compiles a number literal, e.g. "2" and stores the machine code.
func (b *buf) emitBasicLit(e *ast.BasicLit) {
	switch e.Kind {
	default:
		panic(err(e.Pos(), "syntax error:", e.Value, "(", typ(e), ")"))
	case token.FLOAT, token.INT:
		v, err := strconv.ParseFloat(e.Value, 64)
		if err != nil {
			panic(err)
		}
		b.emit(mov_float_rax(v), mov_rax_xmm0)
	}
}

// emitBinaryExpr compiles a binary expression, e.g. x+1, and stores the machine code.
func (b *buf) emitBinaryExpr(n *ast.BinaryExpr) {
	b.emitExpr(n.Y) // y in xmm0

	// stash result
	reg := b.allocReg()
	if reg == -1 {
		b.emit(mov_xmm0_rax, push_rax)
	} else {
		b.emit(mov_xmm(0, reg))
	}

	b.emitExpr(n.X) // x in xmm0
	if reg == -1 {
		b.emit(pop_rax, mov_rax_xmm1) // x in xmm1
	} else {
		b.emit(mov_xmm(reg, 1))
	}
	b.freeReg(reg)

	switch n.Op {
	default:
		panic(err(n.Pos(), "syntax error:", n.Op))
	case token.ADD:
		b.emit(add_xmm1_xmm0)
	case token.SUB:
		b.emit(sub_xmm1_xmm0)
	case token.MUL:
		b.emit(mul_xmm1_xmm0)
	case token.QUO:
		b.emit(div_xmm1_xmm0)
	}

	// result in xmm0
}

func typ(x interface{}) reflect.Type {
	return reflect.TypeOf(x)
}
func err(pos token.Pos, msg ...interface{}) error {
	return fmt.Errorf("%v: %v", pos, fmt.Sprintln(msg...))
}

func (b *buf) dump(fname string) {
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(b.Bytes())
}
