package main

import (
	"github.com/faiface/pixel/imdraw"
)

type button struct {
	label         string
	width, height int
	onClick       func()
}

func (b *button) Render(imd *imdraw.IMDraw) {

}
