package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	port := ":" + "8000"

	// Create a listening socket.
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		// Accept new connections.
		c, err := ln.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		// Process newly accepted connection.
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	for {
		// Read what has been sent from the client.
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		cdata := strings.TrimSpace(netData)
		if len(cdata) > 0 {
			_, _ = c.Write([]byte("GopherAcademy Advent 2019!\r\n"))
		}

		fmt.Println("--------> ", cdata)

		if cdata == "EXIT" {
			break
		}
	}
	_ = c.Close()
}
