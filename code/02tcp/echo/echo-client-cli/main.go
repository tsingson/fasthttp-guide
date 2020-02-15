package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	host = flag.String("host", "localhost", "host")
	port = flag.String("port", "3001", "port")
)

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", *host+":"+*port)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connecting to " + *host + ":" + *port)
	var wg sync.WaitGroup
	wg.Add(2)
	go handleWrite(conn, &wg)
	go handleRead(conn, &wg)
	wg.Wait()
}

func handleWrite(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	tid := strconv.FormatInt(int64(time.Now().Unix()), 10)

	for i := 100; i > 0; i-- {
		_, e := conn.Write([]byte(tid + " " + strconv.Itoa(i) + "\r\n"))
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			break
		}
		time.Sleep(3 * time.Second)
	}
}

func handleRead(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	reader := bufio.NewReader(conn)
	for i := 1; i <= 100; i++ {
		line, err := reader.ReadString(byte('\n'))
		if err != nil {
			fmt.Print("Error to read message because of ", err)
			return
		}
		fmt.Print(line)
	}
}
