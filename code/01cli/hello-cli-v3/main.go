package main

import (
	"fmt"
	"time"

	"github.com/integrii/flaggy"
	"github.com/tsingson/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	log := logger.New(
		logger.WithDebug(),
		logger.WithDays(31),
		logger.WithLevel(zapcore.DebugLevel))
	defer log.Sync()
	// logger.SetLevel(zap.DebugLevel)

	who := "中国"

	// Add a flag
	flaggy.String(&who, "w", "who", "input your name")

	// Parse the flag
	flaggy.Parse()

	fmt.Println("Hello World, ", who)
	log.Info("hello-01cli-v3", zap.String("input", who))

	time.Sleep(3 * time.Millisecond)
}
