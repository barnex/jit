package jit

import "fmt"

// FoldConst returns a new expression where all constant subexpressions have been replaced by numbers.
// E.g.:
// 	1+1 -> 2
func FoldConst(e expr) expr {
	switch e := e.(type) {
	default:
		return e
	case binexpr:
		return foldBinexpr(e)
	case callexpr:
		return foldCallexpr(e)
	}
}

func isConst(e expr) bool {
	_, ok := e.(*constant)
	return ok
}

func foldBinexpr(e binexpr) expr {
	x := FoldConst(e.x)
	y := FoldConst(e.y)

	if isConst(x) && isConst(y) {
		x := x.(*constant).value
		y := y.(*constant).value
		var v float64
		switch e.op {
		default:
			panic(fmt.Sprintf("foldBinexpr %v", e.op))
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

func foldCallexpr(e callexpr) expr {
	arg := FoldConst(e.arg)
	if isConst(arg) {
		a := arg.(*constant).value
		f := funcs[e.fun]
		v := callCFunc(f, a)
		return constant{v}
	}
	return callexpr{fun: e.fun, arg: arg}
}
