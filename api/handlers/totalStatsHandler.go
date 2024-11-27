package handlers

import (
	"fmt"
	"net/http"
	"sync"
)

func TotalStatsHandler(w http.ResponseWriter, r *http.Request) {
	err := processURL(w, r)
	if err != nil {
		handleProcessURLError(w, err)
		return
	}

	dbConfigFilePath := "config/db-config.json"

	// get server names from the config file r further iteration and
	// initialization of serverStats structure
	serverNames, err := getServerNames(dbConfigFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check if no server names are found in the config file
	if len(serverNames) == 0 {
		http.Error(w, "No server names found in config", http.StatusInternalServerError)
		return
	}

	resultChan := make(chan *serverStats, len(serverNames)) // channel to send serverStats structs
	errChan := make(chan error, len(serverNames))           // channel to send errors

	var wg sync.WaitGroup // wait group to wait for all goroutines to finish

	// create goroutines to process server statistics and send them to the resultChan and errChan
	for _, srvNm := range serverNames {
		wg.Add(1)
		go processServerStats(srvNm, dbConfigFilePath, resultChan, errChan, &wg)
	}

	// wait for all goroutines to finish and close the result and error channels
	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	var serverStatList []*serverStats // list to store serverStats structs

	// receive serverStats structs from the resultChan and handle errors from the errChan
	for i := 0; i < len(serverNames); i++ {
		select {
		case result := <-resultChan:
			serverStatList = append(serverStatList, result)
		case err := <-errChan:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = writeServerStatsJSON(w, serverStatList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// debugPrint(serverStatList, w)
}

func processServerStats(srvNm, dbConfigFilePath string, resultChan chan<- *serverStats, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	srv := new(serverStats)
	srv.serverName = srvNm

	srv.initializeSQLqueries()
	err := srv.setConfigAndConnectDB(dbConfigFilePath)
	if err != nil {
		errChan <- fmt.Errorf("error connecting to database for server %s: %v", srvNm, err)
		return
	}
	defer srv.DB.Close()

	err = srv.makeQueries()
	if err != nil {
		errChan <- fmt.Errorf("error making queries for server %s: %v", srvNm, err)
		return
	}

	srv.readyStatIndicators, err = extractStatIndicators(srv.statIndicatorsRows)
	if err != nil {
		errChan <- fmt.Errorf("error extracting statistics for server %s: %v", srvNm, err)
		return
	}

	err = srv.readyStatIndicators.fillIpAndCountryMap()
	if err != nil {
		errChan <- fmt.Errorf("error filling IP and Country map for server %s: %v", srvNm, err)
		return
	}

	resultChan <- srv
}

func handleProcessURLError(w http.ResponseWriter, err error) {
	switch err {
	case http.ErrNotSupported:
		http.Error(w, "only GET requests are supported", http.StatusMethodNotAllowed)
	case ErrParsingForm:
		http.Error(w, "error parsing form", http.StatusBadRequest)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// The debugPrint function remains unchanged

// func debugPrint(serverStatsList []*serverStats, w http.ResponseWriter) {
// 	fmt.Println("-------------------")
// 	for _, srv := range serverStatsList {
// 		fmt.Fprintf(w, "<h1>Server: %s ", srv.serverName)
// 		// fmt.Fprintf(w, "Host: %s</h1>", srv.Host)
// 		// fmt.Fprintf(w, "<p>DSN: %s</p>", srv.Dsn)
// 		fmt.Fprintf(w, "<h3>Total Stats:</h3>")
// 		fmt.Fprintf(w, "<p>Total records: %d</p>", srv.readyStatIndicators.TotalRecords)
// 		fmt.Fprintf(w, "<p>Unique IP Count: %d</p>", srv.readyStatIndicators.UniqueIPCount)
// 		fmt.Fprintf(w, "<p>Records per day: %.2f</p>", srv.readyStatIndicators.RecordsPerDay)
// 		fmt.Fprintf(w, "<h4>Top 10 IPs:</h4>")
// 		for ip, num := range srv.readyStatIndicators.TopTenIPs {
// 			fmt.Fprintf(w, "<p>%s : %d [%s]</p>", ip, num, srv.readyStatIndicators.MapIpCountries[ip])
// 		}
// 		fmt.Fprintf(w, "<h4>Top 10 Department ports:</h4>")
// 		for dpt, num := range srv.readyStatIndicators.TopTenDpt {
// 			fmt.Fprintf(w, "<p>%s : %d</p>", dpt, num)
// 		}

// 	}

// }
