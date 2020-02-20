package main

import (
	"fmt"
	"net"
	"os"
)

// only needed below for sample processing

func main() {
	fmt.Println("Launching httpserver http://127.0.0.1:3001")

	// listen on all interfaces
	ln, err := net.Listen("tcp", ":3001")
	if err != nil {
		os.Exit(-1)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
}
