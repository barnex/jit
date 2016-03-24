package jit

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

type expr interface {
	compile(b *buf)
}

type variable struct {
	name string
}

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

type constant struct{ value float64 }

func (e *constant) compile(b *buf) {
	b.emit(mov_float_rax(e.value), mov_rax_xmm0)
}

type binexpr struct{ x, y expr }

func (e *binexpr) compileArgs(b *buf) {
	e.y.compile(b) // y in xmm0

	// stash result
	reg := b.allocReg()
	if reg == -1 {
		b.emit(mov_xmm0_rax, push_rax)
	} else {
		b.emit(mov_xmm(0, reg))
	}

	e.x.compile(b) // x in xmm0
	if reg == -1 {
		b.emit(pop_rax, mov_rax_xmm1) // x in xmm1
	} else {
		b.emit(mov_xmm(reg, 1))
	}
	b.freeReg(reg)
}

type add struct{ binexpr }

func (e *add) compile(b *buf) {
	e.compileArgs(b)
	b.emit(add_xmm1_xmm0)
}

type sub struct{ binexpr }

func (e *sub) compile(b *buf) {
	e.compileArgs(b)
	b.emit(sub_xmm1_xmm0)
}

type mul struct{ binexpr }

func (e *mul) compile(b *buf) {
	e.compileArgs(b)
	b.emit(mul_xmm1_xmm0)
}

type quo struct{ binexpr }

func (e *quo) compile(b *buf) {
	e.compileArgs(b)
	b.emit(div_xmm1_xmm0)
}

type callexpr struct {
	fun  string
	args []expr
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

func Parse(expr string) (root expr, e error) {
	node, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %v", expr, err)
	}
	defer func() {
		if err := recover(); err != nil {
			root = nil
			e = fmt.Errorf("parse %q: %v", expr, err)
		}
	}()
	return parseExpr(node), nil
}

func parseExpr(node ast.Expr) expr {
	switch node := node.(type) {
	default:
		panic(fmt.Sprintf("syntax error: %T", node))
	case *ast.BasicLit:
		return parseBasicLit(node)
	case *ast.BinaryExpr:
		return parseBinaryExpr(node)
	case *ast.CallExpr:
		return parseCallExpr(node)
	case *ast.Ident:
		return parseIdent(node)
	case *ast.ParenExpr:
		return parseExpr(node.X)
	case *ast.UnaryExpr:
		return parseUnaryExpr(node)
	}
}

func parseBasicLit(node *ast.BasicLit) expr {
	switch node.Kind {
	default:
		panic(fmt.Sprintf("syntax error: %v (%T)", node.Value, node))
	case token.FLOAT, token.INT:
		v, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			panic(err)
		}
		return &constant{v}
	}
}

func parseBinaryExpr(node *ast.BinaryExpr) expr {
	x := parseExpr(node.X)
	y := parseExpr(node.Y)
	switch node.Op {
	default:
		panic(fmt.Sprintf("syntax error:", node.Op))
	case token.ADD:
		return &add{binexpr{x, y}}
	case token.SUB:
		return &sub{binexpr{x, y}}
	case token.MUL:
		return &mul{binexpr{x, y}}
	case token.QUO:
		return &quo{binexpr{x, y}}
	}
}

func parseCallExpr(node *ast.CallExpr) expr {
	args := make([]expr, 0, len(node.Args))
	for _, a := range node.Args {
		args = append(args, parseExpr(a))
	}
	fun := node.Fun.(*ast.Ident).Name
	if funcs[fun] == 0 {
		panic(fmt.Sprintf("undefined:", fun))
	}
	if len(args) != 1 {
		panic(fmt.Sprintf("%v needs 1 argument, have %v", fun, len(args)))
	}
	return &callexpr{fun, args}
}

func parseIdent(node *ast.Ident) expr {
	switch node.Name {
	default:
		panic(fmt.Sprintf("undefined: %v", node.Name))
	case "x", "y":
		return &variable{node.Name}
	}
}

func parseUnaryExpr(node *ast.UnaryExpr) expr {
	switch node.Op {
	default:
		panic(fmt.Sprintf("syntax error: %v", node.Op))
	case token.ADD:
		return parseExpr(node.X)
	case token.SUB:
		return &sub{binexpr{&constant{0}, parseExpr(node.X)}}
	}
}
