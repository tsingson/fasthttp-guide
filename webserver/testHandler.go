package webserver

import (
	"strconv"

	"github.com/savsgio/gotils"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func (ws *webServer) testGet() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		var tid = strconv.FormatInt(int64(ctx.ID()), 10)
		l := ws.Log.Named(tid)
		l.Debug("testGet")

		if ws.debug {
			l.Debug(tid, zap.String("request", ctx.String()))
			ctx.Request.Header.VisitAll(func(key, value []byte) {
				// l.Info("requestHeader", zap.String("key", gotils.B2S(key)), zap.String("value", gotils.B2S(value)))
				l.Debug(tid, zap.String("key", gotils.B2S(key)), zap.String("value", gotils.B2S(value)))
			})

			l.Debug(tid, zap.String("http payload", gotils.B2S(ctx.Request.Body())))

		}

		ctx.SetContentType(ContentRest)
		ctx.SetStatusCode(200)
		ctx.SetBody([]byte(`{"id":2101127497763529765,"plan_start_date":"2019-07-29","plan_end_date":"2019-02-12","actual_start_date":"2019-07-29","actual_end_date":"2019-07-29","is_done":false,"last_updated":"2019-08-01T14:12:17.983236","is_deleted":false,"user_id":2098735545843717147,"title":"00002"}`))
		return
	}

}

func (ws *webServer) testPost() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		var tid = strconv.FormatInt(int64(ctx.ID()), 10)
		l := ws.Log.Named(tid)
		l.Debug("testPost")

		if ws.debug {
			l.Debug(tid, zap.String("request", ctx.String()))

			ctx.Request.Header.VisitAll(func(key, value []byte) {
				// l.Info("requestHeader", zap.String("key", gotils.B2S(key)), zap.String("value", gotils.B2S(value)))
				l.Debug(tid, zap.String("key", gotils.B2S(key)), zap.String("value", gotils.B2S(value)))
			})

			l.Debug(tid, zap.String("http payload", gotils.B2S(ctx.Request.Body())))

		}

		ctx.SetContentType(ContentRest)
		ctx.SetStatusCode(200)
		ctx.SetBody([]byte(`{"id":2101127497763529765,"plan_start_date":"2019-07-29","plan_end_date":"2019-02-12","actual_start_date":"2019-07-29","actual_end_date":"2019-07-29","is_done":false,"last_updated":"2019-08-01T14:12:17.983236","is_deleted":false,"user_id":2098735545843717147,"title":"00002"}`))
		return
	}

}
