package main

import "net/http"

func totalStatsHandler(w http.ResponseWriter, r *http.Request) {
	err := preProcessResponse(w, &r)
	if err != nil {
		if err == http.ErrNotSupported {
			http.Error(w, "Only GET requests are supported", http.StatusMethodNotAllowed)
		}
		if err == ErrParsingForm {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
		}
	}
}
