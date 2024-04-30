package parser

import (
	"errors"
	"io"
	"log"
	"os"
	"strconv"
)

var ErrIncorrectFilePosition = errors.New("the file position is larger than file size")

var VarLogFPEnvVarName string = "VARLOGFP"

// TODO: consider how to set the environment variable differently

type FilePosition struct {
	filePosition int64
}

// FindFP() method takes a pointer to an open file and finds the current file position
func (fp *FilePosition) FindFP(file *os.File) error {
	filePosition, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	fp.filePosition = filePosition

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

	if fp.filePosition > fileSize {
		return false, ErrIncorrectFilePosition
	}

	return true, nil
}

// GetFPFromEnv() method reads the value of the environment variable VARLOGFP
// func (fp *FilePosition) GetFPFromEnv() int64 {
// 	fpStr := os.Getenv(VarLogFPEnvVarName)
// 	if fpStr == "" {
// 		log.Printf("the %s environment variable was not found. The file position was set to 0", VarLogFPEnvVarName)
// 		return 0
// 	}

// 	// convert fpStr (env var value) from string to int64
// 	fpInt64, err := strconv.ParseInt(fpStr, 10, 64)
// 	if err != nil {
// 		log.Println("can't convert string go integer, the file position was set to 0")
// 		return 0
// 	}

//		return fpInt64
//	}
func (fp *FilePosition) GetFPFromEnv() error {
	fpStr := os.Getenv(VarLogFPEnvVarName)
	if fpStr == "" {
		log.Printf("the %s environment variable was not found. The file position was set to 0", VarLogFPEnvVarName)
		return nil
	}

	// convert fpStr (env var value) from string to int64
	fpInt64, err := strconv.ParseInt(fpStr, 10, 64)
	if err != nil {
		log.Println("can't convert string go integer, the file position was set to 0")
		return err
	}

	fp.filePosition = fpInt64

	return nil
}

// WriteFPToEnv writes file position to environment to VARLOGFP
func (fp *FilePosition) WriteFPToEnv() error {
	err := os.Setenv(VarLogFPEnvVarName, strconv.Itoa(int(fp.filePosition)))
	if err != nil {
		return err
	}

	return nil
}
