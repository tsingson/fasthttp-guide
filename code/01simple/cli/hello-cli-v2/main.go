package main

import (
	"flag"
	"fmt"
)

func main() {

	var who = flag.String("who", "Golang", "TCP address to listen to")

	flag.Parse()

	fmt.Println("Hello World, ", who)
}
