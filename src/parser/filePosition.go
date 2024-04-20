package parser

import (
	"io"
	"os"
)

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

}
