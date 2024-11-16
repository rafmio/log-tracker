package handlers

import (
	"api/dbops"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type totalStatsParams struct {
	statIndicatorsNames   map[string]string    // map[statistic_name]"Statistic Print Name"
	statIndicatorsQueries map[string]string    // map[statistic_name]"SQL Query"
	statIndicatorsRows    map[string]*sql.Rows // the result of the SQL-query
}

type serverStat struct {
	serverName string // server name: 'cute_ganymede', 'black_oxygenium', etc
	totalStatsParams
	dbops.DBConfig
	err error
}

func TotalStatsHandler(w http.ResponseWriter, r *http.Request) {
	// processURL makes:
	// 	1. checking whether the HTTP request method is a GET method
	// 	2. setting headers
	// 	3. parsing request's parameters and populate r.Form
	err := processURL(w, r)
	if err != nil {
		if err == http.ErrNotSupported {
			http.Error(w, "only GET requests are supported", http.StatusMethodNotAllowed)
		}
		if err == ErrParsingForm {
			http.Error(w, "error parsing form", http.StatusBadRequest)
		}
	}

	dbConfigFilePath := "config/db-config.json"

	var serverStatList []*serverStat

	// filling serverStatList
	err = newServerStatList(dbConfigFilePath, serverStatList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// filling totalStatsParams
	// filling database configs
	for _, srv := range serverStatList {
		srv.totalStatsParams = *newTotalStatsParams()
		srv.DBConfig, err := dbops.NewDBConfig(dbConfigFilePath, srv.serverName)
		if err != nil {
			srv.err = err
		}
	}
}

func newServerStatList(dbConfigFilePath string, serverStatList []*serverStat) error {
	// open config file
	file, err := os.ReadFile(dbConfigFilePath)
	if err != nil {
		log.Println("opening config file:", err)
		return err
	}

	// unmarshalling JSON data to struct
	dbConfigs := make(map[string]dbops.DBConfig) // variable for storing unmarshalled data
	err = json.Unmarshal(file, &dbConfigs)
	if err != nil {
		log.Println("Unmarshalling JSON:", err)
		return err
	}

	for srvName := range dbConfigs { // see 'Simplify range: https://pkg.go.dev/golang.org/x/tools/gopls/internal/analysis/simplifyrange'
		if strings.Contains(srvName, "test") {
			continue
		} else {
			newServerStat := new(serverStat)
			newServerStat.serverName = srvName
			serverStatList = append(serverStatList, newServerStat)
		}
	}

	return nil
}

func newTotalStatsParams() *totalStatsParams {
	tsp := new(totalStatsParams)
	tsp.statIndicatorsNames = make(map[string]string)
	tsp.statIndicatorsQueries = make(map[string]string)

	// add your custom statistics here
	tsp.statIndicatorsNames["total_records"] = "Total Number of Records"
	tsp.statIndicatorsNames["unique_ip_count"] = "Unique IP Count"
	tsp.statIndicatorsNames["records_per_day"] = "Records Per Day"
	tsp.statIndicatorsNames["top_10_ips"] = "Top 10 Most Frequent IP Addresses"

	// set SQL queries for statistics
	tsp.statIndicatorsQueries["total_records"] = `SELECT COUNT(*) FROM lg_tab`
	tsp.statIndicatorsQueries["unique_ip_count"] = `SELECT COUNT(DISTINCT srcip) FROM lg_tab`

	// returned value: average records per day
	tsp.statIndicatorsQueries["records_per_day"] = `
	   SELECT AVG(daily_count) AS average_records_per_day
		FROM (
		    SELECT COUNT(*) AS daily_count
    		FROM lg_tab
    		GROUP BY DATE(tmstmp)
			) AS subquery;
    `

	// returned value: ip:count
	tsp.statIndicatorsQueries["top_10_ips"] = `
        SELECT srcip, COUNT(*) AS ip_count
        FROM lg_tab
        GROUP BY srcip
        ORDER BY ip_count DESC
        LIMIT 10
    `

	// returned value: destination_port:count
	tsp.statIndicatorsQueries["top_10_dpt"] = `
	SELECT dpt, COUNT(*) AS count
	FROM lg_tab
	GROUP BY dpt
	ORDER BY count DESC
	LIMIT 10;
	`

	return tsp
}
