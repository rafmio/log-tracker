package parser

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type StubFileConfig struct {
	Pattern        string `json:"pattern"`
	ExcludePattern string `json:"excludePattern"`
	Directory      string `json:"directory"`
	FilePosition   string `json:"filePosition"`
}

func TestReadFileConfig(t *testing.T) {
	// creating temp json config file
	fileConfigToWrite := StubFileConfig{
		Pattern:        "*.log",
		ExcludePattern: "*.gz",
		Directory:      "/home/raf/log-tracker/log-files",
		FilePosition:   "0",
	}

	file, err := os.CreateTemp(".", "fileConfig.json")
	if err != nil {
		log.Println("error creating file")
	} else {
		log.Println("temp file has been created")
	}
	// defer os.Remove(file.Name())
	defer file.Close()

	jsonData, err := json.MarshalIndent(fileConfigToWrite, "", "    ")
	if err != nil {
		log.Println("error marshaling json")
	}

	_, err = file.Write(jsonData)
	if err != nil {
		log.Println("error writing to file")
	}

	t.Run("run reading config of filesystem", func(t *testing.T) {
		// reading the config file
		fileConfigToRead, err := ReadFileConfig("fileConfig.json")
		if err != nil {
			t.Error("fail to read file:", err.Error())
		}

		if fileConfigToRead.Pattern != fileConfigToWrite.Pattern {
			t.Error("fail to read file pattern")
		}
		if fileConfigToRead.ExcludePattern != fileConfigToWrite.ExcludePattern {
			t.Error("fail to read file exclude pattern")
		}
		if fileConfigToRead.Directory != fileConfigToWrite.Directory {
			t.Error("fail to read file directory")
		}
		if fileConfigToRead.FilePosition != fileConfigToWrite.FilePosition {
			t.Error("fail to read file position")
		}

	})
}
