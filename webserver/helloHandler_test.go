package webserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"

	"github.com/tsingson/fasthttp-example/logger"
)

func TestWebServer_hello(t *testing.T) {
	log := logger.Console()
	ws := &webServer{
		Addr:  ":3000",
		Log:   log,
		debug: true,
	}

	flog := logger.InitZapLogger(ws.Log)
	s := &fasthttp.Server{
		Handler: ws.hello(),
		Logger:  flog,
	}

	ln, _ := reuseport.Listen("tcp4", ws.Addr)
	defer func() {
		_ = ln.Close()
	}()

	go func() {
		_ = s.Serve(ln)
	}()

	// -------------------------------------------------------

	c := &fasthttp.HostClient{
		Addr: "localhost:3000",
	}

	// Fetch google page via local proxy.
	statusCode, body, err := c.Get(nil, "http://google.com/hello")
	assert.NoError(t, err)
	assert.Equal(t, statusCode, 200)
	assert.Equal(t, body, []byte(`hello world`))
}

func TestHello(t *testing.T) {
	log := logger.Console()
	ws := &webServer{
		Addr:  ":3000",
		Log:   log,
		debug: true,
	}

	flog := logger.InitZapLogger(ws.Log)
	s := &fasthttp.Server{
		Handler: Hello(),
		Logger:  flog,
	}

	ln, _ := reuseport.Listen("tcp4", ws.Addr)

	defer func() {
		_ = ln.Close()
	}()

	go func() {
		_ = s.Serve(ln)
	}()

	// -------------------------------------------------------

	c := &fasthttp.HostClient{
		Addr: "localhost:3000",
	}

	// Fetch google page via local proxy.
	statusCode, body, err := c.Get(nil, "http://google.com/hello")
	assert.NoError(t, err)
	assert.Equal(t, statusCode, 200)
	assert.Equal(t, body, []byte(`hello world`))
}
