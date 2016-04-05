/*
Package jit provides a toy just-in-time compiler for arithmetic expressions of variables x and y. E.g.:
	code, err := Compile("(x+1) * (y+2) / 3")
	x, y := 1.0, 2.0
	z := code.Eval(x, y)

Works on 64-bit linux only.
*/
package jit
