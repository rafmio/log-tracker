package main

import (
	"encoding/json"
	"log"
	"os"
)

type FileConfig struct {
	Pattern        string `json:"pattern"`
	ExcludePattern string `json:"excludePattern"`
	Directory      string `json:"directory"`
	FilePosition   string `json:"filePosition"`
}

func main() {
	fileConfig := FileConfig{
		Pattern:        "*.log",
		ExcludePattern: "*.gz",
		Directory:      "/home/raf/log-tracker/log-files",
		FilePosition:   "0",
	}

	file, err := os.Create("fileConfig.json")
	if err != nil {
		log.Println("error creating file")
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(fileConfig, "", "    ")
	if err != nil {
		log.Println("error marshaling json")
	}

	_, err = file.Write(jsonData)
	if err != nil {
		log.Println("error writing to file")
	}
}
