package handlers

import "net/http"

type totalStatsParams struct {
	statIndicatorsNames   map[string]string // map[statistic_name]"Statistic Print Name"
	statIndicatorsQueries map[string]string // map[statistic_name]"SQL Query"

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
}

func newTotalStatsParams() *totalStatsParams {
	g := new(totalStatsParams)
	g.statIndicatorsNames = make(map[string]string)
	g.statIndicatorsQueries = make(map[string]string)

	// add your custom statistics here
	g.statIndicatorsNames["total_records"] = "Total Number of Records"
	g.statIndicatorsNames["unique_ip_count"] = "Unique IP Count"
	g.statIndicatorsNames["records_per_day"] = "Records Per Day"
	g.statIndicatorsNames["top_10_ips"] = "Top 10 Most Frequent IP Addresses"

	// set SQL queries for statistics
	g.statIndicatorsQueries["total_records"] = `SELECT COUNT(*) FROM lg_tab`
	g.statIndicatorsQueries["unique_ip_count"] = `SELECT COUNT(DISTINCT srcip) FROM lg_tab`

	// returned value: average records per day
	g.statIndicatorsQueries["records_per_day"] = `
	   SELECT AVG(daily_count) AS average_records_per_day
		FROM (
		    SELECT COUNT(*) AS daily_count
    		FROM lg_tab
    		GROUP BY DATE(tmstmp)
			) AS subquery;
    `

	// returned value: ip:count
	g.statIndicatorsQueries["top_10_ips"] = `
        SELECT srcip, COUNT(*) AS ip_count
        FROM lg_tab
        GROUP BY srcip
        ORDER BY ip_count DESC
        LIMIT 10
    `

	// returned value: destination_port:count
	g.statIndicatorsQueries["top_10_dpt"] = `
	SELECT dpt, COUNT(*) AS count
	FROM lg_tab
	GROUP BY dpt
	ORDER BY count DESC
	LIMIT 10;
	`

	return g
}

// QUERIES FOR STATISTICS
// g.queryStatIndicators["totalNumberEntries"] = `SELECT COUNT(*) FROM lg_tab`
// g.queryStatIndicators["uniqueIpCount"] = `SELECT COUNT(DISTINCT srcip) FROM lg_tab`
// g.queryStatIndicators["entriesPerDay"] = `
// 	SELECT AVG(daily_count) AS average_records_per_day
// 	FROM (
// 		SELECT COUNT(*) AS daily_count
// 		FROM lg_tab
// 		GROUP BY DATE(tmstmp)
// 		`
