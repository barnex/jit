package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
)

func (b *Buf) Compile(expr string) error {
	root, err := parser.ParseExpr(expr)
	if err != nil {
		return fmt.Errorf(`parse "%s": %v`, expr, err)
	}

	b.emit(push_rbp, mov_rsp_rbp) // function preamble
	b.emitExpr(root)              // function body (jit code)
	b.emit(pop_rax, mov_rax_xmm0) // result from stack returned via xmm0
	b.emit(pop_rbp, ret)          // return from function

	b.instr, err = makeExecutable(b.Bytes())
	return err
}

func (b *Buf) emitExpr(e ast.Expr) {
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
		//case *ast.CallExpr:
		//	return w.compileCallExpr(e)
	}
}

func (b *Buf) emitIdent(e *ast.Ident) {
	switch e.Name {
	default:
		panic(err(e.Pos(), "undefined variable:", e.Name))
	case "x":
		b.emit(mov_xmm0_rax, push_rax)
	case "y":
		b.emit(mov_xmm1_rax, push_rax)
	}
}

func (b *Buf) emitBasicLit(e *ast.BasicLit) {
	switch e.Kind {
	default:
		panic(err(e.Pos(), "syntax error:", e.Value, "(", typ(e), ")"))
	case token.FLOAT, token.INT:
		v, err := strconv.ParseFloat(e.Value, 64)
		if err != nil {
			panic(err)
		}
		b.emit(mov_imm_rax(v), push_rax)
	}
}

func (b *Buf) emitBinaryExpr(n *ast.BinaryExpr) {
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
