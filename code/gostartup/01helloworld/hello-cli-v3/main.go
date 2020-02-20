// fasthttp-guide
// build:
// cd ......./fasthttp-guide
// go install ./code/01helloworld/hello-cli-v3
// running:
// hello-cli-v3 -who=tsingson

package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/integrii/flaggy"
	"github.com/tsingson/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// exitFail is the exit code if the program
	// fails.
	exitFail = 1
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

func run(args []string, stdout io.Writer) error {
	os.Args = args

	zlog := logger.New(
		logger.WithDebug(),
		logger.WithDays(31),
		logger.WithLevel(zapcore.DebugLevel))
	defer zlog.Sync()

	log := zlog.Named("main")

	who := "Golang"

	// Add a flag
	flaggy.String(&who, "w", "who", "input your name")

	// Parse the flag
	flaggy.Parse()

	_, _ = fmt.Fprint(stdout, "Hello, ", who)
	log.Info("hello", zap.String("input", who))

	time.Sleep(3 * time.Millisecond)

	return nil
}
