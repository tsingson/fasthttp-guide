package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("input:   ")
	scanner()
}

func scanner() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println("input is : ", scanner.Text())
	}
}
