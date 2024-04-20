package parser

import (
	"errors"
	"io"
	"os"
)

var ErrIncorrectFilePosition = errors.New("the file position is larger than file size")

type FilePosition struct {
	filePosition int64
}

// the FindFP() method takes a pointer to an open file and finds the current file position
func (fp *FilePosition) FindFP(file *os.File) error {
	filePosition, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	fp.filePosition = filePosition

	return nil
}

// the IfFPCorrect() method compares the file position with the file size
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
