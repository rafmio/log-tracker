package logtracker

import (
	"log"
	"logtracker/parser"
	"os"
)

func main() {
	err := os.Setenv("VARLOGPATH", "/home/raf/log-tracker/log-files")
	if err != nil {
		log.Println("can't set env variable")
	}

	path := os.Getenv("VARLOGPATH")

	file, err := parser.SelectAndOpen(path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

}
