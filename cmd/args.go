package main

import (
	"flag"
)

var (
	filePath string
)

func init() {
	flag.StringVar(&filePath, "config", "config.json", "Path to the configuration file")
	flag.Parse()
}
