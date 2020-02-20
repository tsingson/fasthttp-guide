package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/tsingson/logger"
	"github.com/valyala/fasthttp"
)

func main() {
	log := logger.New(logger.WithDays(31), logger.WithDebug())
	defer log.Sync()
	log.Info("fasthttp header parse testing")

	requestHeader := `GET /get/path?a=b?c=d HTTP/1.1
Host: httpbin.org
Accept-Language: fr

`

	htmlForm := `<form action="" method="get" class="form-example">
  <div class="form-example">
    <label for="name">Enter your name: </label>
    <input type="text" name="name" id="name" required>
  </div>
  <div class="form-example">
    <label for="email">Enter your email: </label>
    <input type="email" name="email" id="email" required>
  </div>
  <div class="form-example">
    <input type="submit" value="Subscribe!">
  </div>
</form>`

	requestHTTPMessage := requestHeader + htmlForm
	var req fasthttp.RequestHeader
	br := bufio.NewReader(bytes.NewBufferString(requestHTTPMessage))
	if err := req.Read(br); err != nil {
		fmt.Printf("%v", err)
		os.Exit(-1)
	}

	fmt.Println(req.String())
	fmt.Println(string(req.Header()))
	fmt.Println(string(req.Method()))
	fmt.Println(string(req.RequestURI()))
	fmt.Println(string(req.RequestURI()))
	fmt.Println(string(req.Host()))
	fmt.Println(string(req.Peek("Accept-Language")))

	responseHeader := `HTTP/1.1 200 OK
Date: Sat, 09 Oct 2010 14:28:02 GMT
Server: Apache
Last-Modified: Tue, 01 Dec 2009 20:18:22 GMT
ETag: "51142bc1-7449-479b075b2891b"
Accept-Ranges: bytes
Content-Length: 29769
Content-Type: text/html

`

	html := `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8"/>
    <title>React App</title>
</head>
<body>
<div id="root">
    <h1> golang HTML testing</h1>
    <p>example html for go HTTP payload ( HTTP message for responseHeader ) </p>
</div>
</body>
</html>
`
	responseHTTPMessage := responseHeader + html

	var h fasthttp.ResponseHeader
	bw := bufio.NewReader(bytes.NewBufferString(responseHTTPMessage))
	if err := h.Read(bw); err != nil {
		fmt.Printf("%v", err)
		os.Exit(-1)
	}
	body, err := readBodyIdentity(bw, 1024*4, nil)
	if err == nil {
		fmt.Println("---------------\n", string(body))
	}

	fmt.Println(string(h.Header()))
	log.Info("status code: ", h.StatusCode())
	log.Info("content type: ", string(h.Peek("Content-Type")))

	var v fasthttp.ResponseHeader
	v.SetStatusCode(200)
	v.SetContentType("text/html")
	fmt.Println(string(v.Header()))

	fmt.Println("==========================================")
	var r fasthttp.Request
	r.SetRequestURI("/get/Get")
	r.SetHost("httpbin.org")
	r.SetBodyString("ok")

	size, er2 := r.WriteTo(os.Stdout)
	if er2 != nil {
		fmt.Printf("%v", er2)
	} else {
		fmt.Println("------------size > ", size)
	}
}

// code copy from fasthttp
func readBodyIdentity(r *bufio.Reader, maxBodySize int, dst []byte) ([]byte, error) {
	dst = dst[:cap(dst)]
	if len(dst) == 0 {
		dst = make([]byte, 1024)
	}
	offset := 0
	for {
		nn, err := r.Read(dst[offset:])
		if nn <= 0 {
			if err != nil {
				if err == io.EOF {
					return dst[:offset], nil
				}
				return dst[:offset], err
			}
			panic(fmt.Sprintf("BUG: bufio.Read() returned (%d, nil)", nn))
		}
		offset += nn
		if maxBodySize > 0 && offset > maxBodySize {
			return dst[:offset], ErrBodyTooLarge
		}
		if len(dst) == offset {
			n := round2(2 * offset)
			if maxBodySize > 0 && n > maxBodySize {
				n = maxBodySize + 1
			}
			b := make([]byte, n)
			copy(b, dst)
			dst = b
		}
	}
}

// code copy from fasthttp
func round2(n int) int {
	if n <= 0 {
		return 0
	}
	n--
	x := uint(0)
	for n > 0 {
		n >>= 1
		x++
	}
	return 1 << x
}

// code copy from fasthttp
var ErrBodyTooLarge = errors.New("body size exceeds the given limit")
