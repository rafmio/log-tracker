package parser

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
)

var ErrIncorrectFilePosition = errors.New("the file position is larger than file size")

type FilePosition struct {
	Fp int64
}

// FindFP() method takes a pointer to an open file and finds the current file position
func (fp *FilePosition) FindFP(file *os.File) error {
	filePosition, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	fp.Fp = filePosition

	return nil
}

// IfFPCorrect() method compares the file position with the file size.
// It returns false if file position greater than file size
func (fp *FilePosition) IfFPCorrect(file *os.File) (bool, error) {
	fi, err := file.Stat()
	if err != nil {
		return false, err
	}
	fileSize := fi.Size()

	if fp.Fp > fileSize {
		return false, ErrIncorrectFilePosition
	}

	return true, nil
}

// the method reads the file position from the configuration file
// and writes it to the FilePosition structure
func (fp *FilePosition) ReadFPFromFile(fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("reading file position:", err.Error())
		return err
	}

	var fileConfig FileConfig

	err = json.Unmarshal(data, &fileConfig)
	if err != nil {
		log.Println("unmarshal data:", err.Error())
		return err
	}

	fp.Fp, err = strconv.ParseInt(fileConfig.FilePosition, 10, 64)

	return nil
}

func (fp *FilePosition) WriteFPToFile(fileConfig FileConfig, fileName string) error {

	file, err := os.Open(fileName)
	if err != nil {
		log.Println("WriteFPToFile(), error opening file", err.Error())
		return err
	}
	defer file.Close()

	fileConfig.FilePosition = strconv.Itoa(int(fp.Fp))
	jsonData, err := json.MarshalIndent(fileConfig, "", "    ")
	if err != nil {
		log.Println("error marshaling json")
		return err
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		log.Println("WriteFPToFile(), error writing to file", err.Error())
		return err
	}

	return nil
}
