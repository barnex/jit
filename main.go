package main

import (
	"fmt"
	"os"
)

func main() {
	exprs := []string{
		"x",
		"y",
	}
	x := 1.0
	y := 2.0
	for _, expr := range exprs {
		var b Buf
		err := b.Compile(expr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		res := b.call(x, y)
		b.Free()
		fmt.Printf("with x=%v, y=%v: %v=%v\n", x, y, expr, res)
	}

}

//func main() {
//	var code codeBuffer
//
//	code.emit(openFunc, mul10, closeFunc)
//
//	mem, err := makeExecutable(code.Bytes())
//	if err != nil {
//		fatal(err)
//	}
//	defer unix.Munmap(mem)
//	result := call(mem, 3, 30)
//	fmt.Println(result)
//}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
