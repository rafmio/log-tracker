package main

import (
	"log"
	"os"
)

func main() {
	err := ParserRunner()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
