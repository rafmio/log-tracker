package main

type generalStatsPerServer struct {
	totalNumberEntries string
	uniqueIpCount      string
	entriesPerDay      string
}

type generalStatQueryParams struct {
	// DB column names
	tmstmpColumnName string
	srcIpColumnName  string

	// internal and external names of statistical indicators,
	// SQL queries for every statistical indicator
	statIndicators      map[string]string // names of statistical indicators: map["internalIndicatorName"]"Name For Displaying"
	queryStatIndicators map[string]string // SQL queries itself: map["internalIndicatorName"]"SQL query"

	// queryResults
	generalStatResults map[string]generalStatsPerServer // map["server_name"]map["internalIndicatorName"]float64
}

func (g *generalStatQueryParams) setColumnNames() {
	g.tmstmpColumnName = "tmstmp"
	g.srcIpColumnName = "srcip"
}

func (g *generalStatQueryParams) setStatIndicators(dbConns *DBConnections) {
	// internal and external names of statistical indicators,
	g.statIndicators = make(map[string]string)
	g.statIndicators["totalNumberEntries"] = "Total number of entries"
	g.statIndicators["uniqueIpCount"] = "Unique IP count"
	g.statIndicators["entriesPerDay"] = "Average number of entries per day"

	// SQL queries for every statistical indicator
	g.queryStatIndicators = make(map[string]string)
	g.queryStatIndicators["totalNumberEntries"] = `SELECT COUNT(*) FROM lg_tab;`
	g.queryStatIndicators["uniqueIpCount"] = `SELECT COUNT(DISTINCT srcip) FROM lg_tab;`
	g.queryStatIndicators["entriesPerDay"] = `
		SELECT AVG(daily_count) AS average_records_per_day
		FROM (
			SELECT COUNT(*) AS daily_count
			FROM lg_tab
			GROUP BY DATE(tmstmp)
			`
}
