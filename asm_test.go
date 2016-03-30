package jit

import (
	"fmt"
	"testing"
)

func TestMovXmm(t *testing.T) {
	tests := []struct {
		r1, r2 int
		want   []byte
	}{
		// reference values obtained with gcc and objdump.
		{0, 0, []byte{0xf3, 0x0f, 0x7e, 0xc0}},
		{0, 1, []byte{0xf3, 0x0f, 0x7e, 0xc8}},
		{0, 2, []byte{0xf3, 0x0f, 0x7e, 0xd0}},
		{0, 3, []byte{0xf3, 0x0f, 0x7e, 0xd8}},
		{0, 4, []byte{0xf3, 0x0f, 0x7e, 0xe0}},
		{0, 5, []byte{0xf3, 0x0f, 0x7e, 0xe8}},
		{0, 6, []byte{0xf3, 0x0f, 0x7e, 0xf0}},
		{0, 7, []byte{0xf3, 0x0f, 0x7e, 0xf8}},
		{1, 0, []byte{0xf3, 0x0f, 0x7e, 0xc1}},
		{1, 1, []byte{0xf3, 0x0f, 0x7e, 0xc9}},
		{1, 2, []byte{0xf3, 0x0f, 0x7e, 0xd1}},
		{1, 3, []byte{0xf3, 0x0f, 0x7e, 0xd9}},
		{1, 4, []byte{0xf3, 0x0f, 0x7e, 0xe1}},
		{1, 5, []byte{0xf3, 0x0f, 0x7e, 0xe9}},
		{1, 6, []byte{0xf3, 0x0f, 0x7e, 0xf1}},
		{1, 7, []byte{0xf3, 0x0f, 0x7e, 0xf9}},
		{2, 0, []byte{0xf3, 0x0f, 0x7e, 0xc2}},
		{2, 1, []byte{0xf3, 0x0f, 0x7e, 0xca}},
		{2, 2, []byte{0xf3, 0x0f, 0x7e, 0xd2}},
		{2, 3, []byte{0xf3, 0x0f, 0x7e, 0xda}},
		{2, 4, []byte{0xf3, 0x0f, 0x7e, 0xe2}},
		{2, 5, []byte{0xf3, 0x0f, 0x7e, 0xea}},
		{2, 6, []byte{0xf3, 0x0f, 0x7e, 0xf2}},
		{2, 7, []byte{0xf3, 0x0f, 0x7e, 0xfa}},
		{3, 0, []byte{0xf3, 0x0f, 0x7e, 0xc3}},
		{3, 1, []byte{0xf3, 0x0f, 0x7e, 0xcb}},
		{3, 2, []byte{0xf3, 0x0f, 0x7e, 0xd3}},
		{3, 3, []byte{0xf3, 0x0f, 0x7e, 0xdb}},
		{3, 4, []byte{0xf3, 0x0f, 0x7e, 0xe3}},
		{3, 5, []byte{0xf3, 0x0f, 0x7e, 0xeb}},
		{3, 6, []byte{0xf3, 0x0f, 0x7e, 0xf3}},
		{3, 7, []byte{0xf3, 0x0f, 0x7e, 0xfb}},
		{4, 0, []byte{0xf3, 0x0f, 0x7e, 0xc4}},
		{4, 1, []byte{0xf3, 0x0f, 0x7e, 0xcc}},
		{4, 2, []byte{0xf3, 0x0f, 0x7e, 0xd4}},
		{4, 3, []byte{0xf3, 0x0f, 0x7e, 0xdc}},
		{4, 4, []byte{0xf3, 0x0f, 0x7e, 0xe4}},
		{4, 5, []byte{0xf3, 0x0f, 0x7e, 0xec}},
		{4, 6, []byte{0xf3, 0x0f, 0x7e, 0xf4}},
		{4, 7, []byte{0xf3, 0x0f, 0x7e, 0xfc}},
		{5, 0, []byte{0xf3, 0x0f, 0x7e, 0xc5}},
		{5, 1, []byte{0xf3, 0x0f, 0x7e, 0xcd}},
		{5, 2, []byte{0xf3, 0x0f, 0x7e, 0xd5}},
		{5, 3, []byte{0xf3, 0x0f, 0x7e, 0xdd}},
		{5, 4, []byte{0xf3, 0x0f, 0x7e, 0xe5}},
		{5, 5, []byte{0xf3, 0x0f, 0x7e, 0xed}},
		{5, 6, []byte{0xf3, 0x0f, 0x7e, 0xf5}},
		{5, 7, []byte{0xf3, 0x0f, 0x7e, 0xfd}},
		{6, 0, []byte{0xf3, 0x0f, 0x7e, 0xc6}},
		{6, 1, []byte{0xf3, 0x0f, 0x7e, 0xce}},
		{6, 2, []byte{0xf3, 0x0f, 0x7e, 0xd6}},
		{6, 3, []byte{0xf3, 0x0f, 0x7e, 0xde}},
		{6, 4, []byte{0xf3, 0x0f, 0x7e, 0xe6}},
		{6, 5, []byte{0xf3, 0x0f, 0x7e, 0xee}},
		{6, 6, []byte{0xf3, 0x0f, 0x7e, 0xf6}},
		{6, 7, []byte{0xf3, 0x0f, 0x7e, 0xfe}},
		{7, 0, []byte{0xf3, 0x0f, 0x7e, 0xc7}},
		{7, 1, []byte{0xf3, 0x0f, 0x7e, 0xcf}},
		{7, 2, []byte{0xf3, 0x0f, 0x7e, 0xd7}},
		{7, 3, []byte{0xf3, 0x0f, 0x7e, 0xdf}},
		{7, 4, []byte{0xf3, 0x0f, 0x7e, 0xe7}},
		{7, 5, []byte{0xf3, 0x0f, 0x7e, 0xef}},
		{7, 6, []byte{0xf3, 0x0f, 0x7e, 0xf7}},
		{7, 7, []byte{0xf3, 0x0f, 0x7e, 0xff}},
	}
	for _, test := range tests {
		have := fmt.Sprintf("%x", mov_xmm(test.r1, test.r2))
		want := fmt.Sprintf("%x", test.want)
		if have != want {
			t.Errorf("movq %%xmm%v,%%xmm%v: have %v, want %v", test.r1, test.r2, have, want)
		}
	}
}
