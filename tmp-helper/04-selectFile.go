package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	files, err := filepath.Glob("/home/raf/log-tracker/log-files/ufw.log*")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file)
	}

	fmt.Println("--------------")

	for _, filename := range files {
		fi, err := os.Stat(filename) // getting information about the file
		if err != nil {
			log.Fatal(err)
		}

		if strings.Contains(filename, ".gz") {
			fmt.Println("Hoba!")
			continue
		}
		fmt.Println(filename, "|", "fi.ModTime():", fi.ModTime())
	}
}
