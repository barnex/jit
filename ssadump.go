// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/barnex/jit"
)

func main() {
	expr := os.Args[1]
	root, err := jit.Parse(expr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	jit.SSADump(root)
}
