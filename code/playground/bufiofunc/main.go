package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

func main() {
	// An artificial input source.
	const input = "1234-----------------------5678\r\n\r\n1234567901234567890"
	reader  := bufio.NewReader(strings.NewReader(input))
	scanner := bufio.NewScanner(reader)
	// Create a custom split function by wrapping the existing ScanWords function.
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// Skip leading spaces.
		if len(data) <= 4 {
			return 0, nil , nil
		}
		start := 0
		// for width := 0; start < len(data); start += width {
		// 	var r rune
		// 	r, width = utf8.DecodeRune(data[start:])
		// 	if !isSpace(r) {
		// 		break
		// 	}
		// }
		// Scan until space, marking end of word.

		for   i := 0; i < len(data)-4; i +=4 {

			   if bytes.Equal(data[start + i:start + i+4] , []byte("\r\n\r\n")) {
				return start + i, data[start:i], nil
			}
		}
		// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
		if atEOF && len(data) > start {
			return len(data), data[start:], nil
		}
		// Request more data.
		return start, nil, nil
	}
	// Set the split function for the scanning operation.
	scanner.Split(split)
	// Validate the input
	// if scanner.Scan() {
	// 	fmt.Printf("---------- %s\n", scanner.Text())
	// }

	if scanner.Scan() {
		l := len( scanner.Bytes())
		fmt.Println(">> ", l )
		bo, e1 := reader.Peek( l )

		if e1  ==nil {
			fmt.Println(string(bo ))
	}
	// if err := scanner.Err(); err != nil {
	// 	fmt.Fprintln(os.Stderr, "reading input:", err)
	// }
	// fmt.Printf("%d\n", count)
	// if err := scanner.Err(); err != nil {
	// 	fmt.Printf("Invalid input: %s", err)
	// }
}


func isSpace(r rune) bool {
	if r <= '\u00FF' {
		// Obvious ASCII ones: \t through \r plus space. Plus two Latin-1 oddballs.
		switch r {
		case ' ', '\t', '\n', '\v', '\f', '\r':
			return true
		case '\u0085', '\u00A0':
			return true
		}
		return false
	}
	// High-valued ones.
	if '\u2000' <= r && r <= '\u200a' {
		return true
	}
	switch r {
	case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}


func funx() {
	requestHeader := `GET /get/path?a=b?c=d HTTP/1.1\r\nHost: httpbin.org\r\nAccept-Language: fr\r\n\r\n`
	s := []byte(requestHeader)

	reader := bufio.NewReader(bytes.NewReader(s))
	fmt.Println("size:", reader.Buffered())
	// rb := make([]byte, 2)

	// reader.Peek(2)
	fmt.Println("size:", reader.Buffered())
}

func tryRead(r *bufio.Reader, n int) error {

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
			return err
		}

		// // n == 1 on the first read for the request.
		// if n == 1 {
		// 	// We didn't read a single byte.
		// 	return errNothingRead{err}
		// }

		return fmt.Errorf("error when reading request headers: %s", err)
	}
	b = mustPeekBuffered(r)
	headersLen := bytes.IndexAny(b, "\r\n")

	mustDiscard(r, headersLen)
	return nil
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
