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

	// the result of sending SQL queries to the database
	queryResults map[string]float64 // map[indicatorName]float64
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

	// for _, dbConfig := range lt.dbConfigs {
	// 	db, err := sql.Open(dbConfig.DriverName, lt.dsns[dbConfig.Name])
	// 	if err != nil {
	// 		return err
	// 	}

	// 	lt.dbs[dbConfig.Name] = db // don't forget to db.Close()!
	// }

	return nil
}

func (lt *ltGeneralStats) generalStatQueryParams() {
	lt.statIndicators = make(map[string]string)
	lt.statIndicators["totalNumberEntries"] = "Total number of entries"
	lt.statIndicators["uniqueIpCount"] = "Unique IP count"
	lt.statIndicators["entriesPerDay"] = "Average number of entries per day"

	lt.queryStatIndicators = make(map[string]string)
	lt.queryStatIndicators["totalNumberEntries"] = `SELECT COUNT(*) FROM lg_tab;`
	lt.queryStatIndicators["uniqueIpCount"] = `SELECT COUNT(DISTINCT srcip) FROM lg_tab;`
	lt.queryStatIndicators["entriesPerDay"] = `
		SELECT AVG(daily_count) AS average_records_per_day
		FROM (
    		SELECT COUNT(*) AS daily_count
    		FROM lg_tab
    		GROUP BY DATE(tmstmp)
		) AS daily_counts;
		 `
}

func (lt *ltGeneralStats) makeQuery(queryResults map[string]float64) error {

	return nil
}
