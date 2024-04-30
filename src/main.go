package logtracker

import (
	"log"
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

	fp := parser.FilePosition{filePosition: 0}
	fp.GetFPFromEnv()


	// open the file
	file, err := parser.SelectAndOpen(path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	// read file since exact file position
	logLines, err := parser.FileReader(file, ???)
	if err != nil {
		log.Println(err)
	}

}
