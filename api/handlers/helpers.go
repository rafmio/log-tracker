package main

import (
	"errors"
	"net/http"
	"time"
)

var (
	ErrParsingForm        = errors.New("Error parsing form")
	ErrIncorrectSetDate   = errors.New("Invalid start or end date")
	ErrIncorrectDateRange = errors.New("Interval is more than 48 hours")
)

func preProcessResponse(w http.ResponseWriter, r *http.Request) error {
	// checking whether the HTTP request method is a GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are supported", http.StatusMethodNotAllowed)

		return http.ErrNotSupported
	}

	// SET HEADERS ----------------------------------------------
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")             // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // methods 'PUT', 'PATCH' and 'DELETE' has been deleted
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, hx-request, hx-target, hx-current-url")

	// parse form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return ErrParsingForm
	}

	return nil
}

func parseAndValidateDateRange(startDateParam, endDateParam, layoutDateTime string) ([2]time.Time, error) {
	var dateRange [2]time.Time

	startDate, err := time.Parse(layoutDateTime, r.FormValue(startDateParam))
	if err != nil {
		return dateRange, ErrIncorrectSetDate
	}
	endDate, err := time.Parse(layoutDateTime, r.FormValue(endDateParam))
	if err != nil {
		return dateRange, ErrIncorrectSetDate
	}

	// check if interval is valid (less than 48 hours)
	if endDate.Sub(startDate) > time.Hour*48 {
		return dateRange, ErrIncorrectDateRange
	}

	dateRange = [2]time.Time{startDate, endDate}
	return dateRange, nil
}
