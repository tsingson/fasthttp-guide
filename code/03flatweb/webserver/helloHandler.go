package webserver

import (
	"strconv"

	"github.com/valyala/fasthttp"

	"github.com/tsingson/fasthttp-guide/pkg/utils"
)

func (ws *webServer) helloWorldGetHandler() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		tid := strconv.FormatInt(int64(ctx.ID()), 10)
		log := ws.Log.Named(tid)

		if ws.debug {
			utils.RequestCtxDebug(ctx, log, false)
		}

		ctx.SetContentType(ContentText)
		ctx.SetStatusCode(200)
		ctx.SetBody([]byte(`hello world`))
	}
}

func Hello() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		tid := strconv.FormatInt(int64(ctx.ID()), 10)

		ctx.Request.Header.Add("tid", tid)
		ctx.SetContentType(ContentText)
		ctx.SetStatusCode(200)
		ctx.SetBody([]byte(`hello world`))
	}
}
