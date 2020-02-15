package webserver

import (
	"fmt"
	"net"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"

	"github.com/tsingson/fasthttp-guide/logger"
)

func TestWebServer_simplePostHandler(t *testing.T) {
	// setup logger that output to console
	log := logger.Console()

	// init a webServer to use the post handler method
	ws := &webServer{
		Addr:  ":3000",
		Log:   log,
		debug: false,
	}

	// setup listener , it's fasthttp in memory listener for TESTING only
	ln := fasthttputil.NewInmemoryListener()

	// now running fasthttp server in a goroutine
	handler := ws.simplePostHandler()
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

	postPayloadByte := []byte(`{"actual_start_date":"2019-07-29","actual_end_date":"2019-07-29","plan_start_date":"2019-07-29","plan_end_date":"2019-02-12","title":"养殖计划00002","user_id":2098735545843717147}`)

	// req.SetRequestURI("http://localhost:3000/")
	req.SetHost("localhost:3000")
	req.Header.Add("Accept", "application/json")
	req.Header.SetMethod("POST")

	req.SetBody(postPayloadByte)

	err := client.Do(req, resp)
	assert.NoError(t, err)

	payload := []byte(`{"id":2101127497763529765,"plan_start_date":"2019-07-29","plan_end_date":"2019-02-12","actual_start_date":"2019-07-29","actual_end_date":"2019-07-29","is_done":false,"last_updated":"2019-08-01T14:12:17.983236","is_deleted":false,"user_id":2098735545843717147,"title":"00002"}`)

	assert.Equal(t, resp.StatusCode(), 200)
	assert.Equal(t, resp.Body(), payload)
}

func BenchmarkWebServer_simplePostHandler(b *testing.B) {
	b.ReportAllocs()
	// setup logger that output to console
	log := logger.Console()

	// init a webServer to use the post handler method
	ws := &webServer{
		Addr:  ":3000",
		Log:   log,
		debug: false,
	}
	handler := ws.simplePostHandler()

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
	postPayloadByte := []byte(`{"actual_start_date":"2019-07-29","actual_end_date":"2019-07-29","plan_start_date":"2019-07-29","plan_end_date":"2019-02-12","title":"养殖计划00002","user_id":2098735545843717147}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer func() {
			// 用完需要释放资源
			fasthttp.ReleaseResponse(resp)
			fasthttp.ReleaseRequest(req)
		}()

		// req.SetRequestURI("http://localhost:3000/")
		req.SetHost("localhost:3000")
		req.Header.Add("Accept", "application/json")
		req.Header.SetMethod("POST")

		req.SetBody(postPayloadByte)

		_ = client.Do(req, resp)
	}
}

func BenchmarkWebServer_simplePostHandler2(b *testing.B) {
	// setup logger that output to console
	log := logger.Console()

	// init a webServer to use the post handler method
	ws := &webServer{
		Addr:  ":3000",
		Log:   log,
		debug: false,
	}
	handler := ws.simplePostHandler()

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

	postPayloadByte := []byte(`{"actual_start_date":"2019-07-29","actual_end_date":"2019-07-29","plan_start_date":"2019-07-29","plan_end_date":"2019-02-12","title":"养殖计划00002","user_id":2098735545843717147}`)

	namProc := runtime.GOMAXPROCS(runtime.NumCPU() * 4)

	b.SetParallelism(namProc)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
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

			// req.SetRequestURI("http://localhost:3000/")
			req.SetHost("localhost:3000")
			req.Header.Add("Accept", "application/json")
			req.Header.SetMethod("POST")

			req.SetBody(postPayloadByte)

			_ = client.Do(req, resp)
		}
	})
}
