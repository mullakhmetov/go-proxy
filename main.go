package main

import (
	"log"
	"os"
)

func main() {
	configFile, err := os.Open("settings.json")
	if err != nil {
		log.Fatalf("Error during opening config file:", err.Error())
	}
	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf(err.Error())
	}
	StartProxy(config)
}
