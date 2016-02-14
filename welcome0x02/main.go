package main

import (
	"fmt"

	"github.com/Lealen/engi"
)

const width = 1920
const height = 1080

//GameName
const GameName = "Welcome 0x02"

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
		ScaleOnResize:          false,
		KeepAspectRatioOnScale: false,
	}

	//engi.Open(options, &Game{})
	//engi.Open(options, &Menu{})
	engi.Open(options, &Gui{})
}
