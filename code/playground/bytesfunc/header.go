package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)
type noCopy struct{} //nolint:unused

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type RequestHeader struct {
	noCopy noCopy //nolint:unused,structcheck

	disableNormalizing bool
	noHTTP11           bool
	connectionClose    bool

	// These two fields have been moved close to other bool fields
	// for reducing RequestHeader object size.
	cookiesCollected bool
	rawHeadersParsed bool

	contentLength      int
	contentLengthBytes []byte

	method      []byte
	requestURI  []byte
	host        []byte
	contentType []byte
	userAgent   []byte

	// h     []argsKV
	// bufKV argsKV
	//
	// cookies []argsKV

	rawHeaders []byte

	// stores an immutable copy of headers as they were received from the
	// wire.
	rawHeadersCopy []byte
}


// // Read reads request header from r.
// //
// // io.EOF is returned if r is closed before reading the first header byte.
// func (h *RequestHeader) Read(r *bufio.Reader) error {
// 	n := 1
// 	for {
// 		err := h.tryRead(r, n)
// 		if err == nil {
// 			return nil
// 		}
// 		if err != errNeedMore {
// 			h.resetSkipNormalize()
// 			return err
// 		}
// 		n = r.Buffered() + 1
// 	}
// }

func (h *RequestHeader) tryRead(r *bufio.Reader, n int) error {
	h.resetSkipNormalize()
	b, err := r.Peek(n)
	if len(b) == 0 {
		if err == io.EOF {
			return err
		}

		if err == nil {
			panic("bufio.Reader.Peek() returned nil, nil")
		}

		// This is for go 1.6 bug. See https://github.com/golang/go/issues/14121 .
		if err == bufio.ErrBufferFull {
			return &ErrSmallBuffer{
				error: fmt.Errorf("error when reading request headers: %s", errSmallBuffer),
			}
		}

		// n == 1 on the first read for the request.
		if n == 1 {
			// We didn't read a single byte.
			return errNothingRead{err}
		}

		return fmt.Errorf("error when reading request headers: %s", err)
	}
	b = mustPeekBuffered(r)
	headersLen, errParse := h.parse(b)
	if errParse != nil {
		return headerError("request", err, errParse, b)
	}
	mustDiscard(r, headersLen)
	return nil
}

//
// func (h *RequestHeader) parse(buf []byte) (int, error) {
// 	m, err := h.parseFirstLine(buf)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	var rawHeaders []byte
// 	rawHeaders, _, err = readRawHeaders(h.rawHeaders[:0], buf[m:])
// 	if err != nil {
// 		return 0, err
// 	}
// 	h.rawHeadersCopy = append(h.rawHeadersCopy[:0], rawHeaders...)
// 	var n int
// 	n, err = h.parseHeaders(buf[m:])
// 	if err != nil {
// 		return 0, err
// 	}
// 	h.rawHeaders = append(h.rawHeaders[:0], buf[m:m+n]...)
// 	h.rawHeadersParsed = true
// 	return m + n, nil
// }


// func (h *RequestHeader) parseFirstLine(buf []byte) (int, error) {
// 	bNext := buf
// 	var b []byte
// 	var err error
// 	for len(b) == 0 {
// 		if b, bNext, err = nextLine(bNext); err != nil {
// 			return 0, err
// 		}
// 	}
//
// 	// parse method
// 	n := bytes.IndexByte(b, ' ')
// 	if n <= 0 {
// 		return 0, fmt.Errorf("cannot find http request method in %q", buf)
// 	}
// 	h.method = append(h.method[:0], b[:n]...)
// 	b = b[n+1:]
//
// 	// parse requestURI
// 	n = bytes.LastIndexByte(b, ' ')
// 	if n < 0 {
// 		h.noHTTP11 = true
// 		n = len(b)
// 	} else if n == 0 {
// 		return 0, fmt.Errorf("requestURI cannot be empty in %q", buf)
// 	} else if !bytes.Equal(b[n+1:], strHTTP11) {
// 		h.noHTTP11 = true
// 	}
// 	h.requestURI = append(h.requestURI[:0], b[:n]...)
//
// 	return len(buf) - len(bNext), nil
// }


func (h *RequestHeader) resetSkipNormalize() {
	h.noHTTP11 = false
	h.connectionClose = false

	h.contentLength = 0
	h.contentLengthBytes = h.contentLengthBytes[:0]

	h.method = h.method[:0]
	h.requestURI = h.requestURI[:0]
	h.host = h.host[:0]
	h.contentType = h.contentType[:0]
	h.userAgent = h.userAgent[:0]

	// h.h = h.h[:0]
	// h.cookies = h.cookies[:0]
	h.cookiesCollected = false

	h.rawHeaders = h.rawHeaders[:0]
	h.rawHeadersParsed = false
}


func readRawHeaders(dst, buf []byte) ([]byte, int, error) {
	n := bytes.IndexByte(buf, '\n')
	if n < 0 {
		return nil, 0, errNeedMore
	}
	if (n == 1 && buf[0] == '\r') || n == 0 {
		// empty headers
		return dst, n + 1, nil
	}

	n++
	b := buf
	m := n
	for {
		b = b[m:]
		m = bytes.IndexByte(b, '\n')
		if m < 0 {
			return nil, 0, errNeedMore
		}
		m++
		n += m
		if (m == 2 && b[0] == '\r') || m == 1 {
			dst = append(dst, buf[:n]...)
			return dst, n, nil
		}
	}
}



const (
	hex2intTable                = "\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x00\x01\x02\x03\x04\x05\x06\a\b\t\x10\x10\x10\x10\x10\x10\x10\n\v\f\r\x0e\x0f\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\n\v\f\r\x0e\x0f\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10"
	toLowerTable                = "\x00\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f !\"#$%&'()*+,-./0123456789:;<=>?@abcdefghijklmnopqrstuvwxyz[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\u007f\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f\xa0\xa1\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xab\xac\xad\xae\xaf\xb0\xb1\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xbb\xbc\xbd\xbe\xbf\xc0\xc1\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xca\xcb\xcc\xcd\xce\xcf\xd0\xd1\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf\xe0\xe1\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9\xea\xeb\xec\xed\xee\xef\xf0\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xfb\xfc\xfd\xfe\xff"
	toUpperTable                = "\x00\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`ABCDEFGHIJKLMNOPQRSTUVWXYZ{|}~\u007f\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f\xa0\xa1\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xab\xac\xad\xae\xaf\xb0\xb1\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xbb\xbc\xbd\xbe\xbf\xc0\xc1\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xca\xcb\xcc\xcd\xce\xcf\xd0\xd1\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf\xe0\xe1\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9\xea\xeb\xec\xed\xee\xef\xf0\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xfb\xfc\xfd\xfe\xff"
	quotedArgShouldEscapeTable  = "\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x01\x01\x01\x01\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x01\x01\x01\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x01\x01\x00\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01"
	quotedPathShouldEscapeTable = "\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x00\x01\x00\x01\x01\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x01\x01\x01\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x01\x01\x00\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01"
)


func normalizeHeaderValue(ov, ob []byte, headerLength int) (nv, nb []byte, nhl int) {
	nv = ov
	length := len(ov)
	if length <= 0 {
		return
	}
	write := 0
	shrunk := 0
	lineStart := false
	for read := 0; read < length; read++ {
		c := ov[read]
		if c == '\r' || c == '\n' {
			shrunk++
			if c == '\n' {
				lineStart = true
			}
			continue
		} else if lineStart && c == '\t' {
			c = ' '
		} else {
			lineStart = false
		}
		nv[write] = c
		write++
	}

	nv = nv[:write]
	copy(ob[write:], ob[write+shrunk:])

	// Check if we need to skip \r\n or just \n
	skip := 0
	if ob[write] == '\r' {
		if ob[write+1] == '\n' {
			skip += 2
		} else {
			skip++
		}
	} else if ob[write] == '\n' {
		skip++
	}

	nb = ob[write+skip : len(ob)-shrunk]
	nhl = headerLength - shrunk
	return
}

func normalizeHeaderKey(b []byte, disableNormalizing bool) {
	if disableNormalizing {
		return
	}

	n := len(b)
	if n == 0 {
		return
	}

	b[0] = toUpperTable[b[0]]
	for i := 1; i < n; i++ {
		p := &b[i]
		if *p == '-' {
			i++
			if i < n {
				b[i] = toUpperTable[b[i]]
			}
			continue
		}
		*p = toLowerTable[*p]
	}
}

// AppendNormalizedHeaderKey appends normalized header key (name) to dst
// and returns the resulting dst.
//
// Normalized header key starts with uppercase letter. The first letters
// after dashes are also uppercased. All the other letters are lowercased.
// Examples:
//
//   * coNTENT-TYPe -> Content-Type
//   * HOST -> Host
//   * foo-bar-baz -> Foo-Bar-Baz
func AppendNormalizedHeaderKey(dst []byte, key string) []byte {
	dst = append(dst, key...)
	normalizeHeaderKey(dst[len(dst)-len(key):], false)
	return dst
}

// AppendNormalizedHeaderKeyBytes appends normalized header key (name) to dst
// and returns the resulting dst.
//
// Normalized header key starts with uppercase letter. The first letters
// after dashes are also uppercased. All the other letters are lowercased.
// Examples:
//
//   * coNTENT-TYPe -> Content-Type
//   * HOST -> Host
//   * foo-bar-baz -> Foo-Bar-Baz
func AppendNormalizedHeaderKeyBytes(dst, key []byte) []byte {
	return AppendNormalizedHeaderKey(dst, b2s(key))
}

var (
	errNeedMore    = errors.New("need more data: cannot find trailing lf")
	errInvalidName = errors.New("invalid header name")
	errSmallBuffer = errors.New("small read buffer. Increase ReadBufferSize")
)

type errNothingRead struct {
	error
}

// ErrSmallBuffer is returned when the provided buffer size is too small
// for reading request and/or response headers.
//
// ReadBufferSize value from Server or clients should reduce the number
// of such errors.
type ErrSmallBuffer struct {
	error
}

func mustPeekBuffered(r *bufio.Reader) []byte {
	buf, err := r.Peek(r.Buffered())
	if len(buf) == 0 || err != nil {
		panic(fmt.Sprintf("bufio.Reader.Peek() returned unexpected data (%q, %v)", buf, err))
	}
	return buf
}

func mustDiscard(r *bufio.Reader, n int) {
	if _, err := r.Discard(n); err != nil {
		panic(fmt.Sprintf("bufio.Reader.Discard(%d) failed: %s", n, err))
	}
}


// b2s converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// s2b converts string to a byte slice without memory allocation.
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func s2b(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}
