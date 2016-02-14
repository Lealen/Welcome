package main

import (
	"encoding/json"
	"log"
)

var config struct {
	Fullscreen,
	VSync bool
	FPSLimit        int
	KeepAspectRatio bool
}

func loadDefaultConfig() {
	config.Fullscreen = false
	config.VSync = false
	config.FPSLimit = 120
	config.KeepAspectRatio = true
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
		if err.Error() == "open config.ini: no such file or directory" {
			saveConfig()
			return
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
