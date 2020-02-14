package main

import (
	"os"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/tsingson/fasthttp-guide/logger"
	"github.com/tsingson/fasthttp-guide/pkg/utils"
)

func main() {
	var log *zap.Logger = logger.Console()

	// setup a full URL
	Protocol := "http"
	IPAddress := "127.0.0.1"
	Port := "3001"
	Path := "/uri/hello"
	Parameters := "id=1&who=tsingson&where=china"
	Anchor := "#anchor"

	fullURL := Protocol + "://" + IPAddress + ":" + Port + Path + "?" + Parameters + Anchor

	// fasthttp web client initial
	// get response / request from pool
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	// release back to pool
	defer func() {
		fasthttp.ReleaseResponse(response)
		fasthttp.ReleaseRequest(request)
	}()

	// -------------------------------------------------------
	//      构造 web client 请求数据
	// -------------------------------------------------------
	// 指定 HTTP 请求的 URL
	request.SetRequestURI(fullURL)
	request.SetHost("127.0.0.1:3001")

	// 指定 HTTP 请求的方法
	request.Header.SetMethod("GET")
	// 设置 HTTP 请求的 HTTP header

	request.Header.SetBytesKV([]byte("Content-Type"), []byte("text/plain; charset=utf8"))
	request.Header.SetBytesKV([]byte("User-Agent"), []byte("fasthttp-code web client"))
	request.Header.SetBytesKV([]byte("Accept"), []byte("text/plain; charset=utf8"))

	request.SetBody([]byte("payload"))

	// 设置 web client 请求的超时时间
	timeOut := 3 * time.Second

	// 计时开始
	t1 := time.Now()

	// DO request
	err := fasthttp.DoTimeout(request, response, timeOut)
	if err != nil {
		log.Error("post request error", zap.Error(err))
		os.Exit(-1)
	}
	// -------------------------------------------------------
	//      处理返回结果
	// -------------------------------------------------------
	elapsed := time.Since(t1)

	log.Debug("fasthttp client", zap.Duration("elapsed", elapsed))

	utils.RequestDebug(request, log.Named("request"), true)
	log.Debug("=====================================================================")
	utils.ResponseDebug(response, log.Named("response"), true)
}
