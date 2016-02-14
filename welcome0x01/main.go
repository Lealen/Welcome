package main

import (
	"fmt"

	"github.com/Lealen/engi"
)

const width = 640
const height = 480

//GameName
const GameName = "Welcome 0x01"

func main() {
	loadConfig()

	fmt.Println(config)

	engi.SetFPSLimit(config.FPSLimit)

	options := engi.RunOptions{
		Title:                  GameName,
		Fullscreen:             config.Fullscreen,
		Width:                  width,
		Height:                 height,
		VSync:                  config.VSync,
		ScaleOnResize:          false, //true,
		KeepAspectRatioOnScale: config.KeepAspectRatio,
	}

	//engi.Open(options, &Game{})
	engi.Open(options, &Menu{})
}
