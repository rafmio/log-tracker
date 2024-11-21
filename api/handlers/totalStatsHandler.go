package handlers

import (
	"fmt"
	"net/http"
)

// type totalStatsParams struct {
// }

// type serverStats struct {
// 	serverName            string               // server name: 'cute_ganymede', 'black_oxygenium', etc
// 	statIndicatorsNames   map[string]string    // map[statistic_name]"Statistic Print Name"
// 	statIndicatorsQueries map[string]string    // map[statistic_name]"SQL Query"
// 	statIndicatorsRows    map[string]*sql.Rows // the result of the SQL-query
// 	*dbops.DBConfig                            // DB configuration
// 	 readyStatIndicators                  *totalStatsResult
// 	err                   error
// }

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

	// initialize serverStatList []*serverStats
	serverStatList, err := newServerStatsList(dbConfigFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// 1. initialize SQL queries
	// 2. establish DB connection
	// 3. make queries
	for _, srv := range serverStatList {
		srv.initializeSQLqueries()                         // initialize SQL queries
		err := srv.setConfigAndConnectDB(dbConfigFilePath) // establish DB connection
		if err != nil {
			http.Error(w, "Error connecting to database", http.StatusInternalServerError)
			srv.err = err
			continue
		}
		defer srv.DB.Close()

		err = srv.makeQueries() // make queries and populate srv.statIndicatorsRows
		if err != nil {
			http.Error(w, "Error making queries", http.StatusInternalServerError)
			srv.err = err
			continue
		}

		srv.readyStatIndicators, err = extractStatIndicators(srv.statIndicatorsRows)
		if err != nil {
			http.Error(w, "Error extracting statistics", http.StatusInternalServerError)
			srv.err = err
			continue
		}
	}

	// for _, srv := range serverStatList {
	// 	srv.readyStatIndicators, err = extractStatIndicators(srv.statIndicatorsRows)
	// 	if err != nil {
	// 		http.Error(w, "Error extracting statistics", http.StatusInternalServerError)
	// 	}
	// }

	// DEBUG print:
	debugPrint(serverStatList, w)
}

func debugPrint(serverStatsList []*serverStats, w http.ResponseWriter) {
	for _, srv := range serverStatsList {
		fmt.Fprintf(w, "<h1>Server: %s ", srv.serverName)
		fmt.Fprintf(w, "Host: %s</h1>", srv.Host)
		fmt.Fprintf(w, "<p>DSN: %s</p>", srv.Dsn)
		fmt.Fprintf(w, "<h3>Total Stats:</h3>")
		fmt.Fprintf(w, "<p>Total records: %d</p>", srv.readyStatIndicators.totalRecords)
		fmt.Fprintf(w, "<p>Unique IP Count: %d</p>", srv.readyStatIndicators.uniqueIPCount)
		fmt.Fprintf(w, "<p>Records per day: %.2f</p>", srv.readyStatIndicators.recordsPerDay)
		fmt.Fprintf(w, "<h4>Top 10 IPs:</h4>")
		for ip, num := range srv.readyStatIndicators.topTenIPs {
			fmt.Fprintf(w, "<p>%s : %d</p>", ip, num)
		}
		fmt.Fprintf(w, "<h4>Top 10 Department ports:</h4>")
		for dpt, num := range srv.readyStatIndicators.topTenDpt {
			fmt.Fprintf(w, "<p>%s : %d</p>", dpt, num)
		}

	}

}
