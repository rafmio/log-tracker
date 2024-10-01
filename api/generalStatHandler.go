package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"
)

type ltGeneralStats struct {
	dbConfigFilePath   string             // the path to the database connection configuration file
	dbConfigs          map[string]Source  // map database connection configurations returned by readConfig()
	dsns               map[string]string  // data source names for every database
	db                 *sql.DB            // connection returned by sql.Open()
	tmpl               *template.Template // html template parsed from a file
	totalNumberEntries map[string]int
	uniqueIpCount      map[string]int
}

func (l *ltGeneralStats) setDbConfigFilePath() {
	l.dbConfigFilePath = "db-config.json"
}

func (l *ltGeneralStats) readConfig(dbConfigFilePath string) error {
	// reading file with configuration for DB connection
	file, err := os.ReadFile(dbConfigFilePath)
	if err != nil {
		log.Println("Opening config file:", err)
	}

	// unmarshalling JSON data to struct
	dbConfig := make(map[string]Source)

	err = json.Unmarshal(file, &dbConfig)
	if err != nil {
		log.Println("Unmarshalling JSON:", err)
	}

	return err
}

func (l *ltGeneralStats) setDsn() {
	dsnStringTemplate := "host=%s port=%s user=%s dbname=%s password=%s sslmode=%s"

	for _, dbConfig := range l.dbConfigs {
		l.dsns[dbConfig.Name] = fmt.Sprintf(dsnStringTemplate,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.User,
			dbConfig.DBName,
			dbConfig.Password,
			dbConfig.SslMode,
		)
	}
}
