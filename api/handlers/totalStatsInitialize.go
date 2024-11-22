package handlers

import (
	"api/dbops"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strings"
)

// serverStats is used to store the configs and result of the SQL queries
type serverStats struct {
	serverName            string               // server name: 'cute_ganymede', 'black_oxygenium', etc
	statIndicatorsNames   map[string]string    // map[statistic_name]"Statistic Print Name"
	statIndicatorsQueries map[string]string    // map[statistic_name]"SQL Query"
	statIndicatorsRows    map[string]*sql.Rows // the result of the SQL-query
	*dbops.DBConfig                            // DB configuration
	readyStatIndicators   *totalStatsResult
	err                   error
}

func getServerNames(dbConfigFilePath string) ([]string, error) {
	// open config file
	serverNames := make([]string, 0)

	file, err := os.ReadFile(dbConfigFilePath)
	if err != nil {
		log.Println("opening config file:", err)
		return nil, err
	}

	// unmarshaling JSON data to struct
	dbConfigs := make(map[string]dbops.DBConfig) // variable for storing unmarshaled data
	err = json.Unmarshal(file, &dbConfigs)
	if err != nil {
		log.Println("Unmarshalling JSON:", err)
		return nil, err
	}

	// unmarshaled configs contain information about each server
	// where key is server name and value is DBConfig
	// (we omitted DBConfig for now)
	for srvName := range dbConfigs { // see 'Simplify range: https://pkg.go.dev/golang.org/x/tools/gopls/internal/analysis/simplifyrange'
		if strings.Contains(srvName, "test") {
			continue // we no need "test" servers in production
		} else {
			serverNames = append(serverNames, srvName)
		}
	}

	return serverNames, nil
}

func (s *serverStats) initializeSQLqueries() {
	// s := new(totalStatsParams)
	s.statIndicatorsNames = make(map[string]string)
	s.statIndicatorsQueries = make(map[string]string)

	// map of stat indicator names
	s.statIndicatorsNames["total_records"] = "Total Number of Records"
	s.statIndicatorsNames["unique_ip_count"] = "Unique IP Count"
	s.statIndicatorsNames["records_per_day"] = "Records Per Day"
	s.statIndicatorsNames["top_10_ips"] = "Top 10 Most Frequent IP Addresses"
	s.statIndicatorsNames["top_10_dpt"] = "Top 10 Most Frequent Destination Ports"

	// set SQL queries for statistics
	s.statIndicatorsQueries["total_records"] = `SELECT COUNT(*) FROM lg_tab`
	s.statIndicatorsQueries["unique_ip_count"] = `SELECT COUNT(DISTINCT srcip) FROM lg_tab`

	// returned value: average records per day
	s.statIndicatorsQueries["records_per_day"] = `
	   SELECT AVG(daily_count) AS average_records_per_day
		FROM (
		    SELECT COUNT(*) AS daily_count
    		FROM lg_tab
    		GROUP BY DATE(tmstmp)
			) AS subquery;
    `

	// returned value: ip:count
	s.statIndicatorsQueries["top_10_ips"] = `
        SELECT srcip, COUNT(*) AS ip_count
        FROM lg_tab
        GROUP BY srcip
        ORDER BY ip_count DESC
        LIMIT 10;
    `

	// returned value: destination_port:count
	s.statIndicatorsQueries["top_10_dpt"] = `
	SELECT dpt, COUNT(*) AS count
	FROM lg_tab
	GROUP BY dpt
	ORDER BY count DESC
	LIMIT 10;
	`
}

func (s *serverStats) setConfigAndConnectDB(dbConfigFilePath string) error {
	var err error
	s.DBConfig, err = dbops.NewDBConfig(dbConfigFilePath, s.serverName)
	if err != nil {
		return err
	}

	s.DBConfig.SetDSN()
	err = s.DBConfig.EstablishDbConnection()
	if err != nil {
		return err
	}

	return nil
}
