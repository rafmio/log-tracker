package main

type ltGeneralStats struct {

	// SQL query parameters
	statIndicators      map[string]string // names of statistical indicators: map["internalName"]"name for displaying"
	queryStatIndicators map[string]string // SQL queries itself: map["internalName"]"SQL query"

	// the result of sending SQL queries to the database
	queryResults map[string]float64 // map[indicatorName]float64
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
