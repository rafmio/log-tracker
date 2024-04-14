package parser

import (
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

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting current dir: %v", err)
	}

	envVarPath := currentDir + "/" + tempDir

	if err = os.Setenv("VARLOGDIR", envVarPath); err != nil {
		t.Fatalf("setting env var VARLOGDIR: %v", err)
	}

	for i := 0; i < 5; i++ {
		fileName := "ufw.log" + strconv.Itoa(i)
		file := filepath.Join(tempDir, fileName)
		if i == 4 {
			break
		}
		if err = os.WriteFile(file, []byte(line), 0644); err != nil {
			t.Fatalf("writing file %s: %v", file, err)
		}

		time.Sleep(time.Second * 2)
	}
}
