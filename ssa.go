package jit

import "fmt"

var(
	ass []ssaexpr
)


type ssaexpr interface{}
type ssavar string
type ssaconst float64
type ssabin struct{ op string; x, y int }
type ssacall struct{ fun string; arg int }

func(s ssabin)String()string { return fmt.Sprintf("x%d %v x%d", s.x, s.op, s.y) }
func(s ssacall)String()string { return fmt.Sprintf("%v(x%d)", s.fun, s.arg) }

func SSADump(e expr) {
	ass = []ssaexpr{
		ssavar("x"),
		ssavar("y"),
	}
	ssaDump(e)

	for i,e := range ass{
		fmt.Printf("x%v = %v\n", i, e)
	}
}

func ssaDump(e expr){
	switch e := e.(type) {
	default:
		panic(fmt.Sprintf("%v: %T", e, e))
case *variable:
		ass = append(ass, ssavar(e.name))
		case *constant:
		ass = append(ass, ssaconst(e.value))
	case *binexpr:
		ssaDump(e.x)
		lhs := currAss()
		ssaDump(e.y)
		rhs := currAss()
		ass = append(ass, ssabin{op: e.op, x:lhs, y:rhs})
	case *callexpr:
		ssaDump(e.args[0])
		arg := currAss()
		ass = append(ass, ssacall{fun:e.fun, arg:arg})
	}
}

func currAss() int{return len(ass)-1}
