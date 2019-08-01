package webserver

import (
	"github.com/valyala/fasthttp"
)

func recovery(next func(ctx *fasthttp.RequestCtx)) func(ctx *fasthttp.RequestCtx) {
	fn := func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if rvr := recover(); rvr != nil {
				/**
				  logEntry := GetLogEntry(r)
				  if logEntry != nil {
				  	logEntry.Panic(rvr, debug.Stack())
				  } else {
				  	fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
				  	debug.PrintStack()
				  }
				*/

				ctx.Error("recover", 500)
			}
		}()
		next(ctx)
	}
	return fn
}
