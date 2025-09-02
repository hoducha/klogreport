package main

import "github.com/fatih/color"

type ColorPalette struct {
	Name   string
	Colors []*color.Color
}

var colorPalettes = map[string]ColorPalette{
	"colorbrewer": {
		Colors: []*color.Color{
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(115)), // rgb(141,211,199)
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(230)), // rgb(255,255,179)
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(146)), // rgb(190,186,218)
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(210)), // rgb(251,128,114)
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(110)), // rgb(128,177,211)
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(216)), // rgb(253,180,98)
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(150)), // rgb(179,222,105)
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(218)), // rgb(252,205,229)
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(188)), // rgb(217,217,217)
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(139)), // rgb(188,128,189)
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(194)), // rgb(204,235,197)
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(228)), // rgb(255,237,111)
		},
	},
	"tableau10": {
		Colors: []*color.Color{
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(68)),  // #729ece
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(215)), // #ff9e4a
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(71)),  // #67bf5c
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(210)), // #ed665d
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(140)), // #ad8bc9
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(137)), // #a8786e
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(212)), // #ed97ca
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(248)), // #a2a2a2
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(186)), // #cdcc5d
			color.New(color.Attribute(38), color.Attribute(5), color.Attribute(116)), // #6dccda
		},
	},
}

func generateOverflowColor(index int, basePaletteSize int) *color.Color {
	// Strategy: Use a spread of colors across the 256-color spectrum
	// avoiding common terminal colors (0-15) and very dark colors
	
	// Color ranges to use for overflow (avoiding palette conflicts)
	colorRanges := [][]int{
		{52, 88},   // Reds/magentas
		{22, 58},   // Greens
		{17, 53},   // Blues
		{130, 166}, // Oranges/browns
		{89, 125},  // Purples
		{59, 95},   // Cyans/teals
		{160, 196}, // More reds
		{70, 106},  // More greens
	}
	
	overflowIndex := index - basePaletteSize
	rangeIndex := overflowIndex % len(colorRanges)
	colorRange := colorRanges[rangeIndex]
	
	// Pick color within the selected range
	rangeSize := colorRange[1] - colorRange[0]
	colorOffset := (overflowIndex / len(colorRanges)) % rangeSize
	colorCode := colorRange[0] + colorOffset
	
	// Ensure we don't go below 52 (too dark) or above 231 (too bright)
	if colorCode < 52 {
		colorCode = 52 + (colorCode % 20)
	}
	if colorCode > 231 {
		colorCode = 52 + ((colorCode - 231) % 20)
	}
	
	return color.New(color.Attribute(38), color.Attribute(5), color.Attribute(colorCode))
}