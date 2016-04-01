#jit

Toy just-in-time compiler for arithmetic expressions of floating-point variables x and y, like `sqrt(x*x + y*y) - 2*cos(x+1)`.  Intended for fun and learning only.

The generated machine instructions are for x68-64 and are tested on Linux only.

## Parsing & AST

Parsing transforms an expression, like `(x+y)*y` into an AST (Abstract Syntax Tree), like:

```
      *
     /  \ 
    +    y
   / \
  x   1
```

For simplicity, we use Go's built-in parser (package `go/parser`). However, we transform the Go AST into our own representation, to stay independent of `go/ast`'s internal details.



## Code generation

`(x+1)*(2+sqrt(y+1))` compiles to:

```
55                              push   %rbp
48 89 e5                        mov    %rsp,%rbp
48 81 ec 10 00 00 00            sub    $0x10,%rsp
66 48 0f 7e c0                  movq   %xmm0,%rax
48 89 85 f8 ff ff ff            mov    %rax,-0x8(%rbp)
66 48 0f 7e c8                  movq   %xmm1,%rax
48 89 85 f0 ff ff ff            mov    %rax,-0x10(%rbp)
48 b8 00 00 00 00 00 00 f0 3f   movabs $0x3ff0000000000000,%rax
66 48 0f 6e c0                  movq   %rax,%xmm0
f3 0f 7e d0                     movq   %xmm0,%xmm2
48 8b 85 f0 ff ff ff            mov    -0x10(%rbp),%rax
66 48 0f 6e c0                  movq   %rax,%xmm0
f3 0f 7e ca                     movq   %xmm2,%xmm1
f2 0f 58 c1                     addsd  %xmm1,%xmm0
48 b8 30 1c 40 00 00 00 00 00   movabs $0x401c30,%rax
ff d0                           callq  *%rax
f3 0f 7e d0                     movq   %xmm0,%xmm2
48 b8 00 00 00 00 00 00 00 40   movabs $0x4000000000000000,%rax
66 48 0f 6e c0                  movq   %rax,%xmm0
f3 0f 7e ca                     movq   %xmm2,%xmm1
f2 0f 58 c1                     addsd  %xmm1,%xmm0
f3 0f 7e d0                     movq   %xmm0,%xmm2
48 b8 00 00 00 00 00 00 f0 3f   movabs $0x3ff0000000000000,%rax
66 48 0f 6e c0                  movq   %rax,%xmm0
f3 0f 7e d8                     movq   %xmm0,%xmm3
48 8b 85 f8 ff ff ff            mov    -0x8(%rbp),%rax
66 48 0f 6e c0                  movq   %rax,%xmm0
f3 0f 7e cb                     movq   %xmm3,%xmm1
f2 0f 58 c1                     addsd  %xmm1,%xmm0
f3 0f 7e ca                     movq   %xmm2,%xmm1
f2 0f 59 c1                     mulsd  %xmm1,%xmm0
48 81 c4 10 00 00 00            add    $0x10,%rsp
5d                              pop    %rbp
c3                              retq   
```

## Assembler

## Dynamic loading

## Use case: implicit function plotter

An example use case is an implicit function plotter. Here we plot the curve implicitly defined by `(x*x-y*y-y*x-4)*(x*x+y*y-16)`:

![fig](plotter.png)
