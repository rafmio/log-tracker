package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type totalStatsResult struct {
	TotalRecords   int               `json:"total_records"`
	UniqueIPCount  int               `json:"unique_ip_count"`
	RecordsPerDay  float64           `json:"records_per_day"`
	TopTenIPs      map[string]int    `json:"top_10_ips"`
	TopTenDpt      map[string]int    `json:"top_10_dpt"`
	MapIpCountries map[string]string `json:"map_ip_countries"`
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
		tsr.TotalRecords = count
	}

	if rows, ok := statIndicatorsRows["unique_ip_count"]; ok {
		var uniqCount int
		if rows.Next() {
			if err := rows.Scan(&uniqCount); err != nil {
				return tsr, err
			}
		}
		tsr.UniqueIPCount = uniqCount
	}

	if rows, ok := statIndicatorsRows["records_per_day"]; ok {
		var avg float64
		if rows.Next() {
			if err := rows.Scan(&avg); err != nil {
				return tsr, err
			}
		}
		tsr.RecordsPerDay = avg
	}

	if rows, ok := statIndicatorsRows["top_10_ips"]; ok {
		tsr.TopTenIPs = make(map[string]int)
		for rows.Next() {
			var ip string
			var count int
			if err := rows.Scan(&ip, &count); err != nil {
				return tsr, err
			}
			tsr.TopTenIPs[ip] = count
		}
	}

	if rows, ok := statIndicatorsRows["top_10_dpt"]; ok {
		tsr.TopTenDpt = make(map[string]int)
		for rows.Next() {
			var dpt string
			var count int
			if err := rows.Scan(&dpt, &count); err != nil {
				return tsr, err
			}
			tsr.TopTenDpt[dpt] = count
		}
	}

	return tsr, nil
}

// fillIpAndCountryMap() matches IP addresses with their corresponding countries
// makes GET requests to a geolocation API
func (t *totalStatsResult) fillIpAndCountryMap() error {
	t.MapIpCountries = make(map[string]string)

	// data structure for parsing geolocation API response
	type Response struct {
		IP      string `json:"ip"`
		Success bool   `json:"success"`
		Country string `json:"country"`
	}

	// Make API requests to geolocation API here
	for ip := range t.TopTenIPs {
		resp, err := http.Get(fmt.Sprintf("http://ipwho.is/%s", ip))
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		// Parse and extract country from API response
		var response Response
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Println("Error decoding IP information: ", err)
			return err
		}

		if !response.Success {
			fmt.Println("Error retrieving IP information for", ip)
		}

		t.MapIpCountries[ip] = response.Country
		// fmt.Println(ip, ":", response.Country)
	}

	return nil
}
