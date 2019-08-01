package main

import (
	"fmt"
	"runtime"

	"github.com/tsingson/fasthttp-example/webserver"
	"go.uber.org/zap"
)

func main() {

	runtime.MemProfileRate = 0
	runtime.GOMAXPROCS(128)

	fmt.Println("----------------------------------------------------")
	//
	// current, _ := utils.GetCurrentPath()
	//
	// tls(current)

	var s = webserver.DefaultServer()

	stopSignal := make(chan struct{})

	undo := zap.RedirectStdLog(s.Log)
	defer undo()

	s.Run()
	// select {}
	<-stopSignal
}
