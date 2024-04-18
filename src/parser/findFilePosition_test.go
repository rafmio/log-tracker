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
	lines := []string{
		"Apr 11 23:00:15 localhost kernel: [23604659.769285] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=91.240.118.248 DST=194.58.102.129 LEN=40 TOS=0x00 PREC=0x00 TTL=250 ID=34638 PROTO=TCP SPT=41605 DPT=63383 WINDOW=1024 RES=0x00 SYN URGP=0 ",
		"Apr 11 23:00:33 localhost kernel: [23604677.653676] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=185.73.125.150 DST=194.58.102.129 LEN=40 TOS=0x00 PREC=0x00 TTL=246 ID=30104 PROTO=TCP SPT=54319 DPT=42274 WINDOW=1024 RES=0x00 SYN URGP=0 ",
		"Apr 11 23:01:06 localhost kernel: [23604710.957509] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=194.120.230.94 DST=194.58.102.129 LEN=40 TOS=0x00 PREC=0x00 TTL=245 ID=54321 PROTO=TCP SPT=59401 DPT=80 WINDOW=65535 RES=0x00 SYN URGP=0 ",
		"Apr 11 23:01:14 localhost kernel: [23604718.786047] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=185.73.125.150 DST=194.58.102.129 LEN=40 TOS=0x00 PREC=0x00 TTL=246 ID=6305 PROTO=TCP SPT=54319 DPT=42607 WINDOW=1024 RES=0x00 SYN URGP=0 ",
		"Apr 11 23:01:32 localhost kernel: [23604736.897883] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=79.124.59.82 DST=194.58.102.129 LEN=40 TOS=0x00 PREC=0x00 TTL=245 ID=22924 PROTO=TCP SPT=43726 DPT=7851 WINDOW=1024 RES=0x00 SYN URGP=0 ",
		"Apr 11 23:01:52 localhost kernel: [23604756.660266] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=87.251.67.173 DST=194.58.102.129 LEN=40 TOS=0x00 PREC=0x00 TTL=246 ID=20220 PROTO=TCP SPT=47854 DPT=41554 WINDOW=1024 RES=0x00 SYN URGP=0 ",
		"Apr 11 23:02:16 localhost kernel: [23604781.065927] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=83.97.73.250 DST=194.58.102.129 LEN=44 TOS=0x00 PREC=0x00 TTL=249 ID=19810 PROTO=TCP SPT=42798 DPT=21977 WINDOW=1025 RES=0x00 SYN URGP=0 ",
		"Apr 11 23:02:34 localhost kernel: [23604798.747909] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=172.233.164.157 DST=194.58.102.129 LEN=40 TOS=0x00 PREC=0x00 TTL=238 ID=56842 PROTO=TCP SPT=51268 DPT=55332 WINDOW=1024 RES=0x00 SYN URGP=0 ",
		"Apr 11 23:02:53 localhost kernel: [23604817.717013] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=193.37.69.142 DST=194.58.102.129 LEN=40 TOS=0x00 PREC=0x00 TTL=246 ID=31037 PROTO=TCP SPT=43047 DPT=43580 WINDOW=1024 RES=0x00 SYN URGP=0 ",
		"Apr 11 23:03:12 localhost kernel: [23604836.487395] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=212.192.158.71 DST=194.58.102.129 LEN=44 TOS=0x00 PREC=0x00 TTL=248 ID=54321 PROTO=TCP SPT=44496 DPT=2099 WINDOW=65535 RES=0x00 SYN URGP=0 ",
	}
	// create a temp dir to test
	tempDir, err := os.MkdirTemp(".", "logs")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}

	// set path to directory with files
	// currentDir, err := os.Getwd()
	// if err != nil {
	// t.Fatalf("getting current dir: %v", err)
	// }

	// fullPathToDirWithLogs := filepath.Join(currentDir, tempDir) // а зачем нам полный путь?

	// slice for storing temp file names
	tmpFiles := make([]string, 0)

	for i := 0; i < 5; i++ {
		fileName := filepath.Join("ufw.log", ".", strconv.Itoa(i))
		tmpFile := filepath.Join(tempDir, fileName)
		tmpFiles = append(tmpFiles, tmpFile)

	}

	return tempDir, tmpFiles
}
