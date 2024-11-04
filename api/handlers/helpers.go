package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrParsingForm        = errors.New("Error parsing form")
	ErrIncorrectSetDate   = errors.New("Invalid start or end date")
	ErrIncorrectDateRange = errors.New("Interval is more than 48 hours")
)

type fetchEntriesParams struct {
	w                   http.ResponseWriter
	r                   *http.Request
	sourceNameParam     string
	startDateParam      string
	endDateParam        string
	layoutDateTime      string
	startDate           time.Time
	endDate             time.Time
	rawQueryStr         string
	dbName              string
	tableName           string
	timeStampColumnName string
	preparedQuery       string
}

func processURL(w http.ResponseWriter, r *http.Request) error {
	// checking whether the HTTP request method is a GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are supported", http.StatusMethodNotAllowed)

		return http.ErrNotSupported
	}

	// set headers
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

func newFetchEntriesParams(w http.ResponseWriter, r *http.Request) *fetchEntriesParams {
	fep := new(fetchEntriesParams)
	fep.w = w
	fep.r = r
	fep.sourceNameParam = "source_name"
	fep.startDateParam = "start_date"
	fep.endDateParam = "end_date"
	fep.layoutDateTime = "2006-01-02T15:04"
	// fep.rawQueryStr = `SELECT * FROM %s WHERE %s >= $1 AND %s <= $2`
	fep.rawQueryStr = `SELECT * FROM %s WHERE %s BETWEEN %s AND %s`
	fep.dbName = r.FormValue(fep.sourceNameParam)
	fep.tableName = "lg_tab"
	fep.timeStampColumnName = "tmstmp"

	preparedQuery := fmt.Sprintf(
		fep.rawQueryStr,
		fep.tableName,
		fep.timeStampColumnName,
		fep.startDate.String(),
		fep.endDate.String(),
	)
	fep.preparedQuery = preparedQuery

	return fep
}

func (f *fetchEntriesParams) parseAndValidateDateRange() error {

	startDate, err := time.Parse(f.layoutDateTime, f.r.FormValue(f.startDateParam))
	if err != nil {
		return ErrIncorrectSetDate
	}
	endDate, err := time.Parse(f.layoutDateTime, f.r.FormValue(f.endDateParam))
	if err != nil {
		return ErrIncorrectSetDate
	}

	// check if interval is valid (less than 48 hours)
	if endDate.Sub(startDate) > time.Hour*48 {
		return ErrIncorrectDateRange
	}

	f.startDate = startDate
	f.endDate = endDate

	return nil
}
