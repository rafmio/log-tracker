package handlers

import (
	"encoding/json"
	"net/http"
)

func writeServerStatsJSON(w http.ResponseWriter, serverStatList []*serverStats) error {
	// create a slice to hold the JSON-friendly version of serverStats
	type jsonServerStats struct {
		ServerName    string            `json:"server_name"`
		TotalRecords  int               `json:"total_records"`
		UniqueIPCount int               `json:"unique_ip_count"`
		RecordsPerDay float64           `json:"records_per_day"`
		TopTenIPs     map[string]int    `json:"top_10_ips"`
		TopTenDpt     map[string]int    `json:"top_10_dpt"`
		IPCountryMap  map[string]string `json:"map_ip_countries"`
	}

	jsonStats := make([]jsonServerStats, 0, len(serverStatList))

	// Convert serverStats to JSON-friendly struct
	for _, srv := range serverStatList {
		jsonStat := jsonServerStats{
			ServerName:    srv.serverName,
			TotalRecords:  srv.readyStatIndicators.TotalRecords,
			UniqueIPCount: srv.readyStatIndicators.UniqueIPCount,
			RecordsPerDay: srv.readyStatIndicators.RecordsPerDay,
			TopTenIPs:     srv.readyStatIndicators.TopTenIPs,
			TopTenDpt:     srv.readyStatIndicators.TopTenDpt,
			IPCountryMap:  srv.readyStatIndicators.MapIpCountries,
		}
		jsonStats = append(jsonStats, jsonStat)
	}

	// Set the content type header
	// w.Header().Set("Content-Type", "application/json")

	// Encode the JSON and write it to the ResponseWriter
	encoder := json.NewEncoder(w)
	err := encoder.Encode(jsonStats)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return err
	}

	return nil
}
