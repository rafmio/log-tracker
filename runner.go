package main

import (
	"log"
	"logtracker/dbhandler"
	"logtracker/parser"
	"strconv"
)

func ParserRunner() error {
	// set path to config file
	fileConfigName := "config/fileConfig.json"
	fileDBConfigName := "config/databaseConfig.json"

	// extracting configuration
	fileConfig, err := parser.ReadFileConfig(fileConfigName)
	if err != nil {
		return err
	}

	// select and open file
	file, err := parser.SelectAndOpen(fileConfig)
	if err != nil {
		return err
	} else {
		log.Printf("the %s file has been selected\n", file.Name())
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
	}
	if !correct {
		log.Println("incorrect file position. The file position has been set to 0")
		fp.Fp = int64(0)
	}

	// read the log-file
	logLinesSls, err := parser.FileReader(file, fp.Fp)
	if err != nil {
		return err
	} else {
		log.Printf("there are %d lines for reading\n", len(logLinesSls))
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
