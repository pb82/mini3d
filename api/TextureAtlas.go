package api

import "image/color"

type TextureAtlas interface {
	W() int
	H() int
	ColorAt(x, y int) color.Color
}
