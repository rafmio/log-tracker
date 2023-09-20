package main

import (
	"os"
)

var LogFileName string = "ufw"

func main() {
	file, err := os.OpenFile(LogFileName, os.O_RDONLY, 0)
}
