package main

import (
	"log"
	"logtracker/dbhandler"
	"logtracker/parser"
	"os"
)

func main() {
	// уставливаем переменные окружения для пути, где лежат файлы
	// и для файловой позиции
	err := os.Setenv("VARLOGPATH", "/home/raf/log-tracker/log-files")
	if err != nil {
		log.Println("can't set env variable")
	}

	err = os.Setenv(parser.VarLogFPEnvVarName, "")
	if err != nil {
		log.Println("can't set env variable")
	}

	path := os.Getenv("VARLOGPATH")

	fp := new(parser.FilePosition)
	err = fp.GetFPFromEnv() // get fp from env
	if err != nil {
		log.Println(err)
	}

	// open the file
	file, err := parser.SelectAndOpen(path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	// check if the file position is correct
	correct, err := fp.IfFPCorrect(file)
	if err != nil {
		log.Println(err)
	}
	if !correct {
		log.Println("incorrect file position")
		fp.filePosition = int64(0)
	}

	// read file since exact file position
	logLines, err := parser.FileReader(file, fp.filePosition)
	if err != nil {
		log.Println(err)
	}

	err = fp.FindFP(file)
	if err != nil {
		log.Println(err)
	}

	err = fp.WriteFPToEnv()
	if err != nil {
		log.Println(err)
	}

	for _, logLine := range logLines {
		logEntry, err := parser.ParseLog(logLine)
		if err != nil {
			log.Println(err)
		}

		err = dbhandler.InsertToDB(logEntry)
		if err != nil {
			log.Println(err)
		}
	}
}
