package main

import (
	"fmt"

	"github.com/integrii/flaggy"
)

func main() {
	who := "中国"

	// Add a flag
	flaggy.String(&who, "w", "who", "input your name")

	// Parse the flag
	flaggy.Parse()

	fmt.Println("Hello World, ", who)
}
