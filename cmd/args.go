package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

var (
	filePath string
	debug    bool
)

func init() {
	flag.StringVar(&filePath, "config", "config.json", "Path to the configuration file")
	flag.BoolVar(&debug, "debug", false, "Enable debug output")
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug output is enabled")
	}
}
