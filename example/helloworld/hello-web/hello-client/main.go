package main

import (
	"net/url"
	"os"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/tsingson/fasthttp-guide/logger"

	"github.com/tsingson/fasthttp-guide/pkg/vtils"
)

func main() {
	var log *zap.Logger = logger.Console()
	baseURL := "http://127.0.0.1:3001"

	// 随便指定一个字串做为 web 请求的事务ID , 用来打印多条日志时, 区分是否来自同一个 web 请求事务
	tid := "12345678"

	// -------------------------------------------------------
	//      构造 web client 请求的 URL
	// -------------------------------------------------------

	var fullURL string
	{
		relativeUrl := "/uri/"
		u, err := url.Parse(relativeUrl)
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

		base, err := url.Parse(baseURL)
		if err != nil {
			log.Fatal("error", zap.Error(err))
			os.Exit(-1)
		}

		fullURL = base.ResolveReference(u).String()

		log.Debug("---------------- HTTP 请求 URL -------------")

		log.Debug(tid, zap.String("http request URL > ", fullURL))

	}
	// -------------------------------------------------------
	//      fasthttp web client 的初始化, 与清理
	// -------------------------------------------------------
	//  fasthttp 从缓存池中申请 request / response 对象
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	// 释放申请的对象到池中
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()
	// -------------------------------------------------------
	//      构造 web client 请求数据
	// -------------------------------------------------------
	// 指定 HTTP 请求的 URL
	req.SetRequestURI(fullURL)

	// 指定 HTTP 请求的方法
	req.Header.SetMethod("GET")
	// 设置 HTTP 请求的 HTTP header

	req.Header.SetBytesKV([]byte("Content-Type"), []byte("text/plain; charset=utf8"))
	req.Header.SetBytesKV([]byte("User-Agent"), []byte("fasthttp-guide web client"))
	req.Header.SetBytesKV([]byte("Accept"), []byte("text/plain; charset=utf8"))
	req.Header.SetBytesKV([]byte("TransactionID"), []byte(tid))

	// 设置 web client 请求的超时时间
	timeOut := 3 * time.Second

	// 计时开始
	t1 := time.Now()

	// DO request
	err := fasthttp.DoTimeout(req, resp, timeOut)
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

	log.Debug(tid, zap.Int("http status code", resp.StatusCode()))
	log.Debug("---------------- HTTP 响应 header 与 payload -------------")

	// -------------------------------------------------------
	// 注意对比一下, 下面的代码段, 与 web server  中几乎一样
	// -------------------------------------------------------
	{
		// 取出 web client 请求中的 HTTP header
		{
			log.Debug("---------------- HTTP header 每一个键值对-------------")
			resp.Header.VisitAll(func(key, value []byte) {

				log.Debug(tid, zap.String("key", vtils.B2S(key)), zap.String("value", vtils.B2S(value)))
			})

		}
		// 取出 web client 请求中的 HTTP payload
		{
			log.Debug("---------------- HTTP payload -------------")
			log.Debug(tid, zap.String("http payload", vtils.B2S(resp.Body())))
		}
	}
}
