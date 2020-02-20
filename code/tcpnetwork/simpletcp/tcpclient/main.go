package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

func main() {
	err := tcpclient()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func tcpclient() error {
	target := "localhost:8000"

	addr, err := net.ResolveTCPAddr("tcp", target)
	if err != nil {
		return err
	}

	// Establish a connection with the server.
	var conn *net.TCPConn
	conn, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	// _ = conn.SetNoDelay(false) // Disable TCP_NODELAY; Nagle's Algorithm takes action.

	fmt.Println("Sending Gophers down the pipe...")

	for i := 1; i <= 5; i++ {
		// Send the word "GOPHER" to the open connection.
		_, err = conn.Write([]byte(string(strconv.Itoa(i)) + "GOPHER\n"))
		if err != nil {
			// TODO: handle error
			continue
		}
		time.Sleep(1 * time.Millisecond)
	}
	_, _ = conn.Write([]byte("EXIT\n"))
	var message []byte = make([]byte, 1024)
	n, er1 := bufio.NewReader(conn).Read(message)
	if er1 != nil {
		if err == io.EOF {
			// TODO: handle io.EOF
			fmt.Println(err)
		}
		fmt.Println(err)

	}
	if n > 0 {
		fmt.Println(string(message))
	}
	time.Sleep(100 * time.Millisecond)
	return nil
}
