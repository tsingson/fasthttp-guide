package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	ips, err := net.LookupIP("httpbin.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		os.Exit(1)
	}
	for _, ip := range ips {
		fmt.Printf("IN A %s\n", ip.String())
	}
	ip := "35.170.216.115"
	if len(ips) > 0 {
		ip = ips[0].String()
	}
	// connect to this socket
	conn, err := net.Dial("tcp", ip+":80")
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	// read in input from stdin
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Text to send: ")
	// text, _ := reader.ReadString('\n')
	request := `GET / HTTP/1.1
Host: httpbin.org
Accept-Language: fr

`
	// send to socket
	fmt.Fprintf(conn, request+"\r\n")
	// listen for reply
	var message []byte = make([]byte, 1024, 1024)
	n, _ := bufio.NewReader(conn).Read(message)
	if n > 0 {
		fmt.Print(string(message))
	}
}
