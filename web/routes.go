package cmd

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	// registration of static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// registration of dynamic routes
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/about", handleAbout)
	mux.HandleFunc("/logtracker", handleLogtracker)
	mux.HandleFunc("/budgeting", handleBudgeting)
	mux.HandleFunc("/orders", handleOrders)
	mux.HandleFunc("/logtracker/logs", handleLogLogs)
	mux.HandleFunc("/logtracker/statistic", handleLogStatistic)
	mux.HandleFunc("/logtracker/logs/fetch", handleLogFetch)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/about.html")
}

func handleLogtracker(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/logtracker.html")
}

func handleBudgeting(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/budgeting.html")
}

func handleOrders(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/orders.html")
}

func handleLogLogs(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/logtracker-logs.html")
}

func handleLogStatistic(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/logtracker-statistic.html")
}

func handleLogFetch(w http.ResponseWriter, r *http.Request) {
		start := time.Now()                                   // DEBUG

		sourceName := r.URL.Query().Get("source_name")
		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")

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
