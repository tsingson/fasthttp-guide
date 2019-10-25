package main

import (
	"fmt"
	"runtime"

	"go.uber.org/zap"

	"github.com/tsingson/fasthttp-example/webserver"
)

func main() {
	runtime.MemProfileRate = 0
	runtime.GOMAXPROCS(128)

	fmt.Println("----- fasthttp server starting -----")
	cfg := webserver.Default()
	s := webserver.NewServer(cfg)

	stopSignal := make(chan struct{})

	undo := zap.RedirectStdLog(s.Log)
	defer undo()

	err := s.Run()
	if err != nil {
		panic("server start fail")
	}
	// select {}
	<-stopSignal
}
