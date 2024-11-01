package main

import (
	"fmt"
	"net/http"
)

const (
	port = ":8082"
)

func runAPI() error {
	mux := http.NewServeMux() // allocates and returns *http.ServeMux

	fmt.Printf("Listening on port %s\n", port)

	// creating routes
	mux.HandleFunc("/logtracker/fetch_entries", fetchEntriesHandler)
	mux.HandleFunc("/logtracker/statistics/total_stats", totalStatsHandler)

	// running server
	if err := http.ListenAndServe(port, mux); err != nil {
		return err
	}

	return nil
}
