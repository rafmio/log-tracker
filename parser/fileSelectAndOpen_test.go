package parser

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestSelectAndOpen(t *testing.T) {
	// test line for write into temp files
	line := "Apr 11 23:03:12 localhost kernel: [23604836.487395] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=212.192.158.71 DST=194.58.102.129 LEN=44 TOS=0x00 PREC=0x00 TTL=248 ID=54321 PROTO=TCP SPT=44496 DPT=2099 WINDOW=65535 RES=0x00 SYN URGP=0 "

	// create a temp dir to check correctness of the SelectAndOpen func
	tempDir, err := os.MkdirTemp(".", "varlog")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // for delete all temp dirs and files

	// the target directory will be set at env.
	// set environment variable for path
	// currentDir, err := os.Getwd()
	// if err != nil {
	// 	t.Fatalf("getting current dir: %v", err)
	// }

	// construct full path to target test dir
	// envVarPath := filepath.Join(currentDir, tempDir)

	// fill slice with temp filenames (ufw.log*) for range:
	fileNames := make([]string, 0, 10)

	for i := 0; i < 5; i++ {
		ufwFileName := "ufw.log"     // filename base for construct ufw file names
		arbFileName := "someLog.log" // filename base for construct other file names

		// constract exact file names for ufw and other files
		if i != 0 {
			ufwFileName = ufwFileName + "." + strconv.Itoa(i)
			arbFileName = arbFileName + "." + strconv.Itoa(i)

			ufwFileName = filepath.Join(tempDir, ufwFileName) // "varlog/ufw.log
			arbFileName = filepath.Join(tempDir, arbFileName)
			fileNames = append(fileNames, ufwFileName, arbFileName)
		} else {
			ufwFileName = filepath.Join(tempDir, ufwFileName) //"varlog/ufw.log.*
			arbFileName = filepath.Join(tempDir, arbFileName)
			fileNames = append(fileNames, ufwFileName, arbFileName)
		}

		// create temp files, write line there (except for i == 0 'ufw.log')
		if i == 0 {
			if err = os.WriteFile(ufwFileName, []byte(""), 0644); err != nil {
				t.Fatalf("writing file %s: %v", ufwFileName, err)
			}
			if err = os.WriteFile(arbFileName, []byte(""), 0644); err != nil {
				t.Fatalf("writing file %s: %v", arbFileName, err)
			}
		} else {
			if err = os.WriteFile(ufwFileName, []byte(line), 0644); err != nil {
				t.Fatalf("writing file %s: %v", ufwFileName, err)
			}
			if err = os.WriteFile(arbFileName, []byte(line), 0644); err != nil {
				t.Fatalf("writing file %s: %v", arbFileName, err)
			}
		}

		// change time of creation and access (-10 for every ineration)
		timeRange := i * -10
		timeToSet := time.Now().Add(time.Hour * time.Duration(timeRange))
		if err = os.Chtimes(ufwFileName, timeToSet, timeToSet); err != nil {
			t.Fatalf("changing time of file %s: %v", ufwFileName, err)
		}
		if err = os.Chtimes(arbFileName, timeToSet, timeToSet); err != nil {
			t.Fatalf("changing time of file %s: %v", arbFileName, err)
		}
	}

	t.Run("run SelectAndOpen()", func(t *testing.T) {
		var fileConfig = FileConfig{
			Pattern:        "*ufw.log*",
			ExcludePattern: "*.gz",
			// Directory:      "/home/raf/log-tracker/log-files",
			Directory:    tempDir,
			FilePosition: "0",
		}
		// call SelectAndOpen():
		file, _ := SelectAndOpen(fileConfig)
		// want := fileNames[2] // 'ufw.log.1' - latest nonempty file
		// want := filepath.Join(tempDir, fileNames[2])
		want := fileNames[2]
		gotFileName := file.Name()
		if gotFileName != want {
			t.Errorf("got %s, want %s", gotFileName, want)
		}
	})

}

// func createTempFileConfig(t testing.TB) (FileConfig, error) {
// 	t.Helper()
// 	file, err := createTmpJSONFile(t)
// 	if err != nil {
// 		return FileConfig{}, err
// 	}
// 	defer os.Remove(file.Name())
// 	defer file.Close()

// 	fileConfigToRead, err := ReadFileConfig(file.Name())
// 	if err != nil {
// 		return FileConfig{}, err
// 	}

// 	return fileConfigToRead, nil
// }
