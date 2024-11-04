package cmd

import (
	"api/handlers"
	"fmt"
	"net/http"
)

const (
	port = ":8082"
)

func RunAPI() error {
	mux := http.NewServeMux() // allocates and returns *http.ServeMux

	fmt.Printf("Listening on port %s\n", port)

	// creating routes
	mux.HandleFunc("/logtracker/fetch_entries", handlers.FetchEntriesHandler)
	mux.HandleFunc("/logtracker/statistics/total_stats", handlers.TotalStatsHandler)

	// running server
	if err := http.ListenAndServe(port, mux); err != nil {
		return err
	}

	return nil
}
