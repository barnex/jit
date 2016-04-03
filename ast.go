package jit

// this file provides the AST (abstract syntax tree) building blocks
// and analysis functions.

import (
	"fmt"
)

type expr interface {
	children() []expr
	simplify() expr
}

type variable struct {
	name string
}

func (variable) children() []expr {
	return nil
}

func (e variable) simplify() expr {
	return e
}
func (e variable) String() string {
	return e.name
}

type constant struct {
	value float64
}

func (constant) children() []expr {
	return nil
}

func (e constant) simplify() expr {
	return e
}

func (e constant) String() string {
	return fmt.Sprint(e.value)
}

func isConst(e expr) bool {
	_, ok := e.(*constant)
	return ok
}

type binexpr struct {
	op   string
	x, y expr
}

func (e binexpr) children() []expr {
	return []expr{e.x, e.y}
}

func (e binexpr) simplify() expr {
	x := e.x.simplify()
	y := e.y.simplify()

	if isConst(x) && isConst(y) {
		x := x.(*constant).value
		y := y.(*constant).value
		var v float64
		switch e.op {
		default:
			panic("bug")
		case "+":
			v = x + y
		case "-":
			v = x - y
		case "*":
			v = x * y
		case "/":
			v = x / y
		}
		return constant{v}
	}
	return binexpr{op: e.op, x: x, y: y}
}

func (e binexpr) String() string {
	return fmt.Sprintf("(%v%v%v)", e.x, e.op, e.y)
}

type callexpr struct {
	fun string
	arg expr
}

func (e callexpr) children() []expr {
	return []expr{e.arg}
}

func (e callexpr) simplify() expr {
	//args := make([]expr, len(e.args))
	//for i, a := range e.args {
	//	args[i] = a.simplify()
	//}
	//if len(args) == 1 && isConst(args[0]) {
	arg := e.arg.simplify()
	if isConst(arg) {
		a := arg.(*constant).value
		f := funcs[e.fun]
		v := callCFunc(f, a)
		return constant{v}
	}
	return callexpr{fun: e.fun, arg: arg}
}

func (e callexpr) String() string {
	return fmt.Sprintf("%v(%v)", e.fun, e.arg)
}

// recordCalls iterates over the AST with given root
// and records, in m, for each encountered expression whether it contains a function call.
// Used to determine whether evaluating an expression causes the register contents to be destroyed.
func recordCalls(root expr, m map[expr]bool) {
	for _, c := range root.children() {
		recordCalls(c, m)
		if m[c] {
			m[root] = true
		}
	}
	if _, ok := root.(callexpr); ok {
		m[root] = true
	}
}

// recordDepth iteratates over the AST with given root
// and records, in m, the number of binary expressions under each expression encountered.
// Used to decide which side of a binary expression requires least registers.
func recordDepth(root expr, m map[expr]int) {
	for _, c := range root.children() {
		recordDepth(c, m)
		if m[c] > m[root] {
			m[root] = m[c]
		}
	}
	if _, ok := root.(binexpr); ok {
		m[root]++
	}
}
