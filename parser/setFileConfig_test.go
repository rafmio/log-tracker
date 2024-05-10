package parser

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

var fileConfigToWrite = FileConfig{
	Pattern:        "*ufw.log*",
	ExcludePattern: "*.gz",
	Directory:      "/home/raf/log-tracker/log-files",
	FilePosition:   "0",
}

func TestReadFileConfig(t *testing.T) {
	// creating temp json config file
	file, err := createTmpJSONFile(t)
	if err != nil {
		t.Error("creating tmp JSON file:", err.Error())
	}
	defer os.Remove(file.Name())
	defer file.Close()

	t.Run("run reading config of filesystem", func(t *testing.T) {
		// reading the config file
		fileConfigToRead, err := ReadFileConfig(file.Name())
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

func createTmpJSONFile(t testing.TB) (*os.File, error) {
	t.Helper()
	file, err := os.CreateTemp(".", "fileConfig.json")
	if err != nil {
		log.Println("error creating file")
		return nil, err
	}

	jsonData, err := json.MarshalIndent(fileConfigToWrite, "", "    ")
	if err != nil {
		log.Println("error marshaling json")
		return nil, err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		log.Println("error writing to file")
		return nil, err
	}

	return file, nil
}
