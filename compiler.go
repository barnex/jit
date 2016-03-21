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
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strconv"

	"golang.org/x/sys/unix"
)

// Code stores JIT compiled machine code and allows to evaluate it.
type Code struct {
	instr []byte
}

// Compile compiles an arithmetic expression, which may contain the variables x and y. E.g.:
// 	(x+1) * (y-2)
// If no longer needed, the returned code must be explicitly freed with Free().
func Compile(expr string) (c *Code, e error) {
	root, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, fmt.Errorf(`parse "%s": %v`, expr, err)
	}

	// catch bailout panics
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
			c = nil
		}
	}()

	var b buf
	b.emit(push_rbp, mov_rsp_rbp) // function preamble
	//b.emit(sub_rbp(16))
	b.emitExpr(root)              // function body (jit code)
	b.emit(pop_rax, mov_rax_xmm0) // result from stack returned via xmm0
	//b.emit(add_rbp(16))
	b.emit(pop_rbp, ret) // return from function

	b.dump("b.out")

	instr, err := makeExecutable(((*bytes.Buffer)(&b)).Bytes())
	if err != nil {
		return nil, err
	}
	return &Code{instr}, nil
}

func (b *buf) dump(fname string) {
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(((*bytes.Buffer)(b)).Bytes())
}

// Eval executes the code, passing values for the variables x and y,
// and returns the result.
func (c *Code) Eval(x, y float64) float64 {
	return call(c.instr, x, y)
}

// Free unmaps the code, after which Eval cannot be called anymore.
func (c *Code) Free() {
	unix.Munmap(c.instr)
	c.instr = nil
}

// buf accumulates machine code.
type buf bytes.Buffer

// emit writes machine code to the buffer.
func (b *buf) emit(ops ...[]byte) {
	for _, op := range ops {
		((*bytes.Buffer)(b)).Write(op)
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
	if len(e.Args) != 1 {
		panic(err(e.Pos(), "need one argument"))
	}
	b.emitExpr(e.Args[0])

	name := e.Fun.(*ast.Ident).Name
	fptr := funcs[name]
	if fptr == 0 {
		panic(err(e.Pos(), "undefined:", name))
	}

	b.emit(pop_rax, mov_rax_xmm2)
	b.emit(mov_xmm0_rax, push_rax)
	b.emit(mov_xmm1_rax, push_rax)
	b.emit(mov_xmm2_rax)
	b.emit(mov_rax_xmm0)
	b.emit(mov_uint_rax(fptr), call_rax)
	b.emit(mov_xmm0_rax, mov_rax_xmm2)
	b.emit(pop_rax, mov_rax_xmm1)
	b.emit(pop_rax, mov_rax_xmm0)
	b.emit(mov_xmm2_rax)
	b.emit(mov_xmm2_rax, push_rax)
}

// emitIdent compiles an identifier (x or y) and stores the machine code.
func (b *buf) emitIdent(e *ast.Ident) {
	switch e.Name {
	default:
		panic(err(e.Pos(), "undefined variable:", e.Name))
	case "x":
		b.emit(mov_xmm0_rax, push_rax)
	case "y":
		b.emit(mov_xmm1_rax, push_rax)
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
		b.emit(mov_float_rax(v), push_rax)
	}
}

// emitBinaryExpr compiles a binary expression, e.g. x+1, and stores the machine code.
func (b *buf) emitBinaryExpr(n *ast.BinaryExpr) {
	b.emitExpr(n.X)
	b.emitExpr(n.Y)
	b.emit(pop_rax, mov_rax_xmm3) // get right operand
	b.emit(pop_rax, mov_rax_xmm2) // get left operand

	switch n.Op {
	default:
		panic(err(n.Pos(), "syntax error:", n.Op))
	case token.ADD:
		b.emit(add_xmm3_xmm2)
	case token.SUB:
		b.emit(sub_xmm3_xmm2)
	case token.MUL:
		b.emit(mul_xmm3_xmm2)
	case token.QUO:
		b.emit(div_xmm3_xmm2)
	}

	b.emit(mov_xmm2_rax, push_rax)
}

func typ(x interface{}) reflect.Type {
	return reflect.TypeOf(x)
}
func err(pos token.Pos, msg ...interface{}) error {
	return fmt.Errorf("%v: %v", pos, fmt.Sprintln(msg...))
}
