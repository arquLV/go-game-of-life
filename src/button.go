package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type button struct {
	label         string
	width, height int
	onClick       func()
}

func (b *button) Render(imd *imdraw.IMDraw) {
	imd.Color = color.RGBA{0x00, 0xff, 0x00, 0xff}
	imd.Push(pixel.V(200, 200), pixel.V(300, 240))
	imd.Rectangle(0)
}
