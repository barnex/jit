package jit

import "fmt"

var (
	ass    []ssaexpr
	assOf  = make(map[expr]int)
	exprOf []expr
)

type ssaentry struct {
	e ssaexpr
}

//func (s*ssaentry)

type ssaexpr interface{}
type ssavar string
type ssaconst float64
type ssabin struct {
	op   string
	x, y int
}
type ssacall struct {
	fun string
	arg int
}

func (s ssabin) String() string  { return fmt.Sprintf("x%d %v x%d", s.x, s.op, s.y) }
func (s ssacall) String() string { return fmt.Sprintf("%v(x%d)", s.fun, s.arg) }

func SSADump(e expr) {
	emit(variable{"x"}, ssavar("x"))
	emit(variable{"y"}, ssavar("y"))
	ssaDump(e)

	for i, e := range ass {
		fmt.Printf("x%v = %v	// %v\n", i, e, exprOf[i])
	}
}

func emit(e expr, s ssaexpr) int {
	if i, ok := assOf[e]; ok {
		return i
	}

	ass = append(ass, s)
	if p, ok := assOf[e]; ok {
		panic(fmt.Sprint("duplicate assignment of ", e, ", previously:", p))
	}
	exprOf = append(exprOf, e)
	assOf[e] = len(ass) - 1
	return assOf[e]
}

func ssaDump(e expr) int {
	switch e := e.(type) {
	default:
		panic(fmt.Sprintf("%v: %T", e, e))
	case variable:
		return emit(e, ssavar(e.name))
	case constant:
		return emit(e, ssaconst(e.value))
	case binexpr:
		lhs := ssaDump(e.x)
		rhs := ssaDump(e.y)
		return emit(e, ssabin{op: e.op, x: lhs, y: rhs})
	case callexpr:
		arg := ssaDump(e.arg)
		return emit(e, ssacall{fun: e.fun, arg: arg})
	}
}
