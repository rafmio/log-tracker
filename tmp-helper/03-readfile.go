package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.ReadFile("text.txt")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(file))
}
