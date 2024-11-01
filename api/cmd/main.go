package main

import "log"

func main() {
	err := runAPI()
	if err != nil {
		log.Printf("Error running API: %v\n", err)
		return
	}
	log.Println("API running successfully")
}
