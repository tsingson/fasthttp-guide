package webserver

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
	"github.com/valyala/fasthttp/reuseport"

	"github.com/tsingson/fasthttp-example/logger"
)

func TestWebServer_hello(t *testing.T) {
	// setup logger that output to console
	log := logger.Console()
	ws := &webServer{
		Addr:  ":3000",
		Log:   log,
		debug: false,
	}

	handler := ws.helloWorldGetHandler()
	// setup fasthttp logger

	flog := logger.InitZapLogger(ws.Log)

	// setup fasthttp server

	s := &fasthttp.Server{
		Handler: handler,
		Logger:  flog,
	}

	// setup listener
	ln, _ := reuseport.Listen("tcp4", ws.Addr)

	// remember to close listener
	defer func() {
		_ = ln.Close()
	}()

	// now running fasthttp server in a goroutine

	go func() {
		_ = s.Serve(ln)
	}()

	// -------------------------------------------------------
	// now, the real http client what you want
	// -------------------------------------------------------

	c := &fasthttp.HostClient{
		Addr:                          "localhost:3000",
		DisableHeaderNamesNormalizing: true,
	}

	// http client Fetch the testing fasthttp server  via local proxy.
	statusCode, body, err := c.Get(nil, "http://google.com/hello")
	assert.NoError(t, err)
	assert.Equal(t, statusCode, 200)
	assert.Equal(t, body, []byte(`hello world`))
}

func TestWebServer_hello2(t *testing.T) {
	// setup logger that output to console
	log := logger.Console()
	ws := &webServer{
		Addr:  ":3000",
		Log:   log,
		debug: false,
	}

	handler := ws.helloWorldGetHandler()
	// setup fasthttp logger

	// setup fasthttp server
	// setup listener , it's fasthttp in memory listener for TESTING only
	ln := fasthttputil.NewInmemoryListener()

	// now running fasthttp server in a goroutine

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	// remember to close listener
	defer func() {
		_ = ln.Close()
	}()

	// -------------------------------------------------------
	// now, the real http client what you want
	// -------------------------------------------------------

	client := &fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
		DisableHeaderNamesNormalizing: true,
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.SetHost("localhost:3000")
	req.Header.SetMethod("GET")

	// http client Fetch the testing fasthttp server  via local proxy.
	err := client.Do(req, resp)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode(), 200)
	assert.Equal(t, resp.Body(), []byte(`hello world`))
}

func BenchmarkWebServer_hello(b *testing.B) {
	// setup logger that output to console
	log := logger.Console()
	ws := &webServer{
		Addr:  ":3000",
		Log:   log,
		debug: false,
	}

	handler := ws.helloWorldGetHandler()
	// setup fasthttp logger

	// setup fasthttp server
	// setup listener , it's fasthttp in memory listener for TESTING only
	ln := fasthttputil.NewInmemoryListener()

	// now running fasthttp server in a goroutine

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	// remember to close listener
	defer func() {
		_ = ln.Close()
	}()

	// -------------------------------------------------------
	// now, the real http client what you want
	// -------------------------------------------------------

	client := &fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
		DisableHeaderNamesNormalizing: true,
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.SetHost("localhost:3000")
	req.Header.SetMethod("GET")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = client.Do(req, resp)
	}
}

func TestHello(t *testing.T) {
	// setup fasthttp server
	addr := ":3000"
	s := &fasthttp.Server{
		Handler: Hello(),
	}
	// setup listener
	ln, _ := reuseport.Listen("tcp4", addr)

	// remember to close listener
	defer func() {
		_ = ln.Close()
	}()

	// now running fasthttp server in a goroutine

	go func() {
		_ = s.Serve(ln)
	}()

	// -------------------------------------------------------
	// now, the real http client what you want
	// -------------------------------------------------------

	c := &fasthttp.HostClient{
		Addr:                          "localhost:3000",
		DisableHeaderNamesNormalizing: true,
	}

	// http client Fetch the testing fasthttp server  via local proxy.
	statusCode, body, err := c.Get(nil, "http://google.com/hello")
	assert.NoError(t, err)
	assert.Equal(t, statusCode, 200)
	assert.Equal(t, body, []byte(`hello world`))
}

func BenchmarkHello(b *testing.B) {
	b.ReportAllocs()
	addr := ":3000"
	s := &fasthttp.Server{
		Handler: Hello(),
	}
	// setup listener
	ln, _ := reuseport.Listen("tcp4", addr)

	// remember to close listener
	defer func() {
		_ = ln.Close()
	}()

	// now running fasthttp server in a goroutine

	go func() {
		_ = s.Serve(ln)
	}()
}
