package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	requestHeader := `GET /get/path?a=b?c=d HTTP/1.1
Host: httpbin.org
Accept-Language: fr

`
	s := []byte(requestHeader)
	buf := bytes.NewReader(s)
	reader := bufio.NewReader(buf)
	line, _ := reader.ReadSlice('\n')
	fmt.Printf("the line:%s\n", line)

	const input = "This is The Golang Standard Library.\r\nWelcome you!\r\ntsingson"
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	fmt.Println(count)

	fmt.Println("-----------------------")

	in  := "abcdefghijkl"

	sc  := bufio.NewScanner(strings.NewReader(in ))

	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		fmt.Printf("%t\t%d\t%s\n", atEOF, len(data), data)

		return 0, nil, nil

	}

	sc .Split(split)

	bf  := make([]byte, 1)

	sc .Buffer(bf, bufio.MaxScanTokenSize)

	for sc .Scan() {

		fmt.Printf("%s\n", sc .Text())

	}

	// fmt.Println(bytes.Index(s, []byte(" ")))    // 2
	// fmt.Println(bytes.Index(s, []byte("\r\n"))) // 2
	// fmt.Println(bytes.IndexAny(s, "ole"))       // 1
	// fmt.Println(bytes.IndexByte(s, 'l'))        // 2
	// fmt.Println(bytes.IndexRune(s, '界'))        // 9
	// fmt.Println("-------------------")
	// fmt.Println(bytes.Contains(s, []byte("Hello"))) // true
	// fmt.Println(bytes.ContainsAny(s, "llo"))        // true
	// fmt.Println(bytes.ContainsRune(s, '世'))         // true
	// fmt.Println(bytes.Count(s, []byte("llo")))      // 1
	// fmt.Println(bytes.HasPrefix(s, []byte("llo")))  // false
	// fmt.Println(bytes.HasSuffix(s, []byte("世界")))   // true
	//
	// fmt.Println("-------------------")
	// a := []byte("hello")
	// b := []byte("world")
	// fmt.Println(bytes.Equal(a, b))   // false
	// fmt.Println(bytes.Compare(a, b)) // -1
	// fmt.Println(bytes.Compare(b, a)) // 1
}
