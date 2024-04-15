package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestSelectAndOpen(t *testing.T) {
	line := "Apr 11 23:03:12 localhost kernel: [23604836.487395] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=212.192.158.71 DST=194.58.102.129 LEN=44 TOS=0x00 PREC=0x00 TTL=248 ID=54321 PROTO=TCP SPT=44496 DPT=2099 WINDOW=65535 RES=0x00 SYN URGP=0 "

	// создаем временую директорию для временных тестовых файлов
	// чтобы проверить корректно ли функция SelectAndOpen()
	// определяет нужный файл (не последний, не пустой)
	tempDir, err := os.MkdirTemp(".", "varlog")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	// defer os.RemoveAll(tempDir) // for delete all temp dirs and files

	// set environment variable for path
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting current dir: %v", err)
	}

	envVarPath := currentDir + "/" + tempDir

	if err = os.Setenv("VARLOGDIR", envVarPath); err != nil {
		t.Fatalf("setting env var VARLOGDIR: %v", err)
	}

	// fill slice with temp filenames:
	fileNames := make([]string, 0, 10)

	for i := 0; i < 5; i++ {
		ufwFileName := "ufw.log"     // база для создания пути временных файлов ufw
		arbFileName := "someLog.log" // база для создания пути прочих файлов

		// изменяем имена файлов для разнообразия, а также добавляем их в слайс
		// для дальнейшей работы с путями (файлы с именями *.1, *2)
		if i != 0 {
			ufwFileName = ufwFileName + "." + strconv.Itoa(i)
			arbFileName = arbFileName + "." + strconv.Itoa(i)

			ufwFileName = filepath.Join(tempDir, ufwFileName)
			arbFileName = filepath.Join(tempDir, arbFileName)
			fileNames = append(fileNames, ufwFileName, arbFileName)
		} else {
			ufwFileName = filepath.Join(tempDir, ufwFileName)
			arbFileName = filepath.Join(tempDir, arbFileName)
			fileNames = append(fileNames, ufwFileName, arbFileName)
		}

		// создаем временные файлы, записываем туда строку, кроме i == 0
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

		// меняем время файлов (по -10 часов * i каджой итерации)
		timeRange := i * -10
		timeToSet := time.Now().Add(time.Hour * time.Duration(timeRange))
		if err = os.Chtimes(ufwFileName, timeToSet, timeToSet); err != nil {
			t.Fatalf("changing time of file %s: %v", ufwFileName, err)
		}
		if err = os.Chtimes(arbFileName, timeToSet, timeToSet); err != nil {
			t.Fatalf("changing time of file %s: %v", arbFileName, err)
		}
	}

	fmt.Println(fileNames)
	// for i := 0; i < 5; i++ {
	// 	fileName := "ufw.log" + strconv.Itoa(i)
	// 	file := filepath.Join(tempDir, fileName)
	// 	if i == 4 {
	// 		break
	// 	}
	// 	if err = os.WriteFile(file, []byte(line), 0644); err != nil {
	// 		t.Fatalf("writing file %s: %v", file, err)
	// 	}

	// 	time.Sleep(time.Second * 2)
	// }

}
