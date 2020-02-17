// fasthttp-guide
// build:
// cd ......./fasthttp-guide
// go install ./code/01helloworld/hello-cli-v2
// running:
// hello-cli-v2 -who=tsingson

package main

import (
	"flag"
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

	who := flag.String("who", "Golang", "input your name")

	flag.Parse()

	fmt.Fprint(stdout, "Hello, ", *who)
	return nil
}
