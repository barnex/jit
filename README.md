#jit

Toy just-in-time compiler for arithmetic expressions of floating-point variables x and y, like `sqrt(x*x + y*y) - 2*cos(x+1)`.  Intended for fun and learning only.

The generated machine instructions are for x68-64 and are tested on Linux only.

## Parsing and the AST

Parsing transforms an expression, like `(x+1)*y` into an AST (Abstract Syntax Tree), like:

```
      *
     /  \ 
    +    y
   / \
  x   1
```

For simplicity, we use Go's built-in parser (package `go/parser`). However, we transform the Go AST into our own representation, to stay independent of `go/ast`'s internal details.

Our AST's nodes are of type `expr`, an interface implemented by the concrete types:

```
variable
constant
binexpr
callexpr
```

## Constant folding

We employ constant folding on the AST, i.e. replacing constant expressions by their numerical value. E.g.:

```
((x*x)/(1+sqrt(2))) -> ((x*x)/2.414213562373095)
```


## Compilation

### calling convention

We generate code for a function body, following the System V AMD64 ABI calling conention (x, y are passed via `xmm0`, `xmm1` respectively. Result returned in `xmm0`)

### locals

Internally, we keep on using `xmm0` and `xmm1` as the arguments for function calls or binary expressions, and store results in `xmm0`. The original arguments x and y are safely stored as local variables on the stack.

### registerization

We generate code using a stack machine strategy, but registerize into `xmm2`-`xmm7` if possible -- or spill to the stack otherwise.

For binary expressions, we evaluate the "deepest" branch first. This ensures we only need `O(log(N))` registers for `N` binary expression nodes. Otherwise, unbalanced expressions like "(x+(x+(x+(x+(x+(x+(x+(x+(x+(x+y))))))))))" would require `O(N)` registers and quickly start spilling to the stack.

Finally, we must take care that function calls do not destroy the contents of `xmm2`-`xmm7`. When a branch of a binary expression contains a function call, we simply never store the result of the other branch in a register.

The result is quite good, with most expressions of typical length hitting the stack only very few times, if at all:

```
2*(x+y)*(x-y)/2 :                                         4 registers,  0 stack spills
1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1 :         1 register,   0 stack spills
1+2+(3+2*4+((((5+6*2)+7)+(8))+9)+10*sin(2-x+y/3))+11 :    4 registers,  0 stack spill
1+x+(y+2*4+((((5+y*2)+7)+sqrt(x))+y)+10*sin(2-x+y/x))+y : 4 registers,  1 stack spill
```

### putting it all together

Putting it all together, the final code looks relatively OK for a small project like this. There are a few redundant `mov`s, but these are cheap. Immediate values and function calls could have been a bit more elegant, e.g. using `rip`-relative addressing.


As an example, `(x+1)*(2+sqrt(y+1))` compiles to:

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
