package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/tsingson/logger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"go.uber.org/zap"
)

func main() {
	log := logger.New(logger.WithDebug(), logger.WithStoreInDay())
	address := "127.0.0.1:3001"

	// hijackHandler is called on hijacked connection.
	hijackHandler := func(c net.Conn) {
		fmt.Fprintf(c, "This message is sent over a hijacked connection to the client %s\n", c.RemoteAddr())
		fmt.Fprintf(c, "Send me something and I'll echo it to you\n")
		var buf [1]byte
		for {
			if _, err := c.Read(buf[:]); err != nil {
				log.Printf("error when reading from hijacked connection: %s", err)
				return
			}
			fmt.Fprintf(c, "You sent me %q. Waiting for new data\n", buf[:])
			log.Info("hijack", zap.ByteString("read", buf[:]))
		}
	}

	// requestHandler is called for each incoming request.
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		path := ctx.Path()
		switch {
		case string(path) == "/hijack":
			// Note that the connection is hijacked only after
			// returning from requestHandler and sending http response.
			ctx.Hijack(hijackHandler)

			// The connection will be hijacked after sending this response.
			fmt.Fprintf(ctx, "Hijacked the connection!")
		case string(path) == "/":
			fmt.Fprintf(ctx, "Root directory requested")
		default:
			fmt.Fprintf(ctx, "Requested path is %q", path)
		}
	}

	s := &fasthttp.Server{
		Handler: requestHandler,
		Logger:  log,
	}

	// reuse port
	ln, err := reuseport.Listen("tcp4", address) //  s.Cfg.ProxyConfig.ServerPort) // ":8095") //
	// s.ln, err = net.Listen("tcp4", stbmodel.AAAPort) //  s.Cfg.ProxyConfig.ServerPort) // ":8095") //
	if err != nil {
		log.Error("connect error",
			zap.String("addr", address),
			zap.Error(err))
		time.Sleep(1 * time.Second)
		os.Exit(-1)
	}

	err = s.Serve(ln)

	if err != nil {
		log.Error("fasthttp web httpserver start error",
			zap.String("addr", address),
			zap.Error(err))
		time.Sleep(1 * time.Second)
		os.Exit(-1)
	}
}
