//+build ignore

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/barnex/jit"
)

func main() {
	expr := strings.Join(os.Args[1:], " ")
	code, err := jit.Compile(expr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(code.Eval(9, 999))
}
