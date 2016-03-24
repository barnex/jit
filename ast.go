package jit

type expr interface {
	compile(b *buf)
	children() []expr
}

func walk(root expr, f func(expr)) {
	for _, c := range root.children() {
		walk(c, f)
	}
	f(root)
}

func recordCalls(root expr, m map[expr]bool){
	walk(root, func(e expr){
		for _, c := range root.children(){
			recordCalls(c, m)
			if m[c]{
				m[root]	= true
			}
		}
		if _, ok := e.(*callexpr); ok{
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

type constant struct {
	leaf
	value float64
}

type binexpr struct{ x, y expr }

func (e *binexpr) children() []expr {
	return []expr{e.x, e.y}
}

type add struct{ binexpr }
type sub struct{ binexpr }
type mul struct{ binexpr }
type quo struct{ binexpr }

type callexpr struct {
	fun  string
	args []expr
}

func (e *callexpr) children() []expr {
	var c []expr
	for _, a := range e.args {
		c = append(c, a)
	}
	return c
}
