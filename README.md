Toy just-in-time compiler for arithmetic expressions of variables x and y. E.g.:

```
code, err := Compile("(x+1) * (y+2) / 3")
defer code.Free()
x, y := 1.0, 2.0
z := code.Eval(x, y) // executes machine code generated on-the-fly!
```

The generated machine instructions work on x68-64 linux only.

Inspired by the book "The Go Programming Language" by Alan A. A. Donovan and Brian W. Kernighan, section 7.9: Example: Expression Evaluator.
