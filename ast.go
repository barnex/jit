package jit

// this file provides the AST (abstract syntax tree) building blocks
// and analysis functions.

import (
	"fmt"
)

type expr interface {
	compile(b *buf)
	children() []expr
	String() string
}

type leaf struct{}

func (_ leaf) children() []expr { return nil }

type variable struct {
	leaf
	name string
}

func (e *variable) String() string {
	return e.name
}

type constant struct {
	leaf
	value float64
}

func (e *constant) String() string {
	return fmt.Sprint(e.value)
}

type binexpr struct {
	op   string
	x, y expr
}

func (e *binexpr) children() []expr {
	return []expr{e.x, e.y}
}

func (e *binexpr) String() string {
	return fmt.Sprintf("(%v%v%v)", e.x, e.op, e.y)
}

type callexpr struct {
	fun  string
	args []expr
}

func (e *callexpr) String() string {
	return fmt.Sprintf("%v(%v)", e.fun, e.args[0])
}

func (e *callexpr) children() []expr {
	var c []expr
	for _, a := range e.args {
		c = append(c, a)
	}
	return c
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
	if _, ok := root.(*callexpr); ok {
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
	if _, ok := root.(*binexpr); ok {
		m[root]++
	}
}
