package config

import (
	"encoding/json"
	"io"
	"log"
)

// Revisit to implement io reader
func ReadConfig(r io.Reader) Configuration {
	log.Println("Loading in config file")

	decoder := json.NewDecoder(r)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)

	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Loaded in config file")
	return configuration
}
