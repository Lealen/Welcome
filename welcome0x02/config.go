package main

import (
	"encoding/json"
	"log"
	"os"
)

var config struct {
	Width,
	Height int
	Fullscreen,
	VSync bool
	FPSLimit int
}

func loadDefaultConfig() {
	config.Width = 1366
	config.Height = 768
	config.Fullscreen = false
	config.VSync = false
	config.FPSLimit = 120
}

func saveConfig() {
	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		log.Println(err)
		return
	}
	saveToFile("config.ini", string(data))
	if err != nil {
		log.Println(err)
		return
	}
}

func loadConfig() {
	loadDefaultConfig()
	data, err := loadFromFile("config.ini")
	if err != nil {
		if os.IsNotExist(err) {
			saveConfig()
		}
		log.Println(err)
		return
	}
	err = json.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Println(err)
		return
	}
}
