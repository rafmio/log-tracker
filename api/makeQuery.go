package main

type generalStatQueryParams struct {
	statIndicators      map[string]string // names of statistical indicators: map["internalName"]"name for displaying"
	queryStatIndicators map[string]string // SQL queries itself: map["internalName"]"SQL query"
}

func newGeneralStatQueryParams() *generalStatQueryParams {
	newGSQP := new(generalStatQueryParams)

	newGSQP.statIndicators = make(map[string]string)
	newGSQP.statIndicators["totalNumberEntries"] = "Total number of entries"
	newGSQP.statIndicators["uniqueIpCount"] = "Unique IP count"

	newGSQP.queryStatIndicators = make(map[string]string)
	newGSQP.queryStatIndicators["totalNumberEntries"] = `SELECT COUNT(*) FROM lg_tab;`
	newGSQP.queryStatIndicators["uniqueIpCount"] = `SELECT COUNT(DISTINCT srcip) FROM lg_tab;`
	newGSQP.queryStatIndicators["entriesPerDay"] = `
		SELECT AVG(daily_count) AS average_records_per_day
		FROM (
    		SELECT COUNT(*) AS daily_count
    		FROM lg_tab
    		GROUP BY DATE(tmstmp)
		) AS daily_counts;
		 `

	return newGSQP
}
