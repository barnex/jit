package jit

import (
	"fmt"
)

type expr interface {
	compile(b *buf)
	children() []expr
	String() string
}

func walk(root expr, f func(expr)) {
	for _, c := range root.children() {
		walk(c, f)
	}
	f(root)
}

func recordCalls(root expr, m map[expr]bool) {
	walk(root, func(e expr) {
		for _, c := range root.children() {
			recordCalls(c, m)
			if m[c] {
				m[root] = true
			}
		}
		if _, ok := e.(*callexpr); ok {
			m[root] = true
		}
	})
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

type binexpr struct{ x, y expr }

func (e *binexpr) children() []expr {
	return []expr{e.x, e.y}
}

func (e*add)String()string{ return fmt.Sprintf("(%v+%v)", e.x, e.y) }
func (e*sub)String()string{ return fmt.Sprintf("(%v-%v)", e.x, e.y) }
func (e*mul)String()string{ return fmt.Sprintf("(%v*%v)", e.x, e.y) }
func (e*quo)String()string{ return fmt.Sprintf("(%v/%v)", e.x, e.y) }

type add struct{ binexpr }
type sub struct{ binexpr }
type mul struct{ binexpr }
type quo struct{ binexpr }

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
