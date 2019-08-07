package webserver

import (
	"strconv"

	"github.com/savsgio/gotils"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func (ws *webServer) hello() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		tid := strconv.FormatInt(int64(ctx.ID()), 10)
		log := ws.Log.Named(tid)
		log.Debug("hello")

		if ws.debug {
			ctx.Request.Header.VisitAll(func(key, value []byte) {
				// log.Info("requestHeader", zap.String("key", gotils.B2S(key)), zap.String("value", gotils.B2S(value)))
				log.Debug(tid, zap.String("key", gotils.B2S(key)), zap.String("value", gotils.B2S(value)))
			})

			log.Debug(tid, zap.String("http payload", gotils.B2S(ctx.Request.Body())))

		}

		ctx.SetContentType(ContentText)
		ctx.SetStatusCode(200)
		ctx.SetBody([]byte(`hello world`))
		return
	}
}
