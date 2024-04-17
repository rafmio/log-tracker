package parser

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

type stubFilePosition struct {
	filePosition int
}

// the functions under test gets *os.File as an argument,
// evaluates current file position, and save it in env variable.
// This is the test function for it:

func TestFindAndSet(t *testing.T) {

	// set stub env variable
	if err := os.Setenv("VARLOGFP", strconv.Itoa(0)); err != nil {
		t.Fatalf("setting env var VARLOGFP: %v", err)
	}

	// creating a temp directory for temp files to run tests on them
	tmpDirName, tmpFileNames := createTempFiles(t)
	defer os.RemoveAll(tmpDirName) // for delete all temp dirs and files

}

// createTempFiles returns name of directory for further os.RemoveAll(dirname)
func createTempFiles(t testing.TB) (string, []string) {
	t.Helper()
	// test line for write into temp files
	line := "Apr 11 23:03:12 localhost kernel: [23604836.487395] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=212.192.158.71 DST=194.58.102.129 LEN=44 TOS=0x00 PREC=0x00 TTL=248 ID=54321 PROTO=TCP SPT=44496 DPT=2099 WINDOW=65535 RES=0x00 SYN URGP=0 "

	// create a temp dir to test
	tempDir, err := os.MkdirTemp(".", "logs")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}

	// set path to directory with files
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting current dir: %v", err)
	}

	fullPathToDirWithLogs := filepath.Join(currentDir, tempDir) // а зачем нам полный путь?

	// slice for storing temp file names
	tmpFiles := make([]string, 0)

	for i := 0; i < 5; i++ {
		fileName := filepath.Join("ufw.log", ".", strconv.Itoa(i))
	}

	return tempDir, tmpFiles
}
