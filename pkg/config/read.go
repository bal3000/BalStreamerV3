package config

import (
	"encoding/json"
	"log"
	"os"
)

// Revisit to implement io reader

func ReadConfig() Configuration {
	file, _ := os.Open("./config.json")
	defer file.Close()
	log.Println("Loading in config file")

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)

	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Loaded in config file")
	return configuration
}
