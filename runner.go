package main

import (
	"log"
	"logtracker/dbhandler"
	"logtracker/parser"
	"strconv"
)

func ParserRunner() error {
	// set path to config file
	fileConfigName := "/home/raf/log-tracker/config/fileConfig.json"
	fileDBConfigName := "/home/raf/log-tracker/config/databaseConfig.json"

	// extracting configuration
	fileConfig, err := parser.ReadFileConfig(fileConfigName)
	if err != nil {
		return err
	}

	// select and open file
	file, err := parser.SelectAndOpen(fileConfig)
	if err != nil {
		return err
	}
	defer file.Close()

	// initialize the file position
	fp := new(parser.FilePosition)
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
	logLinesSls, err := parser.FileReader(file, fp.Fp)
	if err != nil {
		return err
	}

	// find the file position after reading
	err = fp.FindFP(file)
	if err != nil {
		return err
	}

	// write the value of file position to the configuration file
	err = fp.WriteFPToFile(fileConfig, fileConfigName)
	if err != nil {
		return err
	}

	CDBc, err := dbhandler.LoadDatabaseConfig(fileDBConfigName)
	if err != nil {
		log.Println(err)
	}

	for _, logLine := range logLinesSls {
		logEntry, err := parser.ParseLog(logLine)
		if err != nil {
			log.Println(err)
		}

		err = dbhandler.InsertToDb(logEntry, CDBc)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
