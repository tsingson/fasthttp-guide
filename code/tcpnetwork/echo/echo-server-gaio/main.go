// code copy from https://raw.githubusercontent.com/xtaci/gaio/master/examples/echo-server/main.go

package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/xtaci/gaio"
)

// this goroutine will wait for all io events, and sents back everything it received
// in async way
func echoServer(w *gaio.Watcher) {
	for {
		// loop wait for any IO events
		results, err := w.WaitIO()
		if err != nil {
			log.Println(err)
			return
		}

		for _, res := range results {
			switch res.Operation {
			case gaio.OpRead: // read completion event
				fmt.Println(" Read completion EVENT " + time.Now().Format(time.RFC1123))
				if res.Error != nil {
					_ = w.Free(res.Conn)
					continue
				}
				if res.Size > 0 {
					// send back everything, we won't start to read again until write completes.
					// submit an async write request
					fmt.Println("\n  data size ", res.Size)
					_ = w.Write(nil, res.Conn, res.Buffer[:res.Size])
				}
			case gaio.OpWrite: // write completion event
				fmt.Println(" write completion EVENT " + time.Now().Format(time.RFC1123))
				if res.Error == nil {
					// since write has completed, let's start read on this conn again
					fmt.Println(" try to read from client " + time.Now().Format(time.RFC1123))
					_ = w.Read(nil, res.Conn, res.Buffer[:cap(res.Buffer)])
				}
			}
		}
	}
}

func main() {
	w, err := gaio.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	go echoServer(w)

	var ln net.Listener

	ln, err = net.Listen("tcp", "localhost:3001")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("echo httpserver listening on", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("new httpclient", conn.RemoteAddr())

		// submit the first async read IO request
		err = w.Read(nil, conn, nil )
		if err != nil {
			log.Println(err)
			return
		}
	}
}
