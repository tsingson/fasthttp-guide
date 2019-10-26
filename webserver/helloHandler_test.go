package webserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"

	"github.com/tsingson/fasthttp-example/logger"
)

func TestWebServer_hello(t *testing.T) {
	// setup logger that output to console
	log := logger.Console()
	ws := &webServer{
		Addr:  ":3000",
		Log:   log,
		debug: true,
	}
	// setup fasthttp logger

	flog := logger.InitZapLogger(ws.Log)

	// setup fasthttp server

	s := &fasthttp.Server{
		Handler: ws.helloWorldGetHandler(),
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
		Addr: "localhost:3000",
	}

	// http client Fetch the testing fasthttp server  via local proxy.
	statusCode, body, err := c.Get(nil, "http://google.com/hello")
	assert.NoError(t, err)
	assert.Equal(t, statusCode, 200)
	assert.Equal(t, body, []byte(`hello world`))
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
		Addr: "localhost:3000",
	}

	// http client Fetch the testing fasthttp server  via local proxy.
	statusCode, body, err := c.Get(nil, "http://google.com/hello")
	assert.NoError(t, err)
	assert.Equal(t, statusCode, 200)
	assert.Equal(t, body, []byte(`hello world`))
}
