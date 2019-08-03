package main

import (
	"bytes"
	"strconv"

	"github.com/savsgio/gotils"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/tsingson/fasthttp-example/logger"
)

func main() {

	var log *zap.Logger
	var address = "127.0.0.1:3001"
	// var uriPrefix = []byte("/uri")
	// case bytes.HasPrefix(path, imgPrefix):

	log = logger.Console()
	// handler
	requestHandler := func(ctx *fasthttp.RequestCtx) {

		// -------------------------------------------------------
		// 处理 web client 的请求数据
		// -------------------------------------------------------
		// 取出 web client 请求进行 TCP 连接的连接 ID
		var connID = strconv.FormatUint(ctx.ConnID(), 10)
		// 暂时用 connectIDstring 当成 transaction  ID
		var tid = connID
		log.Debug("HTTP 访问 TCP 连接 ID  " + connID)

		// 取出 web 访问的 URL/URI
		var uriPath = ctx.Path()
		{
			// 取出 URI
			log.Debug("---------------- HTTP URI -------------")
			log.Debug(" HTTP 请求 URL 原始数据 > ", zap.String("request", ctx.String()))
		}

		// 取出 web client 请求的 URL/URI 中的参数部分
		{
			log.Debug("---------------- HTTP URI 参数 -------------")
			var uri = ctx.URI().QueryString()
			log.Debug("在 URI 中的原始数据 > " + string(uri))
			log.Debug("---------------- HTTP URI 每一个键值对 -------------")
			ctx.URI().QueryArgs().VisitAll(func(key, value []byte) {
				log.Debug(tid, zap.String("key", gotils.B2S(key)), zap.String("value", gotils.B2S(value)))
			})
		}
		// -------------------------------------------------------
		// 注意对比一下, 下面的代码段, 与 web client  中几乎一样
		// -------------------------------------------------------
		{
			// 取出 web client 请求中的 HTTP header
			{
				log.Debug("---------------- HTTP header 每一个键值对-------------")
				ctx.Request.Header.VisitAll(func(key, value []byte) {
					// l.Info("requestHeader", zap.String("key", gotils.B2S(key)), zap.String("value", gotils.B2S(value)))
					log.Debug(tid, zap.String("key", gotils.B2S(key)), zap.String("value", gotils.B2S(value)))
				})

			}
			// 取出 web client 请求中的 HTTP payload
			{
				log.Debug("---------------- HTTP payload -------------")
				log.Debug(tid, zap.String("http payload", gotils.B2S(ctx.Request.Body())))
			}
		}
		switch {
		// 如果访问的 URI 路由是 /uri 开头 , 则进行下面这个响应
		case len(uriPath) > 1:
			{
				log.Debug("---------------- HTTP 响应 -------------")

				// -------------------------------------------------------
				// 处理逻辑开始
				// -------------------------------------------------------

				// payload 是 []byte , 是 web response 返回的 HTTP payload
				var payload = bytes.NewBuffer([]byte("Hello, "))

				// 这是从 web client 取数据
				var who = ctx.QueryArgs().PeekBytes([]byte("who"))

				if len(who) > 0 {
					payload.Write(who)
				} else {
					payload.Write([]byte(" 中国 "))
				}

				// -------------------------------------------------------
				// 处理 HTTP 响应数据
				// -------------------------------------------------------
				// HTTP header 构造
				ctx.Response.Header.SetStatusCode(200)
				ctx.Response.Header.SetConnectionClose() // 关闭本次连接, 这就是短连接 HTTP
				ctx.Response.Header.SetBytesKV([]byte("Content-Type"), []byte("text/plain; charset=utf8"))
				// HTTP payload 设置
				// 这里 HTTP payload 是 []byte
				ctx.Response.SetBody(payload.Bytes())
			}

			// 访问路踊不是 /uri 的其他响应
		default:
			{
				log.Debug("---------------- HTTP 响应 -------------")

				var payload = bytes.NewBuffer([]byte("Hello, "))

				var who = ctx.QueryArgs().PeekBytes([]byte("who"))

				if len(who) > 0 {
					payload.Write(who)
				} else {
					payload.Write([]byte(" 中国 "))
				}

				//
				ctx.Response.Header.SetStatusCode(200)
				ctx.Response.Header.SetConnectionClose() // 关闭本次连接, 这就是短连接 HTTP
				ctx.Response.Header.SetBytesKV([]byte("Content-Type"), []byte("text/plain; charset=utf8"))
				//
				ctx.Response.SetBody(payload.Bytes())
			}
		}

		return

	}

	// Create custom server.
	s := &fasthttp.Server{
		Handler: requestHandler,
		// Every response will contain 'Server: My super server' header.
		Name: "hello-world server",
		// Other Server settings may be set here.
	}

	log.Debug("------------------ fasthttp 服务器尝试启动------ ")

	if err := s.ListenAndServe(address); err != nil {
		log.Fatal("error in ListenAndServe", zap.Error(err))
	}
}
