package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func main() {
	var code codeBuffer

	code.emit(openFunc, mul10, closeFunc)

	mem, err := makeExecutable(code.Bytes())
	if err != nil {
		fatal(err)
	}
	defer unix.Munmap(mem)
	result := call(mem, 3, 30)
	fmt.Println(result)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
