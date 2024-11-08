package main

import (
	"database/sql"
	"errors"
	"sync"
)

// type generalStatsPerServer struct {
// 	serverName         string
// 	totalNumberEntries string
// 	uniqueIpCount      string
// 	entriesPerDay      string
// }

var (
	ErrListOfStatIndicatorsIsEmpty = errors.New("list of statistical indicators is empty")
)

type generalStatQueryParams struct {
	// DB column names
	tmstmpColumnName string
	srcIpColumnName  string

	// internal and external names of statistical indicators,
	// SQL queries for every statistical indicator
	statIndicators      map[string]string // names of statistical indicators: map["internalIndicatorName"]"Name For Displaying"
	queryStatIndicators map[string]string // SQL queries itself: map["internalIndicatorName"]"SQL query"

	// queryResults
	generalStatResults map[string]map[string]string // map["server_name"]map["internalIndicatorName"]float64
	generalStatErrors  map[string]map[string]error
}

func (g *generalStatQueryParams) setColumnNames() {
	g.tmstmpColumnName = "tmstmp"
	g.srcIpColumnName = "srcip"
}

func (g *generalStatQueryParams) setStatIndicators() {
	// internal and external names of statistical indicators,
	g.statIndicators = make(map[string]string)
	g.statIndicators["totalNumberEntries"] = "Total number of entries"
	g.statIndicators["uniqueIpCount"] = "Unique IP count"
	g.statIndicators["entriesPerDay"] = "Average number of entries per day"

	// SQL queries for every statistical indicator
	g.queryStatIndicators = make(map[string]string)
	g.queryStatIndicators["totalNumberEntries"] = `SELECT COUNT(*) FROM lg_tab`
	g.queryStatIndicators["uniqueIpCount"] = `SELECT COUNT(DISTINCT srcip) FROM lg_tab`
	g.queryStatIndicators["entriesPerDay"] = `
		SELECT AVG(daily_count) AS average_records_per_day
		FROM (
			SELECT COUNT(*) AS daily_count
			FROM lg_tab
			GROUP BY DATE(tmstmp)
			`
}

// makeQuery() range errors in openDBErrs, if err == nil, make query:
// range queryStatIndicators map and fill generalStatResults map
func (g *generalStatQueryParams) makeGeneralStatQuery(dbs map[string]*sql.DB, openDbErrs map[string]error) error {

	switch {
	case len(g.statIndicators) == 0:
		return ErrListOfStatIndicatorsIsEmpty
	case len(dbs) == 0:
		return ErrDSNMapEmpty
	case len(openDbErrs) == 0:
		return ErrOpenDBErrsMapEmpty
	}

	g.generalStatResults = make(map[string]map[string]string)
	g.generalStatErrors = make(map[string]map[string]error)

	var wg sync.WaitGroup
	var mu sync.Mutex

	// range openDbErrs, if error == nil - try to query, else - skip
	for dbName, dbErr := range openDbErrs { // range DB errors
		if dbErr == nil { // if every certain connections is successful
			for statIndName, query := range g.queryStatIndicators { // range over queryStatIndicators map
				rows, err := dbs[dbName].Query(query) // make query from queryStatIndicators map
				if err != nil {                       // if query failed
					// add error to generalStatErrors map
					g.generalStatErrors[dbName] = make(map[string]error)
					g.generalStatErrors[dbName][statIndName] = err
				} else { // if query success
					// add result to generalStatResults map
					defer rows.Close()
					var result string
					for rows.Next() {
						rows.Scan(&result)
					}
					g.generalStatResults[dbName] = make(map[string]string)
					g.generalStatResults[dbName][statIndName] = result
				}
			}
		} else {
			// add error to generalStatErrors map
			g.generalStatErrors[dbName] = make(map[string]error)
			g.generalStatErrors[dbName][dbName] = dbErr
			continue
		}
	}

	return nil
}
