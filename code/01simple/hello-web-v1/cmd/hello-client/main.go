package main

import (
	"net/url"
	"os"
	"time"

	"github.com/integrii/flaggy"
	"github.com/savsgio/gotils"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/tsingson/fasthttp-guide/logger"
	"github.com/tsingson/fasthttp-guide/pkg/utils"
)

func main() {
	var log *zap.Logger = logger.Console()

	// baseURL to call web server
	baseURL := "http://127.0.0.1:3001"

	// Add a flag for baseURL
	flaggy.String(&baseURL, "b", "base-url", "address for webserver")

	// Parse the flag
	flaggy.Parse()

	tid := "12345"
	// -------------------------------------------------------
	//      构造 web client 请求的 URL
	// -------------------------------------------------------

	var fullURL string
	{
		relativeUrl := "/uri/"
		var u *url.URL
		var err error
		u, err = url.Parse(relativeUrl)
		if err != nil {
			log.Fatal("error", zap.Error(err))
		}

		queryString := u.Query()

		// 这里构造 URI 中的数据, 每一个键值对
		{
			queryString.Set("id", "1")
			queryString.Set("who", "tsingson")
			queryString.Set("where", "中国深圳")
		}

		u.RawQuery = queryString.Encode()

		var base *url.URL

		base, err = url.Parse(baseURL)
		if err != nil {
			log.Fatal("error", zap.Error(err))
			os.Exit(-1)
		}

		fullURL = base.ResolveReference(u).String()

		log.Debug("---------------- HTTP 请求 URL -------------")

		log.Debug(tid, zap.String("http request URL > ", fullURL))

	}

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

	// 指定 HTTP 请求的方法
	request.Header.SetMethod("GET")
	// 设置 HTTP 请求的 HTTP header

	request.Header.SetBytesKV([]byte("Content-Type"), []byte("text/plain; charset=utf8"))
	request.Header.SetBytesKV([]byte("User-Agent"), []byte("fasthttp-code web client"))
	request.Header.SetBytesKV([]byte("Accept"), []byte("text/plain; charset=utf8"))
	request.Header.SetBytesKV([]byte("TransactionID"), []byte(tid))

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
	log.Debug("---------------- HTTP 响应消耗时间-------------")

	log.Debug(tid, zap.Duration("elapsed", elapsed))
	log.Debug("---------------- HTTP 响应状态码 -------------")

	log.Debug(tid, zap.Int("http status code", response.StatusCode()))
	log.Debug("---------------- HTTP 响应 header 与 payload -------------")

	// -------------------------------------------------------
	// 注意对比一下, 下面的代码段, 与 web server  中几乎一样
	// -------------------------------------------------------
	{
		// 取出 web client 请求中的 HTTP header
		{
			log.Debug("---------------- HTTP header 每一个键值对-------------")
			response.Header.VisitAll(func(key, value []byte) {
				// l.Info("requestHeader", zap.String("key", gotils.B2S(key)), zap.String("value", gotils.B2S(value)))
				log.Debug(tid, zap.String("key", utils.B2S(key)), zap.String("value", utils.B2S(value)))
			})

		}
		// 取出 web client 请求中的 HTTP payload
		{
			log.Debug("---------------- HTTP payload -------------")
			log.Debug(tid, zap.String("http payload", gotils.B2S(response.Body())))
		}
	}
}
