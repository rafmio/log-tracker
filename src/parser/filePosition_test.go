package parser

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

type TstCase struct {
	name  string
	f     *os.File
	fSize int64 // fs.FileInfo.Size() - want
}

func TestFindFP(t *testing.T) {

	stubFilePosition := new(FilePosition) // stub new instance for FilePosition

	tmpFiles := make([]string, 10)             // slice for storing 10 temp file names
	tmpDirName := createTempFiles(t, tmpFiles) // temp dir for temp files to run tests
	defer os.RemoveAll(tmpDirName)             // for delete all temp dirs and files

	tstCases := make([]TstCase, len(tmpFiles)) // slice for storing test cases

	// fill test cases struct:
	for i := 0; i < len(tmpFiles); i++ {
		t.Run("test FilePosition.FindFP()", func(t *testing.T) {
			tstCases[i].name = tmpFiles[i]                   // initialize name field
			tstCases[i].f, _ = os.Open(tmpFiles[i])          // initialize f field (opened *os.File)
			defer tstCases[i].f.Close()                      // defer closing opened file
			settedFP, _ := tstCases[i].f.Seek(0, io.SeekEnd) // set fp for pass in func further finding
			fi, _ := tstCases[i].f.Stat()                    // get stat info
			tstCases[i].fSize = fi.Size()                    // initialize fSize field
			_ = stubFilePosition.FindFP(tstCases[i].f)       // run FindFP() itself

			if stubFilePosition.filePosition != settedFP || stubFilePosition.filePosition != tstCases[i].fSize {
				t.Errorf("got %d, want %d, want %d", stubFilePosition.filePosition, settedFP, tstCases[i].fSize)
			}
		})
	}
}

// createTempFiles returns name of directory for further os.RemoveAll(dirname)
func createTempFiles(t testing.TB, tmpFiles []string) string {
	t.Helper()
	// test lines for write into temp files
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
	// create a temp dir for storing temp files to test
	tempDir, err := os.MkdirTemp(".", "logs")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}

	for i := 0; i < len(lines); i++ {
		// generating filenames
		fileName := "ufw.log" + "." + strconv.Itoa(i)
		tmpFile := filepath.Join(tempDir, fileName) // "logs/ufw.log.*"
		tmpFiles[i] = tmpFile                       // filling tmpFiles []string slice

		// create len(lens) files, write line there
		f, err := os.Create(tmpFile)
		if err != nil {
			t.Fatalf("creating file: %v", err)
		}
		defer f.Close()

		// filling files +1 line per each iteration
		for j := 0; j <= i; j++ {
			if i == 0 {
				_, err := f.WriteString("") // empty file
				if err != nil {
					t.Fatalf("writing to file: %v", err)
				}
			} else {
				// non-empty files:
				_, err := f.WriteString(lines[j] + "\n")
				if err != nil {
					t.Fatalf("writing to file: %v", err)
				}
			}
		}
	}

	return tempDir
}

func TestIfFPCorrect(t *testing.T) {
	tmpFiles := make([]string, 10)             // slice for storing 10 temp file names
	tmpDirName := createTempFiles(t, tmpFiles) // temp dir for temp files to run tests
	defer os.RemoveAll(tmpDirName)             // for delete all temp dirs and files

	for i := 0; i < len(tmpFiles); i++ {
		t.Run("check if FP is correct", func(t *testing.T) {
			file, _ := os.Open(tmpFiles[i]) // open temp file
			fi, _ := file.Stat()            // save the file's size into fi var

			if i%2 == 0 {
				fp := FilePosition{ // initialize FiliPosition instance
					filePosition: fi.Size(), // with Size()
				}
				result, _ := fp.IfFPCorrect(file)
				if result != true {
					t.Errorf("IfFPCorrect(): got: %v, want: %v", true, result)
				}
			} else {
				fp := FilePosition{
					filePosition: fi.Size() + int64(i+i*1000), // fake incorrect file position
				}
				result, _ := fp.IfFPCorrect(file)
				if result != false {
					t.Errorf("IfFPCorrect(): got: %v, want: %v", false, result)
				}
			}
		})

	}
}

// func TestGetFPFromEnv(t *testing.T) {
// 	for i := 1; i <= 10; i++ {
// 		t.Run("get file position from env", func(t *testing.T) {
// 			want := i * 1000
// 			err := os.Setenv(VarLogFPEnvVarName, strconv.Itoa(want))
// 			if err != nil {
// 				t.Fatalf("os.Setenv(): %v", err)
// 			}
// 			fp := new(FilePosition)
// 			got := fp.GetFPFromEnv()

//				if got != int64(want) {
//					t.Errorf("GetFPFromEnv(): got %v, want %v", got, want)
//				}
//			})
//		}
//	}
func TestGetFPFromEnv(t *testing.T) {
	for i := 1; i <= 10; i++ {
		t.Run("get the file position from env", func(t *testing.T) {
			want := i * 1000
			err := os.Setenv(VarLogFPEnvVarName, strconv.Itoa(want))
			if err != nil {
				t.Fatalf("os.Setenv(): %v", err)
			}
			fp := new(FilePosition)
			_ = fp.GetFPFromEnv()

			if fp.filePosition != int64(want) {
				t.Errorf("GetFPFromEnv(): got %v, want %v", fp.filePosition, want)
			}
		})
	}
}

func TestWriteFPToEnv(t *testing.T) {
	for i := 1; i <= 10; i++ {
		t.Run("write file position to env", func(t *testing.T) {
			want := i * 1000
			fp := FilePosition{
				filePosition: int64(want),
			}
			err := fp.WriteFPToEnv()
			if err != nil {
				t.Fatalf("WriteFPToEnv(): %v", err)
			}
			got := os.Getenv(VarLogFPEnvVarName)
			if got != strconv.Itoa(want) {
				t.Errorf("WriteFPToEnv(): got %v, want %v", got, want)
			}
		})
	}
}
