#jit

Toy just-in-time compiler for arithmetic expressions of floating-point variables x and y, like `sqrt(x*x + y*y) - 2*cos(x+1)`.

The generated machine instructions are for x68-64 and are tested on linux only.

## Code generation

`(x+1)*(2+sqrt(y+1))` compiles to:

```lang:asm
push   %rbp
mov    %rsp,%rbp
sub    $0x10,%rsp
movq   %xmm0,%rax
mov    %rax,-0x8(%rbp)
movq   %xmm1,%rax
mov    %rax,-0x10(%rbp)
movabs $0x3ff0000000000000,%rax
movq   %rax,%xmm0
movq   %xmm0,%xmm2
mov    -0x10(%rbp),%rax
movq   %rax,%xmm0
movq   %xmm2,%xmm1
addsd  %xmm1,%xmm0
movabs $0x401c30,%rax
callq  *%rax
movq   %xmm0,%xmm2
movabs $0x4000000000000000,%rax
movq   %rax,%xmm0
movq   %xmm2,%xmm1
addsd  %xmm1,%xmm0
movq   %xmm0,%xmm2
movabs $0x3ff0000000000000,%rax
movq   %rax,%xmm0
movq   %xmm0,%xmm3
mov    -0x8(%rbp),%rax
movq   %rax,%xmm0
movq   %xmm3,%xmm1
addsd  %xmm1,%xmm0
movq   %xmm2,%xmm1
mulsd  %xmm1,%xmm0
add    $0x10,%rsp
pop    %rbp
retq   
```


Inspired by the book "The Go Programming Language" by Alan A. A. Donovan and Brian W. Kernighan, section 7.9: Example: Expression Evaluator.
