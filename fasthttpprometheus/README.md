# fasthttp-prometheus

Prometheus metrics exporter for go fasthttp framework

## Installation

`$ go get github.com/flf2ko/fasthttp-prometheus`

## Usage

```go
package main

import (
    "fmt"
    "github.com/buaazp/fasthttprouter"
    "github.com/valyala/fasthttp"
    "log"
    fastp "go-fasthttp-prometheus"
)

func Index(ctx *fasthttp.RequestCtx) {
    fmt.Fprint(ctx, "Welcome!\n")
}

func main() {
    router := fasthttprouter.New()
    APIregist(router)

    p := fastp.NewPrometheus("fasthttp")
    fastpHandler := p.WrapHandler(router)

    log.Fatal(fasthttp.ListenAndServe(":8080", fastpHandler))
}

func APIregist(r *fasthttprouter.Router) {
    r.GET("/", Index)
}
```

## Related Project

* [fasthttp](https://github.com/valyala/fasthttp)
* [fasthttprouter](https://github.com/buaazp/fasthttprouter)

## Inspired by

* [go-gin-prometheus](https://github.com/zsais/go-gin-prometheus)
* [gin-prometheus](https://github.com/DanielHeckrath/gin-prometheus)