package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewConsoleDebug  new zap logger for console
func NewConsoleDebug() zapcore.Core {
	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Optimize the console output for human operators.
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.

	var stderr = zapcore.NewCore(consoleEncoder, consoleErrors, highPriority)
	var stdout = zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority)

	return zapcore.NewTee(stderr, stdout)
}

// ConsoleWithStack  console log for debug
func ConsoleWithStack() *zap.Logger {
	core := NewConsoleDebug()
	// From a zapcore.Core, it's easy to construct a Logger.
	return zap.New(core).WithOptions(zap.AddCaller())
}

// Console  console log for debug
func Console() *zap.Logger {
	core := NewConsoleDebug()
	// From a zapcore.Core, it's easy to construct a Logger.
	return zap.New(core)
}
