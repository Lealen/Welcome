package main

import (
	"fmt"

	"github.com/Lealen/engi"
)

func loadTexture(texturename string) (texture *engi.Texture) {
	texture = engi.Files.Image(texturename)
	if texture == nil {
		fmt.Println("Failed to load", texturename)
		texture = engi.Files.Image("notfound.png")
	}
	if texture == nil {
		panic("I wasn't been able to use texture named " + texturename + " nor notfound.png")
	}
	return
}
