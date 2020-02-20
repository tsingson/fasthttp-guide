// fasthttp-guide
// build:
// cd ......./fasthttp-guide
// go install ./code/01helloworld/hello-cli-v1
// running:
// hello-cli-v1 tsingson

package main

import (
	"fmt"
	"io"
	"os"
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
	who := "Golang"
	if len(os.Args) > 1 && len(os.Args[1]) > 0 {
		who = os.Args[1]
	}
	_, err := fmt.Fprint(stdout, "Hello, ", who)
	return err
}
