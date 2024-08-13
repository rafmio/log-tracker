package main

import (
	"fmt"
	"net/http"
	"strings"
)

func tstHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	servers := r.FormValue("servers")
	fmt.Printf("Type of 'servers': %T\n", servers)
	fmt.Println(servers)

	serversSls := strings.Split(servers, ",")
	for i := 0; i < len(serversSls); i++ {
		w.Write([]byte(serversSls[i]))
		w.Write([]byte("\n"))
	}

	w.Write([]byte(servers))
	w.Write([]byte("\n"))
}

func main() {
	http.HandleFunc("/tst", tstHandler)
	http.ListenAndServe(":8082", nil) // listen on port 8082 and serve requests from handlers defined above.
}
