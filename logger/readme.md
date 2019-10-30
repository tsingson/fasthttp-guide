# fasthttp logger interface

```
type Logger interface {
    // Printf must have the same semantics as log.Printf.
    Printf(format string, args ...interface{})
}
```
