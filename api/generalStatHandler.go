package main

import (
	"database/sql"
	"text/template"
)

type ltGeneralStats struct {
	dbConfigFilePath string                     // the path to the database connection configuration file - setDbConfigFilePath()
	dbConfigs        map[string]ConnectDBConfig // map database connection configurations returned by readConfig()
	dsns             map[string]string          // data source names for every database, returned by setDSNs(). Format: map["server_name"]"formatString"
	dbs              map[string]*sql.DB         // connections returned by sql.Open()

	statIndicators      map[string]string
	queryStatIndicators map[string]string

	tmpl *template.Template // html template parsed from a file
}

func (lt *ltGeneralStats) setDBconnectionConfigs() error {
	// Setting the path from where we will read the configuration file to connect to the database
	lt.dbConfigFilePath = setDBconfigFilePath() // return string with 'path/to/file.json'

	lt.dbConfigs = make(map[string]ConnectDBConfig)
	lt.dbConfigs, err := readConfig(lt.dbConfigFilePath)
	if err != nil {
		return err
	}

	setDSNs(lt.dsns, lt.dbConfigs)

	lt.dbs, err = setDBs(lt.dbs, lt.dsns)

	return nil
}
