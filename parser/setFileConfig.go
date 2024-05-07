package parser

import (
	"encoding/json"
	"os"
)

// Stores the values needed to work with files, directories - paths, file position, name pattern
type FileConfig struct {
	Pattern        string `json:"pattern"`
	ExcludePattern string `json:"excludePattern"`
	Directory      string `json:"directory"`
	FilePosition   string `json:"filePosition"`
}

// Reads the specified configuration file and returns the fileConfig structure
func ReadFileConfig(fileName string) (FileConfig, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return FileConfig{}, err
	}

	var fileConfig FileConfig

	err = json.Unmarshal(data, &fileConfig)
	if err != nil {
		return FileConfig{}, err
	}

	return fileConfig, nil
}
