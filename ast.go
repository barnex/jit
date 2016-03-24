package jit

type expr interface {
	compile(b *buf)
	walkChildren(func(expr))
}

func walk(root expr, f func(expr)) {
	f(root)
	root.walkChildren(f)
}

type leaf struct{}

func (_ leaf) walkChildren(func(expr)) {}

type variable struct {
	leaf
	name string
}

type constant struct {
	leaf
	value float64
}

type binexpr struct{ x, y expr }

func (e *binexpr) walkChildren(f func(expr)) {
	walk(e.x, f)
	walk(e.y, f)
}

type add struct{ binexpr }
type sub struct{ binexpr }
type mul struct{ binexpr }
type quo struct{ binexpr }

type callexpr struct {
	fun  string
	args []expr
}

func (e *callexpr) walkChildren(f func(expr)) {
	for _, a := range e.args {
		walk(a, f)
	}
}
