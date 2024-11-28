package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"web/cmd"
)

var (
	apiAddrFetch = "http://localhost:8082/fetch"
	// apiAddrStat  = "http://localhost:8082/logtracker/statistic/generalstat"
	apiAddrStat = "http://localhost:8082/logtracker/statistic/totalStats"
)

func main() {
	mux := http.NewServeMux()
	cmd.RegisterRoutes(mux)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/about.html")
	})

	http.HandleFunc("/logtracker", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/logtracker.html")
	})

	http.HandleFunc("/budgeting", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/budgeting.html")
	})

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/orders.html")
	})

	http.HandleFunc("/logtracker/logs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/logtracker-logs.html")
	})

	http.HandleFunc("/logtracker/statistic", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/logtracker-statistic.html")
	})

	http.HandleFunc("/logtracker/logs/fetch", func(w http.ResponseWriter, r *http.Request) {
		log.Println("running /logtracker/logs/fetch handler") // DEBUG
		start := time.Now()                                   // DEBUG
		// w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// w.Header().Set("Access-Control-Allow-Origin", "*")             // Разрешить все источники
		// w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // methods 'PUT', 'PATCH' and 'DELETE' has been deleted
		// w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, hx-request, hx-target, hx-current-url")

		sourceName := r.URL.Query().Get("source_name")
		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")

		log.Println("source_name:", sourceName) // DEBUG

		// build URL for the API request
		apiURL := fmt.Sprintf("%s?source_name=%s&start_date=%s&end_date=%s",
			apiAddrFetch,
			sourceName,
			startDate,
			endDate,
		)

		// send request to the API
		resp, err := http.Get(apiURL)
		if err != nil {
			log.Println("http.Get():", err)
			w.WriteHeader(http.StatusInternalServerError) // ?
			return
		}
		defer resp.Body.Close()

		// read response from the API
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("io.ReadAll():", err)
			w.WriteHeader(http.StatusInternalServerError) // ?
			return
		}

		// write response to the client
		w.Write(body)
		elapsed := time.Since(start)
		log.Println("elapsed:", elapsed)
	})

	// http.HandleFunc("/logtracker/statistic/generalstat", func(w http.ResponseWriter, r *http.Request) {
	http.HandleFunc("/logtracker/statistic/totalStats", func(w http.ResponseWriter, r *http.Request) {
		log.Println("DEGUG: running /logtracker/statistic/totalStats handler") // DEGUG
		w.Write([]byte("DEBUG: /logtracker/statistic/totalStats handler"))     // DEGUG

		resp, err := http.Get(apiAddrStat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("http.Get():", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("io.ReadAll():", err)
		}

		w.Write(body)

	})

	log.Println("Starting server on :8084")

	if err := http.ListenAndServe(":8084", nil); err != nil {
		log.Fatal(err)
	}
}
