// +build ignore

package main

import (
	"fmt"

	. "github.com/barnex/jit"
)

func main() {

	code := []byte{
		0xf2, 0x0f, 0x10, 0x44, 0x24, 0x08, // movsd  0x8(%rsp),%xmm0
		0xf2, 0x0f, 0x10, 0x4c, 0x24, 0x10, // movsd  0x10(%rsp),%xmm1
		0xf2, 0x0f, 0x58, 0xc1, // addsd  %xmm1,%xmm0
		0xf2, 0x0f, 0x11, 0x44, 0x24, 0x18, // movsd  %xmm0,0x18(%rsp)
		0xc3, // retq
	}

	code, err := MakeExecutable(code)
	if err != nil {
		panic(err)
	}

	f := MakeFunc(code)
	fmt.Println(f(40, 2))
}

func xxadd(x, y float64) float64 {
	return x + y
}
