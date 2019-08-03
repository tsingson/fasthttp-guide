package main

import (
	"fmt"
	"os"
)

func main() {
	var who = "中国"
	if len(os.Args[1]) > 0 {
		who = os.Args[1]
	}
	fmt.Println("Hello World, ", who)
}
