# Cumin

Cumin is a function currying library with dynamic invocation implemented through runtime function introspection and parameter type coercion. It basically provides support for generic functions in Go: arbitrary functions with arbitrary invocations. 

Cumin allows for go code to accept any function as a callback handler. This is somewhat antithetical to golang's principles, and there's a wealth of debate around generics in go, but this is useful in situations like [exposing go functions to python](https://github.com/damouse/glusnek).

It works by reflecting the type and number of parameters in passed functions, testing those types against incoming parameters, and attempting to cast those types appropriately. 

## Examples

```go
fn := func() int {
  return 1
}

c, _ := NewCurry(fn)

// Returns 1
r, err := c.Invoke([]interface{}{})

// Returns an invocation error: too many arguments
r, err := c.Invoke([]interface{}{1, 2})
```
