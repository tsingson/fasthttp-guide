// code copy from https://colobu.com/2014/12/02/go-socket-programming-TCP/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/tsingson/fasthttp-guide/pkg/utils"
)

var (
	host = flag.String("host", "localhost", "host")
	port = flag.String("port", "3001", "port")
)

func main() {
	flag.Parse()
	fmt.Println("try to connect echo server > ")
	conn, err := net.Dial("tcp", *host+":"+*port)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connecting success : " + *host + ":" + *port)
	var wg sync.WaitGroup
	wg.Add(3)
	go handleWrite(conn, &wg)
	go handleRead(conn, &wg)
	go handler(conn , &wg)
	wg.Wait()
}
func handler( conn net.Conn , wg *sync.WaitGroup)  {
	fmt.Println("==============================")
	fmt.Println("==============================")
	fmt.Println("==============================")
	fmt.Println("==============================")
	fmt.Println("==============================")
	defer wg.Done()

	msg, er1 := Generate("123456789", 1024)
	if er1 != nil {
		wg.Done()
	}

	t1 := time.Now()

	for k := 0; k < 100; k++ {

		fmt.Println(" **************************************************  ", len(msg), " "+strconv.Itoa(k ))
		_, e := conn.Write(msg)
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			os.Exit(-1)
		}
		buf := make([]byte, 1024*2)
		reader := bufio.NewReader(conn)

		fmt.Println(" try to read message from server > ")
		line, er2 := reader.Read(buf)
		if er2 != nil {
			fmt.Print("Error to read message because of ", er2 )
			return
		}
		fmt.Println("\n ---->>>>>>>>>>>>>>> ---> read data size ", line)
		fmt.Println(utils.B2S(buf[:line]))

		fmt.Println("************************", time.Since(t1))
	}
}

func handleWrite(conn net.Conn, wg *sync.WaitGroup) {
	fmt.Println(" >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ")
	defer wg.Done()

	// tid := strconv.FormatInt(int64(time.Now().Unix()), 10)

	msg, err := Generate("123456789", 1024)
	if err != nil {
		runtime.Goexit()
	}

	for i := 100; i > 0; i-- {
		fmt.Println(" try to write message to server >>>>>>>>>>>>>  data size is  ", len(msg), " "+strconv.Itoa(i))
		_, e := conn.Write(msg)
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			break
		}
		time.Sleep(300 * time.Microsecond)
	}

}

func handleRead(conn net.Conn, wg *sync.WaitGroup) {
	fmt.Println(" <<<<<<<<<<<<<<<<<<<<<<<<<<<<< ")
	defer wg.Done()

	buf := make([]byte, 1024*2)
	reader := bufio.NewReader(conn)
	for i := 1; i <= 100; i++ {
		fmt.Println(" try to read message from server > ")
		line, err := reader.Read(buf)
		if err != nil {
			fmt.Print("Error to read message because of ", err)
			return
		}
		fmt.Println("\n ---->>>>>>>>>>>>>>> ---> read data size ", line)
		fmt.Println(utils.B2S(buf[:line]))
	}

}
