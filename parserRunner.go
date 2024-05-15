package parser

import (
	"log"
	"strconv"
)

func ParserRunner() error {
	// set path to config file
	fileConfigName := "/home/raf/log-tracker/config/fileConfig.json"

	// extracting configuration
	fileConfig, err := ReadFileConfig(fileConfigName)
	if err != nil {
		return err
	}

	// select and open file
	file, err := SelectAndOpen(fileConfig)
	if err != nil {
		return err
	}
	defer file.Close()

	// initialize the file position
	fp := new(FilePosition)
	fpInt, err := strconv.Atoi(fileConfig.FilePosition)
	fp.Fp = int64(fpInt)

	// check if the file position is correct
	correct, err := fp.IfFPCorrect(file)
	if err != nil {
		log.Println(err)
		return err
	}
	if !correct {
		log.Println("incorrect file position")
		fp.Fp = int64(0)
	}

	// read the log-file
	logLinesSls, err := FileReader(file, fp.Fp)
	if err != nil {
		return err
	}

	// find the file position after reading
	err = fp.FindFP(file)
	if err != nil {
		return err
	}

	// write the value of file position to the configuration file
}
