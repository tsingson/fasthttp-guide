package main

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/integrii/flaggy"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/tsingson/logger"

	"github.com/tsingson/fasthttp-guide/pkg/utils"
)

/**
121.211.149.176 - - [04/Aug/2019:23:59:01 +0800] "GET /epg?date=eq.2019-08-05&channel_id=eq.20024 HTTP/1.1" 200 2249 "-" "okhttp/3.5.0" "-" "0.060" "0.061" "-" request_body&&&&&-
*/

var log *logger.Logger

func main() {
	log = logger.New(logger.WithDebug(), logger.WithStoreInDay())
	address := "127.0.0.1:3001"

	// Add a flag
	flaggy.String(&address, "addr", "address", "address for webserver")

	// Parse the flag
	flaggy.Parse()

	// -------------------------------------------------------
	// 创建 fasthttp 服务器
	// -------------------------------------------------------
	// Create custom server.
	s := &fasthttp.Server{
		Handler: requestHandler,          // 注意这里
		Name:    "hello-01cli-v1 server", // 服务器名称
	}
	// -------------------------------------------------------
	// 运行服务端程序
	// -------------------------------------------------------
	log.Debug("------------------ fasthttp 服务器尝试启动------ ")

	if err := s.ListenAndServe(address); err != nil {
		log.Fatal("error in ListenAndServe", zap.Error(err))
	}
}

// requestHandler handler for fasthttp
func requestHandler(ctx *fasthttp.RequestCtx) {
	// -------------------------------------------------------
	// 处理 web client 的请求数据
	// -------------------------------------------------------
	// 取出 web client 请求进行 TCP 连接的连接 ID
	connID := strconv.FormatUint(ctx.ConnID(), 10)
	// 取出 web client 请求 HTTP header 中的事务ID
	tid := string(ctx.Request.Header.PeekBytes([]byte("TransactionID")))
	if len(tid) == 0 {
		tid = "12345678"
	}

	log.Debug("HTTP 访问 TCP 连接 ID  " + connID)

	// 取出 web 访问的 URL/URI
	uriPath := ctx.Path()
	{
		// 取出 URI
		log.Debug("---------------- HTTP URI -------------")
		log.Debug(" HTTP 请求 URL 原始数据 > ", zap.String("request", fmt.Sprintf("#%016X - %s<->%s - %s %s", ctx.ID(), ctx.LocalAddr(), ctx.RemoteAddr(), ctx.Request.Header.Method(), ctx.URI().FullURI())))
	}

	// 取出 web client 请求的 URL/URI 中的参数部分
	{
		log.Debug("---------------- HTTP URI 参数 -------------")
		uri := ctx.URI().QueryString()
		log.Debug("在 URI 中的原始数据 > " + string(uri))
		log.Debug("---------------- HTTP URI 每一个键值对 -------------")
		ctx.URI().QueryArgs().VisitAll(func(key, value []byte) {
			log.Debug(tid, zap.String("key", utils.B2S(key)), zap.String("value", utils.B2S(value)))
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
				// l.Info("requestHeader", zap.String("key", utils.B2S(key)), zap.String("value", utils.B2S(value)))
				log.Debug(tid, zap.String("key", utils.B2S(key)), zap.String("value", utils.B2S(value)))
			})
		}
		// 取出 web client 请求中的 HTTP payload
		{
			log.Debug("---------------- request HTTP payload -------------")
			log.Debug(tid, zap.String("http payload", utils.B2S(ctx.Request.Body())))
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
			payload := bytes.NewBuffer([]byte("Hello, "))

			// 这是从 web client 取数据
			who := ctx.QueryArgs().PeekBytes([]byte("who"))

			if len(who) > 0 {
				payload.Write(who)
			} else {
				payload.Write([]byte(" 中国 "))
			}
			where := ctx.QueryArgs().PeekBytes([]byte("where"))
			if len(where) > 0 {
				payload.Write(where)
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
			ctx.Response.Header.SetBytesKV([]byte("TransactionID"), []byte(tid))
			// HTTP payload 设置
			// 这里 HTTP payload 是 []byte
			log.Debug(tid, zap.String("payload", payload.String()))
			ctx.Response.SetBody(payload.Bytes())
		}

		// 访问路踊不是 /uri 的其他响应
	default:
		{
			log.Debug("---------------- HTTP 响应 -------------")

			// -------------------------------------------------------
			// 处理逻辑开始
			// -------------------------------------------------------

			// payload 是 []byte , 是 web response 返回的 HTTP payload
			payload := bytes.NewBuffer([]byte("Hello, "))

			// 这是从 web client 取数据
			who := ctx.QueryArgs().PeekBytes([]byte("who"))

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
			ctx.Response.Header.SetBytesKV([]byte("TransactionID"), []byte(tid))
			// HTTP payload 设置
			// 这里 HTTP payload 是 []byte
			log.Debug(tid, zap.String("payload", payload.String()))
			ctx.Response.SetBody(payload.Bytes())
		}
	}

	return
}
