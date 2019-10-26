package webserver

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/tsingson/fasthttp-example/pkg/goutils"
)

func (ws *webServer) simpleGetHandler() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		// var tid = strconv.FormatInt(int64(ctx.ID()), 10)
		tid := goutils.B2S(ctx.Request.Header.Peek("TransactionID"))
		l := ws.Log.Named(tid)
		l.Debug("simpleGetHandler")

		if ws.debug {
			l.Debug(tid, zap.String("request", ctx.String()))
			ctx.Request.Header.VisitAll(func(key, value []byte) {
				// l.Info("requestHeader", zap.String("key", gotils.B2S(key)), zap.String("value", gotils.B2S(value)))
				l.Debug(tid, zap.String("key", goutils.B2S(key)), zap.String("value", goutils.B2S(value)))
			})

			l.Debug(tid, zap.String("http payload", goutils.B2S(ctx.Request.Body())))

		}

		ctx.SetContentType(ContentRest)
		ctx.SetStatusCode(200)
		ctx.SetBody([]byte(`{"id":2101127497763529765,"plan_start_date":"2019-07-29","plan_end_date":"2019-02-12","actual_start_date":"2019-07-29","actual_end_date":"2019-07-29","is_done":false,"last_updated":"2019-08-01T14:12:17.983236","is_deleted":false,"user_id":2098735545843717147,"title":"00002"}`))
		return
	}
}

func (ws *webServer) simplePostHandler() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		// var tid = strconv.FormatInt(int64(ctx.ID()), 10)
		tid := goutils.B2S(ctx.Request.Header.Peek("TransactionID"))
		l := ws.Log.Named(tid)

		if ws.debug {
			l.Debug("simplePostHandler")
			l.Debug(tid, zap.String("request", ctx.String()))
			ctx.Request.Header.VisitAll(func(key, value []byte) {
				l.Debug(tid, zap.String("key", goutils.B2S(key)), zap.String("value", goutils.B2S(value)))
			})
			l.Debug(tid, zap.String("http payload", goutils.B2S(ctx.Request.Body())))
		}
		payload := []byte(`{"id":2101127497763529765,"plan_start_date":"2019-07-29","plan_end_date":"2019-02-12","actual_start_date":"2019-07-29","actual_end_date":"2019-07-29","is_done":false,"last_updated":"2019-08-01T14:12:17.983236","is_deleted":false,"user_id":2098735545843717147,"title":"00002"}`)
		ctx.SetContentType(ContentRest)
		ctx.SetStatusCode(200)
		ctx.SetBody(payload)
		return
	}
}
