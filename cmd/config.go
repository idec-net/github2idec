package main

import (
	"io/ioutil"
	"log"

	"encoding/json"

	"github.com/idec-net/github2idec/g2i"
)

// loadConfig loads bot configuration from file
func loadConfig(path string) *g2i.Config {
	log.Print("Loading configuration")
	config := &g2i.Config{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}

	return config
}
