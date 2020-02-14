package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/integrii/flaggy"
	"github.com/tsingson/logger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/tsingson/fasthttp-guide/pkg/utils"
)

var log *logger.Logger

func main() {
	log = logger.New(
		logger.WithDebug(),
		// logger.WithStoreInDay(),
		// logger.WithDays(31),
		logger.WithLevel(zapcore.DebugLevel))
	defer log.Sync()

	log.Info("try to start fasthttp web server")

	addr := ":3001"
	compress := true

	flaggy.String(&addr, "addr", "address", "TCP address to listen to")
	flaggy.Bool(&compress, "c", "compress", "Whether to enable transparent response compression")

	flag.Parse()

	h := requestHandler
	if compress {
		h = fasthttp.CompressHandler(h)
	}

	s := &fasthttp.Server{
		Handler: h,
		Logger:  log,
	}

	// reuse port
	ln, err := reuseport.Listen("tcp4", addr) //  s.Cfg.ProxyConfig.ServerPort) // ":8095") //
	// s.ln, err = net.Listen("tcp4", stbmodel.AAAPort) //  s.Cfg.ProxyConfig.ServerPort) // ":8095") //
	if err != nil {
		log.Error("connect error",
			zap.String("addr", addr),
			zap.Error(err))
		time.Sleep(1 * time.Second)
		os.Exit(-1)
	}

	err = s.Serve(ln)

	if err != nil {
		log.Error("fasthttp web server start error",
			zap.String("addr", addr),
			zap.Error(err))
		time.Sleep(1 * time.Second)
		os.Exit(-1)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	utils.RequestCtxDebug(ctx, log.Log, true)
	fmt.Fprintf(ctx, "Hello, world!\n\n")

	fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
	fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
	fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
	fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
	fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
	fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
	fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
	fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
	fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())

	fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)

	ctx.SetContentType("text/plain; charset=utf8")

	// Set arbitrary headers
	ctx.Response.Header.Set("X-My-Header", "my-header-value")

	// Set cookies
	var c fasthttp.Cookie
	c.SetKey("cookie-name")
	c.SetValue("cookie-value")
	ctx.Response.Header.SetCookie(&c)

	fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", ctx.Request.Body())
	// ctx.Response.SetBody(payload)
}
