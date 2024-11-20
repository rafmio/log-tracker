package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

type totalStatsResult struct {
	totalRecords   int               `json:"total_records"`
	uniqueIPCount  int               `json:"unique_ip_count"`
	recordsPerDay  float64           `json:"records_per_day"`
	topTenIPs      map[string]int    `json:"top_10_ips"`
	topTenDpt      map[string]int    `json:"top_10_dpt"`
	mapIpCountries map[string]string // map[ip_string]"Country Name"
}

func extractStatIndicators(statIndicatorsRows map[string]*sql.Rows) (*totalStatsResult, error) {
	tsr := new(totalStatsResult)

	if rows, ok := statIndicatorsRows["total_records"]; ok {
		var count int
		if rows.Next() {
			if err := rows.Scan(&count); err != nil {
				return tsr, err
			}
		}
		tsr.totalRecords = count
	}

	if rows, ok := statIndicatorsRows["unique_ip_count"]; ok {
		var uniqCount int
		if rows.Next() {
			if err := rows.Scan(&uniqCount); err != nil {
				return tsr, err
			}
		}
		tsr.uniqueIPCount = uniqCount
	}

	if rows, ok := statIndicatorsRows["records_per_day"]; ok {
		var avg float64
		if rows.Next() {
			if err := rows.Scan(&avg); err != nil {
				return tsr, err
			}
		}
		tsr.recordsPerDay = avg
	}

	if rows, ok := statIndicatorsRows["top_10_ips"]; ok {
		tsr.topTenIPs = make(map[string]int)
		for rows.Next() {
			var ip string
			var count int
			if err := rows.Scan(&ip, &count); err != nil {
				return tsr, err
			}
			tsr.topTenIPs[ip] = count
		}
	}

	if rows, ok := statIndicatorsRows["top_10_dpt"]; ok {
		tsr.topTenDpt = make(map[string]int)
		for rows.Next() {
			var dpt string
			var count int
			if err := rows.Scan(&dpt, &count); err != nil {
				return tsr, err
			}
			tsr.topTenDpt[dpt] = count
		}
	}

	return tsr, nil
}

// matchIpAndCountry matches IP addresses with their corresponding countries
// makes GET requests to a geolocation API
func matchIpAndCountry(tsr map[string]string) (map[string]string, error) {
	countries := make(map[string]string)
	// Make API requests to geolocation API here
	for ip := range tsr {
		resp, err := http.Get(fmt.Sprintf("http://ipwho.is/%s", ip))
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		// Parse and extract country from API response
		//
	}
}
