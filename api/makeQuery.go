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

	return newGSQP
}
