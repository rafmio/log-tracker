package main

import (
	"api/cmd"
	"log"
)

func main() {
	err := cmd.RunAPI()
	if err != nil {
		log.Printf("Error running API: %v\n", err)
		return
	}
	log.Println("API running successfully")
}
