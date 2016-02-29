package main

import (
	"fmt"
	"image"
	"image/color"

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

var texturecolors = make(map[color.Color]*engi.Texture)

func GetTextureColor(c color.Color) *engi.Texture {
	if texture, ok := texturecolors[c]; ok {
		return texture
	}
	texture := GetSingleUseTextureColor(c)
	texturecolors[c] = texture
	return texture
}

func GetSingleUseTextureColor(c color.Color) *engi.Texture {
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, c)
	return engi.NewTexture(engi.NewImageObject(img))
}

func DeleteTextureColor(c color.Color) {
	delete(texturecolors, c)
}
