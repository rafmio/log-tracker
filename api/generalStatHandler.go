package main

type ltGeneralStats struct {
	*DBConnections
	queryResults map[string]float64 // map[indicatorName]float64
}

func (lt *ltGeneralStats) generalStatQueryParams() {
}
