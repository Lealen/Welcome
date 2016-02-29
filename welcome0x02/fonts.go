package main

import (
	"image/color"

	"github.com/Lealen/engi"
)

type fontAddr struct {
	URL  string
	Size float64
	BG,
	FG color.Color
}

var fonts = make(map[fontAddr]*engi.Font)

func GetFont(URL string, Size float64, BG, FG color.Color) *engi.Font {
	fontaddr := fontAddr{URL: URL, Size: Size, BG: BG, FG: FG}
	if font, ok := fonts[fontaddr]; ok {
		return font
	}
	font := &engi.Font{URL: URL, Size: Size, BG: BG, FG: FG}
	err := font.CreatePreloaded()
	if err != nil {
		panic(err)
		//fmt.Println(err)
	}
	fonts[fontaddr] = font
	return font
}

func DeleteFont(font *engi.Font) {
	delete(fonts, fontAddr{URL: font.URL, Size: font.Size, BG: font.BG, FG: font.FG})
}
