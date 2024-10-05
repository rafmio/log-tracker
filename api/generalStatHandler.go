package main

import (
	"database/sql"
)

type ltGeneralStats struct {
	// DB connection parameters
	dbConfigFilePath string                     // the path to the database connection configuration file - setDbConfigFilePath()
	dbConfigs        map[string]ConnectDBConfig // map database connection configurations returned by readConfig()
	dsns             map[string]string          // data source names for every database, returned by setDSNs(). Format: map["server_name"]"formatString"
	dbs              map[string]*sql.DB         // connections returned by sql.Open() - don't forget to db.Close()!

	// SQL query parameters
	statIndicators      map[string]string // names of statistical indicators: map["internalName"]"name for displaying"
	queryStatIndicators map[string]string // SQL queries itself: map["internalName"]"SQL query"

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

	for _, dbConfig := range lt.dbConfigs {
		db, err := sql.Open(dbConfig.DriverName, lt.dsns[dbConfig.Name])
		if err != nil {
			return err
		}

		lt.dbs[dbConfig.Name] = db // don't forget to db.Close()!
	}

	return nil
}
