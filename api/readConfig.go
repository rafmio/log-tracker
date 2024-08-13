package main

import (
	"encoding/json"
	"log"
	"os"
)

type Server struct {
	Name      string `json:"Name"`
	Host      string `json:"Host"`
	Port      string `json:"Port"`
	DBName    string `json:"DBName"`
	User      string `json:"User"`
	TableName string `json:"TableName"`
	Password  string `json:"Password"`
	SslMode   string `json:"SslMode"`
}

func readConfig(fileName string) map[string]Server {
	// reading file with configuration for DB connection
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("Opening config file:", err)
	}

	// unmarshalling JSON data to struct
	var servers map[string]Server
	err = json.Unmarshal(file, &servers)
	if err != nil {
		log.Println("Unmarshalling JSON:", err)
	}

	return servers
}
