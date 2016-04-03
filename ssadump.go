// +build ignore

package main

import (
	"log"
	"os"

	"github.com/barnex/just-in-time-compiler"
)

func main() {
	expr := os.Args[1]
	_, err := jit.SSADump(expr)
	if err != nil {
		log.Fatal(err)
	}
}
